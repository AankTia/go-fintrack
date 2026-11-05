package repository

import (
	"context"
	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx pgx.Tx, t *model.Transaction) error
	GetById(ctx context.Context, id string, userID string) (*model.Transaction, error)
}

type transactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx pgx.Tx, t *model.Transaction) error {
	query := `INSERT INTO transactions (user_id, account_id, category_id, amount, type, description, transaction_date)
						VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, query, t.UserID, t.AccountId, t.CategoryId, t.Amount, t.Type, t.Description, t.TransactionDate)
	} else {
		row = r.db.QueryRow(ctx, query, t.UserID, t.AccountId, t.CategoryId, t.Amount, t.Type, t.Description, t.TransactionDate)
	}

	return row.Scan(&t.ID)
}

func (r *transactionRepository) GetById(ctx context.Context, id string, userID string) (*model.Transaction, error) {
	query := `SELECT id, user_id, account_id, category_id, amount, type, description, transaction_date
						FROM transactions WHERE id = $1 and user_id = $2`

	var t model.Transaction
	err := r.db.QueryRow(ctx, query, id, userID).Scan(&t.ID, &t.UserID, &t.AccountId, &t.CategoryId, &t.Amount, &t.Type, &t.Description, &t.TransactionDate)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
