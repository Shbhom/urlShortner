package services

import (
	"crypto/rand"
	"math/big"
	"net/url"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (svc *Service) generateShortCode(length int) (string, error) {
	result := make([]byte, length)

	for i := range result {
		n, err := rand.Int(
			rand.Reader,
			big.NewInt(int64(len(alphabet))),
		)

		if err != nil {
			return "", err
		}

		result[i] = alphabet[n.Int64()]
	}

	return string(result), nil
}

func (svc *Service) Addurl(url string) (string, error) {
	key, err := svc.generateShortCode(8)
	if err != nil {
		return "", err
	}
	if err := svc.url.AddUrl(url, key); err != nil {
		return "", err
	}
	return key, nil
}

func (svc *Service) GetUrl(code string) (*url.URL, error) {
	tUrl, err := svc.url.GetUrlByCode(code)
	if err != nil {
		return nil, err
	}
	pUrl, err := url.Parse(tUrl)
	if err != nil {
		return nil, err
	}
	return pUrl, err
}
