package repository

import (
	"context"
	"fmt"

	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx pgx.Tx, t *model.Transaction) error
	GetById(ctx context.Context, id string, userID string) (*model.Transaction, error)
	List(ctx context.Context, userID string, filter model.TransactionListRequest, offset, limit int) ([]*model.TransactionResponse, int64, error)
	Update(ctx context.Context, tx pgx.Tx, t *model.Transaction) error
	Delete(ctx context.Context, tx pgx.Tx, id string) error
	GetMonthlySummary(ctx context.Context, userID string, month int, year int) (income int64, expense int64, err error)
	GetExpenseBreakdown(ctx context.Context, userID string, month int, year int) ([]*model.CategoryBreakdown, error)
	GetCashflowHistory(ctx context.Context, userID string) ([]*model.CashflowHistory, error)
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

func (r *transactionRepository) List(ctx context.Context, userID string, filter model.TransactionListRequest, offset, limit int) ([]*model.TransactionResponse, int64, error) {
	baseQuery := `FROM transactions t
							 	JOIN accounts a ON t.account_id = a.id
							 	JOIN categories c ON t.category_id = c.id
							 	WHERE t.user_id = $1`

	args := []interface{}{userID}
	argCount := 2

	if filter.Type != "" {
		args = append(args, filter.Type)
		baseQuery += fmt.Sprintf(" AND t.type = $%d", argCount)
		argCount++
	}
	if filter.AccountID != "" {
		args = append(args, filter.AccountID)
		baseQuery += fmt.Sprintf(" AND t.account_id = $%d", argCount)
		argCount++
	}
	if filter.StartDate != "" && filter.EndDate != "" {
		args = append(args, filter.StartDate, filter.EndDate)
		baseQuery += fmt.Sprintf(" AND t.transaction_date BETWEEN $%d AND $%d", argCount, argCount+1)
		argCount += 2
	}

	// Count total
	var total int64
	err := r.db.QueryRow(ctx, "SELECT count(*) "+baseQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Fetch data
	query := fmt.Sprintf("SELECT t.id, t.user_id, t.account_id, t.category_id, t.amount, t.type, t.description, t.transaction_date, a.name, c.name %s ORDER BY t.transaction_date DESC LIMIT %d OFFSET %d", baseQuery, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transactions []*model.TransactionResponse
	for rows.Next() {
		var tr model.TransactionResponse
		err := rows.Scan(&tr.ID, &tr.UserID, &tr.AccountId, &tr.CategoryId, &tr.Amount, &tr.Type, &tr.Description, &tr.TransactionDate, &tr.AccountName, &tr.CategoryName)
		if err != nil {
			return nil, 0, err
		}
		transactions = append(transactions, &tr)
	}

	return transactions, total, nil

}

func (r *transactionRepository) Update(ctx context.Context, tx pgx.Tx, t *model.Transaction) error {
	query := `UPDATE transactions SET account_id=$1, category_id=$2, amount=$3, type=$4, description=$5, transaction_date=$6, updated_at=NOW() WHERE id=$7`
	_, err := tx.Exec(ctx, query, t.AccountId, t.CategoryId, t.Amount, t.Type, t.Description, t.TransactionDate, t.ID)
	return err
}

func (r *transactionRepository) Delete(ctx context.Context, tx pgx.Tx, id string) error {
	query := `DELETE FROM transactions WHERE id=$1`
	_, err := tx.Exec(ctx, query, id)
	return err
}

func (r *transactionRepository) GetMonthlySummary(ctx context.Context, userID string, month int, year int) (income int64, expense int64, err error) {
	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0),
    	COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0)
		FROM transactions 
		WHERE user_id = $1 
		AND EXTRACT(MONTH FROM transaction_date) = $2 
		AND EXTRACT(YEAR FROM transaction_date) = $3`

	err = r.db.QueryRow(ctx, query, userID, month, year).Scan(&income, &expense)
	if err != nil {
		return 0, 0, err
	}
	return income, expense, nil
}

func (r *transactionRepository) GetExpenseBreakdown(ctx context.Context, userID string, month int, year int) ([]*model.CategoryBreakdown, error) {
	query := `SELECT
							c.id, c.name, SUM(t.amount) as total
						FROM transactions t
						JOIN categories c ON t.category_id = c.id
						WHERE t.user_id = $1
						AND t.type = 'expense'
						AND EXTRACT(MONTH FROM t.transaction_date) = $2
						AND EXTRACT(YEAR FROM t.transaction_date) = $3
						GROUP BY c.id, c.name
						ORDER BY total DESC`

	rows, err := r.db.Query(ctx, query, userID, month, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var breakdown []*model.CategoryBreakdown
	for rows.Next() {
		var cb model.CategoryBreakdown
		if err := rows.Scan(&cb.CategoryID, &cb.CategoryName, &cb.TotalAmount); err != nil {
			return nil, err
		}
		breakdown = append(breakdown, &cb)
	}

	return breakdown, nil
}

func (r *transactionRepository) GetCashflowHistory(ctx context.Context, userID string) ([]*model.CashflowHistory, error) {
    // Query untuk mengambil total income dan expense per bulan dalam 6 bulan terakhir
    // Menggunakan COALESCE untuk menangani bulan yang tidak ada transaksi
    query := `
        SELECT 
            to_char(transaction_date, 'Mon') as month, 
            SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END) as total_income,
            SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END) as total_expense
        FROM transactions 
        WHERE user_id = $1 
        AND transaction_date >= date_trunc('month', current_date) - interval '6 months'
        GROUP BY date_trunc('month', transaction_date), to_char(transaction_date, 'Mon')
        ORDER BY date_trunc('month', transaction_date) ASC
    `

    rows, err := r.db.Query(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var history []*model.CashflowHistory
    for rows.Next() {
        var h model.CashflowHistory
        if err := rows.Scan(&h.Month, &h.Income, &h.Expense); err != nil {
            return nil, err
        }
        history = append(history, &h)
    }
    return history, nil
}
