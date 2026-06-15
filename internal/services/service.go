package services

import (
	"context"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/shbhom/urlShortner/internal/models"
	"github.com/sqids/sqids-go"
)

type URLRepository interface {
	GetNextSequence(ctx context.Context) (uint64, error)
	AddUrl(ctx context.Context, data models.UrlData) error
	GetUrlByCode(ctx context.Context, short_code string) (string, error)
	BulkUpdateUrlLastInvokation(ctx context.Context, data map[string]string) error
	PurgeOldURLs(ctx context.Context) (int64, error)
}

type CacheRepository interface {
	Get(ctx context.Context, shortCode string) (string, error)
	Set(ctx context.Context, data models.UrlData) error
	Rename(ctx context.Context, oldKey, newKey string) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	Delete(ctx context.Context, key string) error
	RecordInvokation(ctx context.Context, code string) error
}

type Service struct {
	url           URLRepository
	cache         CacheRepository
	sqid          *sqids.Sqids
	validator     *validator.Validate
	AnalyticsChan chan string
	AnalyticsWG   *sync.WaitGroup
}

func NewService(minLen uint8, urlRepo URLRepository, cacheRepo CacheRepository) *Service {
	sq := NewSquid(minLen)
	vali := validator.New()
	wg := &sync.WaitGroup{}
	wg.Add(5)
	srv := &Service{
		url:           urlRepo,
		cache:         cacheRepo,
		sqid:          sq,
		validator:     vali,
		AnalyticsChan: make(chan string, 10000),
		AnalyticsWG:   wg,
	}

	for _ = range 5 {
		go srv.invocationWorker()
	}
	return srv
}
