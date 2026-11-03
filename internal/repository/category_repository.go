package repository

import (
	"context"
	"github.com/AankTia/go-fintrack/internal/model"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error

	GetByUserId(ctx context.Context, userID string) ([]*model.Category, error)

	Delete(ctx context.Context, id string) error
}
