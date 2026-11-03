# 🗓️ WEEK 2: Database Design & Core Entities (Detailed Breakdown)

Goal Utama: Model data lengkap, migrasi schema terpasang, dan CRUD dasar untuk Akun & Kategori selesai.

## 🟦 Day 4: Database Schema & Models

Fokus: Mendefinisikan "kebenaran" data di level database dan memetakannya ke Go structs.

### 4.1 Database Schema (Migrations)

Buat file migrasi baru untuk entitas utama:

- [x] Table accounts: id (UUID), user_id (FK → users), name (string), type (enum: bank, cash, e-wallet), balance (BIGINT/NUMERIC), currency (string), created_at, updated_at.
- [x] Table categories: id (UUID), user_id (FK → users, nullable untuk kategori default), name (string), type (enum: income, expense), icon (string), created_at, updated_at.
- [x] Constraints: Tambahkan INDEX pada user_id di semua tabel untuk optimasi query.

### 4.2 Go Model Mapping

Implementasikan struct di internal/model/:

- [x] Account & Category entity (untuk internal/DB).
- [x] AccountRequest & CategoryRequest DTO (untuk validasi input).
- [x] Pastikan field uang menggunakan int64 (representasi sen/unit terkecil) untuk menghindari floating point error.

### 4.3 Repository Interface

Definisikan contract di internal/repository/:

- [x] AccountRepository interface: Create, GetByID, GetByUserID, UpdateBalance.
- [x] CategoryRepository interface: Create, GetByUserID, Delete.

———

## 🟦 Day 5: Account Management (CRUD)

Fokus: Implementasi fitur pengelolaan akun keuangan pengguna.

### 5.1 Account Repository Implementation

- [ ] Implementasi method repository menggunakan parameterized query (cegah SQL Injection).
- [ ] Wajib: Setiap query SELECT, UPDATE, atau DELETE harus menyertakan WHERE user_id = ? untuk menjamin isolasi data.

### 5.2 Account Service Logic

- [ ] CreateAccount: Validasi nama akun dan tipe akun.
- [ ] GetMyAccounts: Ambil semua akun milik user yang sedang login.
- [ ] Implementasi business rule: Misalnya, satu user minimal harus memiliki satu akun.

### 5.3 Account Handler & Routing

- [ ] POST /api/v1/accounts: Create new account.
- [ ] GET /api/v1/accounts: List all user accounts.
- [ ] Integrasi dengan Auth Middleware untuk mengambil user_id dari JWT token.

### 5.4 Verification

- [ ] Test: Buat akun → Cek apakah tersimpan dengan user_id yang benar.
- [ ] Test: Coba akses akun milik user lain via ID → Harus return 403 Forbidden/404 Not Found.

———

## 🟦 Day 6: Category Management (CRUD)

Fokus: Implementasi pengelompokan transaksi.

### 6.1 Category Repository Implementation

- [ ] Implementasi CRUD kategori.
- [ ] Logic khusus: Kategori bisa bersifat "Global" (user_id IS NULL) atau "Custom" (user_id NOT NULL).

### 6.2 Category Service Logic

- [ ] CreateCategory: Pastikan tidak ada nama kategori duplikat untuk user yang sama.
- [ ] GetCategories: Gabungkan kategori global + kategori custom milik user.

### 6.3 Category Handler & Routing

- [ ] POST /api/v1/categories: Create custom category.
- [ ] GET /api/v1/categories: Get all available categories.
- [ ] DELETE /api/v1/categories/{id}: Hapus custom category.

### 6.4 Verification & Unit Test

- [ ] Unit Test: Tulis test untuk logic GetCategories (memastikan gabungan global & custom benar).
- [ ] Test: Buat kategori → Hapus kategori → Pastikan tidak bisa menghapus kategori global.

———

## ✅ Acceptance Criteria (End of Week 2)

- [ ] Schema Valid: Tabel accounts dan categories ada di DB dengan foreign key ke users.
- [ ] No Floats: Tidak ada penggunaan float32/64 untuk field saldo/nominal uang di seluruh kode.
- [ ] Strict Ownership: Terverifikasi bahwa user tidak bisa memanipulasi akun/kategori milik user lain.
- [ ] API Working: Endpoint CRUD Account dan Category berfungsi dan terproteksi JWT.
- [ ] Code Quality: Repository menggunakan interface, memudahkan unit testing di minggu berikutnya.