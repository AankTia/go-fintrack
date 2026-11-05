package repository

import (
	"context"

	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BudgetRepository interface {
	Upsert(ctx context.Context, budget *model.Budget) error
	GetByPeriod(ctx context.Context, userID string, period string) ([]model.Budget, error)
	Delete(ctx context.Context, userID string, budgetID string) error
	GetActualSpend(ctx context.Context, userID string, categoryID string, period string) (int64, error)
}

type budgetRepository struct {
	db *pgxpool.Pool
}

func NewBudgetRepository(db *pgxpool.Pool) BudgetRepository {
	return &budgetRepository{db: db}
}

func (r *budgetRepository) Upsert(ctx context.Context, b *model.Budget) error {
	query := `INSERT INTO budgets (user_id, category_id, amount, period, updated_at)
						VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
						ON CONFLICT (user_id, category_id, period)
						DO UPDATE SET amount = EXCLUDED.amount, updated_at = CURRENT_TIMESTAMP
						RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query, b.UserID, b.CategoryID, b.Amount, b.Period).Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
}

func (r *budgetRepository) GetByPeriod(ctx context.Context, userId string, period string) ([]model.Budget, error) {
	query := `SELECT id, user_id, category_id, amount, period, created_at, updated_at
						FROM budgets WHERE user_id = $1 AND period = $2`

	rows, err := r.db.Query(ctx, query, userId, period)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []model.Budget
	for rows.Next() {
		var b model.Budget
		if err := rows.Scan(&b.ID, &b.UserID, &b.CategoryID, &b.Amount, &b.Period, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		budgets = append(budgets, b)
	}

	return budgets, nil
}

func (r *budgetRepository) Delete(ctx context.Context, userID string, budgetID string) error {
	query := `DELETE FROM budgets WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, budgetID, userID)
	return err
}

func (r *budgetRepository) GetActualSpend(ctx context.Context, userID string, categoryID string, period string) (int64, error) {
	var total int64

	// Menghitung sum amount dari transaksi dengan tipe 'expense'
	query := `SELECT COALASCE(SUM(amount), 0)
						FROM transactions
						WHERE user_id = $1
						AND category_id = $2
						AND to_char(transaction_date, 'YYYY-MM') = $3
						AND type = 'expense'`

	err := r.db.QueryRow(ctx, query, userID, categoryID, period).Scan(&total)
	return total, err
}
