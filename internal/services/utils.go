package services

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
)

func (svc *Service) ParseBody(req_body io.Reader, v any) error {
	err := json.NewDecoder(req_body).Decode(v)
	if err != nil {
		return err
	}
	return validator.New().Struct(v)
}
