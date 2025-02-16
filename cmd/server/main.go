package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TheCrabilia/chaos-shortener/internal/api"
	"github.com/TheCrabilia/chaos-shortener/internal/chaos"
	"github.com/TheCrabilia/chaos-shortener/internal/db"
	"github.com/TheCrabilia/chaos-shortener/internal/monitoring"
	"github.com/TheCrabilia/chaos-shortener/internal/shortener"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/justinas/alice"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	latencyRate  float64
	errorRate    float64
	connDropRate float64
	outageRate   float64
)

func main() {
	flag.Float64Var(&latencyRate, "latency-rate", 0, "rate of latency injection")
	flag.Float64Var(&errorRate, "error-rate", 0, "rate of error injection")
	flag.Float64Var(&connDropRate, "conn-drop-rate", 0, "rate of connection drop injection")
	flag.Float64Var(&outageRate, "outage-rate", 0, "rate of outage injection")

	flag.Parse()

	metrics := monitoring.NewMetrics()
	injector := chaos.NewInjector()

	injector.SetLatencyRate(latencyRate)
	injector.SetErrorRate(errorRate)
	injector.SetConnDropRate(connDropRate)
	injector.SetOutageRate(outageRate)

	conn, err := pgx.Connect(context.Background(), os.Getenv("CSHORT_DATABASE"))
	if err != nil {
		panic(err)
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", os.Getenv("CSHORT_MIGRATIONS_PATH")),
		os.Getenv("CSHORT_DATABASE"),
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

	db := db.New(conn)
	shortener := shortener.NewShortener(db)
	handlers := api.NewHandlers(shortener, metrics)

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

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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
