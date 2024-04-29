package main

import (
	"fmt"
	"os"
	"testing"
)

var sdk *DatacloudChannelSdk

func TestMain(m *testing.M) {
	// Setup
	secretId := "your secret id"
	secretKey := "your secret key"
	channelId := "your channel id"
	cloudApiHost := "https://api.geovisearth.com"
	dataCloudHost := "https://datacloud.geovisearth.com"

	sdk, _ = NewDatacloudChannelSdk(secretId, secretKey, channelId, cloudApiHost, dataCloudHost)

	// Run the tests
	code := m.Run()

	// Teardown

	os.Exit(code)
}

func TestCreateCertificationSign(t *testing.T) {
	// Setup
	secretId := "your secret id"
	secretKey := "your secret key"
	sign, _ := NewCertificationSign(secretId, secretKey)

	// Run the tests
	headers := sign.Sign(map[string]string{
		"path":   "/datacloud/auth/phone",
		"method": "GET",
	})

	fmt.Println(headers)
}

func TestGetTokenByPhone(t *testing.T) {
	phoneNumber := "19965407629"
	ret, err := sdk.getTokenByPhone(phoneNumber)
	if err != nil {
		t.Errorf("getTokenByPhone returned an error: %v", err)
	}
	fmt.Println(ret)
}

func TestRefreshToken(t *testing.T) {
	phoneNumber := "19965407629"
	appId := "YZJQ4ggozLp0Mhuy"
	ret, err := sdk.refreshToken(phoneNumber, appId)
	if err != nil {
		t.Errorf("refreshToken returned an error: %v", err)
	}
	fmt.Println(ret)
}

func TestGetUsage(t *testing.T) {
	appId := "YZJQ4ggozLp0Mhuy"
	ret, err := sdk.getUsage(appId)
	if err != nil {
		t.Errorf("getUsage returned an error: %v", err)
	}
	fmt.Println(ret)
}

func TestGetApplicationList(t *testing.T) {
	phoneNumber := "19965407629"
	ret, err := sdk.getApplicationList(phoneNumber)
	if err != nil {
		t.Errorf("getApplicationList returned an error: %v", err)
	}
	fmt.Println(ret)
}
