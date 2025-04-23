package postgres

import (
	"context"
	"log/slog"

	"github.com/hell-ecosystem/purchase-service/internal/domain"
	"github.com/hell-ecosystem/purchase-service/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type purchaseRepo struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewPurchaseRepo(db *pgxpool.Pool, logger *slog.Logger) repository.PurchaseRepository {
	return &purchaseRepo{db: db, logger: logger}
}

func (r *purchaseRepo) InsertMany(ctx context.Context, items []domain.Purchase) error {
	for _, item := range items {
		_, err := r.db.Exec(ctx, `
            INSERT INTO purchases (user_id, name, category, estimated_price, actual_price, created_at, deadline, is_active, is_purchased)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        `,
			item.UserID, item.Name, item.Category, item.EstimatedPrice, item.ActualPrice,
			item.CreatedAt, item.Deadline, item.IsActive, item.IsPurchased,
		)
		if err != nil {
			r.logger.Error("failed to insert single purchase", slog.String("user_id", item.UserID), slog.String("name", item.Name), slog.String("error", err.Error()))
			return err
		}
	}
	r.logger.Debug("successfully inserted purchases", slog.Int("count", len(items)))
	return nil
}

func (r *purchaseRepo) FindByUser(ctx context.Context, userID string) ([]domain.Purchase, error) {
	rows, err := r.db.Query(ctx, `
        SELECT id, name, category, estimated_price, actual_price, created_at, deadline, is_active, is_purchased
        FROM purchases WHERE user_id=$1
    `, userID)
	if err != nil {
		r.logger.Error("DB query failed", slog.String("user_id", userID), slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var purchases []domain.Purchase
	for rows.Next() {
		var p domain.Purchase
		p.UserID = userID
		err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.EstimatedPrice, &p.ActualPrice,
			&p.CreatedAt, &p.Deadline, &p.IsActive, &p.IsPurchased)
		if err != nil {
			r.logger.Warn("failed to scan row", slog.String("user_id", userID), slog.String("error", err.Error()))
			continue
		}
		purchases = append(purchases, p)
	}

	r.logger.Debug("query complete", slog.Int("count", len(purchases)), slog.String("user_id", userID))
	return purchases, nil
}
