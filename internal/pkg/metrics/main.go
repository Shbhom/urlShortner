package metrics

import (
	"context"
	"runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var Meter = otel.Meter("url-shortener")

// Request Metrics
var (
	RedirectRequestsTotal    metric.Int64Counter
	RedirectRequestHistogram metric.Float64Histogram
)

// Redirect Application Metrics
var (
	RedirectionSuccessTotal metric.Int64Counter
	RedirectNotFoundTotal   metric.Int64Counter
)

// Database Metrics
var (
	PostgresReadsTotal       metric.Int64Counter
	PostgresReadDuration     metric.Float64Histogram
	ShortUrlsCreatedTotal    metric.Int64Counter
	PostgresCreationDuration metric.Float64Histogram
)

// Redis Metrics
var (
	RedisHitsTotal          metric.Int64Counter
	RedisMissesTotal        metric.Int64Counter
	RedisGetDurationSeconds metric.Float64Histogram
	RedisSetDurationSeconds metric.Float64Histogram
)

func InitMetrics() error {
	var err error

	// Request Metrics
	RedirectRequestsTotal, err = Meter.Int64Counter(
		"redirect_request",
		metric.WithDescription("total No of redirect requests"),
	)
	if err != nil {
		return err
	}

	RedirectRequestHistogram, err = Meter.Float64Histogram(
		"redirect_request_duration_seconds",
		metric.WithDescription("Histogram of redirect request latencies"),
	)
	if err != nil {
		return err
	}

	// Redirect Application Metrics
	RedirectionSuccessTotal, err = Meter.Int64Counter(
		"redirect_success",
		metric.WithDescription("total no of successful redirections"),
	)
	if err != nil {
		return err
	}

	RedirectNotFoundTotal, err = Meter.Int64Counter(
		"redirect_not_found",
		metric.WithDescription("total no of NOT Found redirections"),
	)
	if err != nil {
		return err
	}

	// Database Metrics
	PostgresReadsTotal, err = Meter.Int64Counter(
		"postgres_reads",
		metric.WithDescription("Total no of postgres Reads"),
	)
	if err != nil {
		return err
	}

	PostgresReadDuration, err = Meter.Float64Histogram(
		"postgres_read_duration_seconds",
		metric.WithDescription("Histogram of postgres read latencies"),
	)
	if err != nil {
		return err
	}

	ShortUrlsCreatedTotal, err = Meter.Int64Counter(
		"short_urls_created",
		metric.WithDescription("Total no of short urls created"),
	)
	if err != nil {
		return err
	}

	PostgresCreationDuration, err = Meter.Float64Histogram(
		"postgres_create_duration_seconds",
		metric.WithDescription("Histogram of postgres create latencies"),
	)
	if err != nil {
		return err
	}

	// Redis Metrics
	RedisHitsTotal, err = Meter.Int64Counter(
		"redis_hits",
		metric.WithDescription("Total number of successful redis cache hits"),
	)
	if err != nil {
		return err
	}

	RedisMissesTotal, err = Meter.Int64Counter(
		"redis_misses",
		metric.WithDescription("Total number of redis cache misses"),
	)
	if err != nil {
		return err
	}

	RedisGetDurationSeconds, err = Meter.Float64Histogram(
		"redis_get_duration_seconds",
		metric.WithDescription("Histogram of redis get latencies"),
	)
	if err != nil {
		return err
	}

	RedisSetDurationSeconds, err = Meter.Float64Histogram(
		"redis_set_duration_seconds",
		metric.WithDescription("Histogram of redis set latencies"),
	)
	if err != nil {
		return err
	}

	// Go Runtime Metrics
	_, err = Meter.Int64ObservableGauge(
		"go_goroutines",
		metric.WithDescription("Number of goroutines that currently exist."),
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			o.Observe(int64(runtime.NumGoroutine()))
			return nil
		}),
	)
	if err != nil {
		return err
	}

	return nil
}
