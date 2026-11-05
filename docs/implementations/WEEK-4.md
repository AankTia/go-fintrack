# 🗓️ WEEK 4: Dashboard & Analytics

## Day 10: Dashboard Summary API

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

## Day 11: Category Breakdown & Charts Data

### Step 1: Definisikan Model (Data Contract)

Buat model baru di `internal/model/dashboard.go` untuk menampung hasil breakdown. Kita perlu menyimpan nama kategori agar frontend tidak perlu melakukan lookup manual.

Tambah di `internal/model/dashboard.go`:

```go
type CategoryBreakdown struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	TotalAmount  int64   `json:"total_amount"`
	Percentage   float64 `json:"percentage"` // Persentase terhadap total pengeluaran bulanan
}
```

———

### Step 2: Implementasi Repository (Data Access)

Tambahkan method baru di `TransactionRepository` untuk mengambil total pengeluaran yang dikelompokkan berdasarkan kategori.

1. Update Interface di `internal/repository/transaction_repository.go`:

```go
type TransactionRepository interface {
    // ... method lainnya
    GetExpenseBreakdown(ctx context.Context, userID string, month int, year int) ([]*model.CategoryBreakdown, error)
}
```

2. Implementasi Method di transactionRepository:
Gunakan `GROUP BY` dan `JOIN` dengan tabel `categories` untuk mendapatkan nama kategori.

```go
func (r *transactionRepository) GetExpenseBreakdown(ctx context.Context, userID string, month int, year int) ([]*model.CategoryBreakdown, error) {
	query := `
		SELECT 
			c.id, c.name, SUM(t.amount) as total
		FROM transactions t
		JOIN categories c ON t.category_id = c.id
		WHERE t.user_id = $1 
		AND t.type = 'expense'
		AND EXTRACT(MONTH FROM t.transaction_date) = $2 
		AND EXTRACT(YEAR FROM t.transaction_date) = $3
		GROUP BY c.id, c.name
		ORDER BY total DESC`

	rows, err := r.db.Query(ctx, query, userID, month, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var breakdown []*model.CategoryBreakdown
	for rows.Next() {
		var cb model.CategoryBreakdown
		if err := rows.Scan(&cb.CategoryID, &cb.CategoryName, &cb.TotalAmount); err != nil {
			return nil, err
		}
		breakdown = append(breakdown, &cb)
	}

	return breakdown, nil
}
```

———

### Step 3: Implementasi Service (Business Logic)

Service akan memanggil repository, lalu menghitung persentase untuk setiap kategori berdasarkan total pengeluaran bulan tersebut.

1. Update Interface di `internal/service/dashboard.go`:

```go
type DashboardService interface {
    GetSummary(ctx context.Context, userID string, month int, year int) (*model.DashboardSummary, error)
    GetBreakdown(ctx context.Context, userID string, month int, year int) ([]*model.CategoryBreakdown, error)
}
```

2. Implementasi Method `GetBreakdown`:

```go
func (s *dashboardService) GetBreakdown(ctx context.Context, userID string, month int, year int) ([]*model.CategoryBreakdown, error) {
	// 1. Ambil data breakdown dari repo
	breakdowns, err := s.transactionRepo.GetExpenseBreakdown(ctx, userID, month, year)
	if err != nil {
		return nil, err
	}

	// 2. Hitung total keseluruhan untuk mencari persentase
	var grandTotal int64
	for _, b := range breakdowns {
		grandTotal += b.TotalAmount
	}

	// 3. Hitung persentase untuk tiap kategori
	if grandTotal > 0 {
		for _, b := range breakdowns {
			b.Percentage = (float64(b.TotalAmount) / float64(grandTotal)) * 100
		}
	}

	return breakdowns, nil
}
```

———

### Step 4: Implementasi Handler (API Layer)

Ekspos fungsi service ke endpoint HTTP.

Tambah method di `internal/handler/dashboard.go`:

```go
func (h *DashboardHandler) GetBreakdown(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Logic ambil month & year dari context/query param (sama seperti GetSummary)
	// ... (copy logic month/year parsing dari GetSummary) ...

	breakdown, err := h.service.GetBreakdown(r.Context(), userID, month, year)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondWithJSON(w, http.StatusOK, breakdown)
}
```

———

### Step 5: Registrasi Route

Jangan lupa untuk mendaftarkan handler baru ini di file routing Anda (biasanya di internal/platform/router atau cmd/api/main.go).

```go
// Contoh registrasi route
router.GET("/api/v1/dashboard/breakdown", dashboardHandler.GetBreakdown)
```

———

### ✅ Checklist Verifikasi (Acceptance Criteria)

Setelah implementasi, lakukan pengujian berikut:

1. Data Accuracy: Input beberapa transaksi pengeluaran di kategori berbeda → Panggil /breakdown → Pastikan TotalAmount per kategori sesuai dan Percentage jika dijumlahkan hasilnya 100%.
2. Security (Ownership): Coba panggil endpoint dengan token User A, pastikan tidak ada data User B yang muncul (sudah tercover oleh WHERE user_id = $1 di repo).
3. Empty State: Jika bulan tersebut tidak ada pengeluaran, pastikan API mengembalikan array kosong [], bukan error 500.
4. Sorting: Pastikan kategori dengan nominal terbesar muncul paling atas.

