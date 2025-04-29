package repository

import (
	"context"

	"github.com/hell-ecosystem/purchase-service/internal/domain"
)

type PurchaseRepository interface {
	InsertMany(ctx context.Context, items []domain.Purchase) error
	FindByUser(ctx context.Context, userID string) ([]domain.Purchase, error)
	MarkAsPurchased(ctx context.Context, userID, purchaseID string) error
	Deactivate(ctx context.Context, userID, purchaseID string) error
	Update(ctx context.Context, userID, purchaseID string, update domain.PurchaseUpdate) error
	CollectStats(ctx context.Context, userID string) (domain.PurchaseStats, error)
	InsertReminder(ctx context.Context, userID, purchaseID string, reminder domain.Reminder) error
	AttachTags(ctx context.Context, userID, purchaseID string, tagIDs []string) error
}
