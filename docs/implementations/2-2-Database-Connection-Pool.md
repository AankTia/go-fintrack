## 2.2 Database Connection Pool (gemma4:12b-mlx high)

### Langkah 1: Instal Driver Database

Kita akan menggunakan `pgx`, driver PostgreSQL paling populer dan modern untuk Go.

Yang perlu dilakukan:
Jalankan perintah ini di terminal pada root project:

```bash
go get github.com/jackc/pgx/v5
```

- Penjelasan: `go get` akan mengunduh library ke dalam module kita. Kita menggunakan pgx karena ia mendukung fitur-fitur terbaru PostgreSQL dan memiliki performa yang sangat baik di Go.

———

### Langkah 2: Buat File Koneksi di internal/platform/database.go

Kita perlu tempat untuk menyimpan logika koneksi database agar tidak berantakan.

Yang perlu dilakukan:
Buat file baru di internal/platform/database.go dan masukkan kode berikut:

```go
package platform

import (
"context"
"fmt"
"os"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"

)

// DB adalah variabel global untuk menyimpan pool koneksi
// Kita gunakan variabel global agar bisa diakses oleh repository di seluruh aplikasi
var DB \*pgxpool.Pool

// Connect menginisialisasi koneksi ke database
// Kita mengembalikan error jika koneksi gagal agar aplikasi bisa berhenti (fail-fast)
func Connect() error {
// Ambil URL database dari environment variable (diatur di .env)
connStr := os.Getenv("DATABASE_URL")
if connStr == "" {
return fmt.Errorf("DATABASE_URL is not set in environment")
}

    // Buat pool koneksi
    // pgxpool.Connect akan membuka beberapa koneksi sekaligus ke database
    config, err := pgxpool.ParseConfig(connStr)
    if err != nil {
    	return fmt.Errorf("failed to parse config: %w", err)
    }

    // Buat pool
    pool, err := pgxpool.Connect(context.Background(), config)
    if err != nil {
    	return fmt.Errorf("failed to connect to database: %w", err)
    }

    // Simpan ke variabel global
    DB = pool
    fmt.Println("Successfully connected to the database")
    return nil

}

// Close menutup semua koneksi dalam pool saat aplikasi dimatikan
func Close() {
if DB != nil {
DB.Close()
}
}
```

- Penjelasan:
  - pgxpool.Pool: Ini adalah "kolam" koneksi kita. Ia akan mengelola kapan harus membuka koneksi baru, kapan harus menutupnya, dan kapan harus menggunakannya kembali.
  - os.Getenv("DATABASE_URL"): Kita mengambil alamat database dari .env agar tidak membocorkan password di dalam kode.
  - Global Variable DB: Di Go, cara paling umum untuk akses database di banyak layer adalah dengan menyimpan pool di variabel global (namun tetap dibatasi oleh scope package).

———

### Langkah 3: Tambahkan Fitur Ping (Verifikasi)

Penting untuk memastikan database benar-benar "nyala" saat aplikasi pertama kali dijalankan, bukan saat user pertama kali klik tombol di aplikasi.

Yang perlu dilakukan:
Tambahkan fungsi Ping di dalam file yang sama (internal/platform/database.go):

```go
// Ping memastikan koneksi ke database masih aktif
func Ping() error {
ctx := context.Background()
err := DB.Ping(ctx)
if err != nil {
return fmt.Errorf("database ping failed: %w", err)
}
return nil
}
```

- Penjelasan: Ping mengirim perintah sederhana ke database. Jika gagal, berarti ada masalah pada firewall, password salah, atau database belum siap.

———

### Langkah 4: Inisialisasi di cmd/api/main.go

Sekarang kita harus memanggil fungsi koneksi tersebut saat aplikasi dijalankan.

Yang perlu dilakukan:
Edit file cmd/api/main.go Anda untuk memanggil fungsi Connect dan Ping:
```go
package main

import (
"fmt"
"log"
"net/http"

    "github.com/AankTia/go-fintrack/internal/platform" // Pastikan path benar
    "github.com/joho/godotenv"

)

func main() {
// 1. Load .env
if err := godotenv.Load(); err != nil {
log.Println("Warning: .env file not found, using system env")
}

    // 2. Connect ke Database
    if err := platform.Connect(); err != nil {
    	log.Fatalf("Could not connect to DB: %v", err)
    }

    // 3. Verifikasi koneksi dengan Ping
    if err := platform.Ping(); err != nil {
    	log.Fatalf("Database not ready: %v", err)
    }
    fmt.Println("Database connection verified!")

    // 4. Jalankan Server (contoh sederhana)
    fmt.Println("Server starting on :8080")
    // ... kode router Anda di sini ...

}
```

- Penjelasan:
  - Kita melakukan Connection Testing di awal. Jika database mati, aplikasi akan langsung berhenti (log.Fatalf). Ini jauh lebih baik daripada aplikasi berjalan tapi semua fitur database error (tidak bisa diakses).

———

### Ringkasan yang Anda pelajari di sini:

1. Connection Pooling: Cara Go mengelola banyak koneksi sekaligus tanpa membuat server kewalahan.
2. Dependency Management: Menggunakan go get untuk library pgx.
3. Environment Handling: Cara mengambil konfigurasi sensitif (DB URL) dari .env menggunakan os.Getenv.
4. Fail-Fast Principle: Memastikan aplikasi berhenti segera jika koneksi database gagal di awal startup.
