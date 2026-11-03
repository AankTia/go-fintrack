# Day 4: Database Schema & Models

Untuk mengimplementasikan Day 4: Database Schema & Models, kita akan membangun "fondasi data". Karena Anda sedang belajar Go, saya akan menjelaskan konsepnya terlebih dahulu sebelum memberikan kodenya.

Di Go, kita biasanya memisahkan antara Database Schema (apa yang ada di PostgreSQL) dan Go Models (bagaimana data itu direpresentasikan di dalam kode).

Berikut adalah step-by-step detail untuk Day 4:

———

### Step 1: Membuat Migrasi Database (Schema)

Penjelasan:
Kita tidak membuat tabel secara manual di database (lewat GUI/Terminal), melainkan menggunakan file `.sql` (migrasi). Tujuannya agar setiap developer punya struktur DB yang sama. Kita butuh dua tabel: `accounts` (untuk tempat uang) dan `categories` (untuk pengelompokan transaksi).

Apa yang harus dilakukan:
Buat dua file baru di folder `migrations/`:

1. `000002_create_accounts_table.up.sql` (untuk membuat tabel)
2. `000002_create_accounts_table.down.sql` (untuk menghapus tabel jika terjadi kesalahan/rollback)
3. `000003_create_categories_table.up.sql`
4. `000003_create_categories_table.down.sql`

Implementasi Kode:
(Silakan buat file-file tersebut dengan isi berikut)

`000002_create_accounts_table.up.sql`:

```sql
CREATE TABLE accounts (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
name VARCHAR(100) NOT NULL,
type VARCHAR(20) NOT NULL, -- bank, cash, e-wallet
balance BIGINT NOT NULL DEFAULT 0, -- disimpan dalam unit terkecil (misal: Rupiah)
currency VARCHAR(3) NOT NULL DEFAULT 'IDR',
created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_accounts_user_id ON accounts(user_id);
```

`000002_create_accounts_table.down.sql`:

```sql
DROP TABLE IF EXISTS accounts;
```

`000003_create_categories_table.up.sql`:

```sql
CREATE TABLE categories (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
user_id UUID REFERENCES users(id) ON DELETE CASCADE, -- nullable untuk kategori global
name VARCHAR(100) NOT NULL,
type VARCHAR(20) NOT NULL, -- income, expense
icon VARCHAR(50),
created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_categories_user_id ON categories(user_id);
```

`000003_create_categories_table.down.sql`:

```sql
DROP TABLE IF EXISTS categories;
```

———

### Step 2: Membuat Go Model (Entities)

Penjelasan:
Model adalah representasi tabel database di dalam Go. Kita menggunakan struct.

- Tipe Data: Untuk uang, kita gunakan `int64` (BIGINT di DB). Jangan pernah gunakan `float64` untuk uang karena masalah presisi pembulatan.
- UUID: Kita gunakan string atau library google/uuid (tergantung yang sudah terinstall, namun umumnya string untuk kemudahan di DTO).

Apa yang harus dilakukan:
Buat file baru `internal/model/account.go` dan `internal/model/category.go`.

Implementasi Kode:

`internal/model/account.go`:

```go
package model

import "time"

// Account adalah entity internal yang merepresentasikan tabel accounts
type Account struct {
ID string `json:"id"`
UserID string `json:"user_id"`
Name string `json:"name"`
Type string `json:"type"` // bank, cash, e-wallet
Balance int64 `json:"balance"`
Currency string `json:"currency"`
CreatedAt time.Time `json:"created_at"`
UpdatedAt time.Time `json:"updated_at"`
}

// AccountRequest adalah DTO untuk validasi input saat Create/Update
type AccountRequest struct {
Name string `json:"name" validate:"required"`
Type string `json:"type" validate:"required,oneof=bank cash e-wallet"`
Balance int64 `json:"balance"`
Currency string `json:"currency"`
}
```

`internal/model/category.go`:

```go
package model

import "time"

// Category adalah entity internal yang merepresentasikan tabel categories
type Category struct {
ID string `json:"id"`
UserID \*string `json:"user_id"` // Menggunakan pointer karena boleh null (global category)
Name string `json:"name"`
Type string `json:"type"` // income, expense
Icon string `json:"icon"`
CreatedAt time.Time `json:"created_at"`
UpdatedAt time.Time `json:"updated_at"`
}

// CategoryRequest adalah DTO untuk validasi input
type CategoryRequest struct {
Name string `json:"name" validate:"required"`
Type string `json:"type" validate:"required,oneof=income expense"`
Icon string `json:"icon"`
}
```

———

### Step 3: Mendefinisikan Repository Interface

Penjelasan:
Di Go, kita menggunakan interface agar kode kita tidak kaku. Interface adalah "kontrak".
Misalnya: "Saya tidak peduli bagaimana cara Anda mengambil data dari DB, yang penting Anda punya fungsi GetByUserID". Ini sangat berguna untuk testing (Mocking).

Apa yang harus dilakukan:
Buat file `internal/repository/account_repository.go` dan `internal/repository/category_repository.go`.

Implementasi Kode:

`internal/repository/account_repository.go`:

```go
package repository

import (
  "context"
  "github.com/yourusername/go-fintrack/internal/model" // Sesuaikan dengan module name di go.mod
)

type AccountRepository interface {
  Create(ctx context.Context, account *model.Account) error
  GetByID(ctx context.Context, id string) (*model.Account, error)
  GetByUserID(ctx context.Context, userID string) ([]\*model.Account, error)
  UpdateBalance(ctx context.Context, id string, amount int64) error
}
```

`internal/repository/category_repository.go`:


```go
package repository

import (
  "context"
  "github.com/yourusername/go-fintrack/internal/model" // Sesuaikan dengan module name di go.mod
)

type CategoryRepository interface {
  Create(ctx context.Context, category *model.Category) error
  GetByUserID(ctx context.Context, userID string) ([]*model.Category, error)
  Delete(ctx context.Context, id string) error
}
```

———

### Ringkasan untuk Proses Belajar Anda:

1. Migrations: Mengelola struktur DB → SQL.
2. Model Struct: Mengelola data di memori Go → Go Struct.
3. DTO (Request): Memisahkan data yang masuk dari user dengan data yang tersimpan di DB (untuk keamanan).
4. Interface: Menentukan apa yang bisa dilakukan oleh database tanpa menentukan bagaimana caranya → Abstraksi.
