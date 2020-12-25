package ding

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
)

// getsign 获取消息请求签名
func getsign(key, timestamp string) string {
	hasher := hmac.New(sha256.New, []byte(key))
	_, _ = hasher.Write([]byte(fmt.Sprintf("%s\n%s", timestamp, key)))

	// 签名后需要使用进行网址编码
	sign := base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	return url.QueryEscape(sign)
}

// AccessToken 访问令牌
type AccessToken struct {
	Token string `json:"token" yaml:"token"`                 // 令牌值
	Key   string `json:"key,omitempty" yaml:"key,omitempty"` // 签名密钥，可为空
}
