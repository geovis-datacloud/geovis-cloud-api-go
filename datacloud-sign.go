package main

type DatacloudSign struct {
	*CloudApiSign
}

func NewDatacloudSign(secretId, secretKey string) (*DatacloudSign, error) {
	sign, err := NewCloudApiSign(secretId, secretKey, "geovis-data-cloud", nil)
	if err != nil {
		return nil, err // 返回错误，而不是忽略
	}
	return &DatacloudSign{CloudApiSign: sign}, nil
}
