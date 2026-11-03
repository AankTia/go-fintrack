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

———

## Step 5: Verifikasi & Testing (Praktek)

Lakukan tes manual untuk memastikan implementasi Anda benar:

1. Positive Test: Login → Create Account → List Account. (Harus berhasil).
2. Negative Test (Security):
   - Login sebagai User A.
   - Coba panggil GET /api/v1/accounts/{id_milik_user_B}.
   - Hasil yang benar: Harus 404 Not Found atau 403 Forbidden. Jika keluar data User B, berarti ada bug serius di Repository/Service Anda.
3. Data Integrity: Cek database, pastikan balance tersimpan sebagai BIGINT dan bukan float.

## Ringkasan File yang Akan Anda Sentuh:

| Layer | File | Action |
|:---|:---|:---|
| Repository | internal/repository/account_repository.go | Lengkapi Implementasi SQL |
| Service | internal/service/account.go | Buat Logic Bisnis |
| Handler | internal/handler/account.go | Buat Endpoint HTTP |
| Main | cmd/api/main.go | Hubungkan semua (Wiring) |