## 🎯 2.3 Migration Setup – Langkah‑langkah lengkap (untuk Go + PostgreSQL)

Berikut panduan “step‑by‑step” yang dirancang khusus untuk pemula.
Setiap poin dilengkapi penjelasan singkat (mengapa langkah itu penting) dan contoh perintah / kode yang dapat Anda salin langsung.

> Catatan:
>
> - Simpan semua file migration dalam folder migrations/ (atau db/migrations/ sesuai preferensi).
> - Semua perubahan pada skema database harus melalui migration—ini menjaga history dan memudahkan rollback.

———

## 1️⃣ Pilih & Instal Tool Migration

| Tool | Kelebihan | Cara Instal |
|---|---|---|
| golang-migrate | Populer, dukung banyak DB, CLI & library Go. | go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest |
| goose | Mudah, scriptable, CLI dan Go library. | go install github.com/pressly/goose@latest |

> Saya rekomendasikan golang-migrate (biasanya sudah ada di banyak proyek Go).
> Jika Anda belum pernah instal, jalankan di terminal (pada root proyek):

# Pastikan GOPATH/bin ada di PATH
```bash
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

> Periksa instalasi:
>
> migrate -version
# Output: migrate version v4.x.x

———

## 2️⃣ Buat Struktur Folder migrations/
```bash
mkdir migrations
```

> Folder ini akan menampung semua file SQL yang berisi up (buat/ubah) dan down (undo) migration.
> Letakkan folder ini di root proyek (atau db/migrations/) agar mudah di‑commit.

———

## 3️⃣ Buat Migration Pertama (contoh: users table)

1. Generate file template
```bash
migrate create -ext sql -dir migrations -seq create_users_table
```

- `-ext sql` → file berformat .sql.
- `-dir migrations` → folder target.
- `-seq` → gunakan nomor urut (mis. 001_create_users_table.up.sql).

Hasil:
```text
migrations/
├─ 001_create_users_table.up.sql
├─ 001_create_users_table.down.sql
```

2. Edit up dan down

001_create_users_table.up.sql:
```sql
-- +migrate Up
CREATE TABLE IF NOT EXISTS users (
    id          BIGSERIAL PRIMARY KEY,
    email       TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT now()
);
```

001_create_users_table.down.sql:
```sql
-- +migrate Down
DROP TABLE IF EXISTS users;
```
> Penjelasan:
>
> - +migrate Up/Down adalah komentar yang dibaca oleh migrator.
> - BIGSERIAL cocok untuk id yang auto‑increment.
> - Gunakan TIMESTAMP WITH TIME ZONE agar aman lintas zona waktu.

———

## 4️⃣ Tambahkan Config DB ke .env (jika belum ada)

Buat file `.env` (jangan commit, masukkan ke .gitignore):
```conf
DATABASE_URL=postgres://user:pass@localhost:5432/fintrack?sslmode=disable
```

> Pastikan Docker Compose telah mengatur user/password sesuai.

———

## 5️⃣ Jalankan Migration Secara Manual (untuk cek)

# Pastikan variable sudah ter-load (atau export di shell)
```bash
export DATABASE_URL=postgres://admin:rahasia@localhost:5432/fintrack_db?sslmode=disable

migrate -path migrations -database "$DATABASE_URL" up
```

- Output akan menunjukkan success dan number of migrations applied.
- Cek di DB (psql) apakah tabel users sudah ada.

> Rollback (jika salah):
>
> `migrate -path migrations -database "$DATABASE_URL" down`

———

## 6️⃣ Jalankan Migration Otomatis Saat Startup (opsional tapi disarankan)

Tambahkan kode di internal/platform/database.go (atau file init yang cocok) sebelum memulai server:

```go
package platform

import (
    "fmt"
    "log"
    "os"

    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbURL string) error {
    m, err := migrate.New(
        "file://migrations",
        dbURL,
    )
    if err != nil {
        return fmt.Errorf("migration init: %w", err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("migration up: %w", err)
    }
    log.Println("✅ Database migrations applied")
    return nil
}
```

- Panggil `RunMigrations` sebelum server start (cmd/api/main.go).

> Manfaat:
>
> - Semua developer/CI menjalankan migration otomatis.
> - Tidak ada “manual run” yang terlupa.

———

## 7️⃣ Commit & Push
```bash
git add migrations/*.sql .env.example
git commit -m "Add migration: create users table"
git push
```

> `.env.example` hanya contoh yang berisi kunci tanpa password, tidak berisi data sensitif.

———

## 📌 Ringkasan Poin Penting

| Langkah | Tujuan | Hasil Akhir |
|---|---|---|
| 1. Install migrator | Memiliki tool CLI untuk migration | migrate CLI tersedia |
| 2. Buat folder migrations/ | Penyimpanan file migration | Folder kosong |
| 3. Generate migration template | Template up.sql + down.sql | File 001_create_users_table.* |
| 4. Edit SQL | Definisikan skema | Tabel users siap |
| 5. Set environment DATABASE_URL | Koneksi DB | DB URL terconfig |
| 6. Jalankan migrasi manual | Verifikasi migrasi | Tabel sudah ada |
| 7. Tambah migrasi otomatis di Go | Auto‑migration pada startup | Migrate di main.go |
| 8. Commit migration files | Repositori lengkap | Migration versioned |

Dengan mengikuti langkah‑langkah di atas, Anda sudah dapat:

- Menentukan skema database secara terstruktur.
- Menjaga history perubahan skema (versi).
- Menjamin setiap developer/CI memiliki skema yang sama.