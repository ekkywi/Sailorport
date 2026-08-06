# Sailorport — Progress

> Update file ini setiap selesai 1 step. Ini sumber kebenaran saat pindah mesin.

## Status saat ini

- **Step selesai:** 7
- **Step berikutnya:** 8 (Scaffold / golden path — 1 template)
- **Terakhir dikerjakan:** 2026-08-06
- **Mesin terakhir:** rumah / lokal

## Checklist step belajar

- [x] Step 0 — Struktur repo + README + git init
- [x] Step 1 — Go install + endpoint GET `/healthz`
- [x] Step 2 — Refactor `internal/` + config env + POST `/api/v1/echo`
- [x] Step 3 — PostgreSQL via Docker Compose + koneksi DB
- [x] Step 4 — Migrasi tabel `services`
- [x] Step 5 — CRUD catalog API
- [x] Step 6 — Portal web menampilkan catalog (list + create)
- [x] Step 7 — Update & delete service di portal (CRUD UI lengkap)
- [ ] Step 8 — Scaffold / golden path (1 template)
- [ ] Step 9 — Auth OIDC + RBAC dasar
- [ ] Step 10 — Worker registry + agent deploy

## Yang sudah jalan

```bash
# 1. Postgres
cd deploy/compose && docker compose up -d

# 2. API
cd apps/api && go run .

# 3. Portal web
cd apps/web && npm install && npm run dev
```

| Endpoint / UI | Method / aksi | Hasil |
|---------------|---------------|-------|
| `/healthz` | GET | `{"status":"ok","service":"sailorport-api","version":"0.1.0"}` |
| `/api/v1/echo` | POST | `{"reply":"Sailorport received: ..."}` |
| `/api/v1/services` | GET | list services (`[]` jika kosong) |
| `/api/v1/services` | POST | create service → `201` |
| `/api/v1/services/{id}` | GET | get satu service |
| `/api/v1/services/{id}` | PUT | update service |
| `/api/v1/services/{id}` | DELETE | hapus → `204` |
| `http://localhost:5173` | UI | list + create + edit + delete |

Env vars API: `PORT`, `APP_ENV`, `APP_VERSION`, `DATABASE_URL`

Default `DATABASE_URL`:
`postgres://sailorport:sailorport@localhost:5432/sailorport?sslmode=disable`

Portal: Vite proxy `/api` dan `/healthz` → `http://localhost:8080`  
API: middleware CORS untuk `http://localhost:5173`

## Struktur file saat ini

```text
Sailorport/
├── apps/api/
│   ├── main.go
│   ├── go.mod
│   └── internal/
│       ├── config/
│       ├── db/
│       ├── handler/   (health, echo, services, cors)
│       ├── migrate/
│       ├── model/
│       └── store/
├── apps/web/
│   ├── vite.config.ts
│   └── src/
│       ├── App.tsx      (form dual-mode create/edit + list + hapus)
│       ├── api.ts       (list/create/update/delete)
│       ├── types.ts
│       ├── App.css
│       └── main.tsx
├── apps/worker/       ← kosong
├── apps/agent/        ← kosong
├── packages/shared/   ← kosong
├── deploy/compose/    ← Postgres
└── docs/
```

## Catatan belajar

### Backend (Step 0–5)

- Handler → store → model; `$1` parameter query
- Goose + embed untuk migrasi
- CORS middleware untuk browser Vite
- Status: 200, 201, 204, 400, 404, 409

### Frontend (Step 6–7)

- `useState` / `useEffect` / `fetch`
- `import type` untuk tipe saja (`FormEvent`)
- Form dual-mode: `editingId === null` → create, ada id → update
- `window.confirm` sebelum delete
- URL harus exact: `/api/v1/services/{id}` (bukan `/api/api/service/...`)
- Linux case-sensitive: `App.css` ≠ `app.css`
- Blank putih / 404 → cek Console + Network (F12)

## Blocker / pertanyaan

- (kosong)

## Next action (urutan Step 8)

1. Rancang endpoint scaffold (mis. `POST /api/v1/scaffold`)
2. Siapkan 1 template service (API Go atau Node sederhana)
3. Generate file dari template + daftar ke catalog
4. Tombol/form “Create from template” di portal

## Cara lanjut di mesin lain

Lihat `docs/CONTINUE.md` dan paste prompt dari `docs/RESUME-PROMPT.md` ke chat Cursor baru.
