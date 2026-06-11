package metrics

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus"
)

type dbStatsCollector struct {
	db                  *sql.DB
	openConnections     *prometheus.Desc
	inUseConnections    *prometheus.Desc
	idleConnections     *prometheus.Desc
	waitCount           *prometheus.Desc
	waitDurationSeconds *prometheus.Desc
}

func NewDBStatsCollector(db *sql.DB) prometheus.Collector {
	return &dbStatsCollector{
		db: db,
		openConnections: prometheus.NewDesc(
			"go_sql_open_connections",
			"The number of established connections both in use and idle.",
			nil, nil,
		),
		inUseConnections: prometheus.NewDesc(
			"go_sql_in_use_connections",
			"The number of connections currently in use.",
			nil, nil,
		),
		idleConnections: prometheus.NewDesc(
			"go_sql_idle_connections",
			"The number of idle connections.",
			nil, nil,
		),
		waitCount: prometheus.NewDesc(
			"go_sql_wait_count_total",
			"The total number of connections waited for.",
			nil, nil,
		),
		waitDurationSeconds: prometheus.NewDesc(
			"go_sql_wait_duration_seconds_total",
			"The total time blocked waiting for a new connection.",
			nil, nil,
		),
	}
}

func (c *dbStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.openConnections
	ch <- c.inUseConnections
	ch <- c.idleConnections
	ch <- c.waitCount
	ch <- c.waitDurationSeconds
}

func (c *dbStatsCollector) Collect(ch chan<- prometheus.Metric) {
	stats := c.db.Stats()

	ch <- prometheus.MustNewConstMetric(c.openConnections, prometheus.GaugeValue, float64(stats.OpenConnections))
	ch <- prometheus.MustNewConstMetric(c.inUseConnections, prometheus.GaugeValue, float64(stats.InUse))
	ch <- prometheus.MustNewConstMetric(c.idleConnections, prometheus.GaugeValue, float64(stats.Idle))
	ch <- prometheus.MustNewConstMetric(c.waitCount, prometheus.CounterValue, float64(stats.WaitCount))
	ch <- prometheus.MustNewConstMetric(c.waitDurationSeconds, prometheus.CounterValue, stats.WaitDuration.Seconds())
}

func RegisterDBStatsCollector(db *sql.DB) {
	prometheus.MustRegister(NewDBStatsCollector(db))
}
