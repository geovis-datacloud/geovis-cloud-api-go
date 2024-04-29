package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DatacloudChannelSdk struct {
	secretId      string
	secretKey     string
	channelId     string
	cloudApiHost  string
	dataCloudHost string
	geoSign       *DatacloudSign
	client        *http.Client
}

type RefreshTokenPayload struct {
	AppId   string `json:"appId"`
	Channel string `json:"channel"`
	Mobile  string `json:"mobile"`
}

func NewDatacloudChannelSdk(secretId, secretKey, channelId, cloudApiHost, dataCloudHost string) (*DatacloudChannelSdk, error) {
	if cloudApiHost == "" {
		cloudApiHost = "https://api.geovisearth.com"
	}

	if dataCloudHost == "" {
		dataCloudHost = "https://datacloud.geovisearth.com"
	}

	geoSign, err := NewDatacloudSign(secretId, secretKey)
	if err != nil {
		return nil, err
	}

	return &DatacloudChannelSdk{
		secretId:      secretId,
		secretKey:     secretKey,
		channelId:     channelId,
		cloudApiHost:  cloudApiHost,
		dataCloudHost: dataCloudHost,
		geoSign:       geoSign,
		client:        &http.Client{},
	}, nil
}

func (d *DatacloudChannelSdk) getTokenByPhone(phone string) (map[string]interface{}, error) {
	params := map[string]string{
		"path":        "/datacloud/auth/phone",
		"method":      "GET",
		"time":        fmt.Sprintf("%d", time.Now().Unix()),
		"body":        "",
		"queryString": fmt.Sprintf("channel=%s&phone=%s", d.channelId, phone),
	}
	headers := d.geoSign.Sign(params)
	url := fmt.Sprintf("%s/passport/datacloud/auth/phone?channel=%s&phone=%s", d.cloudApiHost, d.channelId, phone)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("auth", headers["authorization"])
	req.Header.Set("referer", "https://datacloud.geovisearth.com")

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *DatacloudChannelSdk) refreshToken(phone, appId string) (map[string]interface{}, error) {
	payload := RefreshTokenPayload{
		AppId:   appId,
		Channel: d.channelId,
		Mobile:  phone,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	params := map[string]string{
		"path":        fmt.Sprintf("/v1/cloudapi/apps/public/%s/refresh", appId),
		"method":      "POST",
		"time":        fmt.Sprintf("%d", time.Now().Unix()),
		"body":        string(payloadBytes),
		"queryString": "",
	}
	headers := d.geoSign.Sign(params)
	url := fmt.Sprintf("%s/v1/cloudapi/apps/public/%s/refresh", d.dataCloudHost, appId)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("authorization", headers["authorization"])
	req.Header.Set("referer", "https://datacloud.geovisearth.com")
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *DatacloudChannelSdk) getUsage(appId string) (map[string]interface{}, error) {
	params := map[string]string{
		"path":        "/v1/cloudapi/developer/devDataPackUsage",
		"method":      "GET",
		"time":        fmt.Sprintf("%d", time.Now().Unix()),
		"body":        "",
		"queryString": fmt.Sprintf("appId=%s&channel=%s", appId, d.channelId),
	}
	headers := d.geoSign.Sign(params)
	url := fmt.Sprintf("%s/v1/cloudapi/developer/devDataPackUsage?appId=%s&channel=%s", d.dataCloudHost, appId, d.channelId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("authorization", headers["authorization"])
	req.Header.Set("referer", "https://datacloud.geovisearth.com")

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *DatacloudChannelSdk) getApplicationList(phone string) (map[string]interface{}, error) {
	params := map[string]string{
		"path":        "/v1/cloudapi/application/myPublicAppList",
		"method":      "GET",
		"time":        fmt.Sprintf("%d", time.Now().Unix()),
		"body":        "",
		"queryString": fmt.Sprintf("channel=%s&mobile=%s", d.channelId, phone),
	}
	headers := d.geoSign.Sign(params)
	url := fmt.Sprintf("%s/v1/cloudapi/application/myPublicAppList?channel=%s&mobile=%s", d.dataCloudHost, d.channelId, phone)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("authorization", headers["authorization"])
	req.Header.Set("referer", "https://datacloud.geovisearth.com")

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
