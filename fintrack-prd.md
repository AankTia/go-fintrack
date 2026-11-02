# Product Requirements Document (PRD)

## FinTrack — Personal Finance SaaS

**Status:** Draft MVP  
**Versi:** 1.0  
**Tanggal:** 27 Oktober 2026  
**Pemilik produk:** FinTrack Team  
**Target MVP:** 8 minggu, dengan waktu pengembangan 1–2 jam per hari

## 1. Ringkasan Eksekutif

FinTrack adalah aplikasi SaaS untuk membantu individu dan rumah tangga kecil mencatat transaksi, memahami cash flow, mengatur anggaran, dan mengambil keputusan finansial dengan lebih mudah.

Versi MVP berfokus pada pengalaman inti yang sederhana: pengguna membuat akun, menambahkan akun keuangan, mencatat pemasukan/pengeluaran, lalu melihat ringkasan kondisi keuangannya melalui dashboard.

FinTrack juga dirancang sebagai project portfolio end-to-end untuk menunjukkan kemampuan membangun produk nyata menggunakan Next.js, Go, dan PostgreSQL.

## 2. Latar Belakang dan Masalah

Banyak orang ingin mengatur keuangan, tetapi:

- Pencatatan transaksi sering dilakukan secara tidak konsisten.
- Aplikasi keuangan yang ada terasa terlalu kompleks.
- Pengguna sulit memahami kategori pengeluaran terbesar.
- Pengguna tidak memiliki gambaran cash flow bulanan yang jelas.
- Spreadsheet fleksibel, tetapi membutuhkan setup dan pemeliharaan manual.

### Problem statement

> Pengguna membutuhkan cara yang sederhana dan cepat untuk mencatat transaksi serta memahami ke mana uang mereka pergi tanpa harus menggunakan aplikasi akuntansi yang kompleks.

## 3. Tujuan Produk

### Tujuan utama

- Membantu pengguna mencatat transaksi dalam waktu singkat.
- Memberikan gambaran saldo dan cash flow yang mudah dipahami.
- Membantu pengguna mengontrol pengeluaran melalui budget sederhana.
- Membangun fondasi produk yang dapat dikembangkan menjadi SaaS berbayar.

### Tujuan pembelajaran dan portfolio

- Menerapkan Go untuk backend REST API yang terstruktur dan teruji.
- Membangun frontend Next.js yang responsive dan usable.
- Memahami autentikasi, authorization, validasi, migration, testing, dan deployment.
- Mendemonstrasikan pertimbangan security untuk data finansial.

### Bukan tujuan MVP

- Menjadi aplikasi akuntansi profesional.
- Mengintegrasikan rekening bank secara otomatis.
- Menyediakan nasihat investasi atau finansial.
- Menyimpan data kartu, credential bank, atau informasi pembayaran sensitif.
- Menyediakan rekonsiliasi rekening otomatis.

## 4. Target Pengguna

### Persona utama: pekerja muda atau freelancer

- Memiliki pemasukan dari satu atau beberapa sumber.
- Menggunakan rekening bank dan e-wallet.
- Ingin tahu pengeluaran terbesar setiap bulan.
- Tidak ingin mempelajari software akuntansi.
- Bersedia mencoba produk gratis yang mudah digunakan.

### Persona sekunder: rumah tangga kecil

- Membutuhkan gambaran pengeluaran bersama.
- Saat MVP, dapat menggunakan satu akun bersama secara manual.
- Fitur multi-user dan role akan dipertimbangkan setelah validasi awal.

## 5. Jobs to Be Done

- Ketika menerima pemasukan, saya ingin mencatatnya agar saldo saya akurat.
- Ketika melakukan pembelian, saya ingin mencatatnya dengan cepat agar tidak lupa.
- Ketika mengevaluasi bulan berjalan, saya ingin melihat kategori pengeluaran terbesar.
- Ketika membuat rencana bulanan, saya ingin menetapkan batas pengeluaran per kategori.
- Ketika akan mengambil keputusan finansial, saya ingin melihat kondisi cash flow terkini.

## 6. User Journey MVP

```text
Landing page
  → Register
  → Login
  → Onboarding singkat
  → Buat akun keuangan
  → Tambah kategori/transaksi pertama
  → Lihat dashboard
  → Tinjau budget dan histori transaksi
```

### Momen aktivasi

Pengguna dianggap aktif apabila telah:

1. Membuat akun.
2. Membuat minimal satu akun keuangan.
3. Mencatat transaksi pertama.

## 7. Ruang Lingkup MVP

### P0 — Wajib tersedia

#### Autentikasi

