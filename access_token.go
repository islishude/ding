package ding

// AccessToken 访问令牌
type AccessToken struct {
	Token string // 令牌
	Key   string // 签名密钥，可为空
}
