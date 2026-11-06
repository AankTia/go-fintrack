package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/AankTia/go-fintrack/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrResourceNotFound = errors.New("resource not found")
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
		return nil, ErrResourceNotFound
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
	// FIX: update_at -> updated_at
	query := `UPDATE accounts SET balance = balance + $1, updated_at = NOW() WHERE id = $2 AND user_id = $3`

	_, err := tx.Exec(ctx, query, amount, accountID, userID)
	return err
}

func (s *TransactionService) ListTransactions(ctx context.Context, userID string, req model.TransactionListRequest) ([]*model.TransactionResponse, model.PaginationMeta, error) {
	if req.Page <= 0 { req.Page = 1 }
	if req.Limit <= 0 { req.Limit = 10 }

	offset := (req.Page - 1) * req.Limit

	data, total, err := s.txRepo.List(ctx, userID, req, offset, req.Limit)
	if err != nil {
		return nil, model.PaginationMeta{}, err
	}

	meta := model.PaginationMeta{
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}

	return data, meta, nil
}

func (s *TransactionService) UpdateTransaction(ctx context.Context, userID string, transactionID string, req model.CreateTransactionRequest) (*model.Transaction, error) {
	// 1. Get old transaction
	oldTx, err := s.txRepo.GetById(ctx, transactionID, userID)
	if err != nil {
		return nil, ErrResourceNotFound
	}

	tx, err := s.db.Begin(ctx)
	if err != nil { return nil, err }
	defer tx.Rollback(ctx)

	// 2. REVERSE old balance
	oldBalanceChange := oldTx.Amount
	if oldTx.Type == "expense" { oldBalanceChange = -oldTx.Amount }
	if err := s.updateBalanceWithTx(ctx, tx, oldTx.AccountId, -oldBalanceChange, userID); err != nil {
		return nil, err
	}

	// 3. Update transaction record
	tDate, _ := time.Parse("2006-01-02", req.Date)
	updatedTx := &model.Transaction{
		ID:              transactionID,
		UserID:          userID,
		AccountId:       req.AccountId,
		CategoryId:      req.CategoryId,
		Amount:          req.Amount,
		Type:            req.Type,
		Description:     req.Description,
		TransactionDate: tDate,
	}
	if err := s.txRepo.Update(ctx, tx, updatedTx); err != nil {
		return nil, err
	}

	// 4. APPLY new balance
	newBalanceChange := req.Amount
	if req.Type == "expense" { newBalanceChange = -req.Amount }
	if err := s.updateBalanceWithTx(ctx, tx, req.AccountId, newBalanceChange, userID); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil { return nil, err }
	return updatedTx, nil
}

func (s *TransactionService) DeleteTransaction(ctx context.Context, userID string, transactionID string) error {
	oldTx, err := s.txRepo.GetById(ctx, transactionID, userID)
	if err != nil { return ErrResourceNotFound }

	tx, err := s.db.Begin(ctx)
	if err != nil { return err }
	defer tx.Rollback(ctx)

	// Reverse balance
	balanceChange := oldTx.Amount
	if oldTx.Type == "expense" { balanceChange = -oldTx.Amount }
	if err := s.updateBalanceWithTx(ctx, tx, oldTx.AccountId, -balanceChange, userID); err != nil {
		return err
	}

	// Delete record
	if err := s.txRepo.Delete(ctx, tx, transactionID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
