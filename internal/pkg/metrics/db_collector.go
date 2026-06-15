package metrics

import (
	"context"
	"database/sql"

	"go.opentelemetry.io/otel/metric"
)

func RegisterDBStatsCollector(db *sql.DB) error {
	openConnections, err := Meter.Float64ObservableGauge(
		"go_sql_open_connections",
		metric.WithDescription("The number of established connections both in use and idle."),
	)
	if err != nil {
		return err
	}

	inUseConnections, err := Meter.Float64ObservableGauge(
		"go_sql_in_use_connections",
		metric.WithDescription("The number of connections currently in use."),
	)
	if err != nil {
		return err
	}

	idleConnections, err := Meter.Float64ObservableGauge(
		"go_sql_idle_connections",
		metric.WithDescription("The number of idle connections."),
	)
	if err != nil {
		return err
	}

	waitCount, err := Meter.Int64ObservableCounter(
		"go_sql_wait_count",
		metric.WithDescription("The total number of connections waited for."),
	)
	if err != nil {
		return err
	}

	waitDurationSeconds, err := Meter.Float64ObservableCounter(
		"go_sql_wait_duration_seconds",
		metric.WithDescription("The total time blocked waiting for a new connection."),
	)
	if err != nil {
		return err
	}

	_, err = Meter.RegisterCallback(func(ctx context.Context, o metric.Observer) error {
		stats := db.Stats()

		o.ObserveFloat64(openConnections, float64(stats.OpenConnections))
		o.ObserveFloat64(inUseConnections, float64(stats.InUse))
		o.ObserveFloat64(idleConnections, float64(stats.Idle))
		o.ObserveInt64(waitCount, stats.WaitCount)
		o.ObserveFloat64(waitDurationSeconds, stats.WaitDuration.Seconds())

		return nil
	}, openConnections, inUseConnections, idleConnections, waitCount, waitDurationSeconds)

	return err
}