- Registrasi dengan email dan password.
- Login dan logout.
- Sesi pengguna yang aman.
- Endpoint profil pengguna.
- Pesan error yang jelas tanpa membocorkan informasi sensitif.

#### Akun keuangan

- Membuat akun: cash, bank, e-wallet, atau tipe custom.
- Melihat daftar akun dan saldo.
- Mengedit nama akun.
- Menonaktifkan atau menghapus akun yang belum memiliki transaksi.

#### Kategori

- Kategori default untuk pemasukan dan pengeluaran.
- Membuat kategori custom.
- Mengedit atau menonaktifkan kategori custom.

#### Transaksi

- Mencatat pemasukan atau pengeluaran.
- Memilih akun dan kategori.
- Mengisi nominal, tanggal, dan catatan opsional.
- Melihat daftar transaksi terbaru.
- Filter berdasarkan tanggal, tipe, akun, dan kategori.
- Mengedit dan menghapus transaksi.

#### Dashboard

- Total saldo.
- Total pemasukan pada periode terpilih.
- Total pengeluaran pada periode terpilih.
- Net cash flow.
- Daftar transaksi terbaru.
- Breakdown pengeluaran per kategori.

#### Budget

- Menetapkan budget bulanan per kategori pengeluaran.
- Melihat pemakaian budget.
- Menampilkan status aman, mendekati batas, atau melebihi batas.

### P1 — Setelah MVP awal

- Export transaksi ke CSV.
- Recurring transaction.
- Household sharing.
- Tema gelap.
- Notifikasi email.
- Target tabungan.

### P2 — Eksplorasi bisnis

- Import CSV dari bank atau e-wallet.
- Integrasi open banking.
- Paket premium dan billing.
- Laporan PDF.
- Multi-currency.
- Insight otomatis.

## 8. Persyaratan Fungsional

### FR-01 — Registrasi

**Deskripsi:** Pengguna dapat membuat akun baru.

**Acceptance criteria:**

- [ ] Email wajib diisi dan memiliki format valid.
- [ ] Email tidak boleh duplikat.
- [ ] Password memenuhi panjang minimum dan aturan keamanan.
- [ ] Password disimpan sebagai hash, bukan plaintext.
- [ ] Setelah berhasil, pengguna dapat login.
- [ ] Error validasi ditampilkan secara aman dan mudah dipahami.

### FR-02 — Login dan logout

**Acceptance criteria:**

- [ ] Pengguna dapat login dengan kredensial valid.
- [ ] Kredensial invalid menghasilkan respons generik.
- [ ] Sesi disimpan menggunakan mekanisme yang aman.
- [ ] Logout mencabut atau menghapus sesi aktif.
- [ ] Halaman yang membutuhkan autentikasi tidak dapat diakses tanpa sesi.

### FR-03 — Mengelola akun keuangan

**Acceptance criteria:**

- [ ] Pengguna dapat membuat akun dengan nama, tipe, dan mata uang.
- [ ] Saldo ditampilkan berdasarkan transaksi yang dimiliki pengguna.
- [ ] User hanya dapat melihat dan mengubah akun miliknya.
- [ ] Akun user lain tidak dapat diakses melalui manipulasi ID.
- [ ] Akun dengan transaksi tidak dihapus secara sembarangan.

### FR-04 — Mengelola transaksi

**Acceptance criteria:**

- [ ] Pengguna dapat memilih tipe income atau expense.
- [ ] Nominal harus lebih besar dari nol.
- [ ] Akun dan kategori harus valid serta dimiliki user.
- [ ] Tanggal transaksi wajib valid.
- [ ] Transaksi tersimpan dan muncul di daftar.
- [ ] Perubahan transaksi memperbarui ringkasan saldo.
- [ ] User tidak dapat membaca, mengedit, atau menghapus transaksi user lain.

### FR-05 — Dashboard

**Acceptance criteria:**

- [ ] Dashboard menampilkan periode yang jelas.
- [ ] Saldo, pemasukan, pengeluaran, dan net cash flow konsisten dengan data transaksi.
- [ ] Data kosong memiliki empty state yang membantu.
- [ ] Chart memiliki label yang dapat dipahami.
- [ ] Dashboard dapat digunakan pada layar mobile.

### FR-06 — Budget

**Acceptance criteria:**

- [ ] Pengguna dapat membuat atau mengubah budget per kategori dan bulan.
- [ ] Satu user tidak dapat memiliki budget duplikat untuk kategori dan bulan yang sama.
- [ ] Penggunaan budget dihitung dari transaksi expense pada periode tersebut.
- [ ] Status budget berubah sesuai penggunaan aktual.

## 9. Persyaratan Non-Fungsional

