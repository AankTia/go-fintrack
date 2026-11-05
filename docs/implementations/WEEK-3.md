Berikut adalah panduan step-by-step implementasi Week 3 (Transaction Core) yang bisa Anda apply secara manual. Saya membaginya per hari sesuai dokumen rencana kerja.

———

### 🟦 DAY 7: Transaction Create & Validation

Goal: Mencatat transaksi dan update saldo akun secara atomik (semua berhasil atau semua gagal).

#### Step 1: Database Migration

Buat file migration baru (misal: `000004_create_transactions_table.up.sql`) untuk membuat tabel transaksi.

- Kolom: `id (UUID), user_id (FK), account_id (FK), category_id (FK), amount (BIGINT), type (enum: income, expense), description (text), transaction_date (date), created_at, updated_at`.
- Index: Tambahkan index pada user_id dan transaction_date.

#### Step 2: Internal Model (`internal/model/transaction.go`)

Buat file baru untuk mendefinisikan entity dan DTO:

- type Transaction (Entity untuk DB)
- type CreateTransactionRequest (DTO untuk input API: account_id, category_id, amount, type, description, date).

#### Step 3: Transaction Repository (`internal/repository/transaction_repository.go`)

Buat repository untuk menangani query database:

- Create(ctx, transaction): Simpan transaksi baru.
- GetById(ctx, id, userID): Ambil satu transaksi dengan filter ownership.
- GetByUserId(ctx, userID, filters): Ambil list transaksi dengan filter tanggal/tipe/akun.

#### Step 4: Transaction Service (`internal/service/transaction.go`)

Bagian paling kritis. Anda harus mengimplementasikan database transaction di sini.

- Metode CreateTransaction:
  1. Gunakan db.Begin(ctx) untuk memulai transaksi.
  2. Panggil transactionRepo.Create.
  3. Tentukan nilai saldo: jika type == "income", amount tetap positif; jika expense, amount dikali -1.
  4. Panggil accountRepo.UpdateBalance menggunakan koneksi transaksi yang sama.
  5. Jika semua sukses → tx.Commit(). Jika ada error → tx.Rollback().

#### Step 5: Transaction Handler (`internal/handler/transaction.go`)

- Buat endpoint POST /api/v1/transactions.
- Validasi: Pastikan amount > 0.
- Ownership Check: Panggil accountService.GetAccountDetail untuk memastikan account_id benar-benar milik user tersebut sebelum memproses transaksi.

———

### 🟦 DAY 8: Transaction List & Filtering

Goal: Mengambil data transaksi dengan efisien.

#### Step 1: Update Repository Query

Modifikasi GetByUserId di TransactionRepository untuk mendukung:

- JOIN ke tabel accounts dan categories agar API bisa mengembalikan nama akun/kategori, bukan hanya ID.
- WHERE clause dinamis untuk filter start_date, end_date, dan type.

#### Step 2: Pagination Logic

- Tambahkan parameter limit dan offset pada query SQL.
- Urutkan berdasarkan transaction_date DESC.

#### Step 3: List Handler

- Buat endpoint GET /api/v1/transactions.
- Parsing query params dari URL.
- Return response dengan format: { "data": [...], "pagination": { "total": 100, "page": 1, "limit": 10 } }.

———

### 🟦 DAY 9: Transaction Update & Delete

Goal: Mengubah data tanpa merusak saldo akun.

#### Step 1: Update Logic (Idempotent)

Di TransactionService, implementasikan metode UpdateTransaction:

1. Begin Transaction.
2. Reverse Old: Ambil data lama → Update saldo akun dengan nilai kebalikan dari transaksi lama (mengembalikan saldo ke posisi sebelum transaksi tersebut ada).
3. Update Record: Update data transaksi dengan nilai baru.
4. Apply New: Update saldo akun dengan nilai transaksi yang baru.
5. Commit.

#### Step 2: Delete Logic

Di TransactionService, implementasikan metode DeleteTransaction:

1. Begin Transaction.
2. Reverse: Ambil data transaksi → Update saldo akun (kembalikan saldo).
3. Delete: Hapus record transaksi.
4. Commit.

