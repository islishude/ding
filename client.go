// Package ding 钉钉机器人 API 实现
// https://ding-doc.dingtalk.com/doc#/serverapi2/qf2nxq

package ding

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// New 创建新的实例，hmkey 为 HMAC key 如没有可以不填
func New(client *http.Client, token string, hmkey ...string) Ding {
	const api = "https://oapi.dingtalk.com/robot/send?access_token="
	if client == nil {
		client = http.DefaultClient
	}

	var secret string
	if len(hmkey) > 0 {
		secret = hmkey[0]
	}

	return &clientimpl{
		endurl: api + token,
		hmkey:  secret,
		bhmkey: []byte(secret),
		client: client,
		now: func() string {
			timestamp := time.Now().UnixNano() / 1e6
			return strconv.FormatInt(timestamp, 10)
		},
	}
}

type clientimpl struct {
	endurl string
	hmkey  string
	bhmkey []byte
	client *http.Client
	now    func() string // 获取当前时间(毫秒)
}

func (d *clientimpl) request(ctx context.Context, reqdata map[string]interface{}) error {
	var reqbuf = bytes.NewBuffer(nil)
	if err := json.NewEncoder(reqbuf).Encode(reqdata); err != nil {
		return err
	}

	var endurl = d.endurl
	if d.hmkey != "" {
		endurl += d.getSign()
	}

	// 需要 1.13 以上版本支持
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endurl, reqbuf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var respdata struct {
		Code int    `json:"errcode"`
		Msg  string `json:"errmsg"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&respdata); err != nil {
		return err
	}

	if respdata.Code != 0 {
		return errors.New(respdata.Msg)
	}

	return nil
}

func (d *clientimpl) getSign() string {
	var timestamp = d.now()

	hash := hmac.New(sha256.New, d.bhmkey)
	dataToSign := []byte(fmt.Sprintf("%s\n%s", timestamp, d.hmkey))
	_, _ = hash.Write(dataToSign)

	// 签名后需要使用进行网址编码
	sign := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	return fmt.Sprintf("&timestamp=%s&sign=%s", timestamp, url.QueryEscape(sign))
}
