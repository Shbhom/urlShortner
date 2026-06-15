package metrics

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// InitProvider initializes the OpenTelemetry MeterProvider
// and configures it to export via OTLP over HTTP.
func InitProvider(ctx context.Context) (func(context.Context) error, error) {
	// The exporter automatically uses the standard environment variables:
	// OTEL_EXPORTER_OTLP_ENDPOINT, OTEL_EXPORTER_OTLP_HEADERS, etc.
	exporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		return nil, err
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			resource.Default().SchemaURL(),
			semconv.ServiceName("url-shortener"),
		),
	)
	if err != nil {
		return nil, err
	}

	histogramView := sdkmetric.NewView(
		sdkmetric.Instrument{Name: "*_duration_seconds"},
		sdkmetric.Stream{
			Aggregation: sdkmetric.AggregationExplicitBucketHistogram{
				Boundaries: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
			},
		},
	)

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
		sdkmetric.WithResource(res),
		sdkmetric.WithView(histogramView),
	)

	// Set the global MeterProvider
	otel.SetMeterProvider(meterProvider)

	return meterProvider.Shutdown, nil
}
