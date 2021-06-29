package ding

import "regexp"

var tokenRegExp = regexp.MustCompile(`^[a-z0-9]{64}$`)

func ValidateToken(token string) bool {
	return tokenRegExp.MatchString(token)
}

var secretKeyRegExp = regexp.MustCompile(`^SEC[a-z0-9]{64}$`)

func ValidateSecretKey(key string) bool {
	return secretKeyRegExp.MatchString(key)
}
