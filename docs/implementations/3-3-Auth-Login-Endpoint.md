# 3.3 Auth: Login Endpoint (POST /auth/login)

### Tahap 1: Persiapan Library & Konfigurasi

Sebelum menulis kode, Anda perlu library untuk menangani JWT.

1. Install Library JWT
Jalankan perintah ini di terminal:

```bash
go get github.com/golang-jwt/jwt/v5
```

- Penjelasan: `jwt/v5` adalah standar industri saat ini untuk Go. Library ini menyediakan fungsi untuk membuat (sign) dan memverifikasi (parse) token.

2. Update `.env`
Tambahkan key untuk menandatangani token agar aman.

```conf
JWT_SECRET=ganti-dengan-string-rahasia-yang-panjang-dan-acak
```

- Penjelasan: `JWT_SECRET` adalah kunci rahasia. Hanya server Anda yang tahu kunci ini. Jika seseorang mencoba mengubah isi token (misal mengubah `user_id`), tanda tangannya tidak akan cocok karena mereka tidak punya kunci ini.

———

### Tahap 2: Implementasi Kode (Apply Manual)

Lakukan perubahan pada file-file berikut secara berurutan.

#### Step A: Update `internal/service/auth.go`

Kita akan mengubah LoginUser agar tidak lagi menggunakan token dummy, melainkan membuat JWT asli.

Yang harus Anda lakukan:

1. Import library jwt dan time.
2. Buat fungsi helper generateJWT di bawah fungsi LoginUser.
3. Update isi fungsi LoginUser.

Implementasi Kode:
(Silakan ganti isi internal/service/auth.go Anda dengan logika ini)

```go
package service

import (
	"context"
	"errors"
	"time" // Tambahkan ini untuk waktu expired token

	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/AankTia/go-fintrack/internal/platform" // Ambil config dari sini
	"github.com/AankTia/go-fintrack/internal/repository"
	"github.com/golang-jwt/jwt/v5" // Tambahkan library JWT
	"golang.org/x/crypto/bcrypt"
)

// ... (fungsi RegisterUser biarkan tetap sama) ...

func LoginUser(ctx context.Context, req *model.LoginRequest) (string, error) {
	// 1. Cari user berdasarkan email
	user, err := repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		// Jangan beri tahu "email tidak ditemukan" agar tidak memudahkan attacker
		return "", errors.New("invalid credentials")
	}

	// 2. Verifikasi password (Input Plaintext vs Hash di DB)
	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// 3. Generate JWT Token menggunakan ID User
	token, err := generateJWT(user.ID)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

// Helper function untuk membuat token
func generateJWT(userID int64) (string, error) {
	// Claims adalah data yang ingin kita simpan di dalam token
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
		"iat":     time.Now().Unix(),                   // Issued At (Kapan dibuat)
	}

	// Buat objek token dengan algoritma HMAC SHA256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tanda tangani token menggunakan Secret Key dari Config
	// platform.Config.JWTSecret harus sudah terisi saat LoadConfig()
	signedToken, err := token.SignedString([]byte(platform.Config.JWTSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
```

Penjelasan untuk Pelajar Go:

- `jwt.MapClaims`: Ini seperti `map[string]interface{}`. Kita menyimpan `user_id` di sini agar saat token dikirim balik ke server, server tahu siapa user-nya tanpa harus query DB lagi.
- `time.Now().Add(...).Unix()`: JWT mewajibkan waktu kadaluarsa dalam format Unix Timestamp (integer detik sejak 1970).
- `SignedString`: Inilah proses enkripsi/tanda tangan. Token yang dihasilkan terdiri dari 3 bagian: Header.Payload.Signature.

———

#### Step B: Update internal/handler/auth.go

Kita akan memperbaiki status code agar lebih semantik (benar secara standar HTTP).

Yang harus Anda lakukan:
Ubah bagian `LoginHandler` untuk membedakan error 401 (Unauthorized) dan 400 (Bad Request).

Implementasi Kode:
(Update fungsi `LoginHandler` Anda menjadi seperti ini)

```go
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	token, err := service.LoginUser(r.Context(), &req)
	if err != nil {
		// Jika error adalah "invalid credentials", kirim 401 Unauthorized
		if err.Error() == "invalid credentials" {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		// Error lainnya kirim 400 atau 500
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Kirim response sukses 200 OK
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
```

Penjelasan untuk Pelajar Go:

- `w.Header().Set(...)`: Sangat penting untuk memberi tahu client (Postman/Frontend) bahwa data yang dikirim adalah JSON.
- `http.StatusUnauthorized`: Menggunakan konstanta bawaan Go (401) lebih baik daripada menulis angka manual agar kode lebih mudah dibaca.

———

### Tahap 3: Verifikasi (Cara Testing)

Setelah Anda apply kode di atas, lakukan pengujian berikut di Postman:

1. Skenario Gagal (Wrong Password):
   - POST /api/v1/auth/login
   - Body: {"email": "user@test.com", "password": "salah-password"}
   - Ekspektasi: Status 401 Unauthorized dengan pesan "invalid credentials".
2. Skenario Gagal (Wrong Email):
   - POST /api/v1/auth/login
   - Body: {"email": "tidakada@test.com", "password": "password123"}
   - Ekspektasi: Status 401 Unauthorized dengan pesan "invalid credentials".
3. Skenario Berhasil:
   - POST /api/v1/auth/login
   - Body: {"email": "user@test.com", "password": "password-benar"}
   - Ekspektasi: Status 200 OK dengan body {"token": "eyJhbG..."}.

### Checklist Belajar untuk Anda:

- [x] Apakah saya sudah menginstal library JWT?
- [x] Apakah saya sudah menambahkan JWT_SECRET di .env?
- [x] Apakah saya paham kenapa kita menggunakan jwt.MapClaims?
- [x] Apakah saya paham perbedaan http.StatusBadRequest (400) dan http.StatusUnauthorized (401)?
