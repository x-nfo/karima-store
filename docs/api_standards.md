# API Standards - Karima Store

## 1. Response Format
Semua response API harus dalam format JSON dengan struktur:
- Success: `{ "status": "success", "data": { ... } }`
- Error: `{ "status": "error", "message": "Pesan kesalahan yang jelas", "code": 400 }`

## 2. Naming Convention
- URL: kebab-case (contoh: `/api/v1/product-bundles`)
- JSON Keys: snake_case (contoh: `is_busui_friendly`)
- Golang Structs: PascalCase (contoh: `ProductSku`)

## 3. HTTP Methods
- GET: Mengambil data
- POST: Membuat data baru
- PUT: Mengupdate seluruh data
- PATCH: Mengupdate sebagian data (contoh: update stok saja)
- DELETE: Menghapus data

## 4. Date Format
Gunakan ISO 8601 (UTC): `YYYY-MM-DDTHH:mm:ssZ`