# AGENTS.md — FinTrack

Instruksi kerja untuk AI assistant yang membantu mengembangkan FinTrack.

## 1. Konteks Proyek

FinTrack adalah SaaS pengelola keuangan pribadi/rumah tangga kecil. Produk membantu pengguna:

- mencatat pemasukan dan pengeluaran;
- melihat saldo dan cash flow;
- mengelompokkan transaksi berdasarkan kategori;
- mengatur budget bulanan;
- memahami kondisi keuangan melalui dashboard sederhana.

Stack utama:

- Frontend: Next.js App Router dan TypeScript.
- UI/styling: Tailwind CSS dan shadcn/ui.
- Icon: lucide-react.
- Form/validation: React Hook Form, Zod, dan `@hookform/resolvers`.
- API state: TanStack Query.
- Chart: Recharts.
- Date utility: date-fns.
- E2E test: Playwright.
- Backend: Go REST API.
- Database: PostgreSQL.
- Local development: Docker Compose.
- API testing: gunakan Postman CLI bila collection atau workflow API tersedia.

Dokumen produk yang menjadi referensi:

- `fintrack-prd.md` — kebutuhan produk dan acceptance criteria.
- `fintrack-roadmap.md` — roadmap belajar dan implementasi 8 minggu.
- `fintrack-task-breakdown.md` — task breakdown per phase.

Jika dokumen tersebut berbeda dengan kode aktual, periksa implementasi dan konfigurasi repository terlebih dahulu. Jangan menganggap contoh di dokumen sudah pasti sama dengan kondisi kode.

Library frontend tersebut merupakan keputusan MVP. Jangan mengganti library tanpa alasan teknis yang jelas atau persetujuan user. Gunakan komponen shadcn/ui yang sudah ada sebelum menambahkan UI library lain.

## 2. Prinsip Utama

1. Utamakan MVP dan scope yang diminta.
2. Pahami kode yang sudah ada sebelum mengubahnya.
3. Buat perubahan sekecil mungkin yang menyelesaikan task.
4. Jangan menambahkan dependency atau abstraksi tanpa alasan yang jelas.
5. Pertahankan behavior yang sudah berjalan kecuali perubahan memang diminta.
6. Data user harus terisolasi secara ketat.
7. Backend adalah sumber kebenaran untuk validasi dan business rules.
8. Setiap perubahan penting harus memiliki verifikasi yang sesuai.
9. Jangan mengklaim test atau command berhasil jika belum dijalankan.
10. Jika requirement ambigu dan pilihan berdampak besar pada arsitektur atau data, jelaskan asumsi dan minta keputusan sebelum implementasi.

## 3. Workflow AI Assistant

### Sebelum mengubah kode

- Baca `AGENTS.md` yang relevan di root dan subdirectory.
- Periksa status perubahan repository.
- Cari file, route, handler, service, model, migration, dan test yang berkaitan.
- Pahami dependency dan pola yang sudah dipakai.
- Identifikasi acceptance criteria dari task.
- Buat rencana singkat jika task menyentuh beberapa layer.

### Saat mengimplementasikan

- Kerjakan satu vertical slice atau unit perubahan yang jelas.
- Ikuti konvensi yang sudah ada.
- Reuse helper, komponen, query, dan pola existing jika sesuai.
- Tambahkan validasi di boundary input dan business layer yang tepat.
- Update migration, API documentation, atau test bila perubahan membutuhkannya.
- Hindari perubahan yang tidak berkaitan dengan task.

### Setelah mengimplementasikan

- Jalankan formatter dan test yang paling relevan.
- Periksa diff untuk perubahan tidak sengaja, secret, debug code, dan dead code.
- Laporkan file yang diubah, verifikasi yang dijalankan, dan hal yang belum diverifikasi.
- Jika ada blocker, jelaskan bukti dan alternatif yang sudah dicoba.

## 4. Aturan Scope dan Prioritas

Prioritas fitur MVP:

1. Authentication dan session.
2. Accounts dan categories.
3. Transactions.
4. Dashboard summary.
5. Monthly budgets.
6. Security, testing, UX, dan deployment.

Jangan membangun fitur berikut tanpa permintaan eksplisit atau bukti validasi user:

- integrasi bank/open banking;
- AI financial advice;
- subscription billing;
- multi-currency kompleks;
- household roles kompleks;
- mobile native app;
- microservices atau event-driven architecture;
- design system besar;
- analytics yang tidak dibutuhkan MVP.

Jika menemukan peluang improvement di luar scope, catat sebagai backlog dan jangan mengubah kode tanpa alasan task yang jelas.

## 5. Struktur Layer yang Diutamakan

### Backend Go

Pertahankan pemisahan tanggung jawab:

