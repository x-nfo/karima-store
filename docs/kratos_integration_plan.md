# Perencanaan Integrasi Authentication Ory Kratos

Dokumen ini berisi panduan teknis untuk frontend developer dalam mengintegrasikan sistem autentikasi Ory Kratos ke dalam aplikasi client (Web/Mobile).

## 1. Informasi Sistem & Versi

| Komponen | Spesifikasi | Keterangan |
| :--- | :--- | :--- |
| **Service** | Ory Kratos | Identity Server |
| **Versi** | `v1.1.0` | Pastikan menggunakan SDK yang kompatibel |
| **SDK** | `@ory/client` | Library JavaScript/TypeScript resmi |
| **Auth Domain** | `https://auth.karimasyari.com` | Base URL untuk semua request Auth |
| **Session Cookie** | `karima_session` | Domain: `karimasyari.com` (Subdomain Sharing) |

> [!NOTE]
> Pastikan aplikasi frontend berjalan di domain yang sama atau subdomain dari `karimasyari.com` (misal: `www.karimasyari.com` atau `admin.karimasyari.com`) agar Session Cookie dapat dibaca.

## 2. Alur Integrasi (Browser Flow)

Integrasi Ory Kratos menggunakan konsep **Self-Service Flows**. Frontend tidak membuat form login sembarangan, tetapi "meminta" form dari Kratos.

### Langkah-langkah Umum:
1. **Inisialisasi**: Browser user diarahkan ke *Init Endpoint* Kratos.
2. **Redirect**: Kratos akan me-redirect browser kembali ke *UI Page* di frontend dengan parameter `?flow=<flow_id>`.
3. **Fetching**: Frontend mengambil detail form (JSON) dari Kratos menggunakan `flow_id` tersebut.
4. **Rendering**: Frontend merender input form berdasarkan data JSON yang diterima.
5. **Submission**: User submit form langsung ke Kratos (Kratos menangani security, hashing, validation).

---

## 3. URL & Endpoints Penting

Berikut adalah daftar URL yang wajib diketahui oleh Frontend Developer.

### A. Memulai Flow (Initialization)
Gunakan URL ini pada tombol seperti "Login", "Register", "Lupa Password".

| Aksi | URL (Link Href) | Deskripsi |
| :--- | :--- | :--- |
| **Login** | `/self-service/login/browser` | Memulai proses login |
| **Register** | `/self-service/registration/browser` | Memulai proses pendaftaran akun |
| **Recovery** | `/self-service/recovery/browser` | Memulai lupa password (reset) |
| **Settings** | `/self-service/settings/browser` | Memulai update profil/password (harus login) |
| **Verification**| `/self-service/verification/browser`| Memulai verifikasi email |
| **Logout** | `/self-service/logout/browser` | Melakukan logout (otomatis redirect) |

*Catatan: Semua URL di atas harus di-prefix dengan `https://auth.karimasyari.com` jika diakses dari domain berbeda, atau proxy path jika dikonfigurasi demikian.*

### B. Halaman UI Frontend (Routing)
Frontend harus menyediakan route/halaman berikut untuk menangani redirect dari Kratos.

| Route Frontend | Kegunaan | Parameter Masuk |
| :--- | :--- | :--- |
| `/login` | Menampilkan form login | `?flow=<id>` |
| `/registration` | Menampilkan form registrasi | `?flow=<id>` |
| `/recovery` | Menampilkan form lupa password | `?flow=<id>` |
| `/settings` | Menampilkan form ganti profil/password | `?flow=<id>` |
| `/verification` | Menampilkan status verifikasi | `?flow=<id>` |
| `/error` | Menampilkan halaman error sistem | `?id=<error_id>` |

---

## 4. Implementasi Teknis (Code Snippet)

### Instalasi SDK
```bash
npm install @ory/client
```

### Konfigurasi SDK
```typescript
import { Configuration, FrontendApi } from "@ory/client"

const frontendApi = new FrontendApi(
  new Configuration({
    basePath: "https://auth.karimasyari.com",
    baseOptions: {
      withCredentials: true, // PENTING: Untuk mengirim/menerima Cookies
    },
  }),
)
```

### Contoh: Halaman Login (React)

```tsx
import { useEffect, useState } from "react"
import { useSearchParams } from "react-router-dom"

export const LoginPage = () => {
  const [searchParams] = useSearchParams()
  const flowId = searchParams.get("flow")
  const [flow, setFlow] = useState(null)

  // 1. Jika tidak ada flowId, redirect user ke Init Endpoint
  useEffect(() => {
    if (!flowId) {
      window.location.href = "https://auth.karimasyari.com/self-service/login/browser"
      return
    }

    // 2. Fetch data flow dari Kratos menggunakan SDK
    frontendApi.getLoginFlow({ id: flowId })
      .then(({ data }) => setFlow(data))
      .catch((err) => {
        // Handle error: misal flow expired
        console.error(err)
      })
  }, [flowId])

  if (!flow) return <div>Loading...</div>

  // 3. Render Form
  return (
    <form action={flow.ui.action} method={flow.ui.method} onSubmit={handleSubmit}>
      {/* 
        Render nodes dinamis dari Kratos.
        flow.ui.nodes berisi field: csrf_token, password, identifier, dll.
      */}
      {flow.ui.nodes.map((node) => (
        <input 
            key={node.attributes.name}
            name={node.attributes.name}
            type={node.attributes.type} 
            defaultValue={node.attributes.value}
        />
      ))}
      <button type="submit">Sign In</button>
    </form>
  )
}
```

## 5. Session Management

Setelah user berhasil login, Kratos akan menanam **Cookie** (`karima_session`) di browser.
Untuk mengecek apakah user sedang login (Session Check) di aplikasi atau backend lain:

**Frontend (Check Session Client-side):**
```typescript
frontendApi.toSession()
  .then(({ data }) => {
    console.log("User Logged In:", data.identity)
  })
  .catch(() => {
    console.log("User Not Logged In")
  })
```

**Backend (Check Session Server-side):**
Backend service dapat memvalidasi session dengan memanggil API Kratos `/sessions/whoami` dengen meneruskan Cookie user.

---

## 6. Testing & Development

Untuk keperluan testing lokal (localhost):

1. **Jalankan Kratos Lokal**: Pastikan Docker container berjalan.
2. **Setup Hosts**:
   Tambahkan entry di `/etc/hosts` agar cookie domain match (opsional tapi recommended):
   ```
   127.0.0.1 auth.karimasyari.local
   127.0.0.1 karimasyari.local
   ```
3. **Akses**: Gunakan URL lokal yang diset di env `KRATOS_PUBLIC_URL`, biasanya `http://localhost:4433`.

> [!IMPORTANT]
> **Common Pitfall - CORS & Cookies**: 
> Isu paling sering terjadi adalah Cookie tidak terkirim.
> 1. Pastikan `withCredentials: true` di setiap request Axios/Fetch.
> 2. Pastikan Frontend dan Kratos berada di _top-level domain_ yang sama (misal `app.site.com` dan `auth.site.com`).
