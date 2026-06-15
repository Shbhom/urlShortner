package services

import (
	"context"
	"log/slog"

	"github.com/robfig/cron/v3"
)

func (s *Service) StartCronJobs(ctx context.Context) {
	c := cron.New()
	c.AddFunc("@daily", func() {
		slog.Info("Running daily URL purge...")
		count, err := s.url.PurgeOldURLs(context.Background())
		if err != nil {
			slog.Error("Failed to purge old URLs", "error", err)
			return
		}
		slog.Info("Successfully purged old URLs", "count", count)
	})
	c.Start()

	// Wait for context cancellation to gracefully stop the cron scheduler
	go func() {
		<-ctx.Done()
		c.Stop()
	}()
}
