package services

import (
	"github.com/shbhom/urlShortner/internal/db"
	"github.com/shbhom/urlShortner/internal/pkg/url"
)

type Service struct {
	url url.Repository
}

func NewService(dbUrl string) *Service {
	db := db.NewPostgres(dbUrl)
	return &Service{
		url: db,
	}
}
