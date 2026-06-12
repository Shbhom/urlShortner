package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shbhom/urlShortner/internal/models"
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
	return c.Client.Get(ctx, fmt.Sprintf("url:%s", shortCode)).Result()
}

func (c *Cache) Set(ctx context.Context, data models.UrlData) error {
	return c.Client.Set(ctx, fmt.Sprintf("url:%s", data.ShortCode), data.TargetUrl, c.TTL).Err()
}
