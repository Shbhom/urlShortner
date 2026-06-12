package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/shbhom/urlShortner/internal/models"
	"github.com/shbhom/urlShortner/internal/pkg/metrics"
)

var ErrURLNotFound = errors.New("url not found")

func (db *DB) GetUrlByCode(ctx context.Context, short_code string) (string, error) {
	start := time.Now()
	// 2. Ensure we record the duration when the function exits
	defer func() {
		metrics.PostgresReadDuration.Observe(time.Since(start).Seconds())
	}()

	// 3. Increment the total reads counter
	metrics.PostgresReadsTotal.Inc()

	var url string
	if err := db.Client.QueryRowContext(ctx, `SELECT url from shornted_url where shortnedkey = $1`, short_code).Scan(&url); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrURLNotFound
		}
		return "", fmt.Errorf("Error while fetching url from shortcode: %w", err)
	}
	return url, nil
}

func (db *DB) AddUrl(ctx context.Context, data models.UrlData) error {
	start := time.Now()
	defer func() {
		metrics.PostgresCreationDuration.Observe(time.Since(start).Seconds())
	}()
	if _, err := db.Client.ExecContext(ctx, `INSERT into shornted_url (shortnedkey,url) VALUES ($1,$2)`, data.ShortCode, data.TargetUrl); err != nil {
		return fmt.Errorf("Erorr while inserting short Url to db: %w", err)
	}
	metrics.ShortUrlsCreatedTotal.Inc()
	return nil
}
