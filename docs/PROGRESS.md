# Sailorport — Progress

> Update file ini setiap selesai 1 step. Ini sumber kebenaran saat pindah mesin.

## Status saat ini

- **Step selesai:** 7.5 (Architecture foundation)
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
- [x] Step 7.5 — Architecture foundation (lapisan service, router, pecah web)
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
| `/healthz` | GET | health JSON |
| `/api/v1/echo` | POST | demo (akan dihapus nanti) |
| `/api/v1/services` | CRUD | lewat `handler → service → store` |
| `http://localhost:5173` | UI | catalog CRUD (`features/catalog`) |

Error API sekarang JSON: `{"error":"..."}`

## Struktur file saat ini

```text
apps/api/
├── main.go                    # wiring saja
└── internal/
    ├── config/ db/ migrate/
    ├── handler/               # HTTP: router, cors, respond, services
    ├── service/               # use case + validasi (+ test)
    ├── store/                 # SQL
    └── model/

apps/web/src/
├── App.tsx                    # shell
├── styles/app.css
└── features/catalog/          # page, form, list, api, types

docs/ARCHITECTURE.md           # aturan lapisan
```

## Catatan belajar (7.5)

- Lapisan: handler (HTTP) → service (bisnis) → store (SQL)
- Interface `Repository` di service = port persistence
- `writeError` / `writeJSON` seragam
- Web: pecah per feature agar `App.tsx` tidak membengkak
- CORS header harus exact: `Access-Control-Allow-Origin`

## Blocker / pertanyaan

- (kosong)

## Next action (urutan Step 8)

1. Baca `docs/ARCHITECTURE.md` — scaffold ikut lapisan yang sama
2. Endpoint scaffold + 1 template
3. Generate + daftar ke catalog
4. UI “Create from template”

## Cara lanjut di mesin lain

Lihat `docs/CONTINUE.md` dan paste prompt dari `docs/RESUME-PROMPT.md`.
