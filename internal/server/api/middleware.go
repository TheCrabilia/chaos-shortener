package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/TheCrabilia/chaos-shortener/internal/server/chaos"
	"github.com/TheCrabilia/chaos-shortener/internal/server/monitoring"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *ResponseWriter) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *ResponseWriter) statusCodeString() string {
	return fmt.Sprintf("%d", w.StatusCode)
}

func NewLoggingMiddleware(logger *slog.Logger, metrics *monitoring.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &ResponseWriter{w, http.StatusOK}

			start := time.Now()
			next.ServeHTTP(rw, r)
			end := time.Since(start).Seconds()

			handlerName := "unknown"
			if name := r.Context().Value(handlerNameKey{}); name != nil {
				handlerName = name.(string)
			}

			metrics.RequestDuration.Record(
				r.Context(),
				end,
				metric.WithAttributes(
					attribute.String("handler", handlerName),
				),
			)

			metrics.ResponsesStatus.Add(
				r.Context(),
				1,
				metric.WithAttributes(
					attribute.String("handler", handlerName),
					semconv.HTTPResponseStatusCode(rw.StatusCode),
				),
			)

			logger.Info(
				"request",
				"method",
				r.Method,
				"url",
				r.URL.String(),
				"status",
				rw.statusCodeString(),
				"latency",
				end,
			)
		})
	}
}

func NewChaosMiddleware(injector *chaos.Injector) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if shouldFail, failureType := injector.ShouldInject(); shouldFail {
				switch failureType {
				case chaos.FailureTypeLatency:
					injector.InjectLatency()
				case chaos.FailureTypeError:
					injector.InjectError(w)
					return
				case chaos.FailureTypeOutage:
					injector.InjectOutage(w)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
