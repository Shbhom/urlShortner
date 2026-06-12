package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shbhom/urlShortner/internal/models"
	"github.com/shbhom/urlShortner/internal/pkg/metrics"
)

type Cache struct {
	Client *redis.Client
	TTL    time.Duration
}

func NewCache(redisAddr string, urlTTL time.Duration) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	ctx := context.Background()
	if cmd := client.Ping(ctx); cmd.Err() != nil {
		log.Fatal("Unable to ping, redis client")
	}
	return &Cache{
		Client: client,
		TTL:    urlTTL,
	}
}

func (c *Cache) Get(ctx context.Context, shortCode string) (string, error) {
	start := time.Now()
	res, err := c.Client.Get(ctx, fmt.Sprintf("url:%s", shortCode)).Result()
	metrics.RedisGetDurationSeconds.Observe(time.Since(start).Seconds())
	switch err {
	case redis.Nil:
		metrics.RedisMissesTotal.Inc()
	case nil:
		metrics.RedisHitsTotal.Inc()
	}
	return res, err
}

func (c *Cache) Set(ctx context.Context, data models.UrlData) error {
	start := time.Now()
	err := c.Client.Set(ctx, fmt.Sprintf("url:%s", data.ShortCode), data.TargetUrl, c.TTL).Err()
	metrics.RedisSetDurationSeconds.Observe(time.Since(start).Seconds())
	return err
}
