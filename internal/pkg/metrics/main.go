package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Request Metrics
var (
	RedirectRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redirect_request_total",
			Help: "total No of redirect requests",
		},
		[]string{
			"endpoint",
			"status",
		},
	)

	RedirectRequestHistogram = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "redirect_request_duration_seconds",
			Help:    "Histogram of redirect request latencies",
			Buckets: prometheus.DefBuckets,
		},
	)
)

// Redirect Application Metrics
var (
	RedirectionSuccessTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "redirect_success_total",
			Help: "total no of successful redirections",
		},
	)

	RedirectNotFoundTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "redirect_not_found_total",
			Help: "total no of NOT Found redirections",
		},
	)
)

// Database Metrics
var (
	PostgresReadsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "postgres_reads_total",
			Help: "Total no of postgres Reads",
		},
	)
	PostgresReadDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "postgres_read_duration_seconds",
			Help:    "Histogram of postgres read latencies",
			Buckets: prometheus.DefBuckets,
		},
	)
	ShortUrlsCreatedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "short_urls_created_total",
			Help: "Total no of short urls created",
		},
	)
	PostgresCreationDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "postgres_create_duration_seconds",
			Help:    "Histogram of postgres create latencies",
			Buckets: prometheus.DefBuckets,
		},
	)
)

// Redis Metrics
var (
	RedisHitsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "redis_hits_total",
			Help: "Total number of successful redis cache hits",
		},
	)
	RedisMissesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "redis_misses_total",
			Help: "Total number of redis cache misses",
		},
	)
	RedisGetDurationSeconds = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "redis_get_duration_seconds",
			Help:    "Histogram of redis get latencies",
			Buckets: prometheus.DefBuckets,
		},
	)
	RedisSetDurationSeconds = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "redis_set_duration_seconds",
			Help:    "Histogram of redis set latencies",
			Buckets: prometheus.DefBuckets,
		},
	)
)
