package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TheCrabilia/chaos-shortener/internal/server/api"
	"github.com/TheCrabilia/chaos-shortener/internal/server/chaos"
	"github.com/TheCrabilia/chaos-shortener/internal/server/db"
	"github.com/TheCrabilia/chaos-shortener/internal/server/monitoring"
	"github.com/TheCrabilia/chaos-shortener/internal/server/shortener"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/justinas/alice"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	ctx := context.Background()

	metrics := monitoring.NewMetrics()
	injector := chaos.NewInjector()

	databaseURL := os.Getenv("CSHORT_DATABASE")
	if databaseURL == "" {
		panic("CSHORT_DATABASE is required")
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", os.Getenv("CSHORT_MIGRATIONS_PATH")),
		databaseURL,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			panic(err)
		}
		slog.Info("No migrations to apply")
	} else {
		slog.Info("Migrations applied successfully")
	}

	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		panic(err)
	}

	poolConfig.MaxConns = 100
	poolConfig.MaxConnLifetime = time.Minute * 5
	poolConfig.HealthCheckPeriod = time.Second * 30

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		panic(err)
	}

	db := db.New(pool)
	shortener := shortener.NewShortener(db)
	handlers := api.NewHandlers(shortener, metrics, injector)

	stdChain := alice.New(
		api.NewLoggingMiddleware(slog.With("component", "api"), metrics),
		api.NewChaosMiddleware(injector),
	)

	mux := http.NewServeMux()
	mux.Handle(
		"POST /shorten",
		api.WithHandlerName("shorten", stdChain.Then(handlers.ShortenURL())),
	)
	mux.Handle(
		"GET /r/{id}",
		api.WithHandlerName("redirect", stdChain.Then(handlers.RedirectURL())),
	)
	mux.Handle("POST /chaos", handlers.ConfigureInjector())
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			slog.Error("server shutdown failed", "error", err)
		}
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
