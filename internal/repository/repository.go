package repository

import (
	"context"

	"github.com/hell-ecosystem/purchase-service/internal/domain"
)

type PurchaseRepository interface {
	InsertMany(ctx context.Context, items []domain.Purchase) error
	FindByUser(ctx context.Context, userID string) ([]domain.Purchase, error)
}
