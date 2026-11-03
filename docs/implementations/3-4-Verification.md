# 3.4 Verification (The "Done" Gate)

### Bagian 1: Persiapan Lingkungan (Setup)

Sebelum menjalankan curl, aplikasi harus menyala dan database harus siap.

Step 1: Jalankan Database
Pastikan container PostgreSQL sudah berjalan.

```bash
docker compose up -d
```

Penjelasan: Kita menggunakan Docker agar database terisolasi dan tidak mengotori OS utama.

Step 2: Siapkan File `.env`
Pastikan Anda memiliki file `.env` di root folder dengan isi yang sesuai dengan `internal/platform/config.go`:

```conf
DB_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
PORT=8080
JWT_SECRET=bebas-isi-apa-saja-yang-panjang-dan-rahasia
```

Penjelasan: Fungsi `platform.LoadConfig()` mencari variabel ini untuk menentukan ke mana aplikasi harus connect dan apa kunci rahasia untuk menandatangani JWT.

Step 3: Jalankan Aplikasi Go

```bash
go run cmd/api/main.go
```

Penjelasan: Perhatikan log di terminal. Anda harus melihat pesan "Database connection verified!" dan "✅ Database migrations applied". Jika ada error di sini, berarti koneksi DB atau migrasi gagal.

———

### Bagian 2: Skenario Pengujian menggunakan curl

Kita akan menguji 4 skenario utama. Gunakan terminal baru untuk menjalankan command berikut.

#### Skenario 1: Register User Baru (Success)

Command:

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
-H "Content-Type: application/json" \
-d '{"email": "belajar@go.com", "password": "password123"}'
```

- Apa yang terjadi di Go?
  1. `RegisterHandler` menerima JSON → Decode ke RegisterRequest.
  2. `RegisterUser` (Service) memvalidasi panjang password → Hash password dengan bcrypt.
  3. `CreateUser` (Repository) menjalankan INSERT INTO users → DB mengembalikan ID.
- Hasil yang diharapkan: Status 201 Created dan response JSON berisi id dan email.

#### Skenario 2: Register Email Duplikat (Error)

Command: (Jalankan command yang sama dengan Skenario 1 sekali lagi)

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
-H "Content-Type: application/json" \
-d '{"email": "belajar@go.com", "password": "password123"}'
```

- Apa yang terjadi di Go?
  1. DB akan menolak insert karena ada UNIQUE constraint pada kolom email.
  2. repository.CreateUser mengembalikan error.
- Hasil yang diharapkan: Status 400 Bad Request.

#### Skenario 3: Login dengan Password Salah (Unauthorized)

Command:

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
-H "Content-Type: application/json" \
-d '{"email": "belajar@go.com", "password": "salahpassword"}'
```

- Apa yang terjadi di Go?
  1. `GetUserByEmail` mencari user.
  2. `bcrypt.CompareHashAndPassword` membandingkan password input dengan hash di DB → Hasilnya tidak cocok.
  3. Service mengembalikan error `"invalid credentials"`.
- Hasil yang diharapkan: Status `401 Unauthorized`.

#### Skenario 4: Login Berhasil (Success)

Command:

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
-H "Content-Type: application/json" \
-d '{"email": "belajar@go.com", "password": "password123"}'
```

- Apa yang terjadi di Go?
  1. User ditemukan → Password cocok.
  2. generateJWT membuat token menggunakan JWT_SECRET dari config.
- Hasil yang diharapkan: Status 200 OK dan response JSON berisi {"token": "eyJhbG..."}.

———

### Tips Belajar Go untuk Implementasi Manual

Jika Anda ingin mempraktekkan perbaikan kode (seperti menangani error duplikat email), berikut adalah pola implementasinya:

1. Cara menangani Error Spesifik di Go:
Saat ini RegisterUser di internal/service/auth.go hanya mengembalikan err apa adanya. Untuk membuat error lebih "manusiawi", Anda bisa memodifikasi bagian repository atau service.

Contoh logika yang bisa Anda terapkan di internal/service/auth.go:

// Logika pseudo-code untuk Anda terapkan manual:
user, err := repository.CreateUser(ctx, req.Email, string(hashed))
if err != nil {
    // Jika error mengandung kata-kata tertentu dari driver pgx/postgres
    // biasanya error unique violation mengandung kode "23505"
    if strings.Contains(err.Error(), "23505") { 
        return nil, errors.New("email already exists")
    }
    return nil, err
}

2. Mengapa menggunakan context.Context?
Anda akan melihat ctx di hampir semua fungsi. Dalam Go, context digunakan untuk:

- Cancellation: Jika user membatalkan request HTTP, Go bisa memberitahu database untuk berhenti memproses query tersebut agar tidak membuang resource.
- Timeout: Membatasi waktu eksekusi query.

3. Mengapa bcrypt?
Jangan pernah menyimpan password dalam plain text. bcrypt melakukan salting secara otomatis sehingga meskipun dua user punya password sama, hash yang disimpan di database akan berbeda.

### 📚 Penjelasan untuk Pembelajaran Go & Security

Agar Anda lebih paham saat mempraktekkannya, berikut penjelasannya:

1. HS256 (Symmetric Signing):
   - Menggunakan satu kunci rahasia yang sama untuk membuat token DAN memverifikasi token.
   - Sangat cocok jika Backend Anda adalah satu kesatuan (monolith) yang mengelola auth dan data.
   - Implementasi: JWT_SECRET → HS256 → Token.
2. ES256/RS256 (Asymmetric Signing):
   - Menggunakan sepasang kunci: Private Key untuk membuat token, dan Public Key untuk verifikasi.
   - Cocok untuk arsitektur Microservices. Service A (Auth) membuat token dengan Private Key, lalu Service B (Payment) bisa memverifikasi token tersebut hanya dengan Public Key tanpa perlu tahu kunci rahasianya.
   - Implementasi: PrivateKey Object → ES256 → Token.

Ringkasan Perbaikan:
Karena Anda menyimpan secret di .env sebagai string sederhana, Anda wajib menggunakan SigningMethodHS256.