package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/shbhom/urlShortner/internal/config"
	"github.com/shbhom/urlShortner/internal/db/postgres"
	"github.com/shbhom/urlShortner/internal/db/redis"
	"github.com/shbhom/urlShortner/internal/pkg/metrics"
	"github.com/shbhom/urlShortner/internal/services"
)

func Run(envType string) {
	config := config.LoadConfig(envType)
	r := mux.NewRouter()

	// Initialize DB and Cache
	db := postgres.NewPostgres(config.DB_URL)
	cacheTTL := time.Duration(config.URL_TTL * int(time.Minute))
	redisCache := redis.NewCache(config.REDIS_ADDR, cacheTTL)

	// Register Metrics
	metrics.RegisterDBStatsCollector(db.Client)
	metrics.RegisterRedisStatsCollector(redisCache.Client)

	svc := services.NewService(config.SHORT_CODE_MIN_LEN, db, redisCache)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 2. Set up WaitGroup for the background worker
	var wg sync.WaitGroup
	wg.Add(1)

	// 3. Pass the cancellable context and WaitGroup
	svc.StartAnalyticsWorker(ctx, &wg)
	serve := &Server{
		Router:   r,
		Config:   config,
		Services: svc,
	}
	serve.routes()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.API_PORT),
		Handler: serve,
	}

	go func() {
		slog.Info(fmt.Sprintf("Server starting on port %d", config.API_PORT))
		if config.API_PORT == 443 {
			if err := srv.ListenAndServeTLS(config.ChainCertPath, config.PemCertPath); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("listen: %s\n", err)
			}
		} else {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("listen: %s\n", err)
			}
		}
	}()
	<-ctx.Done()
	slog.Info("OS signal received. Shutting down gracefully...")

	// 4. Give active HTTP requests 5 seconds to finish
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	close(svc.AnalyticsChan)

	svc.AnalyticsWG.Wait()
	// 5. Wait for the analytics worker to finish its final flush
	wg.Wait()

	slog.Info("Performing final analytics flush...")
	svc.FlushAnalytics(context.Background())

	slog.Info("Server completely shut down")

}
