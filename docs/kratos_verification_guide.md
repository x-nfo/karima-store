# Panduan Verifikasi Integrasi Ory Kratos di Production

Dokumen ini berisi langkah-langkah sistematis untuk memverifikasi apakah integrasi Ory Kratos sudah berjalan dengan benar di lingkungan production.

## 1. Cek Status Container

Langkah pertama adalah memastikan service Kratos berjalan dengan baik di server VPS.

Jalankan command berikut di terminal server:

```bash
docker ps | grep kratos
```

**Ekspektasi Output:**
- Container `kratos` harus dalam status **Up**.
- Container `kratos_migrate` mungkin statusnya **Exited (0)** (ini normal karena hanya berjalan sekali untuk migrasi).

Jika container `kratos` statusnya `Restarting` atau tidak muncul, cek logs:

```bash
docker logs kratos --tail 100
```

## 2. Verifikasi Endpoint Health

Pastikan Kratos dapat diakses dan "sehat". Gunakan `curl` dari lokal komputer atau dari dalam server.

**Dari Server (Internal Network):**
```bash
docker exec -it kratos curl http://localhost:4433/health/ready
```
*Atau jika curl tidak ada di dalam container, gunakan `wget`:*
```bash
docker exec -it kratos wget -qO- http://localhost:4433/health/ready
```

**Dari Public Internet (Browser/Postman):**
Akses URL Kratos Public Anda (misal: `https://auth.karimasyari.com/health/ready`).

**Ekspektasi Output:**
```json
{"status":"ok"}
```

## 3. Verifikasi Konfigurasi (Environment & File)

Pastikan konfigurasi yang termuat sudah benar, terutama terkait URL dan CORS.

### A. Cek Environment Variables
Inspect container untuk melihat variabel lingkungan yang aktif:

```bash
docker inspect kratos --format '{{range .Config.Env}}{{println .}}{{end}}'
```

**Variable Kunci yang Harus Dicek:**
- `SERVE_PUBLIC_BASE_URL`: Harus mengarah ke URL publik (misal `https://auth.karimasyari.com`).
- `SELFSERVICE_ALLOWED_RETURN_URLS`: Harus mencakup domain frontend (`https://karimasyari.com`, dll).
- `DSN`: Koneksi database harus valid.

### B. Cek Konfigurasi File (CORS & Cookie)
Cek file `deploy/kratos/kratos.prod.yml` yang ter-mount (atau baked-in).

```bash
docker exec -it kratos cat /etc/config/kratos/kratos.prod.yml
```

**Pastikan:**
- `serve.public.cors.enabled: true`
- `serve.public.cors.allowed_origins`: Berisi domain frontend Anda.
- `session.cookie.domain`: **PENTING**. Jika frontend ada di `www.karimasyari.com` dan auth di `auth.karimasyari.com`, cookie domain sebaiknya kosong (host-only) atau diset ke main domain jika perlu sharing subdomain, tapi perhatikan `SameSite` policy.

## 4. Test Integrasi Flow (Browser)

Ini adalah tes fungsional utama. Lakukan di Browser (Chrome/Firefox).

### A. Test Login Flow
1. Buka URL inisialisasi Login di browser:
   `https://auth.karimasyari.com/self-service/login/browser`
2. **Ekspektasi**: Browser harus redirect ke Halaman UI Login Frontend Anda (misal `https://karimasyari.com/login?flow=<id>`).
3. Cek URL bar, pastikan ada parameter `?flow=...`.

### B. Test API Fetch Flow
1. Buka Inspector Browser (F12) > Network Tab.
2. Di halaman login frontend, lihat request ke Kratos untuk mengambil data flow.
   - URL: `https://auth.karimasyari.com/self-service/login/flows?id=<flow_id>`
   - Method: `GET`
3. **Ekspektasi**: Status 200 OK dan response JSON berisi field form (csrf_token, password, identifier).
4. Jika Error (401/403/Cors Error): Berarti konfigurasi CORS atau Cookie salah.

### C. Test Submit Form
1. Isi username/password dan submit.
2. Perhatikan request `POST` ke URL action yang ada di flow (biasanya `https://auth.karimasyari.com/self-service/login?flow=<id>`).
3. **Ekspektasi**:
   - Jika sukses: Redirect ke `default_browser_return_url` (misal home page user).
   - Server meresponse dengan header `Set-Cookie` (`karima_session`).

## 5. Cek Session Cookie

Setelah login berhasil:
1. Buka Inspector > Application > Cookies.
2. Cari cookie bernama `karima_session` (atau sesuai config `session.cookie.name`).
3. **Validasi**:
   - **Domain**: Apakah sesuai? (misal `.karimasyari.com` atau `auth.karimasyari.com`). Hal ini menentukan apakah frontend bisa membacanya (jika perlu) atau browser akan mengirimnya.
   - **HttpOnly**: Harusnya `true` (tidak bisa diakses JS).
   - **Secure**: Harusnya `true` (karena HTTPS).
   - **SameSite**: Biasanya `Lax` atau `None` (jika cross-site perlu).

## 6. Test Endpoint "WhoAmI"

Untuk memastikan session valid, panggil endpoint session dari frontend atau browser.

URL: `https://auth.karimasyari.com/sessions/whoami`

**Ekspektasi**:
- Jika login: Return JSON identity user (`200 OK`).
- Jika tidak login: Return Error (`401 Unauthorized`).

## Troubleshooting Checklist

| Gejala | Kemungkinan Penyebab | Solusi |
| :--- | :--- | :--- |
| **502 Bad Gateway** | Container mati atau Port tidak terekspos | Cek `docker logs kratos` dan konfigurasi port/tunnel. |
| **CORS Error** di Browser | `allowed_origins` tidak match | Tambahkan URL frontend lengkap (termasuk protocol & port) di `kratos.prod.yml`. |
| **Redirect Loop** | `default_browser_return_url` salah | Cek konfigurasi env variable `SELFSERVICE_DEFAULT_BROWSER_RETURN_URL`. |
| **Cookie tidak tersimpan** | Domain cookie salah atau Scheme bukan HTTPS | Pastikan akses via HTTPS dan config cookie domain benar. |
| **CSRF Error** | Cookie CSRF hilang/mismatch | Pastikan SameSite setting benar dan gunakan HTTPS. |

---

Jika semua langkah di atas berhasil (Status Container OK, Health OK, Flow Login Berhasil, Session Terbaca), maka integrasi Kratos di Production sudah **BENAR**.
