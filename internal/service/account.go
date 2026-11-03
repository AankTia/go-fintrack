package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/AankTia/go-fintrack/internal/repository"
)

// Custom errors agar Handler bisa menentukan status code HTTP yang tepat
var (
	ErrAccountNotFound = errors.New("account not found")
	ErrInvalidInput    = errors.New("invalid input data")
)

type AccountService struct {
	repo repository.AccountRepository
}

// NewAccountService adalah constructor untuk membuat instance AccountService
func NewAccountService(repo repository.AccountRepository) *AccountService {
	return &AccountService{
		repo: repo,
	}
}

// CreateAccount menangani logika bisnis pembuatan akun baru
func (s *AccountService) CreateAccount(ctx context.Context, userID string, req model.AccountRequest) (*model.Account, error) {
	// Business Validation
	if req.Name == "" {
		return nil, fmt.Errorf("%w: account name is required", ErrInvalidInput)
	}

	// Mapping: Ubah AccountRequest (DTO) menjadi Account (Entity)
	// Menggabungkan data input user dengan userID dari session/token
	account := &model.Account{
		UserID:   userID,
		Name:     req.Name,
		Type:     req.Type,
		Balance:  req.Balance,
		Currency: req.Currency,
	}

	// Set default currency jika kosong
	if account.Currency == "" {
		account.Currency = "IDR"
	}

	// Panggil Repository untuk simpan ke DB
	err := s.repo.Create(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %v", err)
	}

	return account, nil
}

func (s *AccountService) GetMyAccounts(ctx context.Context, userID string) ([]*model.Account, error) {
	accounts, err := s.repo.GetByUserId(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch accounts: %v", err)
	}

	return accounts, nil
}

func (s *AccountService) GetAccountDetail(ctx context.Context, accountID string, userID string) (*model.Account, error) {
	acc, err := s.repo.GetById(ctx, accountID, userID)
	if err != nil {
		// Jika repo mengembalikan sql.ErrNoRows, kita bungkus jadi ErrAccountNotFound
		// agar handler bisa mengembalikan 404 Not Found
		return nil, ErrAccountNotFound
	}

	return acc, nil
}

func (s *AccountService) UpdateAccountBalance(ctx context.Context, accounID string, userID string, amount int64) error {
	// Logic Bisnis: Kita bisa cek dulu apakah saldo akhir akan menjadi negatif
	// (Jika aturan bisnis melarang saldo negatif)

	err := s.repo.UpdateBalance(ctx, accounID, amount, userID)
	if err != nil {
		return fmt.Errorf("failed to update balance: %v", err)
	}

	return nil
}
