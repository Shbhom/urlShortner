package services_test

import (
	"context"
	"testing"
	"time"

	appredis "github.com/shbhom/urlShortner/internal/db/redis"
	"github.com/shbhom/urlShortner/internal/models"
	"github.com/shbhom/urlShortner/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestService_AddAndGetUrl(t *testing.T) {
	// Setup
	fakeDb := services.NewFakeURLRepo()
	fakeCache := services.NewFakeCacheRepo()
	svc := services.NewService(6, fakeDb, fakeCache)

	ctx := context.Background()

	// Action
	shortCode, err := svc.Addurl(ctx, "https://example.com")
	
	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, shortCode)

	targetUrl, err := svc.GetUrl(ctx, shortCode)
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", targetUrl)
}

func TestService_CacheFallback(t *testing.T) {
	// Setup
	fakeDb := services.NewFakeURLRepo()
	fakeCache := services.NewFakeCacheRepo()
	svc := services.NewService(6, fakeDb, fakeCache)

	ctx := context.Background()
	shortCode := "fallback_test_code"
	targetUrl := "https://fallback.example.com"

	// Arrange: Artificially insert a URL into the fakeDb directly
	err := fakeDb.AddUrl(ctx, models.UrlData{ShortCode: shortCode, TargetUrl: targetUrl})
	assert.NoError(t, err)

	// Ensure fakeCache is completely empty for this shortCode
	_, err = fakeCache.Get(ctx, shortCode)
	assert.Error(t, err) // Expect an error or empty depending on fake cache implementation

	// Action
	retrievedUrl, err := svc.GetUrl(ctx, shortCode)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, targetUrl, retrievedUrl)

	// Verify that it populates the cache
	cachedUrl, err := fakeCache.Get(ctx, shortCode)
	assert.NoError(t, err)
	assert.Equal(t, targetUrl, cachedUrl)
}

func TestService_NotFound(t *testing.T) {
	fakeDb := services.NewFakeURLRepo()
	fakeCache := services.NewFakeCacheRepo()
	svc := services.NewService(6, fakeDb, fakeCache)

	ctx := context.Background()

	_, err := svc.GetUrl(ctx, "non_existent_code")
	assert.Error(t, err)
}

func TestService_AnalyticsQueuing(t *testing.T) {
	fakeDb := services.NewFakeURLRepo()
	fakeCache := services.NewFakeCacheRepo()
	svc := services.NewService(6, fakeDb, fakeCache)

	ctx := context.Background()
	shortCode := "test_code"

	err := fakeDb.AddUrl(ctx, models.UrlData{ShortCode: shortCode, TargetUrl: "https://example.com"})
	assert.NoError(t, err)

	_, err = svc.GetUrl(ctx, shortCode)
	assert.NoError(t, err)

	// Wait slightly for the invocationWorker to process the channel message
	time.Sleep(50 * time.Millisecond)

	hash, err := fakeCache.HGetAll(ctx, appredis.ANALYTICS_KEY)
	assert.NoError(t, err)
	assert.Contains(t, hash, shortCode)
}

func TestService_FlushAnalytics(t *testing.T) {
	fakeDb := services.NewFakeURLRepo()
	fakeCache := services.NewFakeCacheRepo()
	svc := services.NewService(6, fakeDb, fakeCache)

	ctx := context.Background()

	err := fakeCache.RecordInvokation(ctx, "code_A")
	assert.NoError(t, err)
	err = fakeCache.RecordInvokation(ctx, "code_B")
	assert.NoError(t, err)

	svc.FlushAnalytics(ctx)

	bulkUpdates := fakeDb.GetBulkUpdates()
	assert.Len(t, bulkUpdates, 1)
	assert.Contains(t, bulkUpdates[0], "code_A")
	assert.Contains(t, bulkUpdates[0], "code_B")

	hash, err := fakeCache.HGetAll(ctx, appredis.ANALYTICS_KEY)
	assert.NoError(t, err)
	assert.Empty(t, hash)
}
