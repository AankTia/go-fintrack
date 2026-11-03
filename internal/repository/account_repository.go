package repository

import (
	"context"

	"github.com/AankTia/go-fintrack/internal/model"
)

type AccountRepository interface {
	Create(ctx context.Context, account *model.Account) error

	GetById(ctx context.Context, id string) (*model.Account, error)

	GetByUserId(ctx context.Context, userID string) ([]*model.Account, error)

	UpdateBalance(ctx context.Context, id string, amount int64) error
}
