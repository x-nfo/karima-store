Subject: Backend Refactoring & Production-Readiness Task - Karima Store
Context: Kita telah melakukan audit awal pada codebase karima-store. Secara arsitektur sistem sudah baik, namun ada beberapa celah kritikal pada aspek stabilitas koneksi, integritas data transaksi, dan keamanan yang harus diperbaiki sebelum deployment ke lingkungan produksi.

Tasks & Technical Requirements:

1. Implementasi Database Transaction pada Checkout
Masalah: Proses pembuatan order, pemotongan stok, dan pencatatan log stok saat ini berjalan secara terpisah.

Tugas: Bungkus seluruh logika bisnis di dalam CheckoutService (terutama pada fungsi pembuatan order) menggunakan GORM Transaction.

Requirement: Pastikan jika salah satu proses (misal: update stok) gagal, maka data order tidak akan tersimpan ke database (rollback). Gunakan db.Transaction(func(tx *gorm.DB) error { ... }).

2. Implementasi Graceful Shutdown
Masalah: Server Fiber saat ini langsung berhenti tanpa menunggu proses background selesai.

Tugas: Tambahkan mekanisme penanganan sinyal OS (SIGTERM, SIGINT) di main.go.

Requirement: Gunakan app.Shutdown() untuk memastikan semua koneksi database dan Redis ditutup dengan rapi, serta memastikan request yang sedang berjalan diselesaikan terlebih dahulu sebelum proses mati.

3. Penguatan Keamanan JWT & Konfigurasi
Tugas: * Hapus hardcoded JWT_SECRET default di config.go. Tambahkan validasi agar aplikasi memberikan log.Fatal jika JWT_SECRET tidak diatur di .env.

Pastikan CORS_ORIGIN di .env produksi hanya berisi domain spesifik frontend, bukan asterisk *.

Ubah AppEnv menjadi production di environment untuk mengaktifkan mode Silent pada logger GORM guna menghindari kebocoran data sensitif di log.

4. Optimalisasi Koneksi Redis & DB
Tugas: * Lakukan pengecekan ulang di seluruh middleware. Pastikan tidak ada inisialisasi NewClient() di dalam fungsi middleware atau handler.

Gunakan shared instance yang sudah diinisialisasi di main.go melalui Dependency Injection.

5. Validasi & Dokumentasi API
Tugas: * Pastikan semua endpoint yang mengubah data (POST/PUT/DELETE) dilindungi oleh authMiddleware.Authenticate() dan memiliki pengecekan RequireRole jika diperlukan.

Update dokumentasi Swagger jika ada perubahan pada skema request/response.

Definition of Done (DoD):

Kode berhasil di-compile tanpa error.

Unit test untuk CheckoutService berhasil dijalankan dengan skenario rollback.

Server dapat dimatikan secara halus (Graceful Shutdown) saat menerima sinyal SIGINT/SIGTERM.

Kredensial sensitif tidak lagi ada yang bersifat hardcoded.