# Day 6: Category Management (CRUD)

### 🟦 Step 1: Implementasi Konkret Category Repository

File: `internal/repository/category_repository.go`

Saat ini file ini hanya berisi interface. Anda perlu membuat struct implementasi yang berinteraksi langsung dengan database menggunakan pgxpool.

Apa yang harus dilakukan:

1. Buat struct `categoryRepository` (private/lowercase) yang menampung *pgxpool.Pool.
2. Buat fungsi constructor `NewCategoryRepository` (public) yang mengembalikan interface `CategoryRepository`.
3. Implementasikan method `Create`, `GetByUserId`, dan `Delete`.

Detail Implementasi Kode:

- Create: Gunakan `INSERT INTO categories`. Jangan lupa gunakan `QueryRow` dan `Scan` untuk mendapatkan `ID` yang baru dibuat.
- GetByUserId: Gunakan query `SELECT ... FROM categories WHERE user_id = $1 OR user_id IS NULL`.
  - Penjelasan: `user_id IS NULL` digunakan agar kategori Global (bawaan sistem) juga muncul untuk semua user.
- Delete: Gunakan `DELETE FROM categories WHERE id = $1 AND user_id = $2`.
  - Penjelasan: Menambahkan `user_id = $2` sangat penting agar user tidak bisa menghapus kategori milik orang lain atau menghapus kategori global.

———

### 🟦 Step 2: Implementasi Category Service

File Baru: `internal/service/category.go`

Service layer adalah tempat "otak" aplikasi berada. Di sini kita melakukan validasi bisnis sebelum data dikirim ke repository.

Apa yang harus dilakukan:

1. Buat `struct CategoryService` yang memiliki dependency `repository.CategoryRepository`.
2. Buat `constructor NewCategoryService`.
3. Implementasikan method:
   - CreateCategory:
     - Validasi: Pastikan Name tidak kosong.
     - Mapping: Ubah model.CategoryRequest (DTO) menjadi model.Category (Entity).
     - Panggil repo.Create.
   - GetMyCategories:
     - Panggil repo.GetByUserId.
     - Di sini Anda bisa menambahkan logika pengurutan (sorting) jika diperlukan.
   - DeleteCategory:
     - Panggil repo.Delete.
     - Jika repository mengembalikan error (data tidak ditemukan atau user tidak punya akses), bungkus menjadi error yang bisa dimengerti oleh Handler (misal: ErrCategoryNotFound).

Penjelasan untuk Pembelajar Go:

- Dependency Injection: Perhatikan bagaimana CategoryService tidak membuat repository-nya sendiri, tetapi "menerimanya" lewat constructor. Ini memudahkan kita jika nanti ingin membuat Unit Test menggunakan Mock Repository.

———

### 🟦 Step 3: Implementasi Category Handler

File Baru: `internal/handler/category.go`

Handler bertugas mengurusi protokol HTTP: membaca request body, mengambil parameter URL, dan menentukan status code response.

Apa yang harus dilakukan:

1. Buat struct CategoryHandler yang memiliki dependency *service.CategoryService.
2. Buat constructor NewCategoryHandler.
3. Implementasikan method:
   - Create (POST /api/v1/categories):
     - Ambil user_id dari context.
     - Decode JSON body ke model.CategoryRequest.
     - Panggil svc.CreateCategory.
     - Return 201 Created.
   - List (GET /api/v1/categories):
     - Ambil user_id dari context.
     - Panggil svc.GetMyCategories.
     - Return 200 OK dengan list kategori.
   - Delete (DELETE /api/v1/categories/{id}):
     - Ambil user_id dari context dan id dari chi.URLParam.
     - Panggil svc.DeleteCategory.
     - Return 204 No Content jika berhasil.

Penjelasan untuk Pembelajar Go:

- Context: Penggunaan `r.Context().Value("user_id")` adalah cara standar di Go untuk meneruskan data (seperti identitas user yang sudah divalidasi oleh middleware) antar layer handler.

———

### 🟦 Step 4: Pendaftaran di Main (Wiring)

File: `cmd/api/main.go` (atau tempat Anda menginisialisasi dependency)

Setelah membuat tiga layer di atas, Anda harus menghubungkannya agar aplikasi tahu cara menggunakannya.

Urutan Wiring:

1. `repo := repository.NewCategoryRepository(db)`
2. `svc := service.NewCategoryService(repo)`
3. `handler := handler.NewCategoryHandler(svc)`
4. Daftarkan route di chi router:
   - `r.Post("/categories", handler.Create)`
   - `r.Get("/categories", handler.List)`
   - `r.Delete("/categories/{id}", handler.Delete)`

———

### 🛠️ Checklist Verifikasi (Setelah Apply)

Setelah Anda mengimplementasikan kode di atas, coba lakukan test manual berikut untuk memastikan tidak ada bug:

1. Create Category: Coba buat kategori baru → Cek di DB apakah user_id terisi dengan benar.
2. List Category: Panggil endpoint GET → Pastikan kategori yang Anda buat DAN kategori global (yang user_id-nya NULL) muncul semua.
3. Security Test (Penting):
   - Coba hapus kategori milik User A menggunakan token User B → Harus Gagal (403/404).
   - Coba hapus kategori Global → Harus Gagal (karena kategori global tidak punya user_id yang cocok dengan user).
4. No Floats: Pastikan tidak ada satu pun variabel float32 atau float64 yang terselip di logika kategori ini.