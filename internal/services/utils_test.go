package services_test

import (
	"strings"
	"testing"

	"github.com/shbhom/urlShortner/internal/models"
	"github.com/shbhom/urlShortner/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestParseBody(t *testing.T) {
	// Setup service instance
	fakeDb := services.NewFakeURLRepo()
	fakeCache := services.NewFakeCacheRepo()
	svc := services.NewService(6, fakeDb, fakeCache)

	t.Run("Invalid JSON", func(t *testing.T) {
		body := strings.NewReader("{bad json")
		var req models.CreateUrlDTO
		err := svc.ParseBody(body, &req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid character")
	})

	t.Run("Missing Required Fields", func(t *testing.T) {
		body := strings.NewReader(`{}`)
		var req models.CreateUrlDTO
		err := svc.ParseBody(body, &req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Field validation for 'Url' failed on the 'https_url' tag")
	})

	t.Run("Success", func(t *testing.T) {
		body := strings.NewReader(`{"url": "https://example.com"}`)
		var req models.CreateUrlDTO
		err := svc.ParseBody(body, &req)
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com", req.Url)
	})
}