#### Step 3: Update & Delete Handlers

- PATCH /api/v1/transactions/{id}
- DELETE /api/v1/transactions/{id}
- Wajib: Lakukan ownership check pada ID transaksi sebelum melakukan operasi.

———

### 🚀 Checklist Verifikasi untuk Anda:

Setelah apply manual, pastikan Anda mengetes hal berikut:

- [ ] Skenario Income: Tambah transaksi income → Saldo akun harus bertambah.
- [ ] Skenario Expense: Tambah transaksi expense → Saldo akun harus berkurang.
- [ ] Skenario Rollback: Simulasi error saat update saldo → Pastikan transaksi tidak tersimpan di tabel transactions.
- [ ] Skenario Security: Coba akses/ubah transaksi milik user lain menggunakan token user yang berbeda → Harus return 403 Forbidden atau 404 Not Found.
- [ ] Skenario Update: Ubah transaksi Expense 10rb menjadi Income 20rb → Total saldo akun harus berubah sebesar +30rb.
  Worked for 9s • 11:05

kasih detail code yang harus diimplementasikan bagaimana

Karena Anda akan mengaplikasikannya secara manual, saya akan memberikan detail kode yang spesifik dan terstruktur. Kita akan mulai dengan Day 7 (Create & Validation) karena ini adalah fondasi utamanya.

### 1. Database Migration

Buat file `migrations/000004_create_transactions_table.up.sql`:

```sql
CREATE TYPE transaction_type AS ENUM ('income', 'expense');

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    amount BIGINT NOT NULL,
    type transaction_type NOT NULL,
    description TEXT,
    transaction_date DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_transactions_user_date ON transactions(user_id, transaction_date);
CREATE INDEX idx_transactions_account ON transactions(account_id);
```

———

### 2. Internal Model (internal/model/transaction.go)

Buat file baru untuk entity dan DTO.

package model

import "time"

type Transaction struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	AccountID       string    `json:"account_id"`
	CategoryID      string    `json:"category_id"`
	Amount          int64     `json:"amount"`
	Type            string    `json:"type"` // income, expense
	Description     string    `json:"description"`
	TransactionDate time.Time `json:"transaction_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateTransactionRequest struct {
	AccountID   string `json:"account_id" validate:"required"`
	CategoryID  string `json:"category_id" validate:"required"`
	Amount      int64  `json:"amount" validate:"required,gt=0"`
	Type        string `json:"type" validate:"required,oneof=income expense"`
	Description string `json:"description"`
	Date        string `json:"date"` // Format YYYY-MM-DD
}

———

### 3. Transaction Repository (internal/repository/transaction_repository.go)

Implementasikan interface untuk akses DB. Perhatikan penggunaan pgx.Tx agar bisa mendukung transaksi atomik dari layer service.

package repository

import (
	"context"
	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx pgx.Tx, t *model.Transaction) error
	GetById(ctx context.Context, id string, userID string) (*model.Transaction, error)
}

type transactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx pgx.Tx, t *model.Transaction) error {
	query := `INSERT INTO transactions (user_id, account_id, category_id, amount, type, description, transaction_date)
	          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	// Gunakan 'tx' jika ada, jika tidak gunakan 'r.db' (untuk fleksibilitas)
	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, query, t.UserID, t.AccountID, t.CategoryID, t.Amount, t.Type, t.Description, t.TransactionDate)
	} else {
		row = r.db.QueryRow(ctx, query, t.UserID, t.AccountID, t.CategoryID, t.Amount, t.Type, t.Description, t.TransactionDate)
	}

	return row.Scan(&t.ID)
}

