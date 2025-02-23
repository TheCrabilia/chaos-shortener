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

const (
	DBPoolMaxConns          = 100
	DBPoolMaxConnLifetime   = time.Minute * 5
	DBPoolHealthCheckPeriod = time.Second * 30
)

func main() {
	ctx := context.Background()

	metrics := monitoring.NewMetrics()
	injector := chaos.NewInjector()

	dbHost := os.Getenv("CSHORT_DB_HOST")
	if dbHost == "" {
		panic("CSHORT_DB_HOST is required")
	}

	dbName := os.Getenv("CSHORT_DB_NAME")
	if dbName == "" {
		panic("CSHORT_DB_NAME is required")
	}

	dbUser := os.Getenv("CSHORT_DB_USERNAME")
	if dbUser == "" {
		panic("CSHORT_DB_USERNAME is required")
	}

	dbPass := os.Getenv("CSHORT_DB_PASSWORD")
	if dbPass == "" {
		panic("CSHORT_DB_PASSWORD is required")
	}

	databaseURL := fmt.Sprintf("postgresql://%s:%s@%s:5432/%s?sslmode=disable", dbUser, dbPass, dbHost, dbName)

	migrationsPath := os.Getenv("CSHORT_MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "db/migrations"
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
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

	poolConfig.MaxConns = DBPoolMaxConns
	poolConfig.MaxConnLifetime = DBPoolMaxConnLifetime
	poolConfig.HealthCheckPeriod = DBPoolHealthCheckPeriod

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
	mux.Handle("GET /healthz", handlers.Health())
	mux.Handle("GET /metrics", promhttp.Handler())

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
