# 🚀 Panduan Setup CI/CD dengan Portainer & GitHub

Panduan ini akan membantu Anda menghubungkan Portainer dengan Repository GitHub Anda untuk deployment otomatis.

## Prasyarat
- [x] Code sudah di-push ke GitHub (Branch: `production`)
- [x] `docker-compose.prod.yml` sudah tersedia di repo
- [ ] Anda memiliki akses admin ke Portainer
- [ ] Anda sudah menyiapkan environment variables (isi `.env.production`)

## Langkah 1: Buat Stack Baru di Portainer

1. Login ke Portainer Web UI.
2. Klik menu **Stacks** di sidebar kiri.
3. Klik tombol **+ Add stack** di pojok kanan atas.
4. Pilih opsi **Repository** (kotak kedua dari kiri).

## Langkah 2: Konfigurasi Repository

Isi form dengan detail berikut:

- **Name**: `karima-prod` (atau nama lain, gunakan huruf kecil)
- **Repository URL**: `https://github.com/x-nfo/karima-store.git`
  *(Jika repo private, centang "Authentication" dan masukkan Access Token)*
- **Repository Reference**: `refs/heads/production`
  *(Atau cukup ketik `production`)*
- **Compose path**: `docker-compose.prod.yml`
  *(PENTING: Jangan gunakan default `docker-compose.yml`)*

## Langkah 3: Konfigurasi Environment Variables

1. Scroll ke bawah ke bagian **Environment variables**.
2. Buka file `.env.production` di komputer lokal Anda.
3. Copy **SEMUA** isinya.
4. Paste ke dalam kotak Environment variables di Portainer (klik "Switch to advanced mode" jika ingin paste sekaligus).

> ⚠️ **PENTING**: Pastikan variable `DB_PASSWORD`, `REDIS_PASSWORD`, dll sudah terisi nilai asli, bukan placeholder!

## Langkah 4: Aktifkan Automatic Updates (CI/CD)

1. Di bawah section Environment variables, cari bagian **Automatic updates**.
2. Aktifkan toggle **Enable automatic updates**.
3. **Fetch interval**: Bisa pilih `5 minutes` (Polling) ATAU gunakan Webhook.
   - **Rekomendasi**: Gunakan **Webhook** untuk update instan.
   - Copy **Webhook URL** yang muncul (contoh: `https://portainer.yourdomain.com/api/stacks/webhooks/...`).

## Langkah 5: Deploy Stack

1. Klik tombol **Deploy the stack** di bagian paling bawah.
2. Tunggu proses build (bisa memakan waktu beberapa menit karena akan download base image dan compile Go app).
3. Jika sukses, status stack akan berubah menjadi "Active".

---

## Langkah 6: Setup GitHub Webhook (Opsional tapi Recommended)

Jika Anda memilih menggunakan Webhook di Langkah 4:

1. Buka Repository GitHub Anda -> **Settings**.
2. Klik **Webhooks** di menu kiri -> **Add webhook**.
3. **Payload URL**: Paste URL webhook dari Portainer tadi.
4. **Content type**: `application/json` (tidak terlalu berpengaruh untuk Portainer, tapi standard bagus).
5. **Secret**: Kosongkan (kecuali Portainer Anda butuh).
6. **Just the push event**: Pilih ini.
7. Klik **Add webhook**.

🎉 **Selesai!**
Sekarang, setiap kali Anda push code perubahan ke branch `production` di GitHub:
1. GitHub memberitahu Portainer.
2. Portainer pull code terbaru.
3. Portainer re-build image dan restart container.

## Troubleshooting

### "Failed to deploy: env file not found"
Solusi: Pastikan Anda **TIDAK** menggunakan `env_file:` di `docker-compose.prod.yml`. (Sudah kita hapus di langkah sebelumnya, jadi aman).

### "Stack creation failed"
Cek log error. Jika masalah network, pastikan network `karima-network` sudah ada di Portainer (Networks -> Add network).

### Container Restarting Terus
Cek logs container:
1. Klik nama stack `karima-prod`.
2. Klik icon logs (dokumen kecil) di sebelah container `backend`.
3. Lihat pesan error (biasanya koneksi database atau config salah).
