# 3.2 Auth: Register Endpoint

## 📊 Analisis Kondisi Saat Ini

### Yang Sudah Ada (✅):

1. ✅ Model RegisterRequest dan UserResponse sudah di internal/model/auth.go
   2 ✅ Handler RegisterHandler sudah ada, memanggil service layer
2. ✅ Repository pattern sudah ada, tapi belum ada fungsi CreateUser yang terintegrasi
3. ✅ Service RegisterUser sudah basic-nya (bcrypt + call repo)

### Yang Belum Ada (❌):

1. ❌ Repository CreateUser belum diimplementasi - ini yang bikin error!
2. ❌ Password hashing butuh library golang.org/x/crypto/bcrypt (ada tapi belum import)
3. ❌ Service perlu import bcrypt untuk hashing
4. ❌ Repository belum punya koneksi ke DB pool

———

## 🎯 Step-by-Step Implementasi: 3.2 Auth: Register Endpoint

### STEP 1: Update Repository Layer

File: `internal/repository/user.go`
Tujuan: Implementasi fungsi CreateUser yang query ke database

```go
// File: internal/repository/user.go
package repository

import (
    "context"
    "fmt"

    "github.com/AankTia/go-fintrack/internal/platform" // DB pool global
)

type User struct {
    ID          int64
    Email        string
    PasswordHash string
}

// CreateUser men-insert user baru ke database dan return user yang tersimpan
// Parameter: email, hashedPassword (password sudah di-hash oleh service layer)
// Return: object User dengan ID baru, atau error jika gagal insert
func CreateUser(ctx context.Context, email, hashedPassword string) (*User, error) {
    // 1. Query SQL dengan parameterized query (cegah SQL Injection)
    //    - Menggunakan $1, $2 untuk placeholder parameter
    //    - RETURNING clause ambil id & email dari DB sebagai primary response
    query := `INSERT INTO users (email, password_hash, created_at, updated_at)
              VALUES ($1, $2, NOW(), NOW())
              RETURNING id, email`

    // 2. Siapkan variabel result
    var user User

    // 3. Execute query ke DB pool global (platform.DB)
    //    - QueryRow return single row atau error
    if err := platform.DB.QueryRow(ctx, query, email, hashedPassword).Scan(
        &user.ID, &user.Email,
    ); err != nil {
        // Wrap error agar mudah debugging: "cannot create user: <details>"
        return nil, fmt.Errorf("cannot create user: %w", err)
    }

    // 4. Return success: user baru dengan ID yang baru saja di-assign DB (BIGSERIAL)
    return &user, nil
}

// GetUserByEmail - SUDAH ADA (tidak perlu ubah)
func GetUserByEmail(ctx context.Context, email string) (*User, error) {
    // ... query SELECT ...
    return &user, nil
}
```

Penjelasan:

- Layer Responsibility: Repository HANYA akses data, tidak ada business logic
- Parameterized Query: $1, $2 mencegah SQL Injection (penting untuk security!)
- RETURNING Clause: Kita minta DB return id dan email setelah insert selesai
- Error Wrapping: %w preserve error chain untuk debugging

———

### STEP 2: Update Service Layer

File: `internal/service/auth.go`
Tujuan: Hashing password dengan bcrypt sebelum dikirim ke repository

```go
// File: internal/service/auth.go
package service

import (
    "context"
    "errors"
    "time"

    "github.com/AankTia/go-fintrack/internal/model"
    "github.com/AankTia/go-fintrack/internal/repository"
    "golang.org/x/crypto/bcrypt" // ⚠️ IMPORT bcrypt library!
)

// RegisterUser memproses pendaftaran user baru
// Flow: Hash Password → Insert ke DB → Return User ID
func RegisterUser(ctx context.Context, req *model.RegisterRequest) (*model.UserResponse, error) {
    // ====================================================
    // STEP 2.1: VALIDASI INPUT (Boundary Validation)
    // ====================================================

    // Email validation - pastikan email valid & tidak kosong
    if len(req.Email) == 0 {
        return nil, errors.New("email is required")
    }

    // Password validation - minimal panjang 8 karakter
    if len(req.Password) < 8 {
        return nil, errors.New("password must be at least 8 characters")
    }

    // ====================================================
    // STEP 2.2: HASH PASSWORD (Security Layer)
    // ====================================================

    // Bcrypt hashing untuk password
    // bcrypt.DefaultCost = 14 (level keamanan default)
    // Generate hash yang kuat dari plain text password
    hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        // Hashing gagal (sangat jarang, tapi handle saja)
        return nil, errors.New("failed to hash password")
    }

    // Convert []byte ke string untuk dikirim ke repository layer
    hashedStr := string(hashed)

    // ====================================================
    // STEP 2.3: CALL REPOSITORY (Data Access Layer)
    // ====================================================

    // Pasang hashed password (bukan plain text!) ke repository
    // Repository akan insert ke DB dengan parameterized query
    user, err := repository.CreateUser(ctx, req.Email, hashedStr)
    if err != nil {
        // Check apakah error dari database karena email duplicate?
        // PostgreSQL: unique_violation = "23505"
        if errors.Is(err, psqlErrorDuplicateEmail) {
            return nil, errors.New("email already exists") // Return 400 Duplicate
        }
        return nil, err
    }

    // ====================================================
    // STEP 2.4: RETURN SUCCESS
    // ====================================================
    // Kembalikan UserResponse (HANYA ID & EMAIL)
    // Password hash TIDAK dikirim ke client (security!)
    return &model.UserResponse{
        ID:    user.ID,
        Email: user.Email,
    }, nil
}

// LoginUser memverifikasi login user
func LoginUser(ctx context.Context, req *model.LoginRequest) (string, error) {
    // 1. Cari user dari database berdasarkan email
    user, err := repository.GetUserByEmail(ctx, req.Email)
    if err != nil {
        return "", errors.New("invalid credentials")
    }

    // 2. Compare password hash dengan input password (bcrypt.CompareHashAndPassword)
    if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        return "", errors.New("invalid credentials")
    }

    // 3. Generate JWT token (TODO - nanti implementasi proper JWT)
    //    Untuk sekarang gunakan mock token
    token := "mocked-jwt-token-placeholder"

    return token, nil
}
```

