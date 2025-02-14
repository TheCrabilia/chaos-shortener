package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/TheCrabilia/chaos-shortener/internal/monitoring"
	"github.com/TheCrabilia/chaos-shortener/internal/shortener"
)

type Handlers struct {
	shortener *shortener.Shortener
	metrics   *monitoring.Metrics
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

func NewHandlers(shortener *shortener.Shortener, metrics *monitoring.Metrics) *Handlers {
	return &Handlers{
		shortener: shortener,
		metrics:   metrics,
		log:       slog.With("component", "api"),
	}
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

func (h *Handlers) ShortenURL() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ShortenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.metrics.ErrorsTotal.WithLabelValues("shorten", "bad_request").Inc()
			errorJSON(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if req.URL == "" {
			h.metrics.ErrorsTotal.WithLabelValues("shorten", "bad_request").Inc()
			errorJSON(w, "url is required", http.StatusBadRequest)
			return
		}

		shortURL, err := h.shortener.Shorten(r.Context(), fullURL(r), req.URL)
		if err != nil {
			h.log.Error("failed to create shortened url", "error", err)
			h.metrics.ErrorsTotal.WithLabelValues("shorten", "shortening").Inc()
			errorJSON(w, "failed to shorten url", http.StatusInternalServerError)
			return
		}

		h.log.Info("shortened", "url", req.URL, "short_url", shortURL)
		h.metrics.URLsCreated.Inc()

		resp := ShortenResponse{
			ShortURL: shortURL,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			h.log.Error("failed to encode response", "error", err)
			h.metrics.ErrorsTotal.WithLabelValues("shorten", "response_encoding").Inc()
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
			h.metrics.ErrorsTotal.WithLabelValues("redirect", "getting_url").Inc()
			return
		}

		h.log.Info("redirecting", "url", fullURL(r)+r.URL.Path, "redirect_url", origURL)
		h.metrics.RedirectsTotal.Inc()

		http.Redirect(w, r, origURL, http.StatusFound)
	})
}
