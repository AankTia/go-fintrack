# Day 5: Account Management (CRUD)

## 🏗️ Arsitektur Aliran Data

Sebelum mulai, pahami alur ini:
Client (HTTP) → Handler → Service → Repository → Database

———

## Step 1: Implementasi Repository (Data Access)

File: `internal/repository/account_repository.go` (Lengkapi file ini)

Penjelasan untuk Pembelajar Go:
Repository bertugas melakukan query ke database.

1. Penting: Kita menggunakan sql.DB untuk eksekusi query.
2. Security: Gunakan placeholder $1, $2 untuk mencegah SQL Injection.
3. Ownership: Di setiap query SELECT atau UPDATE, kita WAJIB menambahkan WHERE user_id = $X. Ini memastikan User A tidak bisa memodifikasi data User B hanya dengan mengganti ID di URL.

Detail Implementasi:

- Buat struct accountRepository yang mengimplementasikan AccountRepository interface.
- Implementasikan method:
  - Create: INSERT INTO accounts (user_id, name, type, balance, currency) VALUES ($1, $2, $3, $4, $5) RETURNING id.
  - GetById: SELECT ... FROM accounts WHERE id = $1 AND user_id = $2. (Cek user_id untuk keamanan).
  - GetByUserId: SELECT ... FROM accounts WHERE user_id = $1.
  - UpdateBalance: UPDATE accounts SET balance = balance + $1, updated_at = NOW() WHERE id = $2 AND user_id = $3.

Di Go, implementasi interface bersifat implicit. Artinya, selama struct Anda memiliki semua method yang didefinisikan di AccountRepository interface dengan tanda tangan (parameter dan return value) yang persis sama, maka struct tersebut otomatis dianggap mengimplementasikan interface tersebut.

Berikut adalah contoh implementasinya di `file internal/repository/account_repository.go`:

```go
package repository

import (
	"context"
	"database/sql"
	"github.com/AankTia/go-fintrack/internal/model"
)

// 1. Definisikan struct yang menyimpan dependency (database connection)
// Kita beri nama kecil (unexported) agar hanya bisa diinisialisasi melalui NewAccountRepository
type accountRepository struct {
	db *sql.DB
}

// 2. Buat Constructor (Factory Function)
// Fungsi ini mengembalikan tipe interface AccountRepository, bukan tipe struct accountRepository.
// Ini adalah praktek Dependency Injection agar layer atas tidak terikat pada implementasi konkret.
func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepository{
		db: db,
	}
}

// 3. Implementasi Method Create
func (r *accountRepository) Create(ctx context.Context, account *model.Account) error {
	// Gunakan query parameterized ($1, $2, dst) untuk mencegah SQL Injection
	query := `
		INSERT INTO accounts (user_id, name, type, balance, currency) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id`

	// ExecQuery digunakan karena kita butuh return value (id)
	err := r.db.QueryRowContext(ctx, query, 
		account.UserID, 
		account.Name, 
		account.Type, 
		account.Balance, 
		account.Currency,
	).Scan(&account.ID) // Mengambil ID yang baru dibuat dari DB dan memasukkannya ke struct account

	return err
}

// 4. Implementasi Method GetById (DENGAN SECURITY CHECK)
func (r *accountRepository) GetById(ctx context.Context, id string, userID string) (*model.Account, error) {
	// PERHATIKAN: Kita menambahkan 'user_id = $2'. 
	// User tidak boleh bisa mengambil data hanya dengan ID akun, dia harus pemiliknya.
	query := `SELECT id, user_id, name, type, balance, currency, created_at, updated_at 
	          FROM accounts WHERE id = $1 AND user_id = $2`

	row := r.db.QueryRowContext(ctx, query, id, userID)

	var acc model.Account
	err := row.Scan(&acc.ID, &acc.UserID, &acc.Name, &acc.Type, &acc.Balance, &acc.Currency, &acc.CreatedAt, &acc.UpdatedAt)

	if err != nil {
		return nil, err // Akan return sql.ErrNoRows jika data tidak ditemukan
	}

	return &acc, nil
}

// 5. Implementasi Method GetByUserId
func (r *accountRepository) GetByUserId(ctx context.Context, userID string) ([]*model.Account, error) {
	query := `SELECT id, user_id, name, type, balance, currency, created_at, updated_at 
	          FROM accounts WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // WAJIB tutup rows untuk mencegah memory leak

	var accounts []*model.Account
	for rows.Next() {
		var acc model.Account
		if err := rows.Scan(&acc.ID, &acc.UserID, &acc.Name, &acc.Type, &acc.Balance, &acc.Currency, &acc.CreatedAt, &acc.UpdatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, &acc)
	}

	return accounts, nil
}

