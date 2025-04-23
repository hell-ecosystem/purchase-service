package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/hell-ecosystem/purchase-service/internal/app"
	httpDelivery "github.com/hell-ecosystem/purchase-service/internal/delivery/http"
	"github.com/hell-ecosystem/purchase-service/internal/repository/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	repo := postgres.NewPurchaseRepo(db)
	service := app.NewPurchaseService(repo)
	handler := httpDelivery.NewHandler(service)

	log.Println("Starting HTTP server on :8080")
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: handler.Routes(),
	}

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("could not listen on :8080: %v", err)
	}
}
