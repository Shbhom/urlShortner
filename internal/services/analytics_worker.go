package services

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/shbhom/urlShortner/internal/db/redis"
)

func (svc *Service) StartAnalyticsWorker(ctx context.Context, wg *sync.WaitGroup) {
	ticker := time.NewTicker(60 * time.Second)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				slog.Info("Shutting Down Analytics worker")
				ticker.Stop()
				return
			case <-ticker.C:
				svc.FlushAnalytics(ctx)
			}
		}
	}()
}

func (svc *Service) FlushAnalytics(ctx context.Context) {
	newKey := fmt.Sprintf(redis.ANALYTICS_NEW_KEY, time.Now().String())
	if err := svc.cache.Rename(ctx, redis.ANALYTICS_KEY, newKey); err != nil {
		if err.Error() == "ERR no such key" {
			slog.Warn("no analytics hash set exists!!")
			return
		}
		slog.Error("Error while renaming analytics batch", "error", err)
		return
	}
	data, err := svc.cache.HGetAll(ctx, newKey)
	if err != nil {
		slog.Error("Error while fetching data from analytics batch", "error", err)
		return
	}
	if err := svc.url.BulkUpdateUrlLastInvokation(ctx, data); err != nil {
		slog.Error("Error while ingesting Analytics data from redis to postgres", "error", err)
		return
	}
	if err := svc.cache.Delete(ctx, newKey); err != nil {
		slog.Error("Error while flushing data from redis", "error", err)
		return
	}
	slog.Info("Successfully flushed analytics batch", "count", len(data))
}

func (svc *Service) invocationWorker() {
	defer svc.AnalyticsWG.Done()
	for code := range svc.AnalyticsChan {
		err := svc.cache.RecordInvokation(context.Background(), code)
		if err != nil {
			slog.Error("Failed to record analytics batch", "code", code, "error", err)
		}
	}
}
