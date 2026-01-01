## 0. General Rules for AI
- Semua perintah terminal untuk container harus menggunakan perintah **Podman**.
- Jika membuat `Dockerfile`, pastikan tidak menggunakan instruksi yang hanya spesifik untuk Docker daemon (seperti mounting docker.sock tanpa penyesuaian).
- Gunakan format `docker-compose.yml` versi 3.8 yang kompatibel dengan Podman.
- Masukan Setiap Pekerjaan dalam Folder Progress_flows dengan subfolder bernama sesuai modul, misal: `Progress_flows/MODUL_1_Fondasi_Infrastruktur_Database`.

## Modul 1: Fondasi Infrastruktur & Database (Tahap 1)
- Tujuan: Menyiapkan server, database, dan sistem migrasi agar tabel terbentuk otomatis.
- Prompt untuk AI Coding: "Saya ingin membangun backend toko online fashion menggunakan Golang Fiber dan PostgreSQL. Ikuti struktur folder standar Go (cmd/internal/pkg).
- Gunakan GORM untuk ORM dan golang-migrate untuk migrasi database.
- Implementasikan skema database berikut dalam file migrasi SQL: [docs/blueprint_db.dbml].
- Buat file docker-compose.yml yang berisi PostgreSQL dan Redis.
- Pastikan ada file .env untuk konfigurasi koneksi DB.
- Buat endpoint GET /health untuk memastikan koneksi DB berhasil."

## Modul 2: Autentikasi & Integrasi Ory Kratos (Tahap 2)
- Tujuan: Mengamankan akses user dan admin.
- Prompt untuk AI Coding: "Hubungkan backend Golang saya dengan Ory Kratos sebagai sistem autentikasi.
- Buat middleware di Go untuk memvalidasi sesi dari Ory Kratos.
- Setiap kali user berhasil register di Kratos, simpan identity_id mereka ke tabel users di database lokal saya.
- Bedakan hak akses "RBAC" antara Customer dan Admin "Role-based access".
- Pastikan session cookie bekerja secara cross-domain antara karima.com dan ks-backend.cloud."

## Modul 3: Manajemen Katalog & Media (Tahap 3)
- Tujuan: Admin bisa input baju, warna, ukuran, dan foto.
-Prompt untuk AI Coding: "Buat API CRUD untuk manajemen produk.
- Implementasikan endpoint untuk Produk, SKU "Warna & Size", dan Kategori.
- Integrasikan upload gambar ke Cloudflare R2 atau AWS S3. Simpan URL-nya di tabel product_media.
- Buat logic Slug Generator otomatis saat nama produk diinput.
- Pastikan admin bisa mengupdate stok secara manual melalui API."

## Modul 4: Pricing Engine & Shipping (Tahap 4 & 5)
- Tujuan: Menghitung harga (Reseller/Flash Sale) dan ongkir (RajaOngkir).
Prompt untuk AI Coding: "Bangun Pricing Engine dan Integrasi Shipping di Golang.
Buat fungsi CalculatePrice yang mengecek: Harga Retail, Harga Tiering (Reseller), dan Flash Sale aktif.
Integrasikan API RajaOngkir/Komerce untuk fitur cek ongkir.
Gunakan subdistrict_id untuk akurasi biaya kirim.
Buat logic untuk menghitung total berat belanjaan (Total Weight = SKU Weight * Qty)."

## Modul 5: Transaksi & Payment Gateway (Tahap 6)
Tujuan: Checkout dan Pembayaran Otomatis.
Prompt untuk AI Coding: "Implementasikan sistem Checkout dan Midtrans Payment Gateway.
Buat endpoint /checkout yang menghasilkan Midtrans Snap Token.
Implementasikan Webhook Listener untuk menerima notifikasi pembayaran dari Midtrans (Success/Expired).
Gunakan Database Transaction (Atomicity): Saat bayar sukses, stok berkurang, status pesanan berubah, dan kirim log ke stock_logs secara bersamaan."

## Modul 6: Notifikasi & Caching (Tahap 7)
Tujuan: Kecepatan website dan pengiriman WA.
Prompt untuk AI Coding: "Optimasi backend dengan Redis dan WhatsApp Gateway.
Gunakan Redis Caching (Pattern: Cache Aside) untuk endpoint katalog produk utama.
Implementasikan integrasi WhatsApp Gateway API untuk mengirim notifikasi otomatis saat order dibuat dan dibayar.
Pastikan pengiriman WA berjalan secara Asynchronous menggunakan Goroutine/Queue."