### Security

- Password menggunakan Argon2id atau bcrypt.
- Cookie/session menggunakan konfigurasi aman.
- Semua endpoint privat memiliki authorization check.
- Validasi dilakukan di frontend dan backend, dengan backend sebagai sumber kebenaran.
- Rate limit diterapkan pada login dan endpoint sensitif.
- Secret tidak disimpan dalam repository.
- Log tidak mencatat password, token, atau data finansial yang tidak diperlukan.

### Performance

- Endpoint dashboard merespons cepat untuk dataset MVP.
- Query transaksi memakai index pada `user_id`, tanggal, akun, dan kategori sesuai kebutuhan.
- Daftar transaksi menggunakan pagination atau limit.
- Frontend tidak mengambil data lebih banyak dari yang dibutuhkan.

### Reliability

- Migration database dapat dijalankan berulang dengan aman.
- Error database ditangani tanpa mematikan seluruh server.
- Health check tersedia.
- Backup dan restore database didokumentasikan sebelum production.

### Accessibility dan usability

- Form memiliki label yang jelas.
- Navigasi keyboard dasar berfungsi.
- Warna tidak menjadi satu-satunya pembeda status.
- Nominal dan tanggal memiliki format konsisten.
- Loading, error, dan empty state tersedia.

## 10. Model Data Tingkat Tinggi

```text
User 1 ─── * Account
User 1 ─── * Category
User 1 ─── * Transaction
Account 1 ─── * Transaction
Category 1 ─── * Transaction
User 1 ─── * Budget
Category 1 ─── * Budget
```

### Entitas inti

| Entitas | Field penting | Catatan |
|---|---|---|
| User | id, email, password_hash, created_at | Email unik |
| Account | id, user_id, name, type, currency | Milik satu user |
| Category | id, user_id, name, kind | `income` atau `expense` |
| Transaction | id, user_id, account_id, category_id, amount, kind, occurred_at | Nominal integer/unit terkecil |
| Budget | id, user_id, category_id, month, limit_amount | Kombinasi user-category-month unik |
| Session | id, user_id, token_hash, expires_at | Dapat dicabut |

### Aturan finansial

- Jangan gunakan floating point untuk nominal uang.
- Gunakan integer unit terkecil untuk MVP IDR.
- Simpan mata uang secara eksplisit.
- Gunakan UTC untuk penyimpanan waktu.
- Semua kalkulasi agregat harus berasal dari data server.
- Jangan menyimpan credential bank atau nomor kartu.

## 11. Desain Teknis Tingkat Tinggi

```text
Next.js Web
  ├── Halaman auth
  ├── Dashboard
  ├── Transaksi
  ├── Akun
  └── Budget
          │ HTTPS/JSON
          ▼
Go REST API
  ├── HTTP handlers
  ├── Middleware auth/CORS/rate limit
  ├── Application services
  ├── Repositories
  └── PostgreSQL
```

### Keputusan Library Frontend

- Next.js App Router dan TypeScript sebagai fondasi aplikasi.
- Tailwind CSS untuk styling responsive.
- shadcn/ui untuk komponen UI yang source code-nya berada di repository dan dapat dikustomisasi.
- lucide-react untuk icon.
- React Hook Form dan Zod dengan `@hookform/resolvers` untuk form serta validasi.
- TanStack Query untuk fetching, caching, mutation, dan invalidation data dari Go API.
- Recharts untuk chart cash flow dan breakdown kategori.
- date-fns untuk format tanggal dan perhitungan periode.
- Playwright untuk smoke test dan end-to-end flow.

Untuk MVP, belum perlu menambahkan Redux atau Zustand. Gunakan TanStack Query untuk server state dan React state untuk state lokal komponen.

### Kontrak API minimum

```text
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/logout
GET    /api/v1/me
GET    /api/v1/dashboard/summary
GET    /api/v1/accounts
POST   /api/v1/accounts
PATCH  /api/v1/accounts/:id
GET    /api/v1/categories
POST   /api/v1/categories
GET    /api/v1/transactions
POST   /api/v1/transactions
PATCH  /api/v1/transactions/:id
DELETE /api/v1/transactions/:id
GET    /api/v1/budgets?month=YYYY-MM
PUT    /api/v1/budgets/:categoryID
```

## 12. UX dan Halaman

### Halaman publik

- Landing page.
- Login.
- Register.
- Privacy/security note sederhana.

### Halaman privat

- Dashboard.
- Transactions.
- Accounts.
- Categories.
- Budgets.
- Settings/profile.

### Prinsip UX

