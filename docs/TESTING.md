# Testing Guide - Karima Store

Dokumen ini menjelaskan cara menjalankan testing pada backend Karima Store, termasuk unit tests, integration tests, dan concurrency tests.

## 1. Prerequisites

Sebelum menjalankan test, pastikan Anda memiliki:

*   **Go 1.24+** terinstall.
*   **Docker** berjalan (untuk database test).

### Setup Test Database

Kita perlu membuat database khusus untuk testing agar tidak mengganggu database development atau production.

```bash
# Masuk ke container postgres yang sudah jalan
docker exec -it karima_postgres psql -U karima_store -d karima_db

# Jalankan query CREATE DATABASE di dalam psql
CREATE DATABASE karima_store_test;

# Keluar dari psql
\q
```

Atau one-liner command:
```bash
docker exec karima_postgres psql -U karima_store -d karima_db -c "CREATE DATABASE karima_store_test;"
```

## 2. Configuration (`.env.test`)

Konfigurasi untuk testing biasanya di-override melalui environment variables atau file khusus.

Untuk menjalankan test yang membutuhkan koneksi database (seperti concurrency test), kita perlu memastikan environment variables diset dengan benar. Perhatikan bahwa `TEST_DB_HOST` harus `localhost` jika dijalankan dari host machine (bukan dari dalam docker container).

## 3. Menjalankan Tests

### Menjalankan Semua Unit Tests

Untuk menjalankan seluruh unit test di project:

```bash
go test ./...
```

Atau dengan verbose output:

```bash
go test -v ./...
```

### Menjalankan Concurrency / Race Condition Tests

Test ini dirancang untuk memastikan sistem aman menangani request yang bersamaan (concurrent), misalnya saat checkout produk limited stock oleh banyak user sekaligus.

**Lokasi File:** `internal/services/race_condition_test.go`

**Scenario:**
1. Ada 1 produk dengan stok 10.
2. 20 request checkout masuk secara bersamaan.
3. **Ekspektasi:** Hanya 10 request yang berhasil, 10 sisanya gagal karena stok habis. Stok akhir produk harus 0.

**Cara Menjalankan:**

Karena test ini membutuhkan koneksi ke DB real (bukan mock DB), kita harus menyuplai kredensial database test via environment variable.

```bash
TEST_DB_HOST=localhost \
TEST_DB_PORT=5432 \
TEST_DB_USER=karima_store \
TEST_DB_PASSWORD=lokal \
TEST_DB_NAME=karima_store_test \
go test -v -count=1 ./internal/services -run TestCheckoutConcurrency
```

**Penjelasan Flag:**
*   `-count=1`: Memastikan test tidak dicache oleh Go, jadi benar-benar dijalankan ulang setiap kali.
*   `-v`: Verbose output (menampilkan log sukses/gagal).
*   `-run TestCheckoutConcurrency`: Hanya menjalankan fungsi test dengan nama spesifik ini.

## 4. Troubleshooting

*   **Error connection refused**: Pastikan `TEST_DB_HOST` diset ke `localhost` (bukan `db` atau nama service docker) jika menjalankan `go test` dari terminal laptop/vps langsung.
*   **Database does not exist**: Pastikan Anda sudah menjalankan langkah "Setup Test Database" di atas.
*   **Foreign Key Violation**: Test setup biasanya menangani pembuatan data dependensi (seperti User), namun jika ada perubahan schema, pastikan `test_setup` juga diupdate.

## 5. Laporan Hasil Testing: Concurrency Checkout

**Judul:** Verifikasi Keamanan Stok pada Transaksi Checkout Bersamaan

**Tanggal Pengujian:** 12 Januari 2026

**Tujuan Testing:**
1.  Memastikan sistem tidak menjual barang melebihi stok yang tersedia (overselling) ketika banyak permintaan checkout masuk secara bersamaan (millisecond yang sama).
2.  Memastikan integritas data stok produk dan log stok akurat setelah transaksi paralel.
3.  Memastikan order number yang digenerate unik dan tidak terjadi duplikasi key database.

**Skenario Pengujian:**
*   **Stok Awal:** 10 unit.
*   **Jumlah Request:** 20 concurrent request (goroutines).
*   **Metode:** Menjalankan `TestCheckoutConcurrency` yang memanggil fungsi `Checkout` secara paralel.

**Hasil Testing:**

| Metric | Ekspektasi | Hasil Aktual | Status |
| :--- | :--- | :--- | :--- |
| Jumlah Order Berhasil | 10 | 10 | ✅ PASS |
| Jumlah Order Gagal | 10 | 10 | ✅ PASS |
| Stok Akhir | 0 | 0 | ✅ PASS |
| Duplikasi Order Number | 0 | 0 | ✅ PASS |

**Catatan Temuan & Perbaikan:**
*   **Isu Awal:** Ditemukan bug di mana order number duplikat ketika request masuk di nanodetik yang hampir sama, menyebabkan error `unique constraint violations` pada database.
*   **Perbaikan:** Format `generateOrderNumber` diperbarui untuk menyertakan precision waktu hingga microsecond/nanosecond.
*   **Kesimpulan:** Logic `stock reservation` menggunakan database transaction (`tx`) berfungsi dengan baik untuk mencegah race condition. Sistem aman digunakan untuk high-concurrency sales (misal: Flash Sale).

