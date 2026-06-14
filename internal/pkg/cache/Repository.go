package cache

import (
	"context"

	"github.com/shbhom/urlShortner/internal/models"
)

type Repository interface {
	Get(ctx context.Context, shortCode string) (string, error)
	Set(ctx context.Context, data models.UrlData) error
	RecordInvokation(ctx context.Context, code string) error
	Rename(ctx context.Context, oldKey, newKey string) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	Delete(ctx context.Context, key string) error
}