// 6. Implementasi Method UpdateBalance
func (r *accountRepository) UpdateBalance(ctx context.Context, id string, amount int64, userID string) error {
	// Update saldo relatif (balance = balance + amount)
	// Lagi-lagi, tambahkan filter user_id untuk keamanan
	query := `UPDATE accounts SET balance = balance + $1, updated_at = NOW() 
	          WHERE id = $2 AND user_id = $3`

	_, err := r.db.ExecContext(ctx, query, amount, id, userID)
	return err
}
```

### 💡 Penjelasan Penting untuk Anda (Pembelajar Go):

1. `context.Context`: Anda melihat ctx di setiap method. Di Go, context digunakan untuk menangani timeout atau pembatalan request. Jika user menutup browser sebelum query selesai, Go bisa membatalkan query database-nya agar tidak membebani server.
2. `Scan()`: Database mengembalikan data dalam bentuk baris/kolom. Scan bertugas memindahkan nilai dari kolom database ke variabel di Go. Urutan variabel di Scan harus sama persis dengan urutan kolom di SELECT.
3. `defer rows.Close()`: defer artinya "jalankan perintah ini di akhir fungsi". Menutup rows sangat krusial di Go; jika lupa, aplikasi Anda akan mengalami connection leak dan akhirnya crash karena kehabisan koneksi ke database.
4. `*sql.DB` vs `AccountRepository`:
   - `accountRepository` (huruf kecil) adalah implementasi internal.
   - `AccountRepository` (huruf besar) adalah interface publik.
   - Dengan mengembalikan interface di `NewAccountRepository`, Anda bisa membuat Mock Repository saat melakukan Unit Testing nanti tanpa perlu koneksi database asli.
5. Security: Perhatikan bahwa saya menambahkan `userID` sebagai parameter di `GetById` dan `UpdateBalance`. Ini adalah implementasi dari Strict Ownership yang diminta di AGENTS.md.

———

## Step 2: Implementasi Service (Business Logic)

File: `internal/service/account.go` (Buat file baru)

Penjelasan untuk Pembelajar Go:
Service adalah "otak" aplikasi. Handler tidak boleh tahu cara database bekerja, dan Repository tidak boleh tahu aturan bisnis.

- Mapping: Service bertugas mengubah AccountRequest (data dari user) menjadi Account (entity untuk DB).
- Validation: Di sinilah Anda mengecek aturan bisnis, misalnya: "Apakah nama akun terlalu pendek?" atau "Apakah saldo awal negatif?".

Detail Implementasi:

- Buat struct AccountService yang memiliki field repo repository.AccountRepository.
- Implementasikan method:
  - CreateAccount: Terima userID dan req model.AccountRequest → Validasi → Panggil repo.Create.
  - GetMyAccounts: Panggil repo.GetByUserId.
  - GetAccountDetail: Panggil repo.GetById.

Kunci utama di Service layer:

1. Tidak boleh ada logic HTTP (seperti http.ResponseWriter atau http.Request).
2. Tidak boleh ada logic SQL (seperti SELECT atau INSERT).
3. Tugas utamanya adalah: Validasi bisnis, konversi data (Request → Entity), dan koordinasi pemanggilan repository.

Berikut adalah detail implementasi untuk file `internal/service/account.go`:

```go
package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/AankTia/go-fintrack/internal/repository"
)