## Day 12: Security & Performance Testing

### Step 1: Ownership Security Audit (Pengetesan Celah Data)

Tujuannya adalah memastikan Strict Isolation. User A tidak boleh bisa melihat summary atau breakdown User B, meskipun mereka tahu ID resource-nya.

Langkah Verifikasi:

1. Uji Token Validasi: Pastikan semua endpoint di /api/v1/dashboard/* menggunakan middleware Auth yang menyuntikkan user_id ke dalam context.
2. Audit Query SQL: Periksa kembali semua query di internal/repository/dashboard.go atau repository terkait.
   - Checklist: Apakah semua query SELECT memiliki klausa WHERE user_id = $1?
   - Anti-Pattern: Hindari query yang menerima userID langsung dari query parameter ?user_id=... tanpa divalidasi terhadap session token.
3. Manual Pengetesan (Postman):
   - Login sebagai User A → Ambil Token A.
   - Login sebagai User B → Ambil Token B.
   - Coba panggil endpoint dashboard menggunakan Token A, tetapi coba manipulasi parameter (jika ada) untuk meminta data User B. Hasilnya harus tetap data User A atau 403 Forbidden.

———

### Step 2: Query Optimization (Analisis Performa DB)

Query agregasi (SUM, GROUP BY) bisa menjadi sangat lambat jika tabel transactions memiliki jutaan baris.

Langkah Analisis:

1. Gunakan EXPLAIN ANALYZE:
   Jalankan query dashboard langsung di database (psql/pgAdmin) dengan prefix EXPLAIN ANALYZE.

   EXPLAIN ANALYZE 
SELECT SUM(amount) FROM transactions 
WHERE user_id = 'some-uuid' 
AND EXTRACT(MONTH FROM transaction_date) = 11 
AND EXTRACT(YEAR FROM transaction_date) = 2026;
2. Identifikasi "Seq Scan":
   Jika hasil EXPLAIN menunjukkan "Seq Scan" (Sequential Scan), artinya database membaca seluruh tabel dari awal sampai akhir. Ini adalah tanda bahaya untuk performa.

Langkah Perbaikan (Indexing):
Jika terjadi Seq Scan, tambahkan Composite Index. Index ini mempercepat filter berdasarkan user dan tanggal secara bersamaan.

-- Tambahkan migration baru untuk index ini
CREATE INDEX idx_transactions_user_date ON transactions (user_id, transaction_date);

Setelah menambahkan index, jalankan kembali EXPLAIN ANALYZE dan pastikan statusnya berubah menjadi "Index Scan".

———

### Step 3: Basic Load Testing (Stress Test)

Memastikan endpoint tetap responsif saat user memiliki riwayat transaksi yang panjang (misal: user aktif selama 2 tahun dengan 10 transaksi per hari).

Langkah Eksekusi:

1. Generate Data Dummy:
   Buat script sederhana atau gunakan SQL untuk memasukkan 1.000 - 5.000 transaksi dummy untuk satu user.

   -- Contoh pseudo-SQL untuk generate data
INSERT INTO transactions (user_id, account_id, category_id, amount, type, transaction_date)
SELECT 'user-uuid', 'acc-uuid', 'cat-uuid', 1000, 'expense', 
       CURRENT_DATE - (n || ' days')::interval
FROM generate_series(1, 5000) n;
2. Ukur Latency:
   Panggil endpoint /api/v1/dashboard/summary dan /api/v1/dashboard/breakdown berulang kali menggunakan Postman atau tool seperti k6 atau hey.
3. Target:
   Respons time harus < 200ms. Jika > 500ms, kembali ke Step 2 untuk optimasi query.

———

### Step 4: Financial Data Integrity Final Check

Melakukan audit terakhir terhadap standar data finansial yang ditetapkan di AGENTS.md.

Checklist Audit:

- [ ] No Floating Point: Periksa kembali internal/model/dashboard.go. Pastikan TotalBalance, TotalIncome, dan TotalExpense menggunakan int64. (Floating point hanya boleh ada di Percentage untuk keperluan display).
- [ ] Zero-Value Handling: Pastikan jika user baru (tidak ada transaksi), API mengembalikan angka 0, bukan null atau error 500. (Gunakan COALESCE(SUM(...), 0) di SQL).
- [ ] Timezone Consistency: Pastikan filter tanggal menggunakan UTC atau timezone yang konsisten agar tidak ada transaksi yang "bergeser" bulan karena perbedaan waktu.

———

### 🏁 Ringkasan Output Day 12

Setelah menyelesaikan Day 12, Anda harus bisa melaporkan:

1. Security: "Verified: User A cannot access User B's dashboard data."
2. Performance: "Verified: Dashboard queries use Index Scan. Latency for 5k records is X ms."
3. Integrity: "Verified: No floating point used for currency calculations."