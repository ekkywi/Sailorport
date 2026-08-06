# Sailorport — Progress

> Update file ini setiap selesai 1 step. Ini sumber kebenaran saat pindah mesin.

## Status saat ini

- **Step selesai:** 8 (Scaffold / golden path — template `go-api`)
- **Step berikutnya:** 9 (Auth OIDC + RBAC dasar) — atau environments bila ingin dulu
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
- [x] Step 8 — Scaffold / golden path (1 template `go-api`)
- [ ] Step 9 — Auth OIDC + RBAC dasar
- [ ] Step 10 — Worker registry + agent deploy

## Yang sudah jalan

```bash
# 1. Postgres
cd deploy/compose && docker compose up -d

# 2. API (dari apps/api; templates auto-detect ../../templates)
cd apps/api && go run .

# opsional override:
# export SAILORPORT_TEMPLATES=~/Projects/Sailorport/templates
# export SAILORPORT_WORKSPACE=~/sailorport-workspace

# 3. Portal
cd apps/web && npm run dev
```

| Endpoint / UI | Method | Hasil |
|---------------|--------|-------|
| `/api/v1/templates` | GET | daftar template |
| `/api/v1/scaffold` | POST | generate folder + create catalog |
| `http://localhost:5173` | UI | form “Create from template” + catalog CRUD |

## Struktur penting Step 8

```text
templates/go-api/          manifest + *.tmpl
apps/api/internal/template registry + generate
apps/api/internal/service/scaffold.go
apps/web/src/features/scaffold/
```

## Catatan belajar

- Scaffold: validasi → generate files → insert catalog
- `{{.Name}}` di `*.tmpl` diganti `text/template`
- Nama service: kebab-case lowercase (`payments-api`)
- Env: `SAILORPORT_TEMPLATES`, `SAILORPORT_WORKSPACE`

## Next action (Step 9 arah)

Auth OIDC + RBAC, atau environments (dev/staging/prod) sebelum deploy agent.