func (r *transactionRepository) GetById(ctx context.Context, id string, userID string) (*model.Transaction, error) {
	query := `SELECT id, user_id, account_id, category_id, amount, type, description, transaction_date 
	          FROM transactions WHERE id = $1 AND user_id = $2`

	var t model.Transaction
	err := r.db.QueryRow(ctx, query, id, userID).Scan(
		&t.ID, &t.UserID, &t.AccountID, &t.CategoryID, &t.Amount, &t.Type, &t.Description, &t.TransactionDate,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

———

### 4. Transaction Service (internal/service/transaction.go)

Ini adalah bagian terpenting. Service harus mengelola pgx.Tx untuk menjamin atomicity.

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/AankTia/go-fintrack/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionService struct {
	db             *pgxpool.Pool
	txRepo         repository.TransactionRepository
	accountRepo    repository.AccountRepository
	categoryRepo   repository.CategoryRepository // Asumsi ada
}

func NewTransactionService(db *pgxpool.Pool, txRepo repository.TransactionRepository, accRepo repository.AccountRepository, catRepo repository.CategoryRepository) *TransactionService {
	return &TransactionService{
		db:           db,
		txRepo:       txRepo,
		accountRepo:  accRepo,
		categoryRepo:  catRepo,
	}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, userID string, req model.CreateTransactionRequest) (*model.Transaction, error) {
	// 1. Parse Date
	tDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		tDate = time.Now()
	}

	// 2. Ownership Check (Pastikan akun milik user)
	_, err = s.accountRepo.GetById(ctx, req.AccountID, userID)
	if err != nil {
		return nil, fmt.Errorf("account not found or unauthorized")
	}

	// 3. Start DB Transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx) // Rollback otomatis jika terjadi panic atau error sebelum Commit

	// 4. Create Transaction Record
	transaction := &model.Transaction{
		UserID:          userID,
		AccountID:       req.AccountID,
		CategoryID:      req.CategoryID,
		Amount:          req.Amount,
		Type:            req.Type,
		Description:     req.Description,
		TransactionDate: tDate,
	}

	if err := s.txRepo.Create(ctx, tx, transaction); err != nil {
		return nil, fmt.Errorf("failed to save transaction: %v", err)
	}

	// 5. Update Account Balance
	balanceChange := req.Amount
	if req.Type == "expense" {
		balanceChange = -req.Amount
	}

	// Note: UpdateBalance di repository harus bisa menerima tx atau menggunakan logic update relatif
	// Karena current AccountRepository.UpdateBalance menggunakan r.db, kita perlu modifikasi repo 
	// atau membuat method baru yang menerima pgx.Tx.
	// Untuk sementara, kita asumsikan ada method UpdateBalanceTx di repository.
	if err := s.updateBalanceWithTx(ctx, tx, req.AccountID, balanceChange, userID); err != nil {
		return nil, fmt.Errorf("failed to update balance: %v", err)
	}

	// 6. Commit
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return transaction, nil
}

// Helper untuk update saldo di dalam transaksi
func (s *TransactionService) updateBalanceWithTx(ctx context.Context, tx pgx.Tx, accountID string, amount int64, userID string) error {
	query := `UPDATE accounts SET balance = balance + $1, updated_at = NOW() WHERE id = $2 AND user_id = $3`
	_, err := tx.Exec(ctx, query, amount, accountID, userID)
	return err
}

———

### 5. Transaction Handler (internal/handler/transaction.go)

package handler

import (
	"net/http"
	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/AankTia/go-fintrack/internal/service"
	"github.com/gin-gonic/gin" // Asumsi menggunakan Gin, sesuaikan jika pakai chi/echo
)

type TransactionHandler struct {
	svc *service.TransactionService
}

func NewTransactionHandler(svc *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: svc}
}

func (h *TransactionHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id") // Didapat dari auth middleware
	var req model.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.svc.CreateTransaction(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tx)
}

### 💡 Catatan Penting untuk Implementasi Manual:

1. Modifikasi AccountRepository: Anda mungkin perlu menambahkan metode yang menerima pgx.Tx di AccountRepository agar UpdateBalance bisa berjalan dalam satu transaksi yang sama dengan CreateTransaction.
2. Dependency Injection: Jangan lupa daftarkan TransactionRepository, TransactionService, dan TransactionHandler di main.go agar mereka terhubung.
3. Order of Application:
   Migration → Model → Repository → Service → Handler.