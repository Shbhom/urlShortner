package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"github.com/shbhom/urlShortner/internal/models"
)

func (svc *Service) generateShortCode(id uint64) (string, error) {
	return svc.sqid.Encode([]uint64{id})
}

func (svc *Service) Addurl(ctx context.Context, url string) (string, error) {
	nextId, err := svc.url.GetNextSequence(ctx)
	if err != nil {
		return "", fmt.Errorf("Error while getting the counter sequence from postgres: %w", err)
	}
	key, err := svc.generateShortCode(nextId)
	if err != nil {
		return "", err
	}
	if err := svc.url.AddUrl(ctx, models.UrlData{ShortCode: key, TargetUrl: url}); err != nil {
		return "", err
	}
	return key, nil
}

func (svc *Service) GetUrl(ctx context.Context, code string) (string, error) {
	recordhit := func() {
		select {
		case svc.AnalyticsChan <- code:
			// Successfully queued to be processed by a worker
		default:
			// The queue is full; drop the event so we don't block the user's redirect
			slog.Warn("Analytics channel full, dropping hit", "code", code)
		}
	}
	cacheUrl, err := svc.cache.Get(ctx, code)
	if err == nil {
		slog.Debug(fmt.Sprintf("Cache hit for code: %s", code))
		recordhit()
		return cacheUrl, nil
	} else if !errors.Is(err, redis.Nil) {
		slog.Error("Redis cache error", "error", err)
	} else {
		slog.Debug(fmt.Sprintf("No cache found for code: %s", code))
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
	recordhit()
	return DBUrl, nil
}
