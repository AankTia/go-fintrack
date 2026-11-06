package service

import (
	"context"
	"fmt"
	"time"

	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/AankTia/go-fintrack/internal/repository"
)

type DashboardService interface {
	GetSummary(ctx context.Context, userID string, mont int, year int) (*model.DashboardSummary, error)
	GetBreakdown(ctx context.Context, userID string, month int, year int) ([]*model.CategoryBreakdown, error)
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

func (s *dashboardService) GetSummary(ctx context.Context, userID string, month int, year int) (*model.DashboardSummary, error) {
	fmt.Println("DEBUG: Testing accountRepo.GetTotalBalance...")
	balance, err := s.accountRepo.GetTotalBalance(ctx, userID)
	if err != nil {
		fmt.Printf("DEBUG: accountRepo ERROR: %v\n", err)
		return nil, err
	}
	fmt.Println("DEBUG: accountRepo SUCCESS")

	fmt.Printf("DEBUG: Calling transactionRepo.GetMonthlySummary with userID: %s, month: %d, year: %d\n", userID, month, year)
	income, expense, err := s.transactionRepo.GetMonthlySummary(ctx, userID, month, year)
	if err != nil {
		fmt.Printf("DEBUG: transactionRepo ERROR: %v\n", err)
		return nil, err
	}
	fmt.Println("DEBUG: transactionRepo SUCCESS")

	return &model.DashboardSummary{
		TotalBalance: balance,
		TotalIncome:  income,
		TotalExpense: expense,
		NetCashFlow:  income - expense,
		ReportMonth:  time.Month(month),
		ReportYear:   year,
	}, nil
}

func (s *dashboardService) GetBreakdown(ctx context.Context, userID string, month int, year int) ([]*model.CategoryBreakdown, error) {
	// Ambil data breakdown dari repo
	breakdowns, err := s.transactionRepo.GetExpenseBreakdown(ctx, userID, month, year)
	if err != nil {
		return nil, err
	}

	// Hitung total keseluruhan untuk mencari persentase
	var grandTotal int64
	for _, b := range breakdowns {
		grandTotal += b.TotalAmount
	}

	// Hitung persentase untuk tiap kategori
	if grandTotal > 0 {
		for _, b := range breakdowns {
			b.Percentage = (float64(b.TotalAmount) / float64(grandTotal)) * 100
		}
	}

	return breakdowns, nil
}