// Custom errors agar Handler bisa menentukan status code HTTP yang tepat
var (
	ErrAccountNotFound = errors.New("account not found")
	ErrInvalidInput    = errors.New("invalid input data")
)

type AccountService struct {
	repo repository.AccountRepository
}

// NewAccountService adalah constructor untuk membuat instance AccountService
func NewAccountService(repo repository.AccountRepository) *AccountService {
	return &AccountService{
		repo: repo,
	}
}

// CreateAccount menangani logika bisnis pembuatan akun baru
func (s *AccountService) CreateAccount(ctx context.Context, userID string, req model.AccountRequest) (*model.Account, error) {
	// 1. Business Validation
	// Contoh: Nama akun tidak boleh kosong
	if req.Name == "" {
		return nil, fmt.Errorf("%w: account name is required", ErrInvalidInput)
	}

	// 2. Mapping: Ubah AccountRequest (DTO) menjadi Account (Entity)
	// Di sini kita menggabungkan data input user dengan userID dari session/token
	account := &model.Account{
		UserID:   userID,
		Name:     req.Name,
		Type:     req.Type,
		Balance:  req.Balance,
		Currency: req.Currency,
	}

	// Set default currency jika kosong
	if account.Currency == "" {
		account.Currency = "IDR"
	}

	// 3. Panggil Repository untuk simpan ke DB
	err := s.repo.Create(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %v", err)
	}

	return account, nil
}

// GetMyAccounts mengambil semua akun milik user tertentu
func (s *AccountService) GetMyAccounts(ctx context.Context, userID string) ([]*model.Account, error) {
	accounts, err := s.repo.GetByUserId(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch accounts: %v", err)
	}

	return accounts, nil
}

// GetAccountDetail mengambil detail satu akun dengan validasi ownership
func (s *AccountService) GetAccountDetail(ctx context.Context, accountID string, userID string) (*model.Account, error) {
	acc, err := s.repo.GetById(ctx, accountID, userID)
	if err != nil {
		// Jika repo mengembalikan sql.ErrNoRows, kita bungkus jadi ErrAccountNotFound
		// agar handler bisa mengembalikan 404 Not Found
		return nil, ErrAccountNotFound
	}

	return acc, nil
}

