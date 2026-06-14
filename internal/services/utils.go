package services

import (
	"encoding/json"
	"io"
)

func (svc *Service) ParseBody(req_body io.Reader, v any) error {
	err := json.NewDecoder(req_body).Decode(v)
	if err != nil {
		return err
	}
	return svc.validator.Struct(v)
}
