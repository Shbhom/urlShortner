package cache

import (
	"context"

	"github.com/shbhom/urlShortner/internal/models"
)

type Repository interface {
	Get(ctx context.Context, shortCode string) (string, error)
	Set(ctx context.Context, data models.UrlData) error
}
