package services

import (
	"time"

	"github.com/shbhom/urlShortner/internal/db/postgres"
	"github.com/shbhom/urlShortner/internal/db/redis"
	"github.com/shbhom/urlShortner/internal/pkg/cache"
	"github.com/shbhom/urlShortner/internal/pkg/metrics"
	"github.com/shbhom/urlShortner/internal/pkg/url"
)

type Service struct {
	url   url.Repository
	cache cache.Repository
}

func NewService(dbUrl, redisAddr string, ttl int) *Service {
	db := postgres.NewPostgres(dbUrl)
	cacheTTL := time.Duration(ttl * int(time.Minute))
	redis := redis.NewCache(redisAddr, cacheTTL)
	metrics.RegisterDBStatsCollector(db.Client)
	metrics.RegisterRedisStatsCollector(redis.Client)
	return &Service{
		url:   db,
		cache: redis,
	}
}
