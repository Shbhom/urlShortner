package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/shbhom/urlShortner/internal/models"
	"github.com/shbhom/urlShortner/internal/pkg/metrics"
)

var ErrURLNotFound = errors.New("url not found")

func (db *DB) GetUrlByCode(ctx context.Context, short_code string) (string, error) {
	start := time.Now()
	// 2. Ensure we record the duration when the function exits
	defer func() {
		metrics.PostgresReadDuration.Record(ctx, time.Since(start).Seconds())
	}()

	// 3. Increment the total reads counter
	metrics.PostgresReadsTotal.Add(ctx, 1)

	var url string
	if err := db.Client.QueryRowContext(ctx, `SELECT url from short_urls where shortnedkey = $1`, short_code).Scan(&url); err != nil {
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
		metrics.PostgresCreationDuration.Record(ctx, time.Since(start).Seconds())
	}()
	if _, err := db.Client.ExecContext(ctx, `INSERT into short_urls (shortnedkey,url) VALUES ($1,$2)`, data.ShortCode, data.TargetUrl); err != nil {
		return fmt.Errorf("Error while inserting short Url to db: %w", err)
	}
	metrics.ShortUrlsCreatedTotal.Add(ctx, 1)
	return nil
}

func (db *DB) GetNextSequence(ctx context.Context) (uint64, error) {
	var nextID uint64
	// We use QueryRow because we expect exactly one integer returned
	err := db.Client.QueryRowContext(ctx, "SELECT nextval('url_counter_seq')").Scan(&nextID)
	if err != nil {
		return 0, err
	}
	return nextID, nil
}

func (db *DB) BulkUpdateUrlLastInvokation(ctx context.Context, data map[string]string) error {
	tx, err := db.Client.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `UPDATE short_urls SET last_invokation = to_timestamp($1) where shortnedkey = $2`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for code, time := range data {
		timeInt, err := strconv.ParseInt(time, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid timestamp %s for code %s", time, code)
		}
		if _, err := stmt.ExecContext(ctx, timeInt, code); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (db *DB) PurgeOldURLs(ctx context.Context) (int64, error) {
	res, err := db.Client.ExecContext(ctx, "DELETE FROM short_urls WHERE last_invokation < NOW() - INTERVAL '1 year'")
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
