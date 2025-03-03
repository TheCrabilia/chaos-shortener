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
)

type Handlers struct {
	shortener *shortener.Shortener
	metrics   *monitoring.Metrics
	injector  *chaos.Injector
	log       *slog.Logger
}

type handlerNameKey struct{}

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
			errorJSON(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if req.URL == "" {
			errorJSON(w, "url is required", http.StatusBadRequest)
			return
		}

		urlID, err := h.shortener.Shorten(r.Context(), req.URL)
		if err != nil {
			h.log.Error("failed to create shortened url", "error", err)
			errorJSON(w, "failed to shorten url", http.StatusInternalServerError)
			return
		}

		shortURL := fmt.Sprintf("%s/r/%s", fullURL(r), urlID)

		h.log.Info("shortened", "url", req.URL, "short_url", shortURL)
		h.metrics.URLsCreated.Add(r.Context(), 1)

		resp := ShortenResponse{
			ID:       urlID,
			ShortURL: shortURL,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			h.log.Error("failed to encode response", "error", err)
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
