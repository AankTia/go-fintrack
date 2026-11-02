# 🗓️ WEEK 6: Frontend Integration (Detailed Breakdown)

Goal Utama: Mengintegrasikan seluruh fitur Backend ke dalam UI yang responsif, intuitif, dan aman.

## 🟦 Day 16: Next.js Setup & Auth UI

Fokus: Menyiapkan boilerplate frontend dan mengamankan akses aplikasi.

### 16.1 Frontend Project Scaffold

- [ ] Initialize Next.js App Router project: npx create-next-app@latest.
- [ ] Install core dependencies: axios, @tanstack/react-query, zod, react-hook-form.
- [ ] Setup shadcn/ui (Button, Input, Card, Form, Toast, Dialog).
- [ ] Konfigurasi .env.local untuk NEXT_PUBLIC_API_URL.

### 16.2 API Client Wrapper

- [ ] Buat instance axios yang terpusat.
- [ ] Implementasikan Interceptor:
  - Request: Otomatis menambahkan Authorization: Bearer <token> dari localStorage/cookie.
  - Response: Handle error 401 (Unauthorized) dengan mengarahkan user kembali ke halaman login.

### 16.3 Auth Pages (UI & Logic)

- [ ] Login Page: Form email & password → Simpan token → Redirect ke Dashboard.
- [ ] Register Page: Form lengkap → Redirect ke Login.
- [ ] Implementasikan Zod Validation pada form untuk feedback error instan sebelum request dikirim.

### 16.4 Verification

- [ ] Test: Alur Register → Login → Masuk Dashboard.
- [ ] Test: Coba akses /dashboard tanpa token → Harus redirect ke /login.

———

## 🟦 Day 17: Accounts & Transactions Frontend

Fokus: Implementasi manajemen data keuangan inti.

### 17.1 Account Management UI

- [ ] Account List: Tampilan Card yang menunjukkan Nama Akun, Tipe, dan Saldo Saat Ini.
- [ ] Create Account Modal: Form untuk tambah akun baru.
- [ ] Gunakan useQuery untuk fetch list akun dan useMutation untuk create akun (dengan invalidateQueries agar list terupdate otomatis).

### 17.2 Transaction List UI

- [ ] Transaction Table: List transaksi dengan kolom: Tanggal, Kategori, Keterangan, Nominal, dan Akun.
- [ ] Filter Bar: Input untuk filter berdasarkan bulan, tipe (income/expense), dan kategori.
- [ ] Implementasikan Pagination UI yang terhubung dengan query params di URL.

### 17.3 Transaction Form

- [ ] Create/Edit Modal:
  - Dropdown pilihan Akun (fetch dari API).
  - Dropdown pilihan Kategori (fetch dari API).
  - Input nominal (dengan format currency).
- [ ] Logic: Saat transaksi berhasil disimpan → Trigger refresh data Akun (karena saldo berubah) dan data Transaksi.

### 17.4 Verification

- [ ] Test: Tambah transaksi pengeluaran → Cek apakah saldo di Account Card berkurang secara real-time.
- [ ] Test: Filter transaksi → Pastikan tabel terupdate sesuai filter.

———

## 🟦 Day 18: Dashboard & Budget Frontend

Fokus: Visualisasi data dan kontrol budget.

### 18.1 Dashboard Summary UI

- [ ] KPI Cards: Tampilan besar untuk Total Balance, Monthly Income, Monthly Expense, dan Net Cash Flow.
- [ ] Cash Flow Chart: Implementasi Recharts (Line atau Bar chart) untuk tren pengeluaran vs pemasukan bulanan.

### 18.2 Category Breakdown Visualization

- [ ] Pie Chart: Visualisasi pengeluaran per kategori.
- [ ] Legend List: Daftar kategori dengan persentase dan nominal.

### 18.3 Budget Tracking UI

- [ ] Budget Progress Bars: List kategori budget dengan progress bar (Warna berubah jadi merah jika is_overlimit === true).
- [ ] Budget Management Form: Modal untuk set/update limit budget bulanan.

### 18.4 Verification

- [ ] Test: Ubah budget → Cek apakah progress bar berubah.
- [ ] Test: Buat transaksi → Cek apakah progress bar budget kategori tersebut bertambah.
- [ ] Test: Cek responsivitas UI di tampilan Mobile (menggunakan Tailwind responsive classes).

———

## ✅ Acceptance Criteria (End of Week 6)

- [ ] Full Integration: Semua fitur Backend (Auth, Account, Transaction, Dashboard, Budget) dapat diakses via UI.
- [ ] State Management: Data terupdate otomatis setelah mutasi (menggunakan TanStack Query invalidation).
- [ ] Security: Token JWT dikelola dengan aman dan interceptor menangani sesi expired.
- [ ] UX: Form memiliki validasi Zod dan user mendapatkan feedback (Toast) saat aksi berhasil/gagal.
- [ ] Responsive: Dashboard dan Tabel transaksi nyaman digunakan di layar desktop maupun mobile.