## 6. Daftar Test Case (Internal Services)

Berikut adalah daftar lengkap unit test dan concurrency test yang tersedia di layer `services`:

| Service | Test Function | Deskripsi |
| :--- | :--- | :--- |
| **AuthService** | `TestAuthService_SyncUser` | Sinkronisasi data user dari provider auth eksternal |
| **CategoryService** | `TestNewCategoryService` | Test konstruktor service |
| | `TestCategoryService_ImplementsInterface` | Verifikasi implementasi interface |
| | `TestCategoryService_GetAllCategories_...` | Test retrieval kategori (Success, Empty, Partial) |
| | `TestCategoryService_GetCategoryStats_...` | Test statistik kategori (Success, Empty, Error) |
| | `TestCategoryService_GetCategoryName_...` | Validasi pengambilan nama kategori (Valid, Invalid, Case Sensitive) |
| | `TestCategoryService_IsValidCategory_...` | Validasi pengecekan kategori valid |
| **CheckoutService** | `TestCheckoutService_TransactionRollback` | Verifikasi rollback DB saat error |
| | `TestCheckoutService_TransactionCommit` | Verifikasi commit DB saat sukses |
| | `TestCheckoutService_SignatureGeneration` | Verifikasi signature Midtrans |
| | `TestCheckoutService_StockDeductionLogic` | Logika pengurangan stok |
| | `TestCheckoutService_PaymentNotificationIdempotency` | Mencegah double process notifikasi pembayaran |
| | `TestCheckoutService_OrderNumberUniqueness` | Verifikasi keunikan nomor order |
| | `TestCheckoutConcurrency` | **Concurrency Test:** Mencegah race condition stok |
| **KomerceService** | `TestKomerceService_SearchDestination` | Pencarian destinasi pengiriman |
| | `TestKomerceService_CalculateShippingCost` | Perhitungan ongkos kirim |
| | `TestKomerceService_CreateOrder` | Pembuatan order pengiriman ke kurir |
| **MediaService** | `TestMediaService_ValidateImageFile_...` | Validasi file gambar (Size, Ext, Security) |
| | `TestMediaService_UploadImage_...` | Upload gambar (Success, Error, Concurrent) |
| | `TestMediaService_DeleteMedia_...` | Penghapusan media |
| | `TestMediaService_SetPrimaryMedia_...` | Set gambar utama produk |
| **NotificationService** | `TestNotificationService_SendWhatsAppMessage_...` | Kirim WA (Success, Fail, API Error) |
| | `TestNotificationService_GetWhatsAppStatus_...` | Cek status koneksi WA gateway |
| | `TestNotificationService_Send...Notification` | Notifikasi Order Created, Payment Success, Shipping |
| | `TestNotificationService_ProcessWhatsAppWebhook` | Handling webhook incoming message |
| **OrderService** | `TestOrderService_GetOrders` | List orders dengan pagination/filter |
| | `TestOrderService_GetOrder` | Get detail order by ID |
| | `TestOrderService_GetOrderByNumber` | Get detail order by Number |
| **PricingService** | `TestPricingService_CalculatePrice` | Perhitungan harga (diskon, pajak, dll) |
| **ProductService** | `TestProductService_CreateProduct_...` | Create produk (Valid, Missing Fields, Duplicate) |
| | `TestProductService_GetProductByID/Update...` | CRUD produk & proteksi SQL Injection |
| | `TestProductService_UpdateProductStock_...` | Update stok (inkremental/dekremental) |
| | `TestProductService_SearchProducts_...` | Pencarian produk |
| | `TestProductService_CacheInvalidation` | Verifikasi hapus cache saat update data |
| **UserService** | `TestNewUserService` | Test konstruktor |
| | `TestUserService_GetUsers_...` | Retrieve users (Pagination, Filters, Limits) |
| | `TestUserService_GetUserByID_...` | Get user detail (Success, NotFound, Admin check) |
| | `TestUserService_UpdateUserRole_...` | Update role (Promote/Demote, Validation) |
| | `TestUserService_Deactivate/ActivateUser_...` | Manajemen status aktif user |
| | `TestUserService_GetUserStats_...` | Statistik user |
| | `TestUserService_ConcurrentRoleUpdates` | **Concurrency:** Update role bersamaan |
| **VariantService** | `TestVariantService_GenerateSKU_...` | Generasi SKU otomatis (Color, Size, Char) |
| | `TestVariantService_CreateVariant_...` | Create varian (Success, Duplicate SKU, Validation) |
| | `TestVariantService_GetVariant...` | Get varian by ID/SKU/Product |
| | `TestVariantService_UpdateVariant...` | Update varian & stok |
| | `TestVariantService_ConcurrentVariantCreation` | **Concurrency:** Create varian bersamaan |

*Catatan: Tabel di atas menampilkan ringkasan scenario. Simbol `...` menandakan terdapat beberapa variasi test case untuk fungsi tersebut (misal: Success, Failed, Edge Cases).*
