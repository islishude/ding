// Package ding 钉钉机器人 API 实现
// https://ding-doc.dingtalk.com/doc#/serverapi2/qf2nxq

package ding

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func defaultNow() string {
	timestamp := time.Now().UnixNano() / 1e6
	return strconv.FormatInt(timestamp, 10)
}

var defaultClient = &http.Client{Timeout: time.Second * 5}

const webhook = "https://oapi.dingtalk.com/robot/send"

// New 创建新的实例，hmkey 为 HMAC key 如没有可以不填
func New(client *http.Client, token string, hmkey ...string) Ding {
	if client == nil {
		client = defaultClient
	}

	var secret string
	if len(hmkey) > 0 {
		secret = hmkey[0]
	}

	return &clientimpl{
		api:    webhook,
		tokens: []AccessToken{{Token: token, Key: secret}},
		client: client,
		now:    defaultNow,
	}
}

// Multi 创建包含访问令牌的实例，每次请求轮换使用下一个访问令牌，防止限流造成的发送失败
func Multi(tokens []AccessToken, client ...*http.Client) Ding {
	impl := &clientimpl{
		api:    webhook,
		tokens: tokens,
		now:    defaultNow,
	}
	if len(client) > 0 && client[0] != nil {
		impl.client = client[0]
	} else {
		impl.client = defaultClient
	}
	return impl
}

type clientimpl struct {
	client *http.Client
	now    func() string // 获取当前时间(毫秒)

	api    string
	tokens []AccessToken

	index int // current access tokens index
	mux   sync.Mutex
}

func (d *clientimpl) request(ctx context.Context, reqdata map[string]interface{}) error {
	var reqbuf = bytes.NewBuffer(nil)
	if err := json.NewEncoder(reqbuf).Encode(reqdata); err != nil {
		return err
	}

	// 需要 1.13 以上版本支持
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.api, reqbuf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()

	accessToken, ok := d.nextAccessToken()
	if !ok {
		return errors.New("no access token")
	}

	q.Set("access_token", accessToken.Token)
	if accessToken.Key != "" {
		timestamp := d.now()
		q.Set("timestamp", timestamp)
		q.Set("sign", getsign(accessToken.Key, timestamp))
	}
	req.URL.RawQuery = q.Encode()

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

// nextAccessToken 获取下一个可用的访问令牌
func (d *clientimpl) nextAccessToken() (AccessToken, bool) {
	if len(d.tokens) == 0 {
		return AccessToken{}, false
	}

	if len(d.tokens) == 1 {
		return d.tokens[0], true
	}

	d.mux.Lock()
	defer d.mux.Unlock()

	if i := d.index; i < len(d.tokens) {
		d.index++
		return d.tokens[i], true
	}

	d.index = 1
	return d.tokens[0], true
}
