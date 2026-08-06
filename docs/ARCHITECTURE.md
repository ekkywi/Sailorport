# Sailorport — Architecture

Modular monolith berlapis. Cukup untuk industri di tahap MVP; hindari over-engineering (hexagonal penuh, CQRS, microservices) sampai domain lebih besar.

## Prinsip

1. **Control plane ≠ runtime** — API/worker orkestrasi; agent di node yang build/run container.
2. **Satu alasan berubah per lapisan** — HTTP ≠ bisnis ≠ SQL ≠ UI.
3. **main/App tipis** — wiring saja; logic di package/feature.
4. **Abstraksi saat perlu** — interface di batas yang diuji (mis. repository), bukan di setiap fungsi.
5. **Nama jujur** — export/path harus match (`listServices`, `/api/v1/services`).

## Lapisan API (`apps/api`)

```text
HTTP request
  → handler     parse request, status code, JSON
  → service     validasi + aturan bisnis
  → store       SQL / Postgres
  → model       entity & DTO
```

| Package | Boleh | Tidak boleh |
|---------|-------|-------------|
| `handler` | Decode JSON, path params, panggil service, map error → HTTP | SQL, aturan bisnis panjang |
| `service` | Validasi, orkestrasi use case | `net/http`, SQL mentah |
| `store` | Query SQL, map `sql.ErrNoRows` | HTTP status, JSON response |
| `model` | Struct data | Logic HTTP/SQL |
| `config` / `db` / `migrate` | Infrastruktur | Domain catalog |

Alur error (contoh catalog):

- input salah → `service.ErrInvalid` → **400**
- tidak ketemu → `service.ErrNotFound` → **404**
- bentrok nama → `service.ErrConflict` → **409**
- lainnya → log + **500**

Response error JSON seragam:

```json
{"error":"name is required"}
```

## Lapisan Web (`apps/web`)

```text
src/
  app/                 shell / routing nanti
  features/<domain>/   UI + api client + types per fitur
  components/ui/       primitif shadcn (button, input, card, …)
  layouts/             layout halaman (AuthLayout, nanti AppShell)
  lib/                 http client, utils (cn)
  styles/              CSS legacy catalog/scaffold (belum Tailwind)
  index.css            Tailwind v4 + theme shadcn
```

| Area | Tanggung jawab |
|------|----------------|
| `features/catalog/*` | Halaman, form, list, `fetch` ke `/api/v1/services` |
| `features/auth/*` | Login/register UI, panggil `/api/v1/auth/*` |
| `layouts/AuthLayout` | Shell split-screen untuk gate login |
| `components/ui/*` | Primitif UI reusable (shadcn) |
| `App.tsx` | Susun layout; jangan menumpuk logic semua fitur |
| `lib/http.ts` | Token + `apiFetch`; tanpa JSX |

## Batas sistem (produk)

```text
Developer → Web → API → Postgres (+ Redis nanti)
                      → Worker → Agent (Docker) → callback
```

## Konvensi

- Route API: `/api/v1/...`
- Health: `GET /healthz`
- Commit: `feat(api):`, `feat(web):`, `refactor(api):`, `docs:`
- Fitur baru: ikut lapisan di atas; jangan bypass `service` dari handler ke store kecuali health/infra

## Yang sengaja ditunda

- Hexagonal/ports penuh per aggregate
- Shared OpenAPI package (nanti di `packages/shared`)
- Migrasi penuh dashboard web ke Tailwind/shadcn (auth sudah)
- State management global di web (Redux/Zustand) — belum perlu
