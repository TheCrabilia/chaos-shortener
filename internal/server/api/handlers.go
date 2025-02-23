package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/TheCrabilia/chaos-shortener/internal/server/chaos"
	"github.com/TheCrabilia/chaos-shortener/internal/server/monitoring"
	"github.com/TheCrabilia/chaos-shortener/internal/server/shortener"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type Handlers struct {
	shortener *shortener.Shortener
	metrics   *monitoring.Metrics
	injector  *chaos.Injector
	log       *slog.Logger
}

type handlerNameKey struct{}

func errorJSON(w http.ResponseWriter, msg string, code int) {
	h := w.Header()
	h.Del("Content-Length")

	h.Set("Content-Type", "application/json")
	h.Set("X-Content-Type-Options", "nosniff")

	error := struct {
		Msg string `json:"message"`
	}{
		Msg: msg,
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&error)
}

func fullURL(r *http.Request) string {
	schema := "http"
	if r.URL.Scheme != "" {
		schema = r.URL.Scheme
	}

	return fmt.Sprintf("%s://%s", schema, r.Host)
}

func WithHandlerName(name string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), handlerNameKey{}, name)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

func NewHandlers(shortener *shortener.Shortener, metrics *monitoring.Metrics, injector *chaos.Injector) *Handlers {
	return &Handlers{
		shortener: shortener,
		metrics:   metrics,
		injector:  injector,
		log:       slog.With("component", "handler"),
	}
}

func (h *Handlers) ShortenURL() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ShortenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.log.Error("failed to decode request", "error", err)
			h.metrics.ErrorsTotal.Add(
				r.Context(),
				1,
				metric.WithAttributes(semconv.HTTPResponseStatusCode(http.StatusBadRequest)),
			)
			errorJSON(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if req.URL == "" {
			h.metrics.ErrorsTotal.Add(
				r.Context(),
				1,
				metric.WithAttributes(semconv.HTTPResponseStatusCode(http.StatusBadRequest)),
			)
			errorJSON(w, "url is required", http.StatusBadRequest)
			return
		}

		shortURL, err := h.shortener.Shorten(r.Context(), fullURL(r), req.URL)
		if err != nil {
			h.log.Error("failed to create shortened url", "error", err)
			h.metrics.ErrorsTotal.Add(
				r.Context(),
				1,
				metric.WithAttributes(semconv.HTTPResponseStatusCode(http.StatusInternalServerError)),
			)
			errorJSON(w, "failed to shorten url", http.StatusInternalServerError)
			return
		}

		h.log.Info("shortened", "url", req.URL, "short_url", shortURL)
		h.metrics.URLsCreated.Add(r.Context(), 1)

		resp := ShortenResponse{
			ShortURL: shortURL,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			h.log.Error("failed to encode response", "error", err)
			h.metrics.ErrorsTotal.Add(
				r.Context(),
				1,
				metric.WithAttributes(semconv.HTTPResponseStatusCode(http.StatusInternalServerError)),
			)
			errorJSON(w, "failed to shorten url", http.StatusInternalServerError)
			return
		}
	})
}

func (h *Handlers) RedirectURL() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		origURL, err := h.shortener.RedirectURL(r.Context(), id)
		if err != nil {
			h.log.Error("failed to get original url", "error", err)
			h.metrics.ErrorsTotal.Add(
				r.Context(),
				1,
				metric.WithAttributes(semconv.HTTPResponseStatusCode(http.StatusInternalServerError)),
			)
			return
		}

		h.log.Info("redirecting", "url", fullURL(r)+r.URL.Path, "redirect_url", origURL)
		h.metrics.RedirectsTotal.Add(r.Context(), 1)

		http.Redirect(w, r, origURL, http.StatusFound)
	})
}

func (h *Handlers) ConfigureInjector() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req InjectorRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.log.Error("failed to decode request", "error", err)
			h.metrics.ErrorsTotal.Add(
				r.Context(),
				1,
				metric.WithAttributes(semconv.HTTPResponseStatusCode(http.StatusBadRequest)),
			)
			errorJSON(w, "invalid request body", http.StatusBadRequest)
			return
		}

		h.injector.SetLatencyRate(req.LatencyRate)
		h.injector.SetErrorRate(req.ErrorRate)
		h.injector.SetOutageRate(req.OutageRate)
		slog.Info(
			"injector configured",
			"latency_rate",
			req.LatencyRate,
			"error_rate",
			req.ErrorRate,
			"outage_rate",
			req.OutageRate,
		)

		resp := InjectorResponse{
			LatencyRate: h.injector.GetLatencyRate(),
			ErrorRate:   h.injector.GetErrorRate(),
			OutageRate:  h.injector.GetOutageRate(),
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			h.log.Error("failed to encode response", "error", err)
			h.metrics.ErrorsTotal.Add(
				r.Context(),
				1,
				metric.WithAttributes(semconv.HTTPResponseStatusCode(http.StatusInternalServerError)),
			)
			errorJSON(w, "failed to configure injector", http.StatusInternalServerError)
			return
		}
	})
}

func (h *Handlers) Health() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})
}