// UpdateAccountBalance menangani penambahan atau pengurangan saldo
func (s *AccountService) UpdateAccountBalance(ctx context.Context, accountID string, userID string, amount int64) error {
	// Logic Bisnis: Kita bisa cek dulu apakah saldo akhir akan menjadi negatif
	// (Jika aturan bisnis melarang saldo negatif)

	err := s.repo.UpdateBalance(ctx, accountID, amount, userID)
	if err != nil {
		return fmt.Errorf("failed to update balance: %v", err)
	}

	return nil
}
```

### 💡 Penjelasan untuk Pembelajar Go:

1. `errors.New` & `%w (Error Wrapping)`:
   Saya membuat `ErrAccountNotFound` dan `ErrInvalidInput`. Di Go, sangat umum untuk mendefinisikan _"Sentinel Errors"_.
   Saat menggunakan `fmt.Errorf("%w: ...", ErrInvalidInput)`, kita melakukan wrapping. Ini memungkinkan Handler nanti melakukan pengecekan: `errors.Is(err, service.ErrInvalidInput)` untuk menentukan apakah harus mengirim status 400 atau 500.
2. Konsep _DTO (Data Transfer Object)_:
   Perhatikan bahwa `CreateAccount` menerima `model.AccountRequest`, bukan `model.Account`.
   - `AccountRequest`: Apa yang user kirim via JSON (hanya Nama, Tipe, Saldo).
   - `Account`: Apa yang ada di database (termasuk ID, UserID, CreatedAt).
     Memisahkan keduanya mencegah user _"menipu"_ sistem dengan mencoba mengirimkan `ID` atau `UserID` sendiri di dalam JSON.
3. Pointer `*model.Account`:
   Kita mengembalikan pointer (*) agar jika terjadi error, kita bisa mengembalikan `nil`. Ini adalah cara standar di Go untuk menunjukkan _"tidak ada data yang ditemukan"_.
4. Dependency Injection (DI):
   `AccountService` tidak membuat repository-nya sendiri (tidak ada `repo := repository.New...` di dalam struct). Ia menerima repository melalui constructor NewAccountService. Ini membuat kode Anda sangat mudah di-test menggunakan Mock.
———

## Step 3: Implementasi Handler (HTTP Layer)

File: internal/handler/account.go (Buat file baru)

Penjelasan untuk Pembelajar Go:
Handler menangani komunikasi HTTP.

- Parsing: Mengubah JSON request menjadi struct menggunakan json.NewDecoder.
- Context: Anda tidak boleh menerima user_id dari body request karena bisa dipalsukan. Ambil user_id dari r.Context() yang sebelumnya sudah diisi oleh AuthMiddleware setelah verifikasi token JWT.

Detail Implementasi:

- Buat struct AccountHandler yang memiliki field svc *service.AccountService.
- Implementasikan method:
  - Create: Ambil userID dari context → Decode JSON body ke model.AccountRequest → Panggil svc.CreateAccount → Kembalikan 201 Created.
  - List: Ambil userID dari context → Panggil svc.GetMyAccounts → Kembalikan 200 OK.

Tugas utama Handler adalah:

1. Menerima Request HTTP: Mengambil data dari body JSON atau URL parameter.
2. Ekstraksi Identitas: Mengambil userID dari context (hasil proses Auth Middleware).
3. Orkestrasi: Memanggil Service yang tepat.
4. Response: Mengirimkan status code HTTP dan JSON response kembali ke client.

Berikut adalah detail implementasi untuk file internal/handler/account.go:

package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/AankTia/go-fintrack/internal/service"
)

type AccountHandler struct {
	svc *service.AccountService
}

func NewAccountHandler(svc *service.AccountService) *AccountHandler {
	return &AccountHandler{
		svc: svc,
	}
}

// Helper function untuk mengirim response JSON agar kode tidak repetitif
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// Helper function untuk mengirim error JSON
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// Create menangani POST /api/v1/accounts
func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	// 1. Ambil userID dari context
	// userID disuntikkan oleh AuthMiddleware setelah JWT diverifikasi
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized: user id not found in context")
		return
	}

	// 2. Decode JSON body ke AccountRequest DTO
	var req model.AccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	// 3. Panggil Service
	account, err := h.svc.CreateAccount(r.Context(), userID, req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			respondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	// 4. Response sukses (201 Created)
	respondWithJSON(w, http.StatusCreated, account)
}

// List menangani GET /api/v1/accounts
func (h *AccountHandler) List(w http.ResponseWriter, r *http.Request) {
	// 1. Ambil userID dari context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// 2. Panggil Service
	accounts, err := h.svc.GetMyAccounts(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to fetch accounts")
		return
	}

	// 3. Response sukses (200 OK)
	respondWithJSON(w, http.StatusOK, accounts)
}

// GetDetail menangani GET /api/v1/accounts/{id}
func (h *AccountHandler) GetDetail(w http.ResponseWriter, r *http.Request) {
	// Ambil userID dari context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Ambil accountID dari URL parameter (asumsi menggunakan chi router)
	// Jika pakai chi, gunakan chi.URLParam(r, "id")
	// Untuk contoh sederhana, kita asumsikan ID diambil dari context/param
	accountID := r.URL.Query().Get("id") 
	if accountID == "" {
		respondWithError(w, http.StatusBadRequest, "account id is required")
		return
	}

	account, err := h.svc.GetAccountDetail(r.Context(), accountID, userID)
	if err != nil {
		if errors.Is(err, service.ErrAccountNotFound) {
			respondWithError(w, http.StatusNotFound, "account not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, account)
}

### 💡 Penjelasan untuk Pembelajar Go:

1. r.Context().Value("user_id"):
   Ini adalah bagian paling penting dalam keamanan API. Anda tidak menerima user_id dari JSON body karena user bisa saja mengirimkan user_id milik orang lain. Dengan mengambilnya dari context, kita menjamin bahwa userID tersebut berasal dari token JWT yang sudah diverifikasi oleh middleware.
   - . (string) disebut sebagai Type Assertion. Karena Value() mengembalikan tipe interface{}, kita harus memberitahu Go bahwa nilai tersebut adalah sebuah string.
2. json.NewDecoder vs json.Unmarshal:
   - Unmarshal digunakan jika data sudah ada dalam bentuk []byte (slice of bytes).
   - NewDecoder lebih efisien untuk membaca data langsung dari stream (r.Body), sehingga menghemat memori karena tidak perlu membaca seluruh body ke dalam buffer terlebih dahulu.
3. Helper respondWithJSON:
   Dalam project nyata, kita tidak ingin menulis w.Header().Set(...) berkali-kali. Membuat helper function membuat kode lebih bersih dan konsisten.
4. Error Handling dengan errors.Is:
   Handler tidak boleh mengembalikan error database mentah (seperti "sql: no rows in result set") ke user karena itu membocorkan detail internal server (security risk).
   Kita mengecek apakah error tersebut adalah service.ErrAccountNotFound, jika ya, kita kirim pesan yang ramah user seperti "account not found" dengan status 404.
———

## Step 4: Dependency Injection (Wiring)

File: cmd/api/main.go (Update file ini)

Penjelasan untuk Pembelajar Go:
Di Go, kita tidak membuat objek secara acak. Kita menggunakan Dependency Injection (menyuntikkan dependency dari bawah ke atas).

1. Database → dimasukkan ke → Repository
2. Repository → dimasukkan ke → Service
3. Service → dimasukkan ke → Handler

Detail Implementasi:
Di dalam main(), setelah koneksi DB berhasil:

// 1. Repo
accountRepo := repository.NewAccountRepository(db) 
// 2. Service
accountSvc := service.NewAccountService(accountRepo)
// 3. Handler
accountHdl := handler.NewAccountHandler(accountSvc)

// 4. Routing (gunakan middleware auth)
r.Group(func(r chi.Router) {
    r.Use(middleware.Auth) // Pastikan route ini terproteksi
    r.Post("/api/v1/accounts", accountHdl.Create)
    r.Get("/api/v1/accounts", accountHdl.List)
})

Dalam Go, kita menggunakan pola Constructor Injection. Artinya, kita membuat objek dari layer terbawah (Database) → Repository → Service → Handler.

Berikut adalah detail implementasi yang harus Anda terapkan di cmd/api/main.go.

### 🛠️ Detail Implementasi cmd/api/main.go

Saya akan memberikan potongan kode yang harus Anda masukkan ke dalam fungsi main(). Saya mengasumsikan Anda menggunakan chi sebagai router (karena sangat populer di ekosistem Go dan ringan).

package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi" // Asumsi menggunakan chi router
	_ "github.com/lib/pq"    // Driver PostgreSQL

	"github.com/AankTia/go-fintrack/internal/handler"
	"github.com/AankTia/go-fintrack/internal/middleware"
	"github.com/AankTia/go-fintrack/internal/repository"
	"github.com/AankTia/go-fintrack/internal/service"
	"github.com/AankTia/go-fintrack/internal/platform" // Asumsi tempat koneksi DB
)

func main() {
	// 1. INISIALISASI DATABASE
	// Pastikan Anda sudah memiliki fungsi untuk koneksi DB di platform/database.go
	db, err := platform.ConnectDB() 
	if err != nil {
		log.Fatal("Could not connect to DB:", err)
	}
	defer db.Close()

	// ---------------------------------------------------------
	// 2. DEPENDENCY INJECTION (THE WIRING)
	// ---------------------------------------------------------

	// Step A: Buat Repository (Layer paling bawah)
	// Kita masukkan 'db' ke dalam repo
	accountRepo := repository.NewAccountRepository(db)

	// Step B: Buat Service
	// Kita masukkan 'accountRepo' ke dalam service
	accountSvc := service.NewAccountService(accountRepo)

	// Step C: Buat Handler
	// Kita masukkan 'accountSvc' ke dalam handler
	accountHdl := handler.NewAccountHandler(accountSvc)

	// ---------------------------------------------------------
	// 3. ROUTING
	// ---------------------------------------------------------
	r := chi.NewRouter()

	// Middleware Global
	r.Use(middleware.Logger) // Contoh middleware logger

	// Grup API v1
	r.Route("/api/v1", func(r chi.Router) {

		// Route Public (misal: Login/Register)
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)

		// Route Private (Wajib Login)
		// Semua route di dalam grup ini akan melewati middleware Auth
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware) // Middleware yang menyuntikkan user_id ke context

			// Endpoint Account Management
			r.Post("/accounts", accountHdl.Create)       // POST /api/v1/accounts
			r.Get("/accounts", accountHdl.List)         // GET /api/v1/accounts
			r.Get("/accounts/{id}", accountHdl.GetDetail) // GET /api/v1/accounts/123
		})
	})

	log.Println("Server starting on :8080...")
	http.ListenAndServe(":8080", r)
}

———

### 💡 Penjelasan untuk Pembelajar Go:

#### 1. Mengapa harus urut (Repo → Svc → Hdl)?

Ini disebut Dependency Graph.

- Handler tidak bisa bekerja tanpa Service.
- Service tidak bisa bekerja tanpa Repository.
- Repository tidak bisa bekerja tanpa DB.
  Jadi, kita harus membuat "pondasinya" dulu (DB), baru membangun lantai di atasnya satu per satu.

#### 2. Apa keuntungan pola ini?

Bayangkan jika suatu saat Anda ingin mengganti database dari PostgreSQL ke MySQL.

- Tanpa DI: Anda harus mengubah kode di Handler, Service, dan Repository.
- Dengan DI: Anda hanya perlu membuat implementasi AccountRepository yang baru untuk MySQL, lalu mengganti satu baris di main.go:
  accountRepo := repository.NewMySQLAccountRepository(db)
  Selesai. Service dan Handler tidak perlu diubah sama sekali karena mereka hanya tahu interface-nya, bukan detail databasenya.

#### 3. Middleware Grouping

Perhatikan r.Group. Ini adalah cara efisien di Go untuk mengelompokkan route.

- Route /auth/login adalah Public (tidak butuh token).
- Route /accounts adalah Private.
  Dengan menaruh r.Use(middleware.AuthMiddleware) di dalam grup, kita tidak perlu memanggil fungsi cek token di setiap handler. Middleware akan mencegat request, memvalidasi JWT, dan jika valid, ia akan memasukkan userID ke dalam context sebelum request sampai ke accountHdl.Create.

### 🚀 Checklist Akhir untuk Anda:

Setelah Anda apply semua langkah (Step 1-4), coba jalankan aplikasinya:

1. go run cmd/api/main.go
2. Buka Postman.
3. Lakukan alur: Register → Login (Simpan Token) → Create Account (Kirim Token di Header) → List Account.
———

## Step 5: Verifikasi & Testing (Praktek)

Lakukan tes manual untuk memastikan implementasi Anda benar:

1. Positive Test: Login → Create Account → List Account. (Harus berhasil).
2. Negative Test (Security):
   - Login sebagai User A.
   - Coba panggil GET /api/v1/accounts/{id_milik_user_B}.
   - Hasil yang benar: Harus 404 Not Found atau 403 Forbidden. Jika keluar data User B, berarti ada bug serius di Repository/Service Anda.
3. Data Integrity: Cek database, pastikan balance tersimpan sebagai BIGINT dan bukan float.

### 1. Pengujian "Happy Path" (Alur Normal)

Tujuannya adalah memastikan fitur berfungsi sebagaimana mestinya.

Skenario A: Membuat Akun Baru

- Method: POST
- URL: /api/v1/accounts
- Header: Authorization: Bearer <TOKEN_ANDA>
- Body (JSON):

  {
  "name": "Tabungan Utama",
  "type": "bank",
  "balance": 1000000,
  "currency": "IDR"
}
- Hasil yang Diharapkan: Status 201 Created dan response JSON berisi objek akun lengkap dengan id.

Skenario B: Melihat Daftar Akun

- Method: GET
- URL: /api/v1/accounts
- Header: Authorization: Bearer <TOKEN_ANDA>
- Hasil yang Diharapkan: Status 200 OK dan array JSON yang berisi akun yang baru saja Anda buat.

———

### 2. Pengujian "Negative Path" (Validasi Input)

Tujuannya untuk memastikan aplikasi tidak crash saat menerima input sampah.

Skenario C: Input Tidak Valid

- Action: Kirim request POST /api/v1/accounts tapi kosongkan field name atau masukkan type yang ngawur (misal: "type": "crypto").
- Hasil yang Diharapkan: Status 400 Bad Request dengan pesan error yang jelas (misal: "account name is required").

Skenario D: Tanpa Token (Unauthorized)

- Action: Panggil GET /api/v1/accounts tanpa mengirimkan header Authorization.
- Hasil yang Diharapkan: Status 401 Unauthorized.

———

### 3. Pengujian "Security Path" (Strict Ownership) 🚨

Ini adalah bagian terpenting. Anda harus mencoba "menjadi hacker" untuk membobol akun user lain.

Persiapan: Anda butuh dua user di database (User A dan User B).

1. Login sebagai User A, buat akun, dan catat ID akun tersebut (misal: ID = 10).
2. Login sebagai User B, simpan tokennya.

Skenario E: Mencoba Mengambil Data User Lain

- Method: GET
- URL: /api/v1/accounts/10 (ID milik User A)
- Header: Authorization: Bearer <TOKEN_USER_B>
- Hasil yang Benar: Status 404 Not Found atau 403 Forbidden.
- Hasil yang SALAH: Status 200 OK dan data akun User A terlihat. (Jika ini terjadi, ada bug di Repository Anda!)

Skenario F: Mencoba Mengupdate Saldo User Lain

- Method: PATCH/PUT (jika sudah buat update)
- URL: /api/v1/accounts/10 (ID milik User A)
- Header: Authorization: Bearer <TOKEN_USER_B>
- Hasil yang Benar: Status 404 Not Found atau 403 Forbidden.

———

### 4. Verifikasi Database (Data Integrity)

Buka database tool Anda (misal: pgAdmin, DBeaver, atau psql) dan jalankan query:

SELECT * FROM accounts;

Cek hal berikut:

1. Tipe Data: Pastikan kolom balance tidak memiliki angka desimal (karena kita pakai BIGINT).
2. Ownership: Pastikan setiap baris memiliki user_id yang benar.
3. Timestamp: Pastikan created_at dan updated_at terisi otomatis.

———

### 💡 Tips untuk Pembelajar Go:

Jika Anda menemukan error saat testing, jangan panik. Gunakan teknik Debug tracing:

1. Cek Handler: Tambahkan fmt.Println("Masuk ke handler Create") untuk memastikan request sampai ke handler.
2. Cek Service: Tambahkan fmt.Println("UserID yang diterima:", userID) untuk memastikan token terbaca benar.
3. Cek Repository: Tambahkan fmt.Println("Executing query for UserID:", userID) sebelum db.Exec.

## Ringkasan File yang Akan Anda Sentuh:

| Layer | File | Action |
|:---|:---|:---|
| Repository | internal/repository/account_repository.go | Lengkapi Implementasi SQL |
| Service | internal/service/account.go | Buat Logic Bisnis |
| Handler | internal/handler/account.go | Buat Endpoint HTTP |
| Main | cmd/api/main.go | Hubungkan semua (Wiring) |