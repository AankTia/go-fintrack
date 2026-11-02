# AGENTS.md — FinTrack

Instruksi kerja untuk AI assistant (Cloud maupun Local Model) yang membantu mengembangkan FinTrack

## 1. Konteks Proyek

FinTrack adalah SaaS pengelola keuangan pribadi/rumah tangga kecil. Produk membantu pengguna:

- mencatat pemasukan dan pengeluaran;
- melihat saldo dan cash flow;
- mengelompokkan transaksi berdasarkan kategori;
- mengatur budget bulanan;
- memahami kondisi keuangan melalui dashboard sederhana.

### Stack Utama:

- **Frontend:** Next.js App Router, TypeScript, Tailwind CSS, shadcn/ui, lucide-react.
- **Frontend Logic:** React Hook Form, Zod, TanStack Query, Recharts, date-fns.
- **Backend:** Go REST API.
- **Database:** PostgreSQL.
- **Infrastruktur:** Docker Compose.
- **Testing:** Playwright (E2E), Postman CLI (API).

### Dokumen Referensi (Wajib Baca):

- `fintrack-prd.md` — Kebutuhan produk dan acceptance criteria.
- `fintrack-roadmap.md` — Roadmap implementasi 8 minggu.
- `fintrack-task-breakdown.md` — Detail task per fase.

**PENTING:** Jika ada perbedaan antara dokumen referensi dan kode aktual, prioritaskan kode aktual. Jangan berasumsi contoh di dokumen adalah kondisi kode saat ini.

---

## 2. Prinsip Utama (Guardrails)

1. **MVP First:** Utamakan scope MVP. Jangan melakukan _over-engineering_.
2. **Minimalist Change:** Buat perubahan sekecil mungkin yang menyelesaikan task.
3. **No New Deps:** Jangan menambah library/dependency tanpa persetujuan user.
4. **Strict Isolation:** Data antar user harus terisolasi total (Strict Ownership).
5. **Backend as Truth:** Validasi bisnis utama wajib ada di Backend.
6. **Evidence-Based:** Jangan mengklaim test/command berhasil jika belum dijalankan.
7. **Ambiguity:** Jika requirement ambigu, jelaskan asumsi dan minta keputusan sebelum menulis kode.

---

## 3. Workflow Implementasi

### Tahap 1: Analisis (Sebelum Coding)

- Baca `AGENTS.md` dan dokumen referensi terkait.
- Cari file, route, handler, service, model, migration, dan test yang berkaitan.
- Pahami pola yang sudah ada (reuse existing patterns).
- Identifikasi Acceptance Criteria (AC).
- Buat rencana singkat (mapping perubahan layer).

### Tahap 2: Eksekusi (Coding)

- Kerjakan satu _vertical slice_ (satu fitur dari DB→ API→ UI).
- Ikuti konvensi penamaan dan struktur yang sudah ada.
- Tambahkan validasi di boundary input (DTO/Request) dan business layer.
- Update migration dan dokumentasi API jika ada perubahan kontrak.

### Tahap 3: Verifikasi (Setelah Coding)

- Jalankan formatter (`gofmt`, `prettier`, dll).
- Jalankan test yang paling relevan.
- Periksa diff: hapus debug code (`fmt.Println`, `console.log`), secret, atau dead code.
- Laporkan: file yang diubah, hasil verifikasi, dan hal yang belum teruji.

---

## 4. Aturan Scope & Prioritas

**Prioritas MVP:**

1. Auth & Session→ 2. Accounts & Categories→ 3. Transactions→ 4. Dashboard→ 5. Monthly Budgets→ 6. Security & UX.

**DILARANG (Kecuali diminta eksplisit):**

- Integrasi bank/open banking, AI financial advice, multi-currency kompleks, microservices, atau design system besar.

---

## 5. Arsitektur & Layer

### Backend Go (Strict Separation)

- `handler`: Parsing HTTP, response status, mapping error. (NO DB QUERY)
- `service`: Business logic, orchestration. (NO HTTP PARSING)
- `repository`: Query database, akses data. (NO BUSINESS LOGIC)
- `model/DTO`: Struktur data internal & kontrak request/response.
- `middleware`: Auth, CORS, logging, rate limit.
- `platform`: Config, DB connection, dependency injection.

### Frontend Next.js

- `route/page`: Komposisi halaman dan data fetching level atas.
- `feature module`: Logika bisnis spesifik domain.
- `component`: UI reusable (prioritaskan shadcn/ui).
- `API client`: Wrapper request ke backend.
- `schema/type`: Validasi Zod dan tipe TypeScript.
- `query/mutation layer`: TanStack Query hooks.

