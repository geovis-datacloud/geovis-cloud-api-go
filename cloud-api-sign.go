package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type CloudApiSign struct {
	secretId    string
	secretKey   []byte
	serviceName string
	options     map[string]string
}

func NewCloudApiSign(secretId, secretKey, serviceName string, options map[string]string) (*CloudApiSign, error) {
	key, err := base64.StdEncoding.DecodeString(secretKey)
	if err != nil {
		return nil, err // 返回错误，而不是忽略
	}
	return &CloudApiSign{
		secretId:    secretId,
		secretKey:   key,
		serviceName: serviceName,
		options:     options,
	}, nil
}

func (c *CloudApiSign) Sign(params map[string]string) map[string]string {
	nonceStr := randStringBytes(21)
	timestamp := time.Now().Unix() * 1000
	timeByte := strconv.AppendInt(nil, timestamp, 10)
	var b bytes.Buffer
	b.WriteString(c.serviceName)
	b.WriteString("\n")
	b.WriteString(params["method"])
	b.WriteString("\n")
	b.WriteString(params["path"])
	b.WriteString("\n")
	b.WriteString(params["queryString"])
	b.WriteString("\n")
	b.WriteString(params["body"])
	b.WriteString("\n")
	b.WriteString(nonceStr)
	b.WriteString("\n")
	b.Write(timeByte)
	toSign := b.Bytes()
	signature := c.createSign(toSign)

	var builder strings.Builder
	// 使用 WriteString 方法拼接字符串
	builder.WriteString("secretId=")
	builder.WriteString(c.secretId)
	builder.WriteString(",")
	builder.WriteString("nonceStr=")
	builder.WriteString(nonceStr)
	builder.WriteString(",")
	builder.WriteString("service=")
	builder.WriteString(c.serviceName)
	builder.WriteString(",")
	builder.WriteString("timestamp=")
	builder.Write(timeByte)
	builder.WriteString(",")
	builder.WriteString("signature=")
	builder.WriteString(signature)

	if c.options != nil && len(c.options) > 0 {
		for k, v := range c.options {
			builder.WriteString(k)
			builder.WriteString("=")
			builder.WriteString(v)
		}
	}

	headers := make(map[string]string)
	headers["authorization"] = builder.String()

	return headers
}

func (c *CloudApiSign) createSign(data []byte) string {
	h := hmac.New(sha256.New, c.secretKey)
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(52)]
	}
	return string(b)
}
