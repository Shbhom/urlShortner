package url

import (
	"context"

	"github.com/shbhom/urlShortner/internal/models"
)

type Repository interface {
	GetUrlByCode(ctx context.Context, short_code string) (string, error)
	AddUrl(ctx context.Context, data models.UrlData) error
}
