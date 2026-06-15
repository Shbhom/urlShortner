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

var (
	ANALYTICS_KEY     = "analytics_batch"
	ANALYTICS_NEW_KEY = "analytics_processing:%s"
)

func NewCache(redisAddr string, urlTTL time.Duration) *Cache {
	var client *redis.Client
	if Opts, err := redis.ParseURL(redisAddr); err != nil {
		log.Fatal("error while parsing redis Opts")
	} else {
		client = redis.NewClient(Opts)
	}
	if client == nil {
		log.Fatal("unable to connect redis")
	}
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
