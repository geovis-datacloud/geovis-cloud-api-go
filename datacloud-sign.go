package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type CloudApiSign struct {
	secretId    string
	secretKey   []byte
	serviceName string
	options     map[string]string
}

func NewCloudApiSign(secretId, secretKey, serviceName string, options map[string]string) *CloudApiSign {
	key, _ := base64.StdEncoding.DecodeString(secretKey)
	return &CloudApiSign{
		secretId:    secretId,
		secretKey:   key,
		serviceName: serviceName,
		options:     options,
	}
}

func (c *CloudApiSign) Sign(params map[string]string) map[string]string {
	nonceStr := randStringBytes(21)
	timestamp := time.Now().Unix() * 1000

	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%d", c.serviceName, params["method"], params["path"], params["queryString"], params["body"], nonceStr, timestamp)
	signature := c.createSign(canonicalRequest)

	authorization := fmt.Sprintf("secretId=%s,nonceStr=%s,service=%s,timestamp=%d,signature=%s", c.secretId, nonceStr, c.serviceName, timestamp, signature)

	if len(c.options) > 0 {
		var extendKeys []string
		for k, v := range c.options {
			extendKeys = append(extendKeys, fmt.Sprintf("%s=%s", k, v))
		}
		authorization = fmt.Sprintf("%s,%s", authorization, strings.Join(extendKeys, ","))
	}

	headers := make(map[string]string)
	headers["authorization"] = authorization

	return headers
}

func (c *CloudApiSign) createSign(data string) string {
	h := hmac.New(sha256.New, []byte(c.secretKey))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func randStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
