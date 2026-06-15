package metrics

import (
	"context"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

type redisStatsCollector struct {
	client *redis.Client

	poolHitsTotal     *prometheus.Desc
	poolMissesTotal   *prometheus.Desc
	poolTimeoutsTotal *prometheus.Desc
	poolTotalConns    *prometheus.Desc
	poolIdleConns     *prometheus.Desc
	poolStaleConns    *prometheus.Desc
	usedMemoryBytes   *prometheus.Desc
	usedMemoryPeak    *prometheus.Desc
	memFragmentation  *prometheus.Desc
}

func NewRedisStatsCollector(client *redis.Client) *redisStatsCollector {
	return &redisStatsCollector{
		client: client,
		poolHitsTotal: prometheus.NewDesc(
			"redis_pool_hits_total",
			"Total number of times a free connection was found in the pool",
			nil, nil,
		),
		poolMissesTotal: prometheus.NewDesc(
			"redis_pool_misses_total",
			"Total number of times a free connection was NOT found in the pool",
			nil, nil,
		),
		poolTimeoutsTotal: prometheus.NewDesc(
			"redis_pool_timeouts_total",
			"Total number of times a wait timeout occurred for a free connection",
			nil, nil,
		),
		poolTotalConns: prometheus.NewDesc(
			"redis_pool_total_connections",
			"Total number of connections in the pool",
			nil, nil,
		),
		poolIdleConns: prometheus.NewDesc(
			"redis_pool_idle_connections",
			"Total number of idle connections in the pool",
			nil, nil,
		),
		poolStaleConns: prometheus.NewDesc(
			"redis_pool_stale_connections",
			"Total number of stale connections removed from the pool",
			nil, nil,
		),
		usedMemoryBytes: prometheus.NewDesc(
			"redis_used_memory_bytes",
			"Total number of bytes allocated by Redis using its allocator",
			nil, nil,
		),
		usedMemoryPeak: prometheus.NewDesc(
			"redis_used_memory_peak_bytes",
			"Peak memory consumed by Redis (in bytes)",
			nil, nil,
		),
		memFragmentation: prometheus.NewDesc(
			"redis_mem_fragmentation_ratio",
			"Ratio between used_memory_rss and used_memory",
			nil, nil,
		),
	}
}

func (c *redisStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.poolHitsTotal
	ch <- c.poolMissesTotal
	ch <- c.poolTimeoutsTotal
	ch <- c.poolTotalConns
	ch <- c.poolIdleConns
	ch <- c.poolStaleConns
	ch <- c.usedMemoryBytes
	ch <- c.usedMemoryPeak
	ch <- c.memFragmentation
}

func (c *redisStatsCollector) Collect(ch chan<- prometheus.Metric) {
	// Pool stats
	stats := c.client.PoolStats()
	ch <- prometheus.MustNewConstMetric(c.poolHitsTotal, prometheus.CounterValue, float64(stats.Hits))
	ch <- prometheus.MustNewConstMetric(c.poolMissesTotal, prometheus.CounterValue, float64(stats.Misses))
	ch <- prometheus.MustNewConstMetric(c.poolTimeoutsTotal, prometheus.CounterValue, float64(stats.Timeouts))
	ch <- prometheus.MustNewConstMetric(c.poolTotalConns, prometheus.GaugeValue, float64(stats.TotalConns))
	ch <- prometheus.MustNewConstMetric(c.poolIdleConns, prometheus.GaugeValue, float64(stats.IdleConns))
	ch <- prometheus.MustNewConstMetric(c.poolStaleConns, prometheus.GaugeValue, float64(stats.StaleConns))

	// Memory stats
	ctx := context.Background()
	infoStr, err := c.client.Info(ctx, "memory").Result()
	if err == nil {
		lines := strings.Split(infoStr, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "used_memory:") {
				if val, err := strconv.ParseFloat(strings.TrimPrefix(line, "used_memory:"), 64); err == nil {
					ch <- prometheus.MustNewConstMetric(c.usedMemoryBytes, prometheus.GaugeValue, val)
				}
			} else if strings.HasPrefix(line, "used_memory_peak:") {
				if val, err := strconv.ParseFloat(strings.TrimPrefix(line, "used_memory_peak:"), 64); err == nil {
					ch <- prometheus.MustNewConstMetric(c.usedMemoryPeak, prometheus.GaugeValue, val)
				}
			} else if strings.HasPrefix(line, "mem_fragmentation_ratio:") {
				if val, err := strconv.ParseFloat(strings.TrimPrefix(line, "mem_fragmentation_ratio:"), 64); err == nil {
					ch <- prometheus.MustNewConstMetric(c.memFragmentation, prometheus.GaugeValue, val)
				}
			}
		}
	}
}

func RegisterRedisStatsCollector(client *redis.Client) {
	prometheus.MustRegister(NewRedisStatsCollector(client))
}
