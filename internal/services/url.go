package services

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/redis/go-redis/v9"
	"github.com/shbhom/urlShortner/internal/models"
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

func (svc *Service) Addurl(ctx context.Context, url string) (string, error) {
	key, err := svc.generateShortCode(8)
	if err != nil {
		return "", err
	}
	if err := svc.url.AddUrl(ctx, models.UrlData{ShortCode: key, TargetUrl: url}); err != nil {
		return "", err
	}
	return key, nil
}

func (svc *Service) GetUrl(ctx context.Context, code string) (string, error) {
	cacheUrl, err := svc.cache.Get(ctx, code)
	if err == nil {
		slog.Info(fmt.Sprintf("Cache hit for code: %s", code))
		return cacheUrl, nil
	} else if !errors.Is(err, redis.Nil) {
		slog.Error("Redis cache error", "error", err)
	} else {
		slog.Warn(fmt.Sprintf("No cache found for code: %s", code))
	}

	//url not found in cache searching DB
	DBUrl, err := svc.url.GetUrlByCode(ctx, code)
	if err != nil {
		return "", err
	}
	
	//After fetching from DB setting url in redis cache
	if err := svc.cache.Set(ctx, models.UrlData{ShortCode: code, TargetUrl: DBUrl}); err != nil {
		slog.WarnContext(ctx, fmt.Sprintf("Error while setting TargetUrl %s for code %s", DBUrl, code))
	}
	
	return DBUrl, nil
}
