package app

import (
	"context"

	"github.com/hell-ecosystem/purchase-service/internal/domain"
)

type PurchaseService interface {
	Create(ctx context.Context, userID string, items []domain.Purchase) error
	GetAll(ctx context.Context, userID string) ([]domain.Purchase, error)
}
