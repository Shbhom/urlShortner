package metrics

import (
	"context"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/metric"
)

func RegisterRedisStatsCollector(client *redis.Client) error {
	poolHitsTotal, err := Meter.Int64ObservableCounter(
		"redis_pool_hits",
		metric.WithDescription("Total number of times a free connection was found in the pool"),
	)
	if err != nil {
		return err
	}

	poolMissesTotal, err := Meter.Int64ObservableCounter(
		"redis_pool_misses",
		metric.WithDescription("Total number of times a free connection was NOT found in the pool"),
	)
	if err != nil {
		return err
	}

	poolTimeoutsTotal, err := Meter.Int64ObservableCounter(
		"redis_pool_timeouts",
		metric.WithDescription("Total number of times a wait timeout occurred for a free connection"),
	)
	if err != nil {
		return err
	}

	poolTotalConns, err := Meter.Int64ObservableGauge(
		"redis_pool_total_connections",
		metric.WithDescription("Total number of connections in the pool"),
	)
	if err != nil {
		return err
	}

	poolIdleConns, err := Meter.Int64ObservableGauge(
		"redis_pool_idle_connections",
		metric.WithDescription("Total number of idle connections in the pool"),
	)
	if err != nil {
		return err
	}

	poolStaleConns, err := Meter.Int64ObservableCounter(
		"redis_pool_stale_connections",
		metric.WithDescription("Total number of stale connections removed from the pool"),
	)
	if err != nil {
		return err
	}

	usedMemoryBytes, err := Meter.Int64ObservableGauge(
		"redis_used_memory_bytes",
		metric.WithDescription("Total number of bytes allocated by Redis using its allocator"),
	)
	if err != nil {
		return err
	}

	usedMemoryPeak, err := Meter.Int64ObservableGauge(
		"redis_used_memory_peak_bytes",
		metric.WithDescription("Peak memory consumed by Redis (in bytes)"),
	)
	if err != nil {
		return err
	}

	memFragmentation, err := Meter.Float64ObservableGauge(
		"redis_mem_fragmentation_ratio",
		metric.WithDescription("Ratio between used_memory_rss and used_memory"),
	)
	if err != nil {
		return err
	}

	_, err = Meter.RegisterCallback(func(ctx context.Context, o metric.Observer) error {
		// Pool stats
		stats := client.PoolStats()
		if stats != nil {
			o.ObserveInt64(poolHitsTotal, int64(stats.Hits))
			o.ObserveInt64(poolMissesTotal, int64(stats.Misses))
			o.ObserveInt64(poolTimeoutsTotal, int64(stats.Timeouts))
			o.ObserveInt64(poolTotalConns, int64(stats.TotalConns))
			o.ObserveInt64(poolIdleConns, int64(stats.IdleConns))
			o.ObserveInt64(poolStaleConns, int64(stats.StaleConns))
		}

		// Server info
		info, err := client.Info(ctx, "memory").Result()
		if err == nil {
			lines := strings.Split(info, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}

				parts := strings.SplitN(line, ":", 2)
				if len(parts) != 2 {
					continue
				}

				key := parts[0]
				value := parts[1]

				switch key {
				case "used_memory":
					if v, err := strconv.ParseInt(value, 10, 64); err == nil {
						o.ObserveInt64(usedMemoryBytes, v)
					}
				case "used_memory_peak":
					if v, err := strconv.ParseInt(value, 10, 64); err == nil {
						o.ObserveInt64(usedMemoryPeak, v)
					}
				case "mem_fragmentation_ratio":
					if v, err := strconv.ParseFloat(value, 64); err == nil {
						o.ObserveFloat64(memFragmentation, v)
					}
				}
			}
		}

		return nil
	}, poolHitsTotal, poolMissesTotal, poolTimeoutsTotal, poolTotalConns, poolIdleConns, poolStaleConns, usedMemoryBytes, usedMemoryPeak, memFragmentation)

	return err
}
