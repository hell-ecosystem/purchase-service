package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/hell-ecosystem/purchase-service/internal/app"
	httpDelivery "github.com/hell-ecosystem/purchase-service/internal/delivery/http"
	"github.com/hell-ecosystem/purchase-service/internal/logger"
	"github.com/hell-ecosystem/purchase-service/internal/repository/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	log := logger.New()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Error("DATABASE_URL is not set")
		os.Exit(1)
	}

	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Error("failed to connect to db", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	repo := postgres.NewPurchaseRepo(db, log)
	service := app.NewPurchaseService(repo, log)
	handler := httpDelivery.NewHandler(service, log)

	log.Info("Starting HTTP server", slog.String("addr", ":8080"))

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: handler.Routes(),
	}

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("could not listen on port", slog.String("error", err.Error()))
	}
}
