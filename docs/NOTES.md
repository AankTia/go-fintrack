# NOTES

## Go Runtime & Module Init

### Cek/Instal Go

Cek versi yang terpasang:

```bash
brew install go
```

> Catatan: Pastikan $PATH mengarah ke /usr/local/bin/go (atau /opt/homebrew/bin/go pada Apple Silicon).
> Verifikasi lagi:

which go # → /usr/local/bin/go
go version

### Inisialisasi Go Module

Inisialisasi module dengan nama yang akan digunakan di GitHub:

```bash
go mod init github.com/AankTia/go-fintrack
```

Ini akan membuat file go.mod berisi:

```text
module github.com/AankTia/go-fintrack

go 1.26.4
```

Pastikan module sudah berfungsi, Buat file sederhana untuk test:

```bash
echo 'package main; import "fmt"; func main(){ fmt.Println("Hello Go!") }' > main.go
go run main.go   # Output: Hello Go!
```

Jika berhasil, module sudah siap.

### Verifikasi Build

Install `golangci-lint`:

```bash
brew install golangci-lint
golangci-lint run
```

Run go fmt for formatting

```bash
go fmt ./...
```

---

## Project Structure

Berikut adalah penjelasan fungsi tiap folder dengan padanan di Laravel/Rails:

| Folder               | Padanan di Laravel/Rails | Fungsi Utama di Go                                                                                                                                     |
| :------------------- | :----------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------- |
| cmd/api/             | public/ (Entry point)    | Tempat aplikasi "hidup". Hanya berisi file main.go yang menjalankan server.                                                                            |
| internal/handler/    | Controllers              | Menangani request HTTP. Tugasnya: Parsing JSON, cek validasi input dasar, dan mengirim response (200 OK, 400 Bad Request, dll).                        |
| internal/service/    | Service Layer            | "Otak" aplikasi. Di sinilah semua aturan bisnis (Business Logic) berada. Misalnya: "Bagaimana cara menghitung saldo?", "Apakah user boleh transaksi?". |
| internal/repository/ | Repository / Eloquent    | Urusan Database. Satu-satunya tempat yang boleh menjalankan Query SQL. Tidak boleh ada logic bisnis di sini.                                           |
| internal/model/      | Models / DTOs            | Mendefinisikan "Bentuk" data. Struktur data untuk tabel database dan struktur request/response.                                                        |
| internal/platform/   | Config / Infrastructure  | Fondasi aplikasi. Tempat koneksi ke Database, Load Config (dari .env), dan helper sistem lainnya.                                                      |
| internal/middleware/ | Middleware               | Filter request. Contoh: Auth check (JWT), Logger, Rate Limiting.                                                                                       |
| migrations/          | Migrations               | File .sql untuk mengubah struktur tabel di database.                                                                                                   |

**Mengapa harus pakai folder internal? (Penting untuk Go)**

Di bahasa Go, folder bernama `internal` memiliki aturan khusus: **hanya kode di dalam project yang sama yang bisa mengaksesnya.** Jika ada orang lain mengimport project Anda sebagai library, mereka tidak bisa mengakses folder `internal`. Ini adalah standar industri di Go untuk menjaga keamanan dan kebersihan kode (encapsulation).

### Visualisasi Alur Data (The Request Flow)

Request (Client) → Handler → Service → Repository → Database

1. Handler: "Ada request POST ke /login. Saya ambil JSON-nya."
2. Service: "Oke, saya cek apakah email ini ada? Saya bandingkan password-nya."
3. Repository: "Saya ambil data dari tabel users berdasarkan email ini."
4. Database: (Mengembalikan data).
