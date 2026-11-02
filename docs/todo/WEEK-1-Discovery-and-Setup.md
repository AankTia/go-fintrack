# 🗓️ WEEK 1: Discovery & Setup (Detailed Breakdown)

Goal Utama: Infrastruktur dasar siap, database terhubung, dan API "Hello World" (Auth) berjalan.

## 🟦 Day 1: Environment Setup & Project Structure

Fokus: Menyiapkan pondasi agar tidak ada konflik path atau dependency di kemudian hari.

### 1.1 Go Runtime & Module Init

- [x] Install/Verify Go version (Latest stable).
- [x] Initialize Go module: go mod init github.com/AankTia/go-fintrack.
- [ ] Setup .gitignore (abaikan .env, binary, dan folder IDE).

### 1.2 Project Structure (Clean Architecture)

Buat struktur folder sesuai standar AGENTS.md (Section 5):

- [ ] cmd/api/main.go (Entry point aplikasi).
- [ ] internal/handler/ (HTTP parsing & response).
- [ ] internal/service/ (Business logic).
- [ ] internal/repository/ (Database queries).
- [ ] internal/model/ (DTOs & DB entities).
- [ ] internal/platform/ (DB connection, config loader).
- [ ] internal/middleware/ (Auth, logging).
- [ ] migrations/ (SQL migration files).

### 1.3 Infrastructure (Docker)

- [ ] Buat docker-compose.yml dengan image postgres:alpine.
- [ ] Konfigurasi environment variable untuk DB (User, Password, DB Name).
- [ ] Verifikasi: Jalankan docker compose up -d dan pastikan DB bisa diakses via tool external (pgAdmin/DBeaver).

———

## 🟦 Day 2: Configuration & Database Connection

Fokus: Menghubungkan aplikasi Go dengan database secara aman dan terstruktur.

### 2.1 Config Loader

- [ ] Install library config (contoh: spf13/viper atau joho/godotenv).
- [ ] Buat struct Config untuk menampung DB_URL, PORT, JWT_SECRET.
- [ ] Implementasi fungsi LoadConfig() di internal/platform/config.go.

### 2.2 Database Connection Pool

- [ ] Install driver PostgreSQL (contoh: jackc/pgx atau lib/pq).
- [ ] Implementasi koneksi singleton/pool di internal/platform/database.go.
- [ ] Tambahkan logic Ping saat startup untuk memastikan DB ready sebelum server jalan.

### 2.3 Migration Setup

- [ ] Pilih & Install tool migrasi (Rekomendasi: golang-migrate atau goose).
- [ ] Buat folder migrations/.
- [ ] First Migration: Buat file .sql untuk tabel users (id, email, password_hash, created_at, updated_at).
- [ ] Verifikasi: Jalankan migrasi dan cek tabel di PostgreSQL.

———

## 🟦 Day 3: Basic REST API (Auth Basics)

Fokus: Memastikan flow request → handler → service → repo bekerja.

### 3.1 HTTP Server & Routing

- [ ] Install router (Rekomendasi: go-chi/chi karena ringan dan std-lib compliant).
- [ ] Setup base router di main.go dengan prefix /api/v1.
- [ ] Implementasi middleware dasar: Logger dan Recovery.

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
