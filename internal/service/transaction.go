package service

import (
	"context"
	"fmt"
	"time"

	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/AankTia/go-fintrack/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionService struct {
	db           *pgxpool.Pool
	txRepo       repository.TransactionRepository
	accountRepo  repository.AccountRepository
	categoryRepo repository.CategoryRepository
}

func NewTransactionService(db *pgxpool.Pool, txRepo repository.TransactionRepository, accountRepo repository.AccountRepository, categoryRepo repository.CategoryRepository) *TransactionService {
	return &TransactionService{
		db:           db,
		txRepo:       txRepo,
		accountRepo:  accountRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, userID string, req model.CreateTransactionRequest) (*model.Transaction, error) {
	// Parse date
	tDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		tDate = time.Now()
	}

	// ownership check
	_, err = s.accountRepo.GetById(ctx, req.AccountId, userID)
	if err != nil {
		return nil, fmt.Errorf("account not found or unauthorized")
	}

	// Start DB transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) // Rollback otomatis jika terjadi panic atau error sebelum Commit

	// Create transaction record
	transaction := &model.Transaction{
		UserID:          userID,
		AccountId:       req.AccountId,
		CategoryId:      req.CategoryId,
		Amount:          req.Amount,
		Type:            req.Type,
		Description:     req.Description,
		TransactionDate: tDate,
	}

	if err := s.txRepo.Create(ctx, tx, transaction); err != nil {
		return nil, fmt.Errorf("failed to save transaction: %v", err)
	}

	// Update account balance
	balanceChange := req.Amount
	if req.Type == "expense" {
		balanceChange = -req.Amount
	}

	// Note: UpdateBalance di repository harus bisa menerima tx atau menggunakan logic update relatif
	// Karena current AccountRepository.UpdateBalance menggunakan r.db, kita perlu modifikasi repo
	// atau membuat method baru yang menerima pgx.Tx.
	// Untuk sementara, kita asumsikan ada method UpdateBalanceTx di repository.
	if err := s.updateBalanceWithTx(ctx, tx, req.AccountId, balanceChange, userID); err != nil {
		return nil, fmt.Errorf("failed to update balance: %v", err)
	}

	// Commit
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *TransactionService) updateBalanceWithTx(ctx context.Context, tx pgx.Tx, accountID string, amount int64, userID string) error {
	query := `UPDATE accounts SET balance = balance + $1, update_at = NOW() WHERE id = $2 AND user_id = $3`

	_, err := tx.Exec(ctx, query, amount, accountID, userID)
	return err
}
