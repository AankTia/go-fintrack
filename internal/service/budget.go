package service

import (
	"context"

	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/AankTia/go-fintrack/internal/repository"
)

type BudgetService interface {
	SetBudget(ctx context.Context, userID string, req model.CreateBudgetRequest) (*model.Budget, error)
	GetBudget(ctx context.Context, userID string, period string) ([]model.Budget, error)
}

type budgetService struct {
	repo repository.BudgetRepository
}

func NewBudgetService(repo repository.BudgetRepository) BudgetService {
	return &budgetService{repo: repo}
}

func (s *budgetService) SetBudget(ctx context.Context, userID string, req model.CreateBudgetRequest) (*model.Budget, error) {
	budget := &model.Budget{
		UserID:     userID,
		CategoryID: req.CategoryID,
		Amount:     req.Amount,
		Period:     req.Period,
	}

	if err := s.repo.Upsert(ctx, budget); err != nil {
		return nil, err
	}

	return budget, nil
}

func (s *budgetService) GetBudget(ctx context.Context, userID string, period string) ([]model.Budget, error) {
	return s.repo.GetByPeriod(ctx, userID, period)
}
