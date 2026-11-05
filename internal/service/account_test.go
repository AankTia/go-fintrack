package service

import (
	"context"
	"testing"
	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAccountRepository adalah mock untuk AccountRepository
type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) Create(ctx context.Context, account *model.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockAccountRepository) GetById(ctx context.Context, id string, userID string) (*model.Account, error) {
	args := m.Called(ctx, id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Account), args.Error(1)
}

func (m *MockAccountRepository) GetByUserId(ctx context.Context, userID string) ([]*model.Account, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*model.Account), args.Error(1)
}

func (m *MockAccountRepository) UpdateBalance(ctx context.Context, id string, amount int64, userID string) error {
	args := m.Called(ctx, id, amount, userID)
	return args.Error(0)
}

func (m *MockAccountRepository) Delete(ctx context.Context, id string, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func TestCreateAccount(t *testing.T) {
	mockRepo := new(MockAccountRepository)
	svc := NewAccountService(mockRepo)
	ctx := context.Background()
	userID := "user-123"

	t.Run("success create account", func(t *testing.T) {
		req := model.AccountRequest{
			Name:     "Bank BCA",
			Type:     "bank",
			Balance:  100000,
			Currency: "IDR",
		}

		// Set expectation: jika Create dipanggil, kembalikan nil (no error)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.Account")).Return(nil).Once()

		acc, err := svc.CreateAccount(ctx, userID, req)

		assert.NoError(t, err)
		assert.Equal(t, "Bank BCA", acc.Name)
		assert.Equal(t, userID, acc.UserID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("fail empty name", func(t *testing.T) {
		req := model.AccountRequest{Name: ""} // Nama kosong

		acc, err := svc.CreateAccount(ctx, userID, req)

		assert.Error(t, err)
		assert.Nil(t, acc)
		assert.Contains(t, err.Error(), "account name is required")
	})
}

func TestDeleteAccount(t *testing.T) {
	mockRepo := new(MockAccountRepository)
	svc := NewAccountService(mockRepo)
	ctx := context.Background()
	userID := "user-123"

	t.Run("fail delete only remaining account", func(t *testing.T) {
		// Mock: User hanya punya 1 akun
		accounts := []*model.Account{{ID: "acc-1", UserID: userID}}
		mockRepo.On("GetByUserId", ctx, userID).Return(accounts, nil).Once()

		err := svc.DeleteAccount(ctx, "acc-1", userID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot delete the only remaining account")
	})

	t.Run("success delete when multiple accounts exist", func(t *testing.T) {
		// Mock: User punya 2 akun
		accounts := []*model.Account{{ID: "acc-1"}, {ID: "acc-2"}}
		mockRepo.On("GetByUserId", ctx, userID).Return(accounts, nil).Once()
		mockRepo.On("Delete", ctx, "acc-1", userID).Return(nil).Once()

		err := svc.DeleteAccount(ctx, "acc-1", userID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}