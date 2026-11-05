# WEEK 5: Budget Management

## Day 13: Budget Model & CRUD

### 1. Database Migration

Buat file migration baru. Saya asumsikan Anda menggunakan penamaan timestamp atau sequential.

SQL (Up):

```sql
CREATE TABLE budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    category_id UUID NOT NULL,
    amount BIGINT NOT NULL,
    period VARCHAR(7) NOT NULL, -- Format: YYYY-MM
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_category FOREIGN KEY(category_id) REFERENCES categories(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_category_period UNIQUE(user_id, category_id, period)
);

CREATE INDEX idx_budgets_user_period ON budgets(user_id, period);
```

SQL (Down):

```sql
DROP TABLE IF EXISTS budgets;
```

———

### 2. Model Layer (internal/model/budget.go)

Buat file baru di `internal/model/budget.go`.

```go
package model

import (
	"time"
	"github.com/google/uuid"
)

type Budget struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	CategoryID uuid.UUID `json:"category_id"`
	Amount     int64     `json:"amount"`
	Period     string    `json:"period"` // Format: YYYY-MM
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateBudgetRequest struct {
	CategoryID uuid.UUID `json:"category_id" validate:"required"`
	Amount     int64     `json:"amount" validate:"required,gt=0"`
	Period     string    `json:"period" validate:"required,len=7"` // YYYY-MM
}
```

———

### 3. Repository Layer (internal/repository/budget_repository.go)

Buat file baru di `internal/repository/budget_repository.go`.

```go
package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"go-fintrack/internal/model"
)

type BudgetRepository interface {
	Upsert(ctx context.Context, budget *model.Budget) error
	GetByPeriod(ctx context.Context, userID uuid.UUID, period string) ([]model.Budget, error)
	Delete(ctx context.Context, userID uuid.UUID, budgetID uuid.UUID) error
}

type budgetRepository struct {
	db *sql.DB
}

func NewBudgetRepository(db *sql.DB) BudgetRepository {
	return &budgetRepository{db: db}
}

func (r *budgetRepository) Upsert(ctx context.Context, b *model.Budget) error {
	query := `
		INSERT INTO budgets (user_id, category_id, amount, period, updated_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id, category_id, period) 
		DO UPDATE SET amount = EXCLUDED.amount, updated_at = CURRENT_TIMESTAMP
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(ctx, query, b.UserID, b.CategoryID, b.Amount, b.Period).
		Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
}

