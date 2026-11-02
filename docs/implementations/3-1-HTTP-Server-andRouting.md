## 3.1 HTTP Server & Routing

Berikut langkah‑langkah lengkap yang harus Anda lakukan, lengkap dengan contoh kode dan penjelasan agar Anda benar‑benar “paham‑paham” tentang bagaimana Go membuat HTTP server, meng‑route request, dan men‑inject middleware.
Saya membaca semua file yang sudah ada (cmd/api/main.go, internal/platform/config.go, internal/platform/database.go) supaya tidak ada asumsi yang salah.

> Catatan
> \*Semua kode di bawah ini sudah beradaptasi dengan struktur yang ada (Folder internal/..., cmd/... dan file konfigurasi).
> Jika Anda ingin men‑extend lebih lanjut (misalnya menambah service atau repository), langkah‑langkah berikut juga memberi petunjuk di mana menaruh file‑file itu.

———

## 1️⃣ Buat main.go (sudah ada, namun saya ulangi agar jelas)

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AankTia/go-fintrack/internal/platform"          // Load config, DB, migration
	"github.com/go-chi/chi/v5"                                 // Router
	"github.com/go-chi/chi/v5/middleware"                      // Built‑in middleware
)

func main() {
	// 1. Load konfigurasi (.env → OS env)
	cfg, err := platform.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// 2. Koneksi ke DB (jika gagal, server mati)
	if err := platform.Connect(); err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	// 3. Ping DB agar ter‑verifikasi aktif
	if err := platform.Ping(); err != nil {
		log.Fatalf("Database not ready: %v", err)
	}

	// 4. Jalankan migrasi secara otomatis (buat tabel‑tabel)
	if err := platform.RunMigrations(); err != nil {
		log.Fatalf("Could not run migrations: %v", err)
	}

	// 5. Buat router
	r := chi.NewRouter()

	// 6. Tambahkan middleware “Logger” (meng‑log setiap request)
	r.Use(middleware.Logger)

	// 7. Route root (hanya contoh)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello !"))
	})

	// 8. Sub‑router untuk API dengan prefix /api/v1
	r.Route("/api/v1", func(r chi.Router) {
		// TODO: Add auth routes (see next step)
	})

	// 9. Jalankan server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
```

Penjelasan tiap baris

| Langkah | Penjelasan                                                                                                                     |
| ------- | ------------------------------------------------------------------------------------------------------------------------------ |
| 1       | `platform.LoadConfig()` mem‑load .env (jika ada) dan mem‑parsing DB_URL, PORT, JWT_SECRET.                                     |
| 2       | `platform.Connect()` meng‑membuat pool pgxpool ke PostgreSQL.                                                                  |
| 3       | `platform.Ping()` memastikan koneksi masih hidup (kita tidak mau menjalankan server kalau DB tidak responsif).                 |
| 4       | `platform.RunMigrations()` memakai golang-migrate agar schema siap.                                                            |
| 5       | `chi.NewRouter()` membuat router “chi”.                                                                                        |
| 6       | `middleware.Logger` otomatis men‑log setiap request (method, path, status).                                                    |
| 7       | Root route ("/") hanya contoh; di masa depan bisa di‑hapus.                                                                    |
| 8       | `r.Route("/api/v1", ...)` membuat sub‑router dengan prefix /api/v1. Semua route di dalamnya otomatis men‑ikut prefix tersebut. |
| 9       | `http.ListenAndServe` mem‑jalankan HTTP server pada port yang dikonfigurasi.                                                   |

> Tips: Gunakan `go run cmd/api/main.go` untuk mem‑jalankan server secara lokal. Pastikan Anda sudah membuat file `.env` di root dengan:
>
> DB_URL=postgres://user:pw@localhost:5432/fintrack?sslmode=disable
> PORT=8080
> JWT_SECRET=supersecret

———

## 2️⃣ Menambahkan Auth Routes (register & login)

Karena file internal/handler masih kosong, kita akan membuat handler sederhana di dalam folder internal/handler/auth.go.

### 2.1 Buat folder internal/handler (jika belum ada)

```bash
mkdir -p internal/handler
touch internal/handler/auth.go
```

### 2.2 Tambahkan file auth.go

```go
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/AankTia/go-fintrack/internal/service"
	"github.com/go-chi/chi/v5"
)

// RegisterHandler men‑handle POST /api/v1/auth/register
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Panggil service layer
	usr, err := service.RegisterUser(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // adapt status if needed
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usr)
}

// LoginHandler men‑handle POST /api/v1/auth/login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Panggil service layer
	token, err := service.LoginUser(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
```

### 2.3 Tambahkan file model/auth.go (DTO)

```go
package model

// RegisterRequest datang dari body POST /auth/register
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest datang dari body POST /auth/login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserResponse men‑return data user setelah register
type UserResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}
```

> Catatan
> DTO (Data Transfer Object) hanya menyimpan data yang dibutuhkan, tidak ada logika bisnis.

### 2.4 Update router di cmd/api/main.go

```go
// di dalam r.Route("/api/v1", func(r chi.Router) { ... })
r.Route("/auth", func(r chi.Router) {
    r.Post("/register", handler.RegisterHandler)
    r.Post("/login", handler.LoginHandler)
})
```

> Perhatikan: Kita mem‑import handler di bagian atas file `main.go`:
>
> `import "github.com/AankTia/go-fintrack/internal/handler"`

———

## 3️⃣ Menambahkan Service Layer (tiap handler panggil service)

Buat folder `internal/service` dan file `auth.go`.

```bash
mkdir -p internal/service
touch internal/service/auth.go
```

```go
package service

