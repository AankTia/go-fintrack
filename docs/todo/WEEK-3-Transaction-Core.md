# 🗓️ WEEK 3: Transaction Core (Detailed Breakdown)

Goal Utama: Implementasi siklus hidup transaksi lengkap (CRUD) dengan update saldo otomatis dan integritas data yang terjamin.

## 🟦 Day 7: Transaction Create & Validation

Fokus: Memastikan setiap transaksi tercatat dengan benar dan mengubah saldo akun secara atomik.

### 7.1 Transaction Model & Schema

- [ ] Table transactions: id (UUID), user_id (FK), account_id (FK), category_id (FK), amount (BIGINT), type (enum: income, expense), description (text), transaction_date (date), created_at, updated_at.
- [ ] Implementasikan struct Transaction dan CreateTransactionRequest DTO di internal/model/.

### 7.2 Atomic Balance Update (The Critical Part)

Implementasikan Database Transaction (BEGIN/COMMIT/ROLLBACK) di level Service:

- [ ] CreateTransaction Flow:
  1. BEGIN Transaction.
  2. Insert data ke tabel transactions.
  3. Update saldo di tabel accounts (Tambah jika income, kurangi jika expense).
  4. COMMIT jika semua berhasil, ROLLBACK jika ada error.
- [ ] Validation: Pastikan saldo tidak menjadi negatif jika user tidak mengizinkan overdraft.

### 7.3 Transaction Handler

- [ ] POST /api/v1/transactions:
  - Validasi input (nominal harus > 0).
  - Verifikasi bahwa account_id dan category_id memang milik user tersebut (Ownership Check).
  - Panggil service → return 201 Created.

### 7.4 Verification

- [ ] Test: Create Income → Cek apakah saldo akun bertambah.
- [ ] Test: Create Expense → Cek apakah saldo akun berkurang.
- [ ] Test: Simulasi error saat update saldo → Pastikan transaksi tidak tersimpan (Rollback).

———

## 🟦 Day 8: Transaction List & Filtering

Fokus: Memberikan akses data transaksi yang efisien dan terfilter.

### 8.1 Transaction Repository (Querying)

- [ ] Implementasi GetTransactionsByUserID dengan join ke tabel accounts dan categories untuk mendapatkan nama akun/kategori.
- [ ] Filtering Logic: Tambahkan filter berdasarkan:
  - Rentang tanggal (start_date s/d end_date).
  - Tipe transaksi (income atau expense).
  - ID Akun tertentu.

### 8.2 Pagination & Sorting

- [ ] Implementasikan limit dan offset untuk mencegah overload data.
- [ ] Default sorting: transaction_date DESC.

### 8.3 List Handler

- [ ] GET /api/v1/transactions:
  - Parsing query parameters (filter, page, limit).
  - Return data dalam bentuk list dengan metadata pagination.

### 8.4 Verification

- [ ] Test: Filter transaksi bulan ini → Pastikan data akurat.
- [ ] Test: Request page 2 → Pastikan data yang muncul berbeda dengan page 1.

———

## 🟦 Day 9: Transaction Update & Delete

Fokus: Menangani perubahan data tanpa merusak konsistensi saldo.

### 9.1 Idempotent Update Logic

Update transaksi lebih kompleks karena mempengaruhi saldo. Flow yang benar:

- [ ] Update Transaction Flow (dalam satu DB Transaction):
  1. Ambil data transaksi lama.
  2. Reverse dampak transaksi lama ke saldo akun (kalau dulu expense, sekarang tambahkan kembali).
  3. Update data transaksi ke nilai baru.
  4. Apply dampak transaksi baru ke saldo akun (kalau sekarang income, tambahkan).
  5. COMMIT.

### 9.2 Transaction Deletion

- [ ] Delete Transaction Flow:
  1. Ambil data transaksi yang akan dihapus.
  2. Kembalikan saldo akun ke posisi semula (Reverse effect).
  3. Hapus record transaksi.
  4. COMMIT.

### 9.3 Update & Delete Handlers

- [ ] PATCH /api/v1/transactions/{id}: Update partial data.
- [ ] DELETE /api/v1/transactions/{id}: Hapus transaksi.
- [ ] Ownership Check: Pastikan user hanya bisa update/delete transaksi miliknya sendiri.

### 9.4 Verification

- [ ] Test: Ubah transaksi Expense Rp 10rb menjadi Income Rp 20rb → Cek apakah saldo akun akhir benar (bertambah Rp 30rb dari posisi sebelumnya).
- [ ] Test: Hapus transaksi → Cek apakah saldo akun kembali normal.

———

## ✅ Acceptance Criteria (End of Week 3)

- [ ] Data Integrity: Saldo akun selalu sinkron dengan total transaksi (tidak ada selisih).
- [ ] Atomicity: Tidak ada transaksi yang tercipta tanpa update saldo (dan sebaliknya).
- [ ] Strict Ownership: User tidak bisa melihat/mengubah transaksi milik user lain.
- [ ] Performance: List transaksi menggunakan pagination dan index pada user_id & transaction_date.
- [ ] Financial Rule: Nominal uang tetap konsisten menggunakan int64/BIGINT.