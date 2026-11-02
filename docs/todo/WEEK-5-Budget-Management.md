# 🗓️ WEEK 5: Budget Management (Detailed Breakdown)

Goal Utama: Pengguna dapat menetapkan batas pengeluaran bulanan per kategori dan memantau sisa budget secara real-time.

## 🟦 Day 13: Budget Model & CRUD

Fokus: Mendefinisikan struktur data budget dan menyediakan API pengelolaannya.

### 13.1 Budget Schema & Model

- [ ] Table budgets: id (UUID), user_id (FK), category_id (FK), amount (BIGINT), period (string, format: YYYY-MM), created_at, updated_at.
- [ ] Constraints: Tambahkan UNIQUE constraint pada (user_id, category_id, period) agar satu kategori hanya punya satu budget per bulan.
- [ ] Implementasikan struct Budget dan CreateBudgetRequest di internal/model/.

### 13.2 Budget Repository

- [ ] Implementasi CreateBudget, GetBudgetByPeriod, UpdateBudget.
- [ ] Ownership Check: Pastikan user_id selalu disertakan dalam query.

### 13.3 Budget Handler & Routing

- [ ] POST /api/v1/budgets: Membuat atau memperbarui budget untuk kategori tertentu di bulan tertentu.
- [ ] GET /api/v1/budgets: List semua budget aktif untuk bulan berjalan.
- [ ] Integrasi Auth Middleware untuk proteksi endpoint.

### 13.4 Verification

- [ ] Test: Buat budget "Makanan" untuk bulan November → Cek apakah tersimpan di DB.
- [ ] Test: Coba buat budget untuk kategori yang sama di bulan yang sama → Pastikan terjadi update (upsert) bukan duplikasi.

———

## 🟦 Day 14: Budget Calculation & Status

Fokus: Menghitung realisasi pengeluaran terhadap budget yang ditetapkan.

### 14.1 Budget Status Logic (Aggregation)

Implementasikan logic perhitungan di Service layer:

- [ ] Actual Spend: Query total transaksi tipe expense untuk category_id dan period tertentu.
- [ ] Remaining Budget: Budget Amount - Actual Spend.
- [ ] Status Percentage: (Actual Spend / Budget Amount) * 100.

### 14.2 Budget Status Endpoint

- [ ] GET /api/v1/budgets/status:
  - Mengembalikan list budget yang mencakup: Nama Kategori, Limit Budget, Total Terpakai, Sisa Budget, dan Persentase.
  - Support query parameter period (misal: ?period=2026-11).

### 14.3 Business Rule: Budget Overlimit

- [ ] Tambahkan flag is_overlimit (boolean) pada response jika Actual Spend > Budget Amount.

### 14.4 Verification

- [ ] Test: Set budget "Transport" Rp 500rb → Buat transaksi Transport Rp 200rb → Cek status: Sisa harus Rp 300rb (60% sisa).
- [ ] Test: Buat transaksi Transport tambahan Rp 400rb → Cek status: is_overlimit harus true.

———

## 🟦 Day 15: Budget Integration & Final Testing

Fokus: Menghubungkan modul budget dengan modul transaksi dan melakukan pengujian menyeluruh.

### 15.1 Transaction-Budget Integration

Meskipun kalkulasi dilakukan saat GET, kita perlu memastikan data konsisten:

- [ ] Review TransactionService: Pastikan tidak ada logic yang menghalangi update budget.
- [ ] (Opsional) Implementasikan trigger atau event sederhana jika ingin mengirim notifikasi saat budget terlampaui (untuk pengembangan selanjutnya).

### 15.2 Integration Testing (End-to-End Flow)

Tulis test case skenario lengkap:

- [ ] Scenario A: User buat budget → User transaksi → Budget berkurang → Selesai.
- [ ] Scenario B: User update budget (menaikkan limit) → Status is_overlimit berubah dari true kembali ke false.
- [ ] Scenario C: User menghapus kategori → Pastikan budget terkait juga terhapus (Cascade Delete) atau ditangani agar tidak terjadi orphaned data.

### 15.3 Security & Performance Audit

- [ ] Verifikasi: Apakah user bisa mengubah budget milik user lain? (Harus: TIDAK).
- [ ] Optimasi: Pastikan query sum(amount) untuk status budget menggunakan index pada category_id dan transaction_date.

———

## ✅ Acceptance Criteria (End of Week 5)

- [ ] Consistency: Status budget (terpakai vs limit) selalu akurat berdasarkan data transaksi terbaru.
- [ ] No Duplicates: Tidak ada double budget untuk kategori yang sama di bulan yang sama.
- [ ] Accuracy: Perhitungan sisa budget menggunakan int64 untuk menghindari error pembulatan.
- [ ] API Ready: Endpoint /budgets dan /budgets/status berfungsi dan terproteksi.
- [ ] Ownership: User hanya bisa mengelola budget milik mereka sendiri.