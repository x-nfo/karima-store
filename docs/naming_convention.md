# Naming Convention - Karima Store

## 1. Database
- Table names: `snake_case` & plural (contoh: `product_skus`).
- Column names: `snake_case` (contoh: `is_active`).

## 2. Golang
- Variables & Functions: `camelCase` (contoh: `calculateTotal`).
- Exported Structs/Functions: `PascalCase` (contoh: `GetOrderById`).
- Interfaces: Diakhiri dengan "er" jika memungkinkan (contoh: `OrderRepository`).

## 3. API Endpoints
- Resource-based: `/products`, `/categories`.
- Action-based: `/checkout`, `/apply-coupon`.
- Versioning: Selalu gunakan prefix `/api/v1/`.