package service

import (
	"context"

	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/AankTia/go-fintrack/internal/repository"
)

type BudgetService interface {
	SetBudget(ctx context.Context, userID string, req model.CreateBudgetRequest) (*model.Budget, error)
	GetBudgets(ctx context.Context, userID string, period string) ([]model.Budget, error)
	GetBudgetStatus(ctx context.Context, userID string, period string) ([]model.BudgetStatus, error)
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

func (s *budgetService) GetBudgets(ctx context.Context, userID string, period string) ([]model.Budget, error) {
	return s.repo.GetByPeriod(ctx, userID, period)
}

func (s *budgetService) GetBudgetStatus(ctx context.Context, userID string, period string) ([]model.BudgetStatus, error) {
	budgets, err := s.repo.GetByPeriod(ctx, userID, period)
	if err != nil {
		return nil, err
	}

	var statuses []model.BudgetStatus

	for _, b := range budgets {
		actual, err := s.repo.GetActualSpend(ctx, userID, b.CategoryID, period)
		if err != nil {
			return nil, err
		}

		remaining := b.Amount - actual
		percentage := 0.0
		if b.Amount > 0 {
			percentage = (float64(actual) / float64(b.Amount)) * 100
		}

		statuses = append(statuses, model.BudgetStatus{
			CategoryID:  b.CategoryID,
			BudgetLimit: b.Amount,
			ActualSpend: actual,
			Remaining:   remaining,
			Percentage:  percentage,
			IsOverlimit: actual > b.Amount,
		})
	}

	return statuses, nil
}