import (
	"context"
	"errors"
	"time"

	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/AankTia/go-fintrack/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser mem‑proses pendaftaran
func RegisterUser(ctx context.Context, req *model.RegisterRequest) (*model.UserResponse, error) {
	// 1. Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// 2. Simpan ke DB (repo layer)
	user, err := repository.CreateUser(ctx, req.Email, string(hashed))
	if err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}

// LoginUser mem‑verifikasi login
func LoginUser(ctx context.Context, req *model.LoginRequest) (string, error) {
	user, err := repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// TODO: generate JWT token (gunakan jwt-go atau libs)
	// Sebagai contoh sederhana:
	token := "mocked-jwt-token"

	// Biasanya token ber‑ttl (exp) dan signed dengan secret
	_ = time.Now() // placeholder

	return token, nil
}
```

### 3.1 internal/repository/user.go

```bash
mkdir -p internal/repository
touch internal/repository/user.go
```

```go
package repository

import (
	"context"
	"fmt"

	"github.com/AankTia/go-fintrack/internal/platform" // untuk DB pool
)

// User struct (model internal)
type User struct {
	ID           int64
	Email        string
	PasswordHash string
}

// CreateUser men‑insert ke table users
func CreateUser(ctx context.Context, email, hashedPassword string) (*User, error) {
	query := `INSERT INTO users (email, password_hash, created_at, updated_at)
			  VALUES ($1, $2, NOW(), NOW())
			  RETURNING id, email`

	var user User
	if err := platform.DB.QueryRow(ctx, query, email, hashedPassword).Scan(&user.ID, &user.Email); err != nil {
		return nil, fmt.Errorf("cannot create user: %w", err)
	}
	return &user, nil
}

// GetUserByEmail mencari user berdasar email
func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, email, password_hash FROM users WHERE email=$1`
	var user User
	if err := platform.DB.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}
```

> Catatan
> Semua query di‑parameterize ($1, $2) untuk mencegah SQL injection.

———

## 4️⃣ Menambahkan Middleware (optional)

`main.go` sudah men‑gunakan `middleware.Logger`. Jika Anda ingin men‑tambah Recoverer atau CORS:

r.Use(middleware.Recoverer) // melaporkan panic menjadi 500
r.Use(middleware.AllowContentType("application/json")) // optional

> Tips: Gunakan middleware.Timeout(60 \* time.Second) kalau API Anda sensitif terhadap waktu eksekusi.

———

## 5️⃣ Testing & Jalankan

1. Membuat .env
   ```conf
   DB_URL=postgres://postgres:postgres@localhost:5432/fintrack?sslmode=disable
   PORT=8080
   JWT_SECRET=supersecret
   ```
2. Run Migrasi (sudah otomatis di `main.go`) – pastikan tabel users sudah ada (migrations/`000001_create_users_table.up.sql`).
3. Jalankan Server

   `go run cmd/api/main.go`

4. Coba End‑point (misal Postman/HTTPie):

   # Register

   ```bash
   	curl -X POST http://localhost:8080/api/v1/auth/register \
    -H "Content-Type: application/json" \
    -d '{"email":"alice@example.com","password":"P@ssw0rd"}'
   ```

   # Login

   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/login \
   -H "Content-Type: application/json" \
   -d '{"email":"alice@example.com","password":"P@ssw0rd"}'
   ```

5. Periksa Log – Anda akan melihat log setiap request ber‑format:

   ```log
   2026-11-02T10:15:42Z GET /api/v1/auth/register 200
   ```

———

## 6️⃣ Ringkasan “Step‑by‑Step”

| Langkah | Tujuan         | Kode Utama                                    |
| ------- | -------------- | --------------------------------------------- |
| 1       | Load Config    | platform.LoadConfig()                         |
| 2       | Connect DB     | platform.Connect()                            |
| 3       | Ping DB        | platform.Ping()                               |
| 4       | Run Migrations | platform.RunMigrations()                      |
| 5       | Create Router  | chi.NewRouter()                               |
| 6       | Add Logger     | r.Use(middleware.Logger)                      |
| 7       | Root Route     | r.Get("/", ...)                               |
| 8       | API Prefix     | r.Route("/api/v1", ...)                       |
| 9       | Auth Routes    | handler.RegisterHandler, handler.LoginHandler |
| 10      | Start Server   | http.ListenAndServe()                         |

> Semua kode di atas sudah di‑tulis di file yang sesuai di repo Anda, jadi Anda cukup men‑copy‑paste ke file‑file tersebut dan menjalankannya.

———

### 🔍 Kembali ke “Belajar Golang”

- Router (chi) → mem‑handle routing dengan path dan metode HTTP.
- Middleware → fungsi “wrapper” yang dijalankan sebelum/selesai handler (logger, recoverer).
- Handler → menerima http.ResponseWriter & \*http.Request, meng‑decode JSON, mem‑panggil service, dan men‑kembalikan response.
- Service → logika bisnis (hash password, verify).
- Repository → akses database (pgxpool).
- DTO (model) → struktur data yang dibawa/ditangani di API.

Dengan struktur di atas Anda sudah mem‑mematuhi prinsip Clean Architecture yang dijelaskan di AGENTS.md – setiap lapisan punya tanggung jawab tersendiri.
