package app

import (
	"context"
	"errors"
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

func (s *purchaseService) MarkAsPurchased(ctx context.Context, userID, purchaseID string) error {
	err := s.db.MarkAsPurchased(ctx, userID, purchaseID)
	if err != nil {
		s.logger.Error("failed to mark purchase as purchased", slog.String("error", err.Error()), slog.String("purchase_id", purchaseID), slog.String("user_id", userID))
		return err
	}

	s.logger.Info("purchase marked as purchased", slog.String("purchase_id", purchaseID), slog.String("user_id", userID))
	return nil
}

func (s *purchaseService) Deactivate(ctx context.Context, userID, purchaseID string) error {
	err := s.db.Deactivate(ctx, userID, purchaseID)
	if err != nil {
		s.logger.Error("failed to deactivate purchase", slog.String("error", err.Error()), slog.String("purchase_id", purchaseID), slog.String("user_id", userID))
		return err
	}

	s.logger.Info("purchase deactivated", slog.String("purchase_id", purchaseID), slog.String("user_id", userID))
	return nil
}

func (s *purchaseService) Update(ctx context.Context, userID, purchaseID string, update domain.PurchaseUpdate) error {
	err := s.db.Update(ctx, userID, purchaseID, update)
	if err != nil {
		s.logger.Error("failed to update purchase", slog.String("error", err.Error()), slog.String("purchase_id", purchaseID), slog.String("user_id", userID))
		return err
	}

	s.logger.Info("purchase updated", slog.String("purchase_id", purchaseID), slog.String("user_id", userID))
	return nil
}

func (s *purchaseService) GetStatistics(ctx context.Context, userID string) (domain.PurchaseStats, error) {
	stats, err := s.db.CollectStats(ctx, userID)
	if err != nil {
		s.logger.Error("failed to collect statistics", slog.String("user_id", userID), slog.String("error", err.Error()))
		return domain.PurchaseStats{}, err
	}

	s.logger.Debug("statistics collected", slog.String("user_id", userID), slog.Int("total", stats.Total))
	return stats, nil
}

func (s *purchaseService) AddReminder(ctx context.Context, userID, purchaseID string, reminder domain.Reminder) error {
	if reminder.RemindAt.IsZero() {
		s.logger.Warn("reminder time not provided", slog.String("purchase_id", purchaseID), slog.String("user_id", userID))
		return errors.New("reminder time is required")
	}

	err := s.db.InsertReminder(ctx, userID, purchaseID, reminder)
	if err != nil {
		s.logger.Error("failed to insert reminder", slog.String("purchase_id", purchaseID), slog.String("user_id", userID), slog.String("error", err.Error()))
		return err
	}

	s.logger.Info("reminder added", slog.String("purchase_id", purchaseID), slog.String("user_id", userID))
	return nil
}

func (s *purchaseService) AddTags(ctx context.Context, userID, purchaseID string, tagIDs []string) error {
	if len(tagIDs) == 0 {
		s.logger.Warn("empty tag list", slog.String("purchase_id", purchaseID), slog.String("user_id", userID))
		return errors.New("tag list cannot be empty")
	}

	err := s.db.AttachTags(ctx, userID, purchaseID, tagIDs)
	if err != nil {
		s.logger.Error("failed to attach tags", slog.String("purchase_id", purchaseID), slog.String("user_id", userID), slog.String("error", err.Error()))
		return err
	}

	s.logger.Info("tags attached", slog.String("purchase_id", purchaseID), slog.String("user_id", userID), slog.Int("tag_count", len(tagIDs)))
	return nil
}