- `handler`: HTTP parsing, response, status code, dan mapping error.
- `service` atau use case: business rules dan orchestration.
- `repository`: query dan akses database.
- `model`/DTO: bentuk data internal dan request/response.
- `middleware`: authentication, authorization, CORS, rate limit, request ID.
- `platform`: config, database, logger, dan dependency wiring.

Jangan menaruh query database di handler. Jangan menaruh parsing HTTP di repository.

### Frontend Next.js

Gunakan pemisahan yang konsisten:

- route/page: composition dan page-level data needs;
- feature module: logic khusus domain;
- component: UI reusable;
- API client: request ke backend;
- schema/type: kontrak dan validasi;
- query/mutation layer: fetching, caching, dan invalidation.

Jangan menduplikasi bentuk response API di banyak tempat. Jika tipe response berubah, cari seluruh consumer-nya.

## 6. Aturan API

- Gunakan prefix versi, misalnya `/api/v1` jika sudah menjadi konvensi repository.
- Gunakan status code HTTP yang semantik.
- Pertahankan format error JSON yang konsisten.
- Validasi input di backend meskipun frontend juga melakukan validasi.
- Jangan mengembalikan password, token mentah, secret, atau data internal yang tidak diperlukan.
- Semua resource privat wajib difilter berdasarkan user yang sedang terautentikasi.
- Jangan percaya `user_id`, ownership, atau role dari request body.
- Gunakan parameterized query.
- Pertahankan backward compatibility bila endpoint sudah digunakan frontend.
- Jika kontrak API berubah, update frontend, test, dan dokumentasi dalam perubahan yang sama.

## 7. Aturan Data Finansial

- Jangan gunakan floating point untuk nominal uang.
- Gunakan integer unit terkecil atau tipe `NUMERIC` sesuai keputusan schema yang sudah ada.
- Simpan currency secara eksplisit.
- Gunakan UTC untuk timestamp penyimpanan.
- Pastikan income dan expense memiliki aturan tanda yang konsisten.
- Saldo dan dashboard summary harus dihitung dari sumber data server.
- Edit dan delete transaksi harus memperbarui aggregate secara benar.
- Tambahkan constraint database untuk aturan yang penting.
- Gunakan migration versioning; jangan mengubah schema production secara manual.
- Jangan menyimpan nomor kartu, PIN, password bank, atau credential pihak ketiga.
- Jangan memasukkan data finansial mentah ke log.
- Jangan memakai data user nyata dalam fixture, screenshot, test, atau demo.

## 8. Security Checklist

Sebelum menyelesaikan perubahan yang menyentuh auth atau data user, pastikan:

- [ ] Password di-hash dengan algoritma yang sesuai.
- [ ] Session/token disimpan dan dicabut dengan aman.
- [ ] Cookie memakai flag aman sesuai environment.
- [ ] CORS tidak terlalu longgar.
- [ ] Login memiliki rate limit atau perlindungan brute-force yang sesuai.
- [ ] Semua input memiliki validasi panjang, format, range, dan enum.
- [ ] Semua query resource memiliki ownership filter.
- [ ] Error tidak membocorkan stack trace atau detail database.
- [ ] Secret hanya berasal dari environment/deployment secret.
- [ ] Log tidak mengandung password, token, atau data finansial sensitif.
- [ ] Test mencoba akses user A terhadap resource user B.

## 9. Database dan Migration

- Setiap perubahan schema harus memiliki migration yang dapat direview.
- Migration harus memiliki arah up dan down bila tooling mendukung.
- Jangan menghapus kolom atau data production tanpa rencana compatibility/rollback.
- Tambahkan index berdasarkan query aktual, bukan spekulasi.
- Gunakan foreign key dan unique constraint untuk invariant penting.
- Pertimbangkan dampak migration terhadap data yang sudah ada.
- Seed data demo tidak boleh menggunakan secret atau data pribadi.
- Uji migration pada database kosong dan database dengan data contoh.

## 10. Testing

### Minimum berdasarkan perubahan

- Perubahan pure function: unit test.
- Perubahan service: unit test untuk business rule.
- Perubahan repository: integration test atau test query yang sesuai.
- Perubahan endpoint: handler/API test.
- Perubahan auth atau ownership: test unauthorized dan cross-user access.
- Perubahan UI form: validation, success, error, dan loading state.
- Perubahan schema: migration up/down dan compatibility check.

### Prioritas test FinTrack

- register/login/logout;
- protected route;
- ownership setiap resource;
- nominal invalid dan tipe transaksi invalid;
- saldo setelah create/update/delete transaksi;
- dashboard aggregate;
- budget period dan status limit;
- data kosong dan error state frontend.

Jalankan command yang didefinisikan repository terlebih dahulu. Jangan mengganti toolchain tanpa kebutuhan. Jika tersedia, gunakan Postman CLI untuk menjalankan collection API setelah server dan environment test siap.

## 11. Konvensi Kode

### Go

