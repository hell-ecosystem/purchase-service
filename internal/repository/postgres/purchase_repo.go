package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

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

func (r *purchaseRepo) MarkAsPurchased(ctx context.Context, userID, purchaseID string) error {
	cmd, err := r.db.Exec(ctx, `
		UPDATE purchases
		SET is_purchased = true, is_active = false
		WHERE id = $1 AND user_id = $2 AND is_active = true
	`, purchaseID, userID)
	if err != nil {
		r.logger.Error("failed to mark purchase as purchased", slog.String("purchase_id", purchaseID), slog.String("user_id", userID), slog.String("error", err.Error()))
		return err
	}
	if cmd.RowsAffected() == 0 {
		r.logger.Warn("no matching active purchase to mark", slog.String("purchase_id", purchaseID), slog.String("user_id", userID))
		return fmt.Errorf("purchase not found or already inactive")
	}
	r.logger.Debug("purchase marked as purchased", slog.String("purchase_id", purchaseID), slog.String("user_id", userID))
	return nil
}

func (r *purchaseRepo) Deactivate(ctx context.Context, userID, purchaseID string) error {
	cmd, err := r.db.Exec(ctx, `
		UPDATE purchases
		SET is_active = false
		WHERE id = $1 AND user_id = $2 AND is_active = true
	`, purchaseID, userID)
	if err != nil {
		r.logger.Error("failed to deactivate purchase", slog.String("purchase_id", purchaseID), slog.String("user_id", userID), slog.String("error", err.Error()))
		return err
	}
	if cmd.RowsAffected() == 0 {
		r.logger.Warn("no matching active purchase to deactivate", slog.String("purchase_id", purchaseID), slog.String("user_id", userID))
		return fmt.Errorf("purchase not found or already inactive")
	}
	r.logger.Debug("purchase deactivated", slog.String("purchase_id", purchaseID), slog.String("user_id", userID))
	return nil
}

func (r *purchaseRepo) Update(ctx context.Context, userID, purchaseID string, update domain.PurchaseUpdate) error {
	set := []string{}
	args := []interface{}{}
	idx := 1

	if update.Name != nil {
		set = append(set, fmt.Sprintf("name = $%d", idx))
		args = append(args, *update.Name)
		idx++
	}
	if update.CategoryID != nil {
		set = append(set, fmt.Sprintf("category_id = $%d", idx))
		args = append(args, *update.CategoryID)
		idx++
	}
	if update.EstimatedPrice != nil {
		set = append(set, fmt.Sprintf("estimated_price = $%d", idx))
		args = append(args, *update.EstimatedPrice)
		idx++
	}
	if update.Deadline != nil {
		set = append(set, fmt.Sprintf("deadline = $%d", idx))
		args = append(args, *update.Deadline)
		idx++
	}
	if len(set) == 0 {
		r.logger.Warn("update called with no fields", slog.String("purchase_id", purchaseID), slog.String("user_id", userID))
		return fmt.Errorf("no update fields provided")
	}

	args = append(args, purchaseID, userID)
	query := fmt.Sprintf(`UPDATE purchases SET %s WHERE id = $%d AND user_id = $%d`, strings.Join(set, ", "), idx, idx+1)

	cmd, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		r.logger.Error("failed to update purchase", slog.String("purchase_id", purchaseID), slog.String("user_id", userID), slog.String("error", err.Error()))
		return err
	}
	if cmd.RowsAffected() == 0 {
		r.logger.Warn("no matching purchase found to update", slog.String("purchase_id", purchaseID), slog.String("user_id", userID))
		return fmt.Errorf("purchase not found")
	}

	r.logger.Debug("purchase updated", slog.String("purchase_id", purchaseID), slog.String("user_id", userID))
	return nil
}

func (r *purchaseRepo) CollectStats(ctx context.Context, userID string) (domain.PurchaseStats, error) {
	var stats domain.PurchaseStats
	row := r.db.QueryRow(ctx, `
		SELECT 
			COUNT(*) FILTER (WHERE user_id = $1) AS total,
			COUNT(*) FILTER (WHERE user_id = $1 AND is_active = true) AS active,
			COUNT(*) FILTER (WHERE user_id = $1 AND is_purchased = true) AS purchased,
			COALESCE(SUM(estimated_price), 0) FILTER (WHERE user_id = $1) AS estimated_total,
			COALESCE(SUM(actual_price), 0) FILTER (WHERE user_id = $1 AND is_purchased = true) AS actual_total
		FROM purchases
	`, userID)

	err := row.Scan(&stats.Total, &stats.Active, &stats.Purchased, &stats.EstimatedTotal, &stats.ActualTotal)
	if err != nil {
		r.logger.Error("failed to collect stats", slog.String("user_id", userID), slog.String("error", err.Error()))
		return domain.PurchaseStats{}, err
	}

	r.logger.Debug("collected stats", slog.String("user_id", userID), slog.Int("total", stats.Total))
	return stats, nil
}

func (r *purchaseRepo) InsertReminder(ctx context.Context, userID, purchaseID string, reminder domain.Reminder) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO reminders (purchase_id, remind_at, message)
		SELECT $1, $2, $3
		WHERE EXISTS (SELECT 1 FROM purchases WHERE id = $1 AND user_id = $4)
	`, purchaseID, reminder.RemindAt, reminder.Message, userID)
	if err != nil {
		r.logger.Error("failed to insert reminder", slog.String("purchase_id", purchaseID), slog.String("user_id", userID), slog.String("error", err.Error()))
		return err
	}
	r.logger.Debug("reminder inserted", slog.String("purchase_id", purchaseID), slog.String("user_id", userID))
	return nil
}

func (r *purchaseRepo) AttachTags(ctx context.Context, userID, purchaseID string, tagIDs []string) error {
	batch := &strings.Builder{}
	for _, tagID := range tagIDs {
		fmt.Fprintf(batch, `INSERT INTO purchase_tags (purchase_id, tag_id) SELECT '%s', '%s' WHERE EXISTS (SELECT 1 FROM purchases WHERE id = '%s' AND user_id = '%s');`, purchaseID, tagID, purchaseID, userID)
	}
	_, err := r.db.Exec(ctx, batch.String())
	if err != nil {
		r.logger.Error("failed to attach tags", slog.String("purchase_id", purchaseID), slog.String("user_id", userID), slog.String("error", err.Error()))
		return err
	}
	r.logger.Debug("tags attached", slog.String("purchase_id", purchaseID), slog.String("user_id", userID), slog.Int("tag_count", len(tagIDs)))
	return nil
}
