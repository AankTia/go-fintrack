# 🗓️ WEEK 8: Deployment & Launch (Detailed Breakdown)

Goal Utama: Aplikasi FinTrack live di server, database terkonfigurasi aman, dan dokumentasi lengkap sebagai bukti kualitas pengerjaan.

## 🟦 Day 22: Dockerize & CI/CD

Fokus: Memastikan aplikasi dapat dijalankan di lingkungan manapun secara konsisten.

### 22.1 Backend Dockerization

- [ ] Buat Dockerfile untuk Go (gunakan Multi-stage Build untuk memperkecil ukuran image: builder → alpine).
- [ ] Optimasi image: Pastikan hanya binary hasil build yang masuk ke image final.

### 22.2 Frontend Dockerization

- [ ] Buat Dockerfile untuk Next.js (gunakan output standalone mode untuk efisiensi memory).
- [ ] Konfigurasi docker-compose.prod.yml yang menggabungkan Backend, Frontend, dan PostgreSQL.

### 22.3 CI/CD Pipeline (GitHub Actions)

- [ ] Buat workflow .github/workflows/main.yml:
  - Lint & Test: Jalankan go test dan eslint setiap ada Push/PR.
  - Build: Build Docker image jika test pass.
  - Deploy (Optional): Auto-deploy ke VPS/Cloud via SSH atau Container Registry.

———

## 🟦 Day 23: Database Migration & Seeding

Fokus: Menyiapkan data awal agar aplikasi tidak terlihat kosong saat pertama kali dibuka oleh reviewer portofolio.

### 23.1 Production DB Setup

- [ ] Konfigurasi database di server (User, Password kuat, dan DB Name).
- [ ] Jalankan seluruh file migrasi dari folder migrations/ ke database produksi.

### 23.2 Demo Seed Data

- [ ] Buat script seed.sql atau tool Go sederhana untuk memasukkan data dummy:
  - 1 User Demo.
  - 3-5 Akun (Bank, Cash, E-Wallet).
  - 10-20 Transaksi (Campuran Income/Expense dalam 3 bulan terakhir).
  - 3 Budget kategori.
- [ ] Verifikasi: Login dengan user demo → Pastikan Dashboard langsung menampilkan grafik dan angka.

———

## 🟦 Day 24: API Documentation & README

Fokus: Mengubah kode menjadi aset portofolio yang bisa dipahami orang lain.

### 24.1 API Documentation

- [ ] Buat dokumentasi endpoint (Rekomendasi: Swagger/OpenAPI atau file API.md yang rapi).
- [ ] Detail per endpoint: Path, Method, Request Body (JSON), Response (Success & Error), dan contoh request.

### 24.2 Professional README.md

Tulis README yang menjual dengan struktur:

- [ ] Project Header: Nama, Logo, dan Deskripsi singkat (The "Why" behind FinTrack).
- [ ] Tech Stack: Daftar teknologi yang digunakan dan alasannya.
- [ ] Key Features: List fitur utama (Auth, Account, Transaction, Dashboard, Budget).
- [ ] Architecture: Jelaskan Clean Architecture (Handler → Service → Repository) yang diterapkan.
- [ ] How to Run: Instruksi docker compose up yang simpel.
- [ ] Learning Outcomes: Apa saja yang dipelajari (Sesuai Section "Learning Goals" di roadmap).

———

## 🟦 Day 25: Final Testing & Launch

Fokus: Quality check terakhir sebelum memberikan link demo kepada orang lain.

### 25.1 Regression Testing

- [ ] Jalankan satu putaran penuh fitur utama di lingkungan produksi:
  - Register → Login → Tambah Akun → Tambah Transaksi → Cek Dashboard.
- [ ] Verifikasi: Pastikan tidak ada hardcoded URL localhost:8080 yang tertinggal di frontend.

### 25.2 User Acceptance Testing (Internal)

- [ ] Coba aplikasi dari berbagai perangkat (Laptop, Tablet, HP).
- [ ] Verifikasi: Pastikan proses loading data terasa smooth dan tidak ada crash saat input data kosong.

### 25.3 GO LIVE! 🚀

- [ ] Deploy aplikasi ke hosting (misal: Railway, Render, DigitalOcean, atau VPS).
- [ ] Setup Domain/Subdomain jika tersedia.
- [ ] Publish link demo dan repository GitHub di LinkedIn/Portofolio.

———

## ✅ Acceptance Criteria (End of Week 8)

- [ ] Live App: Aplikasi dapat diakses via URL publik.
- [ ] Zero-Config Start: Orang lain bisa menjalankan project lokal hanya dengan docker compose up.
- [ ] Documented: API documented dan README.md terlihat profesional.
- [ ] Demo-Ready: Terdapat data dummy yang membuat dashboard terlihat menarik.
- [ ] Stable: Tidak ada bug kritikal pada alur utama (Happy Path).

Selamat! Jika semua checklist ini terpenuhi, Anda telah menyelesaikan build FinTrack dari nol hingga deployment menggunakan standar engineering yang disiplin.