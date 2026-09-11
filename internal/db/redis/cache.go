package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shbhom/urlShortner/internal/models"
	"github.com/shbhom/urlShortner/internal/pkg/metrics"
	"golang.org/x/exp/slog"
)

type Cache struct {
	Client *redis.Client
	TTL    time.Duration
}

var (
	ANALYTICS_KEY     = "analytics_batch"
	ANALYTICS_NEW_KEY = "analytics_processing:%s"
)

func NewCache(redisAddr string, urlTTL time.Duration) (*Cache, error) {
	var client *redis.Client
	if Opts, err := redis.ParseURL(redisAddr); err != nil {
		return nil, fmt.Errorf("error while parsing redis Opts")
	} else {
		client = redis.NewClient(Opts)
	}
	if client == nil {
		return nil, fmt.Errorf("unable to connect redis")
	}
	slog.Info("Successfully connected to redis", "redis Addr", redisAddr)
	ctx := context.Background()
	if cmd := client.Ping(ctx); cmd.Err() != nil {
		return nil, fmt.Errorf("Unable to ping, redis client")
	}
	return &Cache{
		Client: client,
		TTL:    urlTTL,
	}, nil
}

func (c *Cache) Get(ctx context.Context, shortCode string) (string, error) {
	start := time.Now()
	res, err := c.Client.Get(ctx, fmt.Sprintf("url:%s", shortCode)).Result()
	metrics.RedisGetDurationSeconds.Record(ctx, time.Since(start).Seconds())
	switch err {
	case redis.Nil:
		metrics.RedisMissesTotal.Add(ctx, 1)
	case nil:
		metrics.RedisHitsTotal.Add(ctx, 1)
	}
	return res, err
}

func (c *Cache) Set(ctx context.Context, data models.UrlData) error {
	start := time.Now()
	err := c.Client.Set(ctx, fmt.Sprintf("url:%s", data.ShortCode), data.TargetUrl, c.TTL).Err()
	metrics.RedisSetDurationSeconds.Record(ctx, time.Since(start).Seconds())
	return err
}

func (c *Cache) RecordInvokation(ctx context.Context, code string) error {
	now := time.Now().Unix()
	return c.Client.HSet(ctx, ANALYTICS_KEY, code, now).Err()
}

func (c *Cache) Rename(ctx context.Context, oldKey, newKey string) error {
	return c.Client.Rename(ctx, oldKey, newKey).Err()
}

func (c *Cache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.Client.HGetAll(ctx, key).Result()
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.Client.Del(ctx, key).Err()
}