func (r *budgetRepository) GetByPeriod(ctx context.Context, userID uuid.UUID, period string) ([]model.Budget, error) {
	query := `SELECT id, user_id, category_id, amount, period, created_at, updated_at 
	          FROM budgets WHERE user_id = $1 AND period = $2`

	rows, err := r.db.QueryContext(ctx, query, userID, period)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []model.Budget
	for rows.Next() {
		var b model.Budget
		if err := rows.Scan(&b.ID, &b.UserID, &b.CategoryID, &b.Amount, &b.Period, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		budgets = append(budgets, b)
	}
	return budgets, nil
}

func (r *budgetRepository) Delete(ctx context.Context, userID uuid.UUID, budgetID uuid.UUID) error {
	query := `DELETE FROM budgets WHERE id = $1 AND user_id = $2`
	_, err := r.db.ExecContext(ctx, query, budgetID, userID)
	return err
}
```

———

### 4. Service Layer (internal/service/budget.go)

Buat file baru di `internal/service/budget.go`.

```go
package service

import (
	"context"
	"github.com/google/uuid"
	"go-fintrack/internal/model"
	"go-fintrack/internal/repository"
)

type BudgetService interface {
	SetBudget(ctx context.Context, userID uuid.UUID, req model.CreateBudgetRequest) (*model.Budget, error)
	GetBudgets(ctx context.Context, userID uuid.UUID, period string) ([]model.Budget, error)
}

type budgetService struct {
	repo repository.BudgetRepository
}

func NewBudgetService(repo repository.BudgetRepository) BudgetService {
	return &budgetService{repo: repo}
}

func (s *budgetService) SetBudget(ctx context.Context, userID uuid.UUID, req model.CreateBudgetRequest) (*model.Budget, error) {
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

func (s *budgetService) GetBudgets(ctx context.Context, userID uuid.UUID, period string) ([]model.Budget, error) {
	return s.repo.GetByPeriod(ctx, userID, period)
}
```

———

### 5. Handler Layer (internal/handler/budget.go)

Buat file baru di `internal/handler/budget.go`.

```go
package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go-fintrack/internal/model"
	"go-fintrack/internal/service"
)

type BudgetHandler struct {
	svc service.BudgetService
}

func NewBudgetHandler(svc service.BudgetService) *BudgetHandler {
	return &BudgetHandler{svc: svc}
}

func (h *BudgetHandler) CreateOrUpdate(c *gin.Context) {
	userID := uuid.MustParse(c.GetString("user_id"))
	var req model.CreateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	budget, err := h.svc.SetBudget(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set budget"})
		return
	}
	c.JSON(http.StatusOK, budget)
}

func (h *BudgetHandler) List(c *gin.Context) {
	userID := uuid.MustParse(c.GetString("user_id"))
	period := c.Query("period")
	if period == "" {
		// Default ke bulan sekarang jika kosong, misal: 2026-11
		c.JSON(http.StatusBadRequest, gin.H{"error": "period query parameter is required (YYYY-MM)"})
		return
	}

	budgets, err := h.svc.GetBudgets(c.Request.Context(), userID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch budgets"})
		return
	}
	c.JSON(http.StatusOK, budgets)
}
```

### Catatan untuk Registrasi di main.go atau Router:

Jangan lupa untuk melakukan inisialisasi dependency injection:

1. budgetRepo := repository.NewBudgetRepository(db)
2. budgetSvc := service.NewBudgetService(budgetRepo)
3. budgetHdl := handler.NewBudgetHandler(budgetSvc)
4. Tambahkan route:
   - v1.POST("/budgets", budgetHdl.CreateOrUpdate)
   - v1.GET("/budgets", budgetHdl.List)

## Day 14: Budget Calculation & Status

### 1. Model Layer (internal/model/budget.go)

Tambahkan struct baru untuk response status budget. Tambahkan ini di bagian bawah file `internal/model/budget.go` yang sudah Anda buat sebelumnya.

```go
type BudgetStatus struct {
	CategoryID   uuid.UUID `json:"category_id"`
	CategoryName string    `json:"category_name"`
	BudgetLimit  int64     `json:"budget_limit"`
	ActualSpend  int64     `json:"actual_spend"`
	Remaining    int64     `json:"remaining"`
	Percentage   float64   `json:"percentage"`
	IsOverlimit  bool      `json:"is_overlimit"`
}
```

———

### 2. Repository Layer (internal/repository/budget_repository.go)

Kita butuh query untuk mengambil total pengeluaran per kategori dalam satu periode. Tambahkan method berikut ke interface BudgetRepository dan implementasinya.

Tambahkan ke Interface:

```go
GetActualSpend(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID, period string) (int64, error)
```

Tambahkan Implementasi:

```go
func (r *budgetRepository) GetActualSpend(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID, period string) (int64, error) {
	var total int64
	// Menghitung sum amount dari transaksi dengan tipe 'expense'
	// Asumsi: Tabel transactions memiliki kolom amount, category_id, user_id, dan transaction_date
	query := `
		SELECT COALESCE(SUM(amount), 0) 
		FROM transactions 
		WHERE user_id = $1 
		AND category_id = $2 
		AND to_char(transaction_date, 'YYYY-MM') = $3 
		AND type = 'expense'`

	err := r.db.QueryRowContext(ctx, query, userID, categoryID, period).Scan(&total)
	return total, err
}
```

———

### 3. Service Layer (internal/service/budget.go)

Di sinilah "otak" perhitungannya berada. Kita akan menggabungkan data limit budget dengan data pengeluaran aktual.

Tambahkan ke Interface BudgetService:

GetBudgetStatus(ctx context.Context, userID uuid.UUID, period string) ([]model.BudgetStatus, error)

Tambahkan Implementasi:

func (s *budgetService) GetBudgetStatus(ctx context.Context, userID uuid.UUID, period string) ([]model.BudgetStatus, error) {
	// 1. Ambil semua budget yang diset user untuk periode ini
	budgets, err := s.repo.GetByPeriod(ctx, userID, period)
	if err != nil {
		return nil, err
	}

	// Kita butuh nama kategori, asumsikan Anda punya CategoryRepository 
	// Jika belum, Anda bisa join di repository, tapi di sini saya buatkan loop sederhana
	var statuses []model.BudgetStatus

	for _, b := range budgets {
		// 2. Hitung pengeluaran aktual untuk kategori ini
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
			IsOverlimit:  actual > b.Amount,
		})
	}

	return statuses, nil
}

Catatan: Untuk CategoryName, idealnya Anda menginject CategoryRepository ke dalam BudgetService untuk mengambil nama kategori berdasarkan ID.

———

### 4. Handler Layer (internal/handler/budget.go)

Tambahkan endpoint baru untuk menampilkan status budget.

Tambahkan Method Baru:

func (h *BudgetHandler) GetStatus(c *gin.Context) {
	userID := uuid.MustParse(c.GetString("user_id"))
	period := c.Query("period")
	if period == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "period query parameter is required (YYYY-MM)"})
		return
	}

	statuses, err := h.svc.GetBudgetStatus(c.Request.Context(), userID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate budget status"})
		return
	}
	c.JSON(http.StatusOK, statuses)
}

———

### 5. Registrasi Route

Tambahkan route baru ini di router Anda:

v1.GET("/budgets/status", budgetHdl.GetStatus)

### Checklist Verifikasi Manual:

Untuk memastikan implementasi ini benar, coba lakukan skenario berikut:

1. Skenario Normal:
   - Set Budget "Makanan" → Rp 1.000.000.
   - Buat Transaksi "Makanan" → Rp 200.000.
   - Panggil /api/v1/budgets/status?period=YYYY-MM.
   - Hasil yang diharapkan: ActualSpend: 200000, Remaining: 800000, Percentage: 20, IsOverlimit: false.
2. Skenario Overlimit:
   - Buat Transaksi "Makanan" lagi → Rp 900.000.
   - Panggil /api/v1/budgets/status?period=YYYY-MM.
   - Hasil yang diharapkan: ActualSpend: 1100000, Remaining: -100000, Percentage: 110, IsOverlimit: true.
3. Skenario Tanpa Transaksi:
   - Set Budget → Rp 500.000.
   - Jangan buat transaksi apapun.
   - Hasil yang diharapkan: ActualSpend: 0, Remaining: 500000, IsOverlimit: false.

## Day 15: Budget Integration & Final Testing

### 1. Integration Review (Backend)

Secara arsitektural, kalkulasi budget kita sudah bersifat on-the-fly (dihitung saat diminta lewat endpoint /status). Ini adalah pendekatan paling aman untuk mencegah data out-of-sync. Namun, ada satu hal kritikal yang perlu dipastikan: Cascade Delete.

Jika seorang user menghapus kategori, budget yang terkait dengan kategori tersebut harus ikut terhapus agar tidak menjadi data sampah (orphaned data).

Verifikasi Migration:
Pastikan saat membuat tabel budgets di Day 13, Anda menggunakan ON DELETE CASCADE pada foreign key category_id.

-- Cek apakah sudah seperti ini:
CONSTRAINT fk_category FOREIGN KEY(category_id) REFERENCES categories(id) ON DELETE CASCADE

Jika belum, jalankan migration alter table:

ALTER TABLE budgets DROP CONSTRAINT fk_category;
ALTER TABLE budgets ADD CONSTRAINT fk_category FOREIGN KEY(category_id) REFERENCES categories(id) ON DELETE CASCADE;

———

### 2. Security Audit (Ownership Check)

Ini adalah bagian paling krusial dari FinTrack. Anda harus memastikan tidak ada celah bagi User A untuk melihat atau mengubah budget User B.

Audit Point:

1. Upsert: Pastikan user_id diambil dari session/context (c.GetString("user_id")), bukan dari request body. (Kode yang saya berikan sebelumnya sudah melakukan ini).
2. GetByPeriod: Pastikan query mengandung WHERE user_id = $1.
3. GetActualSpend: Pastikan query transaksi menyertakan user_id = $1.

Manual Test Case (Security):

1. Login sebagai User A, ambil category_id milik User B.
2. Coba buat budget menggunakan category_id milik User B tersebut via API.
3. Hasil yang benar: Sistem harus error (karena category tidak ditemukan untuk User A) atau terjadi kegagalan validasi.

———

### 3. Performance Optimization (Indexing)

Karena endpoint /budgets/status melakukan agregasi (SUM) pada tabel transaksi, query ini bisa menjadi lambat jika transaksi sudah ribuan.

Optimasi Database:
Tambahkan index pada kolom yang digunakan di WHERE clause untuk perhitungan pengeluaran.

-- Index untuk mempercepat pencarian transaksi berdasarkan user, kategori, dan tanggal
CREATE INDEX idx_transactions_budget_calc ON transactions(user_id, category_id, transaction_date);

———

### 4. Final Integration Testing (End-to-End Flow)

Lakukan pengujian manual secara berurutan untuk memastikan seluruh vertical slice bekerja:

| Step | Action | Expected Result |
|:---|:---|:---|
| 1 | POST /api/v1/budgets → Limit Rp 500k | Budget tersimpan |
| 2 | GET /api/v1/budgets/status | Actual: 0, Remaining: 500k, Overlimit: false |
| 3 | POST /api/v1/transactions → Expense Rp 200k | Transaksi tersimpan |
| 4 | GET /api/v1/budgets/status | Actual: 200k, Remaining: 300k, Overlimit: false |
| 5 | POST /api/v1/transactions → Expense Rp 400k | Transaksi tersimpan |
| 6 | GET /api/v1/budgets/status | Actual: 600k, Remaining: -100k, Overlimit: true |
| 7 | POST /api/v1/budgets → Update Limit ke 700k | Budget ter-update |
| 8 | GET /api/v1/budgets/status | Actual: 600k, Remaining: 100k, Overlimit: false |

———

### 5. Laporan Akhir (Sesuai Format AGENTS.md)

Jika Anda sudah mengaplikasikan semua kode dari Day 13-15, Anda bisa mencatat progresnya seperti ini:

Selesai:
- Implementasi Budget Management (Week 5)
- Tabel `budgets` dengan UNIQUE constraint (user, category, period)
- Endpoint CRUD Budget & Budget Status (Calculation)
- Optimasi index untuk agregasi pengeluaran

Verifikasi:
- Test Upsert Budget: Berhasil (No duplicates)
- Test Calculation: Berhasil (Actual vs Limit akurat)
- Test Security: Berhasil (Ownership user_id terproteksi)
- Test Performance: Index `idx_transactions_budget_calc` ditambahkan

Catatan:
- Menggunakan BIGINT untuk amount guna menghindari floating point error.
- Menggunakan pendekatan on-the-fly calculation untuk Budget Status guna menjamin konsistensi data.

Selamat! Anda telah menyelesaikan modul Budget Management. Apakah ada bagian yang kurang jelas atau ingin Anda tambahkan?