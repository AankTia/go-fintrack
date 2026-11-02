# 🗓️ WEEK 4: Dashboard & Analytics (Detailed Breakdown)

Goal Utama: Menyediakan endpoint ringkasan keuangan (cash flow) dan breakdown kategori untuk bahan visualisasi dashboard.

## 🟦 Day 10: Dashboard Summary API

Fokus: Membuat "Single Source of Truth" untuk kondisi keuangan pengguna saat ini.

### 10.1 Aggregation Logic (Repository Layer)

Implementasikan query agregasi SQL untuk performa tinggi (jangan tarik semua data ke Go lalu dihitung manual):

- [ ] Total Balance: Sum dari semua saldo akun milik user.
- [ ] Total Income (Monthly): Sum amount dari transactions tipe income pada bulan berjalan.
- [ ] Total Expense (Monthly): Sum amount dari transactions tipe expense pada bulan berjalan.
- [ ] Net Cash Flow: Total Income - Total Expense.

### 10.2 Summary Service & Model

- [ ] Buat model DashboardSummary yang menampung semua angka di atas.
- [ ] Implementasikan service yang mengorkestrasi pemanggilan repository.
- [ ] Tambahkan logic untuk menangani kasus "No Data" (kembalikan 0, bukan null).

### 10.3 Summary Handler

- [ ] GET /api/v1/dashboard/summary:
  - Mengambil user_id dari token.
  - Return summary dalam format JSON.
  - Tambahkan query parameter month & year untuk melihat summary bulan-bulan sebelumnya.

### 10.4 Verification

- [ ] Test: Tambah transaksi income → Refresh summary → Pastikan Total Income dan Total Balance naik.
- [ ] Test: Ubah bulan → Pastikan summary berubah sesuai data bulan tersebut.

———

## 🟦 Day 11: Category Breakdown & Charts Data

Fokus: Menyediakan data untuk grafik (misal: Pie Chart pengeluaran).

### 11.1 Expense by Category Logic

- [ ] Repository Query: Gunakan GROUP BY category_id untuk menghitung total pengeluaran per kategori.
- [ ] Join Query: Gabungkan dengan tabel categories untuk mendapatkan nama kategori dan icon-nya.
- [ ] Sorting: Urutkan dari nominal terbesar ke terkecil.

### 11.2 Breakdown Service & Model

- [ ] Buat model CategoryBreakdown (Kategori Name, Total Amount, Percentage).
- [ ] Implementasikan logic penghitungan persentase terhadap total pengeluaran bulanan.

### 11.3 Breakdown Handler

- [ ] GET /api/v1/dashboard/breakdown:
  - Filter berdasarkan bulan dan tahun.
  - Return list kategori beserta nominal pengeluarannya.

### 11.4 Verification

- [ ] Test: Buat beberapa transaksi di kategori "Makanan" dan "Transport" → Pastikan breakdown menampilkan jumlah yang benar untuk masing-masing kategori.

———

## 🟦 Day 12: Security & Performance Testing

Fokus: Memastikan endpoint dashboard aman dan tidak lambat saat data transaksi sudah banyak.

### 12.1 Ownership Security Audit

- [ ] Verifikasi ulang: Apakah ada celah bagi user untuk melihat summary/breakdown user lain dengan mengganti ID di query? (Harus: TIDAK).
- [ ] Pastikan user_id selalu diambil dari JWT session, bukan dari request body/query.

### 12.2 Query Optimization

- [ ] Explain Analyze: Jalankan EXPLAIN ANALYZE pada query agregasi dashboard.
- [ ] Indexing: Pastikan index sudah terpasang pada transactions(user_id, transaction_date) untuk mempercepat filter bulanan.

### 12.3 Load Test (Basic)

- [ ] Masukkan data dummy transaksi dalam jumlah besar (misal: 1000+ transaksi per user).
- [ ] Ukur latency endpoint /summary dan /breakdown. Target: Response time < 200ms.

### 12.4 Final Review Week 4

- [ ] Review: Apakah ada penggunaan floating point dalam perhitungan persentase? (Jika ada, pastikan hanya digunakan di layer paling akhir untuk display, perhitungan inti tetap integer).

———

## ✅ Acceptance Criteria (End of Week 4)

- [ ] Accuracy: Total balance di dashboard sama persis dengan jumlah saldo di semua akun.
- [ ] Performance: Query agregasi berjalan cepat meskipun data transaksi bertambah.
- [ ] Data Isolation: User hanya bisa melihat statistik keuangan mereka sendiri.
- [ ] API Ready: Endpoint /summary dan /breakdown siap dikonsumsi oleh frontend Next.js di Week 6.