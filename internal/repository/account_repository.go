package repository

import (
	"context"
	"database/sql"
	"github.com/AankTia/go-fintrack/internal/model"
)

type AccountRepository interface {
	Create(ctx context.Context, account *model.Account) error

	GetById(ctx context.Context, id string, userID string) (*model.Account, error)

	GetByUserId(ctx context.Context, userID string) ([]*model.Account, error)

	UpdateBalance(ctx context.Context, id string, amount int64, userID string) error
}

// Definisikan struct yang menyimpan dependency (database connection)
// Kita beri nama kecil (unexported) agar hanya bisa diinisialisasi melalui NewAccountRepositor
type accountRepository struct {
	db *sql.DB
}

// Fungsi ini mengembalikan tipe interface AccountRepository, bukan tipe struct accountRepository.
// Ini adalah praktek Dependency Injection agar layer atas tidak terikat pada implementasi konkret.
func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepository{
		db: db,
	}
}

func (r *accountRepository) Create(ctx context.Context, account *model.Account) error {
	// Gunakan query parameterized ($1, $2, dst) untuk mencegah SQL Injection
	query := `INSERT INTO accounts (user_id, name, type, balance, currency)
						VALUE ($1, $2, $3, $4, $5)
						RETURNING id`

	// ExecQuery digunakan karena kita butuh return value (id)
	err := r.db.QueryRowContext(ctx, query,
		account.UserID,
		account.Name,
		account.Type,
		account.Balance,
		account.Currency,
	).Scan(&account.ID) // Mengambil ID yang baru dibuat dari DB dan memasukkannya ke struct account

	return err
}

func (r *accountRepository) GetById(ctx context.Context, id string, userID string) (*model.Account, error) {
	query := `SELECT id, user_id, name, type, balance, currency, created_at, updated_at
						FROM accounts WHERE id = $1 AND user_id $2`

	row := r.db.QueryRowContext(ctx, query, id, userID)

	var acc model.Account
	err := row.Scan(&acc.ID, &acc.UserID, &acc.Name, &acc.Type, &acc.Balance, &acc.Currency, &acc.CreatedAt, &acc.UpdatedAt)

	if err != nil {
		return nil, err  // Akan return sql.ErrNoRows jika data tidak ditemukan
	}

	return &acc, nil
}

func (r *accountRepository) GetByUserId(ctx context.Context, userID string) ([]*model.Account, error) {
	query := `SELECT id, user_id, name, type, balance, currency, created_at, updated_at 
	          FROM accounts WHERE user_id = $1`
	
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // WAJIB tutup rows untuk mencegah memory leak

	var accounts []*model.Account
	for rows.Next() {
		var acc model.Account
		if err := rows.Scan(&acc.ID, &acc.UserID, &acc.Name, &acc.Type, &acc.Balance, &acc.Currency, &acc.CreatedAt, &acc.UpdatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, &acc)
	}

	return accounts, nil
}

func (r *accountRepository) UpdateBalance(ctx context.Context, id string, amount int64, userId string) error {
	// Update saldo relatif (balance = balance + amount)
	// Lagi-lagi, tambahkan filter user_id untuk keamanan
	query := `UPDATE accounts SET balance = balance + $1, updated_at = NOW() 
	          WHERE id = $2 AND user_id = $3`
	
	_, err := r.db.ExecContext(ctx, query, amount, id, userId)
	return  err
}