Penjelasan:

- Service Layer: Business logic & orchestration (hashing, validation, coordination)
- Bcrypt: Hash password sekali untuk security
- No Plain Text Password: Password hash tidak pernah terekspos ke frontend
- Error Handling: Check duplicate email error dari database

———

### STEP 3: Setup JWT Library

File: `go.mod` (update dependency)

# Command terminal di workspace ini:
```bash
go get golang.org/x/oauth2
go get github.com/golang-jwt/jwt/v4
```

Penjelasan:

- `github.com/golang-jwt/jwt/v4` - library JWT standar untuk Go
- Version v4 karena stable & banyak dipakai
- Akan dipakai di LoginUser untuk generate proper token dengan expiry

———

### STEP 4: Update Main Entry Point (Optional)

File: `cmd/api/main.go`
Tujuan: Pastikan database connection sudah tersedia sebelum run

```go
// File: cmd/api/main.go - bagian platform.Connect() & platform.RunMigrations()

// Pastikan ini ada di main.go:
if err := platform.Connect(); err != nil {
    log.Fatalf("Could not connect to DB: %v", err)
}

// Verify ping
platform.Ping()

// Run migrations before starting server
platform.RunMigrations() // ⚠️ PENTING! Migration harus jalan sebelum server start
```

Penjelasan:

- `RunMigrations()` sudah ada di `internal/platform/database.go`
- Migration tabel users sudah siap di `folder migrations/`
- Pastikan jalankan migration SEBELUM server dimulai

———

## 🔍 VERIFICATION CHECKLIST (Setelah Implementasi)

### 1. Build & Verify
```bash
cd /Users/tiawidi/Developments/Github/go-fintrack
go build -o bin/api ./cmd/api  # Build aplikasi
./bin/api                       # Run aplikasi
```

### 2. Verify Migration Sudah Jalan

Cek di PostgreSQL:
```sql
\d users
```
-- Harus tampil: id (bigint), email (text unique), password_hash (text)

### 3. Test Register Endpoint dengan Postman

```bash
POST http://localhost:8080/api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123"  // Minimal 8 karakter
}
```

Expected Response:
```bash
HTTP Status: 201 Created
{
  "id": 1,
  "email": "user@example.com"
}
```

### 4. Test Duplicate Email
```bash
POST http://localhost:8080/api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",  // Sama seperti sebelumnya
  "password": "SecurePass123"
}
```

Expected Response:
HTTP Status: 400 Bad Request
{
  "error": "email already exists"
}

### 5. Security Verification

- ✅ Password tidak pernah terekspos di response (hanya ID & Email)
- ✅ Database menggunakan parameterized query ($1, $2)
- ✅ Error message tidak membocorkan stack trace/DB details

———

## 📝 Ringkasan Flow Register Endpoint

┌─────────────────┐
│  Client (Frontend) │
│ POST /auth/register │
└────────┬──────────┘
         │ HTTP Request (email, password plain)
         ▼
┌────────────────────────┐
│   Handler Layer        │
│ json.NewDecoder()      │
│ validate input basic   │
│ call service           │
└────────┬───────────────┘
         │ Service Call: RegisterUser()
         ▼
┌────────────────────────┐
│   Service Layer        │
│ 1. Validate input     │
│ 2. Hash password (bcrypt) │
│ 3. Call Repository     │
└────────┬───────────────┘
         │ Repository Call: CreateUser()
         ▼
┌────────────────────────┐
│   Repository Layer     │
│ INSERT INTO users ... │
│ RETURNING id, email    │
└────────┬───────────────┘
         │ Query Result (ID, Email)
         ▼
┌────────────────────────┐
│   Service Returns      │
│ UserResponse {id: 1}  │
└────────┬───────────────┘
         │ HTTP Response (201 Created)
         ▼
┌────────────────────────┐
│   Client Receives      │
│ {"id": 1, "email": "..."} │
└────────────────────────┘

———

## ⚠️ Keamanan Checklist

| Aspek | Status | Implementasi |
|---|---|---|
| Password Hashing | ✅ Done | Bcrypt di Service layer |
| No Plain Text in Response | ✅ Done | UserResponse hanya ID & Email |
| SQL Injection Prevention | ✅ Done | Parameterized query ($1, $2) |
| Error Message Security | ⚠️ Review | Jangan expose stack trace ke client |
| Ownership/Isolation | ⏳ Next Step | Belum ada, nanti setelah Auth selesai |

———

## 🎯 Kesimpulan & Praktek Belajar

### Yang Perlu Kamu Lakukan:

1. Edit internal/repository/user.go → Tambah fungsi CreateUser (Step 1)
2. Edit internal/service/auth.go → Hash password dengan bcrypt (Step 2)
3. Test di Postman - Verify register & duplicate email handling (Step 4)

### Konsep yang Kamu Pelajari dari Ini:

1. Clean Architecture Separation: Handler → Service → Repository (tidak ada mixing)
2. Layer Responsibility: Setiap layer cuma punya satu job, tidak lakukan hal lain
3. Security by Default: Password hash, parameterized query, no sensitive data leaks
4. Error Handling: Wrap errors untuk debugging, return user-friendly message ke client
5. Database Transactions: Insert dengan RETURNING clause (kembali nilai yang baru disisipkan)

———
