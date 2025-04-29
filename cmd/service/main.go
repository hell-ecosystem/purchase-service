package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hell-ecosystem/purchase-service/internal/app"
	"github.com/hell-ecosystem/purchase-service/internal/config"
	httpDelivery "github.com/hell-ecosystem/purchase-service/internal/delivery/http"
	"github.com/hell-ecosystem/purchase-service/internal/logger"
	"github.com/hell-ecosystem/purchase-service/internal/repository/postgres"
	"github.com/hell-ecosystem/purchase-service/internal/tracer"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/getsentry/sentry-go"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	log := logger.New()

	_, cancel := context.WithCancel(context.Background()) //TODO: пока не использую контекст
	defer cancel()

	// Initialize OpenTelemetry Tracer (Jaeger)
	shutdownTracer, err := tracer.InitTracer(cfg.ServiceName, cfg.OtelExporterEndpoint)
	if err != nil {
		log.Error("failed to initialize tracer", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownTracer(ctxTimeout); err != nil {
			log.Error("tracer shutdown failed", slog.String("error", err.Error()))
		}
	}()

	// Initialize Sentry if provided
	if cfg.SentryDSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.SentryDSN,
			Environment:      cfg.SentryEnv,
			TracesSampleRate: cfg.SentrySampleRate,
		})
		if err != nil {
			log.Error("failed to initialize sentry", slog.String("error", err.Error()))
			os.Exit(1)
		}
		defer sentry.Flush(2 * time.Second)
		log.Info("sentry initialized")
	}

	// Connect to Postgres
	db, err := pgxpool.New(context.Background(), "dbURL") //TODO: сделать нормальную строку для подключения бд
	if err != nil {
		log.Error("failed to connect to db", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	// Repositories
	purchaseRepo := postgres.NewPurchaseRepo(db, log)

	// Usecases
	purchaseUseCase := app.NewPurchaseService(purchaseRepo, log)

	// HTTP Handlers
	handler := httpDelivery.NewHandler(purchaseUseCase, log)

	// HTTP Server
	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      handler.Routes(),
		ReadTimeout:  cfg.GetReadTimeout(),
		WriteTimeout: cfg.GetWriteTimeout(),
		IdleTimeout:  cfg.GetIdleTimeout(),
	}

	// Run server
	go func() {
		log.Info("starting HTTP server", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP server error", slog.String("error", err.Error()))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()
	if err := srv.Shutdown(ctxTimeout); err != nil {
		log.Error("server shutdown error", slog.String("error", err.Error()))
	} else {
		log.Info("server gracefully stopped")
	}
}
