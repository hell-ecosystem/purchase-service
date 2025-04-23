package app

import (
	"context"
	"time"

	"github.com/hell-ecosystem/purchase-service/internal/domain"
	"github.com/hell-ecosystem/purchase-service/internal/repository"
)

type purchaseService struct {
	db repository.PurchaseRepository
}

func NewPurchaseService(repo repository.PurchaseRepository) PurchaseService {
	return &purchaseService{db: repo}
}

func (s *purchaseService) Create(ctx context.Context, userID string, items []domain.Purchase) error {
	now := time.Now()
	for i := range items {
		items[i].UserID = userID
		items[i].CreatedAt = now
		items[i].IsActive = true
		items[i].IsPurchased = false
	}
	return s.db.InsertMany(ctx, items)
}

func (s *purchaseService) GetAll(ctx context.Context, userID string) ([]domain.Purchase, error) {
	return s.db.FindByUser(ctx, userID)
}
