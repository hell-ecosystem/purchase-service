package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/hell-ecosystem/purchase-service/internal/domain"
	"github.com/hell-ecosystem/purchase-service/internal/repository"
)

type purchaseService struct {
	db     repository.PurchaseRepository
	logger *slog.Logger
}

func NewPurchaseService(repo repository.PurchaseRepository, logger *slog.Logger) PurchaseService {
	return &purchaseService{db: repo, logger: logger}
}

func (s *purchaseService) Create(ctx context.Context, userID string, items []domain.Purchase) error {
	now := time.Now()

	for i := range items {
		items[i].UserID = userID
		items[i].CreatedAt = now
		items[i].IsActive = true
		items[i].IsPurchased = false
	}

	err := s.db.InsertMany(ctx, items)
	if err != nil {
		s.logger.Error("failed to insert purchases into DB", slog.String("error", err.Error()), slog.Int("count", len(items)))
		return err
	}

	s.logger.Info("purchases inserted into DB", slog.Int("count", len(items)), slog.String("user_id", userID))
	return nil
}

func (s *purchaseService) GetAll(ctx context.Context, userID string) ([]domain.Purchase, error) {
	items, err := s.db.FindByUser(ctx, userID)
	if err != nil {
		s.logger.Error("failed to fetch purchases from DB", slog.String("error", err.Error()), slog.String("user_id", userID))
		return nil, err
	}

	s.logger.Debug("fetched purchases from DB", slog.Int("count", len(items)), slog.String("user_id", userID))
	return items, nil
}
