# Architecture Path - Karima Store

## 1. Domain Strategy
Sistem ini menggunakan pemisahan domain untuk keamanan dan isolasi sesi:

- **Storefront (Customer):** `https://karima.com`
  - Berjalan di Cloudflare Pages.
  - Berinteraksi dengan API untuk belanja.
- **Admin Panel:** `https://admin.ks-backend.cloud`
  - Berjalan di Cloudflare Pages.
  - Dilindungi oleh Cloudflare Zero Trust.
- **API Gateway & Backend:** `https://api.ks-backend.cloud`
  - Berjalan di VPS menggunakan Docker.
  - Endpoint utama untuk semua transaksi.
- **Authentication (Ory Kratos):** `https://auth.ks-backend.cloud`
  - Menangani login, register, dan session management.

## 2. Infrastructure Flow
1. User login melalui `auth.ks-backend.cloud`.
2. Session cookie diberikan untuk domain `.ks-backend.cloud`.
3. Frontend (Admin/Store) mengirimkan request ke `api.ks-backend.cloud`.
4. Backend memverifikasi session ke Ory Kratos secara internal.

## 3. CORS Policy
Backend (Golang) harus mengizinkan (AllowOrigins):
- `https://karima.com`
- `https://admin.ks-backend.cloud`
Wajib mengaktifkan `AllowCredentials: true` untuk mendukung HttpOnly Cookies.

## 4. File: .env.example
Salin ini ke root folder Anda (bukan di dalam folder docs). File ini memberi tahu AI variabel apa saja yang wajib dikonfigurasi.

## 5. Development Environment
- **OS:** Windows dengan WSL2 (Ubuntu/Debian).
- **Container Engine:** **Podman** (bukan Docker Desktop).
- **Orchestration:** `podman-compose`.
- **Note:** AI harus memberikan perintah menggunakan `podman` atau `podman-compose`. Pastikan konfigurasi volume dan networking ramah terhadap Podman Rootless.