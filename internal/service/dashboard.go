package service

import (
	"context"
	"time"

	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/AankTia/go-fintrack/internal/repository"
)

type DashboardService interface {
	GetSummary(ctx context.Context, userID int64, mont int, year int) (*model.DashboardSummary, error)
}

type dashboardService struct {
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
}

func NewDashboardService(accRepo repository.AccountRepository, txRepo repository.TransactionRepository) DashboardService {
	return &dashboardService{
		accountRepo:     accRepo,
		transactionRepo: txRepo,
	}
}

func (s *dashboardService) GetSummary(ctx context.Context, userID int64, month int, year int) (*model.DashboardSummary, error) {
	balance, err := s.accountRepo.GetTotalBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	income, expense, err := s.transactionRepo.GetMonthlySummary(ctx, userID, month, year)
	if err != nil {
		return nil, err
	}

	return &model.DashboardSummary{
		TotalBalance: balance,
		TotalIncome:  income,
		TotalExpense: expense,
		NetCashFlow:  income - expense,
		ReportMonth:  time.Month(month),
		ReportYear:   year,
	}, nil
}
