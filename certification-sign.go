package main

type CertificationSign struct {
	*CloudApiSign
}

func NewCertificationSign(secretId, secretKey string) (*CertificationSign, error) {
	sign, err := NewCloudApiSign(secretId, secretKey, "geovis-certification", nil)
	if err != nil {
		return nil, err // 返回错误，而不是忽略
	}
	return &CertificationSign{CloudApiSign: sign}, nil
}
