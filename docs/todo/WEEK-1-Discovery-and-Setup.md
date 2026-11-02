# 🗓️ WEEK 1: Discovery & Setup (Detailed Breakdown)

Goal Utama: Infrastruktur dasar siap, database terhubung, dan API "Hello World" (Auth) berjalan.

## 🟦 Day 1: Environment Setup & Project Structure

Fokus: Menyiapkan pondasi agar tidak ada konflik path atau dependency di kemudian hari.

### 1.1 Go Runtime & Module Init

- [x] Install/Verify Go version (Latest stable).
- [x] Initialize Go module: go mod init github.com/AankTia/go-fintrack.
- [x] Setup .gitignore (abaikan .env, binary, dan folder IDE).

### 1.2 Project Structure (Clean Architecture)

Buat struktur folder sesuai standar AGENTS.md (Section 5):

- [x] cmd/api/main.go (Entry point aplikasi).
- [x] internal/handler/ (HTTP parsing & response).
- [x] internal/service/ (Business logic).
- [x] internal/repository/ (Database queries).
- [x] internal/model/ (DTOs & DB entities).
- [x] internal/platform/ (DB connection, config loader).
- [x] internal/middleware/ (Auth, logging).
- [x] migrations/ (SQL migration files).

### 1.3 Infrastructure (Docker)

- [x] Buat docker-compose.yml dengan image postgres:alpine.
- [x] Konfigurasi environment variable untuk DB (User, Password, DB Name).
- [x] Verifikasi: Jalankan docker compose up -d dan pastikan DB bisa diakses via tool external (pgAdmin/DBeaver).

———

## 🟦 Day 2: Configuration & Database Connection

Fokus: Menghubungkan aplikasi Go dengan database secara aman dan terstruktur.

### 2.1 Config Loader

- [ ] Install library config (contoh: spf13/viper atau joho/godotenv).
  - [x] `go get github.com/joho/godotenv`; It loads a .env file into os.Getenv. Very lightweight, no extra config needed. (spf13/viper is richer but also more complex; for a simple API we usually pick godotenv.)
  - [x] `go get github.com/go-chi/chi/v5`
  - [x] `go get github.com/go-chi/chi/v5/middleware`

- [x] Buat struct Config untuk menampung DB_URL, PORT, JWT_SECRET.
- [x] Implementasi fungsi LoadConfig() di internal/platform/config.go.
- [x] Implementasi test TestLoadConfig

### 2.2 Database Connection Pool

- [x] Install driver PostgreSQL (contoh: jackc/pgx atau lib/pq).
- [x] Implementasi koneksi singleton/pool di internal/platform/database.go.
- [x] Tambahkan logic Ping saat startup untuk memastikan DB ready sebelum server jalan.

### 2.3 Migration Setup

- [x] Pilih & Install tool migrasi (Rekomendasi: golang-migrate atau goose).
- [x] Buat folder migrations/.
- [x] First Migration: Buat file .sql untuk tabel users (id, email, password_hash, created_at, updated_at).
- [x] Verifikasi: Jalankan migrasi dan cek tabel di PostgreSQL.

———

## 🟦 Day 3: Basic REST API (Auth Basics)

Fokus: Memastikan flow request → handler → service → repo bekerja.

### 3.1 HTTP Server & Routing

- [x] Install router (Rekomendasi: go-chi/chi karena ringan dan std-lib compliant).
- [x] Setup base router di main.go dengan prefix /api/v1.
- [x] Implementasi middleware dasar: Logger dan Recovery.

### 3.2 Auth: Register Endpoint (POST /auth/register)

- [ ] Model: Buat RegisterRequest dan UserResponse DTO.
- [ ] Repository: Implementasi CreateUser (insert ke DB).
- [ ] Service: Implementasi logic hashing password (gunakan bcrypt).
- [ ] Handler: Validasi input → panggil service → return 201 Created.

### 3.3 Auth: Login Endpoint (POST /auth/login)

- [ ] Repository: Implementasi GetUserByEmail.
- [ ] Service: Verifikasi password hash → generate JWT token.
- [ ] Handler: Validasi input → panggil service → return token.

### 3.4 Verification (The "Done" Gate)

- [ ] Buat Postman Collection untuk Register & Login.
- [ ] Test Case:
  - [ ] Register user baru → Berhasil.
  - [ ] Register email yang sama → Error 400 (Duplicate).
  - [ ] Login dengan password salah → Error 401 (Unauthorized).
  - [ ] Login berhasil → Mendapat JWT token.

———

## ✅ Acceptance Criteria (End of Week 1)

- [ ] docker compose up menjalankan DB tanpa error.
- [ ] Migrasi tabel users terpasang otomatis/manual.
- [ ] Endpoint /api/v1/auth/register dan /api/v1/auth/login berfungsi.
- [ ] Password disimpan dalam bentuk Hash (bukan plain text).
- [ ] Kode mengikuti struktur folder yang disepakati di AGENTS.md