---

## 6. Standar API & Data Finansial

### Aturan API:

- Gunakan prefix versi (misal: `/api/v1`).
- Status code HTTP harus semantik (200, 201, 400, 401, 403, 404, 500).
- **Security:** Dilarang mengembalikan password/secret. Wajib filter resource berdasarkan `user_id` dari session (jangan percaya `user_id` dari body).
- Gunakan parameterized query untuk cegah SQL Injection.

### Aturan Data Finansial (Kritikal):

- **TIDAK BOLEH menggunakan Floating Point** untuk uang. Gunakan `integer` (unit terkecil) atau `NUMERIC`.
- Simpan currency secara eksplisit.
- Gunakan UTC untuk semua timestamp.
- Saldo & aggregate harus dihitung di server, bukan di frontend.
- Pastikan konsistensi tanda (+/-) untuk income dan expense.

---

## 7. Security & Database Checklist

### Security:

- [ ] Password di-hash (bcrypt/argon2).
- [ ] Session/Token dikelola dengan aman.
- [ ] Semua resource privat memiliki ownership filter.
- [ ] Input divalidasi (panjang, format, range, enum).
- [ ] Error response tidak membocorkan stack trace.

### Database:

- [ ] Setiap perubahan schema menggunakan migration (Up/Down).
- [ ] Index ditambahkan berdasarkan query aktual.
- [ ] Foreign key & Unique constraint untuk invariant penting.
- [ ] Tidak ada data pribadi/secret dalam seed data.

---

## 8. Strategi Testing

### Minimum Requirement:

- Pure Function→ Unit Test.
- Business Logic→ Service Unit Test.
- DB Query→ Integration Test.
- Endpoint→ API/Handler Test.
- Auth/Ownership→ Test Unauthorized & Cross-user access.

### Prioritas Test FinTrack:

- Flow Register/Login/Logout.
- Ownership resource (User A tidak bisa akses data User B).
- Akurasi saldo setelah CRUD transaksi.
- Dashboard aggregate & Budget limit status.

---

## 9. Konvensi Kode

### Go:

- Ikuti `gofmt`. Gunakan error wrapping.
- Gunakan `context` untuk database & request.
- Hindari `panic` di request normal.
- Interface hanya dibuat jika ada kebutuhan abstraction/testing yang nyata.

### TypeScript/Next.js:

- Hindari `any`. Pisahkan server/client component dengan jelas.
- Jangan simpan secret di env yang terekspos ke browser.
- Pastikan semua form memiliki label dan feedback error yang jelas.

---

## 10. Penanganan Task & Pelaporan

### Task Lintas Layer (Workflow):

1. Ringkasan Requirement→ 2. Mapping Layer→ 3. Migration/Contract→ 4. Backend & Test→ 5. Frontend & Test UI→ 6. E2E Verification.

### Format Laporan Akhir:

```text
Selesai:
- [Ringkasan perubahan singkat]
- [File utama yang diubah]

Verifikasi:
- [Command/test yang dijalankan]
- [Hasilnya]

Catatan:
- [Asumsi/Keputusan penting]
- [Hal yang belum diverifikasi]
```

———

## 11. Stop Condition (Berhenti & Tanya User)

Berhenti jika:

- Perubahan menghapus/memigrasi data production.
- Requirement bertentangan dengan PRD.
- Dibutuhkan secret/credential yang belum diberikan.
- Perubahan menyebabkan Breaking Change pada Public API.
- Scope melebar ke fitur non-MVP tanpa persetujuan.

———

## 12. Local AI Quick Check (Mandatory)

**Jika Anda adalah Local Model, Anda WAJIB melakukan checklist ini sebelum menjawab:**

1. Apakah saya menggunakan floating point untuk uang? → (Harus: **TIDAK**)
2. Apakah saya sudah menambahkan filter ownership `user_id` di query database? → (Harus: **YA**)
3. Apakah saya menambahkan library baru tanpa izin? → (Harus: **TIDAK**)
4. Apakah perubahan ini masuk dalam scope MVP? → (Harus: **YA**)
5. Apakah saya sudah menjalankan formatter/test? → (Harus: **YA**)

**Prinsip Penutup:** Bangun FinTrack sebagai produk kecil yang dapat dipercaya: alur inti sederhana, data finansial terlindungi, dan setiap keputusan teknis dapat dijelaskan.