- Ikuti `gofmt`.
- Gunakan error wrapping yang konsisten dengan pola repository.
- Hindari global mutable state.
- Gunakan context untuk request dan operasi database.
- Tutup resource dan hormati cancellation.
- Hindari panic pada request normal.
- Beri nama fungsi dan error yang jelas.
- Jangan membuat interface hanya untuk menambah abstraksi; buat ketika ada boundary atau kebutuhan testing yang nyata.

### TypeScript/Next.js

- Ikuti formatter dan linter repository.
- Hindari `any` jika tipe yang lebih aman memungkinkan.
- Pisahkan server/client boundary dengan jelas.
- Jangan menyimpan secret di environment variable yang terekspos ke browser.
- Hindari fetch berulang tanpa caching atau invalidation yang jelas.
- Pastikan komponen form memiliki label dan feedback error.

### Naming

- Gunakan istilah domain yang konsisten: `account`, `category`, `transaction`, `budget`.
- Bedakan `income` dan `expense` secara konsisten di API, database, dan UI.
- Jangan mengganti nama field API hanya demi style tanpa migration/compatibility plan.

## 12. Git dan Perubahan File

- Jangan menghapus atau menimpa perubahan user yang tidak berkaitan.
- Jangan melakukan reset, checkout destruktif, atau force push tanpa permintaan eksplisit.
- Hindari commit file environment, credential, build output, dan data lokal.
- Buat commit kecil berdasarkan unit perubahan jika diminta membuat commit.
- Gunakan pesan commit yang menjelaskan intent perubahan.
- Jangan menggabungkan refactor besar dengan feature kecil kecuali diperlukan.

## 13. Dokumentasi

Update dokumentasi jika perubahan memengaruhi:

- cara menjalankan project;
- environment variable;
- endpoint atau payload API;
- schema/migration;
- aturan security;
- alur user;
- command test atau deployment.

Dokumentasi harus menjelaskan hal yang dapat dilakukan sekarang, bukan hanya rencana masa depan.

## 14. Cara Menangani Task dari User

### Task kecil

1. Cari simbol/file terkait.
2. Implementasikan perubahan minimal.
3. Jalankan test relevan.
4. Laporkan hasil dan file yang berubah.

### Task lintas layer

1. Tulis ringkasan requirement dan acceptance criteria.
2. Petakan perubahan database, backend, frontend, test, dan docs.
3. Kerjakan migration/contract terlebih dahulu bila diperlukan.
4. Implementasikan backend dan test.
5. Hubungkan frontend dan test UI.
6. Jalankan verification end-to-end.

### Task ambigu

- Nyatakan asumsi yang paling aman.
- Pilih solusi paling sederhana yang memenuhi tujuan.
- Jangan memperluas scope berdasarkan asumsi.
- Minta klarifikasi jika pilihan tersebut dapat mengubah schema, public API, security, atau biaya operasional.

### Task diagnosis

- Jangan langsung mengubah kode.
- Reproduksi atau kumpulkan bukti terlebih dahulu.
- Identifikasi layer penyebab.
- Jelaskan hipotesis dan bukti.
- Implementasikan fix hanya jika user meminta diagnosis sekaligus perbaikan.

## 15. Definition of Done Umum

Sebuah task dianggap selesai apabila:

- [ ] Acceptance criteria terpenuhi.
- [ ] Perubahan berada di layer yang tepat.
- [ ] Validasi dan authorization sudah dipertimbangkan.
- [ ] Test relevan ditambahkan atau diperbarui.
- [ ] Formatter/linter dijalankan.
- [ ] Tidak ada secret, debug code, atau perubahan tidak sengaja.
- [ ] Dokumentasi diperbarui jika kontrak atau workflow berubah.
- [ ] Hasil verifikasi dilaporkan dengan jujur.

## 16. Format Laporan AI Assistant

Gunakan format ringkas berikut setelah menyelesaikan task:

```text
Selesai:
- Ringkasan perubahan
- File utama yang diubah

Verifikasi:
- Command/test yang dijalankan
- Hasilnya

Catatan:
- Asumsi atau keputusan penting
- Hal yang belum diverifikasi
- Langkah lanjutan opsional
```

## 17. Stop Condition

Berhenti dan minta arahan user jika:

- perubahan dapat menghapus atau memigrasikan data production;
- requirement bertentangan dengan PRD dan tidak ada prioritas yang jelas;
- diperlukan secret, akses akun, atau credential yang belum diberikan;
- perubahan akan mengubah public API secara breaking;
- test gagal karena masalah environment yang tidak dapat dibuktikan dari kode;
- solusi yang aman membutuhkan keputusan bisnis atau privacy policy;
- scope mulai melebar ke fitur P1/P2 tanpa persetujuan.

## Prinsip Penutup

Bangun FinTrack sebagai produk kecil yang dapat dipercaya: alur inti sederhana, data finansial terlindungi, hasil dapat diuji, dan setiap keputusan teknis dapat dijelaskan.
