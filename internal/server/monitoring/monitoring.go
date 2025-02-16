package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	RequestDuration *prometheus.HistogramVec
	RequestsTotal   *prometheus.CounterVec
	ErrorsTotal     *prometheus.CounterVec
	URLsCreated     prometheus.Counter
	RedirectsTotal  prometheus.Counter
}

func NewMetrics() *Metrics {
	m := &Metrics{
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests",
				Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"handler", "method", "status"},
		),
		RequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"handler", "method", "status"},
		),
		ErrorsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_errors_total",
				Help: "Total number of HTTP errors",
			},
			[]string{"handler", "error_type"},
		),
		URLsCreated: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "shortener_urls_created_total",
				Help: "Total number of shortened URLs created",
			}),
		RedirectsTotal: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "shortener_redirects_total",
				Help: "Total number of redirects performed",
			},
		),
	}

	prometheus.MustRegister(m.RequestDuration,
		m.RequestsTotal,
		m.ErrorsTotal,
		m.URLsCreated,
		m.RedirectsTotal,
	)

	return m
}
