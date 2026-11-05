# 🗓️ WEEK 4: Dashboard & Analytics

### 🟦 Day 10: Dashboard Summary API

Goal: Membuat endpoint `GET /api/v1/dashboard/summary` untuk mendapatkan ringkasan keuangan.

#### Step 1: Tambahkan Model

Buat struct untuk menampung data summary.
File: `internal/model/dashboard.go` (Buat file baru)

```go
package model

import "time"

type DashboardSummary struct {
	TotalBalance   int64     `json:"total_balance"`   // Total saldo semua akun
	TotalIncome    int64     `json:"total_income"`     // Pemasukan bulan ini
	TotalExpense   int64     `json:"total_expense"`    // Pengeluaran bulan ini
	NetCashFlow    int64     `json:"net_cash_flow"`    // Income - Expense
	ReportMonth    time.Month `json:"report_month"`
	ReportYear     int       `json:"report_year"`
}
```

#### Step 2: Implementasi Repository (Aggregation Logic)

Kita akan menambahkan method di `TransactionRepository` dan `AccountRepository`.

File: `internal/repository/account_repository.go`
Tambahkan method `GetTotalBalance`:

```go
func (r *accountRepository) GetTotalBalance(ctx context.Context, userID int64) (int64, error) {
	var total int64
	query := `SELECT COALESCE(SUM(balance), 0) FROM accounts WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&total)
	return total, err
}
```

File: `internal/repository/transaction_repository.go`
Tambahkan method `GetMonthlySummary`:

``` go
func (r *transactionRepository) GetMonthlySummary(ctx context.Context, userID int64, month int, year int) (income int64, expense int64, err error) {
	query := `
		SELECT 
			SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END),
			SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END)
		FROM transactions 
		WHERE user_id = $1 
		AND EXTRACT(MONTH FROM transaction_date) = $2 
		AND EXTRACT(YEAR FROM transaction_date) = $3`

	err = r.db.QueryRowContext(ctx, query, userID, month, year).Scan(&income, &expense)
	if err != nil {
		return 0, 0, err
	}
	return income, expense, nil
}
```

#### Step 3: Implementasi Service

Buat service untuk mengorkestrasi data dari kedua repository tersebut.
File: `internal/service/dashboard.go` (Buat file baru)

```go
package service

import (
	"context"
	"time"
	"go-fintrack/internal/model"
	"go-fintrack/internal/repository"
)

type DashboardService interface {
	GetSummary(ctx context.Context, userID int64, month int, year int) (*model.DashboardSummary, error)
}

type dashboardService struct {
	accountRepo    repository.AccountRepository
	transactionRepo repository.TransactionRepository
}

func NewDashboardService(accRepo repository.AccountRepository, txRepo repository.TransactionRepository) DashboardService {
	return &dashboardService{
		accountRepo:    accRepo,
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
		ReportMonth:   time.Month(month),
		ReportYear:    year,
	}, nil
}
```

#### Step 4: Implementasi Handler

Buat handler untuk mengekspos API.
File: `internal/handler/dashboard.go` (Buat file baru)

```go
package handler

import (
	"net/http"
	"strconv"
	"time"
	"go-fintrack/internal/service"
	"github.com/gin-gonic/gin" // Asumsi menggunakan Gin, sesuaikan dengan framework Anda
)

type DashboardHandler struct {
	service service.DashboardService
}

func NewDashboardHandler(s service.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: s}
}

func (h *DashboardHandler) GetSummary(c *gin.Context) {
	userID, _ := strconv.ParseInt(c.GetUint64("userID"), 10, 64) // Ambil dari middleware auth

	monthStr := c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month())))
	yearStr := c.DefaultQuery("year", strconv.Itoa(time.Now().Year()))

	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)

	summary, err := h.service.GetSummary(c.Request.Context(), userID, month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch summary"})
		return
	}

	c.JSON(http.StatusOK, summary)
}
```

#### Step 5: Registrasi di `main.go`

Jangan lupa hubungkan semua layer di file entry point Anda.

File: `cmd/api/main.go` (atau tempat setup DI Anda)

```go
// 1. Init Service
dashboardSvc := service.NewDashboardService(accountRepo, transactionRepo)

// 2. Init Handler
dashboardHdl := handler.NewDashboardHandler(dashboardSvc)

// 3. Route (Dalam group /api/v1 dan middleware Auth)
api.GET("/dashboard/summary", dashboardHdl.GetSummary)
```

———

### 🧪 Cara Verifikasi Manual:

1. Skenario Normal: Panggil GET /api/v1/dashboard/summary. Pastikan total_balance adalah jumlah dari semua saldo akun Anda.
2. Skenario Filter: Panggil GET /api/v1/dashboard/summary?month=1&year=2024. Pastikan income/expense hanya menghitung transaksi di bulan Januari 2024.
3. Skenario No Data: User baru tanpa transaksi harus menerima angka 0, bukan error.

Catatan Keamanan:
Sesuai AGENTS.md, saya memastikan userID diambil dari context (c.GetUint64("userID")) yang diisi oleh middleware auth, bukan dari query parameter, untuk mencegah akses data user lain.