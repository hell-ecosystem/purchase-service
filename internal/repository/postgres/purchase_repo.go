package postgres

import (
	"context"

	"github.com/hell-ecosystem/purchase-service/internal/domain"
	"github.com/hell-ecosystem/purchase-service/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type purchaseRepo struct {
	db *pgxpool.Pool
}

func NewPurchaseRepo(db *pgxpool.Pool) repository.PurchaseRepository {
	return &purchaseRepo{db: db}
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
			return err
		}
	}
	return nil
}

func (r *purchaseRepo) FindByUser(ctx context.Context, userID string) ([]domain.Purchase, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, category, estimated_price, actual_price, created_at, deadline, is_active, is_purchased FROM purchases WHERE user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []domain.Purchase
	for rows.Next() {
		var p domain.Purchase
		p.UserID = userID
		err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.EstimatedPrice, &p.ActualPrice, &p.CreatedAt, &p.Deadline, &p.IsActive, &p.IsPurchased)
		if err != nil {
			continue
		}
		purchases = append(purchases, p)
	}
	return purchases, nil
}