- Tambah transaksi harus dapat dilakukan dalam beberapa langkah.
- Default tanggal adalah hari ini.
- Tipe transaksi dibedakan dengan label dan ikon, bukan warna saja.
- Nominal selalu memiliki format mata uang yang jelas.
- Setelah menyimpan transaksi, pengguna mendapat konfirmasi dan melihat data terbaru.
- Empty state mengarahkan pengguna ke aksi berikutnya.

## 13. Analitik dan Metrik Keberhasilan

### Metrik produk

- **Activation rate:** persentase user yang membuat akun, akun keuangan, dan transaksi pertama.
- **Day-7 retention:** persentase user yang kembali menggunakan produk pada hari ke-7.
- **Weekly active users:** jumlah user aktif mingguan.
- **Transactions per active user:** rata-rata jumlah transaksi per user aktif.
- **Budget adoption:** persentase user aktif yang membuat budget.
- **Task completion:** waktu dan keberhasilan mencatat transaksi.

### Target validasi awal

- Minimal 5 pengguna mencoba MVP.
- Minimal 3 pengguna menyelesaikan alur aktivasi tanpa bantuan langsung.
- Minimal 2 pengguna kembali menggunakan aplikasi dalam 7 hari.
- Tidak ada temuan kritis terkait akses data lintas pengguna.
- Pengguna dapat mencatat transaksi utama dalam waktu kurang dari satu menit setelah memahami UI.

## 14. Rencana Delivery 8 Minggu

| Minggu | Fokus | Deliverable |
|---|---|---|
| 1 | Discovery dan setup | Product brief, repo, frontend/backend skeleton |
| 2 | Database dan backend foundation | Migration, connection, repository, error format |
| 3 | Auth | Register, login, logout, protected route |
| 4 | Core transaction flow | Accounts, categories, transaction CRUD |
| 5 | Insight | Dashboard, chart, budget |
| 6 | Quality | UX states, responsive layout, tests, security review |
| 7 | Release | CI, deployment, API docs, seed data |
| 8 | Validation | User testing, feedback, portfolio case study |

## 15. Risiko dan Mitigasi

| Risiko | Dampak | Mitigasi |
|---|---|---|
| Scope terlalu besar | MVP terlambat | Prioritaskan P0 dan batasi integrasi eksternal |
| Perhitungan saldo tidak konsisten | Hilangnya kepercayaan user | Gunakan server-side aggregation dan test invariants |
| Kebocoran data antar user | Risiko keamanan tinggi | Authorization check pada setiap query resource |
| Data finansial sensitif | Risiko privasi | Minimalkan data, enkripsi transport, secret management |
| Waktu belajar terbatas | Progress tidak konsisten | Task harian kecil dan definition of done jelas |
| Deployment memakan waktu | Demo tertunda | Deploy skeleton lebih awal pada minggu 2–3 |
| Tidak ada validasi kebutuhan | Produk tidak dipakai | Uji dengan pengguna sebelum menambah fitur P1/P2 |

## 16. Definition of Done MVP

MVP dianggap selesai apabila:

- [ ] User baru dapat register dan login.
- [ ] User dapat membuat akun keuangan.
- [ ] User dapat mencatat, mengubah, memfilter, dan menghapus transaksi.
- [ ] Dashboard menghitung saldo dan cash flow dengan benar.
- [ ] User dapat membuat budget bulanan.
- [ ] Authorization test membuktikan data antar user terisolasi.
- [ ] Frontend responsive dan memiliki loading/error/empty state.
- [ ] Backend memiliki unit/integration test yang relevan.
- [ ] Migration dan seed data tersedia.
- [ ] Aplikasi ter-deploy dan dapat dicoba melalui URL demo.
- [ ] README dan dokumentasi API dapat diikuti oleh developer lain.
- [ ] Tidak ada secret atau data pribadi di repository.

## 17. Langkah Berikutnya

1. Review PRD dan tandai fitur yang benar-benar diperlukan untuk versi pertama.
2. Buat `docs/product.md` berdasarkan bagian problem, persona, dan user journey.
3. Ubah FR menjadi issue/task kecil di repository.
4. Mulai dari product brief dan setup repository pada roadmap minggu pertama.
5. Sebelum membangun fitur tambahan, validasikan alur register → transaksi pertama → dashboard.

## Kesimpulan

FinTrack MVP harus membuktikan satu hal utama: pengguna dapat mencatat keuangannya dengan mudah dan memperoleh pemahaman yang lebih baik dari data tersebut. Setiap fitur di luar alur inti perlu menunggu sampai ada bukti kebutuhan dari pengguna, sehingga project tetap realistis untuk diselesaikan sekaligus cukup kuat sebagai portfolio dan fondasi bisnis SaaS.
