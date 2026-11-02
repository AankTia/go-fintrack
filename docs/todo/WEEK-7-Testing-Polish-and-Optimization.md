# 🗓️ WEEK 7: Testing, Polish & Optimization (Detailed Breakdown)

Goal Utama: Mencapai standar kualitas produksi dengan coverage test yang cukup, performa query yang optimal, dan kode yang bersih.

## 🟦 Day 19: Unit Tests (Backend)

Fokus: Menguji logika bisnis secara terisolasi untuk mencegah regresi di masa depan.

### 19.1 Service Layer Unit Tests

Gunakan mock untuk repository agar test tidak bergantung pada database asli:

- [ ] Auth Service: Test validasi password, error saat email duplikat, dan validitas token JWT.
- [ ] Transaction Service:
  - Test logika update saldo (Income vs Expense).
  - Test skenario saldo tidak cukup (jika ada aturan anti-overdraft).
- [ ] Budget Service: Test perhitungan persentase pemakaian budget dan penentuan status is_overlimit.

### 19.2 Repository Integration Tests

Gunakan database test terpisah (atau Docker container sementara):

- [ ] Query Validation: Pastikan query GetTransactions dengan filter tanggal menghasilkan data yang tepat.
- [ ] Ownership Test: Verifikasi bahwa query repository benar-benar menyertakan user_id di setiap WHERE clause.
- [ ] Atomicity Test: Verifikasi bahwa jika update saldo gagal, transaksi tidak tersimpan di DB (test rollback).

### 19.3 Test Execution & Coverage

- [ ] Jalankan go test ./....
- [ ] Cek coverage: Pastikan fungsi-fungsi kritikal di layer service tercover minimal 80%.

———

## 🟦 Day 20: E2E Tests (Playwright)

Fokus: Memastikan "Happy Path" pengguna berjalan mulus dari UI sampai DB.

### 20.1 Playwright Setup

- [ ] Install Playwright di project Next.js.
- [ ] Setup environment test (User dummy untuk login).
- [ ] Buat custom command untuk login agar tidak perlu login manual di setiap test case.

### 20.2 Critical Flow Testing (The "Must-Pass" List)

Tulis skrip E2E untuk alur berikut:

- [ ] Auth Flow: Register → Login → Landing di Dashboard.
- [ ] Transaction Flow: Tambah Akun → Tambah Transaksi → Verifikasi Saldo Akun Berubah.
- [ ] Budget Flow: Set Budget → Tambah Transaksi → Verifikasi Progress Bar Budget Berubah Warna.
- [ ] Dashboard Flow: Verifikasi chart muncul dan angka summary sesuai dengan transaksi yang dibuat.

### 20.3 Performance Baseline (Lighthouse)

- [ ] Jalankan Lighthouse audit pada halaman Dashboard.
- [ ] Optimasi: Pastikan tidak ada unused JS/CSS yang besar dan gambar/icon teroptimasi.

———

## 🟦 Day 21: Code Quality & Optimization

Fokus: Membersihkan "hutang teknis" (technical debt) dan mempercepat aplikasi.

### 21.1 Database Tuning

- [ ] Index Review: Pastikan index sudah ada pada:
  - transactions(user_id, transaction_date)
  - accounts(user_id)
  - budgets(user_id, period)
- [ ] Query Optimization: Cari query yang lambat menggunakan EXPLAIN ANALYZE dan optimalkan join-nya.

### 21.2 Code Polishing

- [ ] Linting: Jalankan golangci-lint di backend dan eslint di frontend.
- [ ] Formatting: Jalankan gofmt dan prettier di seluruh codebase.
- [ ] Cleanup: Hapus semua fmt.Println, console.log, dan komentar kode yang tidak terpakai.

### 21.3 Security Checklist Final Audit

Lakukan review akhir berdasarkan AGENTS.md (Section 7):

- [ ] [ ] Apakah semua password sudah di-hash dengan bcrypt?
- [ ] [ ] Apakah semua endpoint privat memiliki filter user_id?
- [ ] [ ] Apakah error response sudah bersih dari stack trace internal?
- [ ] [ ] Apakah input nominal uang divalidasi (tidak boleh negatif/nol)?

———

## ✅ Acceptance Criteria (End of Week 7)

- [ ] Stability: Semua unit test dan E2E test pass tanpa error.
- [ ] No Regressions: Penambahan fitur budget tidak merusak fitur transaksi.
- [ ] Performance: Response time API agregasi dashboard < 200ms.
- [ ] Clean Code: Tidak ada warning dari linter dan kode terformat rapi.
- [ ] Security Verified: Terbukti secara teknis bahwa User A tidak bisa mengakses data User B melalui API maupun UI.