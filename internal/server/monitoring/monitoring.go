package monitoring

import (
	// "github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type Metrics struct {
	RequestDuration metric.Float64Histogram
	ResponsesStatus metric.Float64Counter
	URLsCreated     metric.Float64Counter
	RedirectsTotal  metric.Float64Counter
}

func NewMetrics() *Metrics {
	m := &Metrics{}

	exporter, err := prometheus.New()
	if err != nil {
		slog.Error("failed to create Prometheus exporter", "error", err)
		os.Exit(1)
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
	)
	otel.SetMeterProvider(provider)

	meter := otel.Meter("cshort")

	m.RequestDuration, _ = meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("Duration of HTTP requests"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10),
	)

	m.ResponsesStatus, _ = meter.Float64Counter(
		"http_responses_total",
		metric.WithDescription("HTTP responses by status code"),
	)

	m.URLsCreated, _ = meter.Float64Counter(
		"shortener_urls_created_total",
		metric.WithDescription("Total number of shortened URLs created"),
	)

	m.RedirectsTotal, _ = meter.Float64Counter(
		"shortener_redirects_total",
		metric.WithDescription("Total number of redirects performed"),
	)

	return m
}
