package app

import (
	"context"

	"github.com/hell-ecosystem/purchase-service/internal/domain"
)

type PurchaseService interface {
	Create(ctx context.Context, userID string, items []domain.Purchase) error
	GetAll(ctx context.Context, userID string) ([]domain.Purchase, error)
	MarkAsPurchased(ctx context.Context, userID, purchaseID string) error
	Deactivate(ctx context.Context, userID, purchaseID string) error
	Update(ctx context.Context, userID, purchaseID string, update domain.PurchaseUpdate) error
	GetStatistics(ctx context.Context, userID string) (domain.PurchaseStats, error)
	AddReminder(ctx context.Context, userID, purchaseID string, reminder domain.Reminder) error
	AddTags(ctx context.Context, userID, purchaseID string, tagIDs []string) error
}