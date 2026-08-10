# Sailorport — Progress

> Update file ini setiap selesai 1 step. Ini sumber kebenaran saat pindah mesin.

## Status saat ini

- **Step selesai:** 10C.1 (Deployments API + agent claim)
- **Step berikutnya:** 10C.2 (Agent poll + docker build/run)
- **Terakhir dikerjakan:** 2026-08-10 — deployments table/API, claim job, status update
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
- [x] Step 9 — Auth lokal + JWT + RBAC
- [x] Step 9 polish — Auth UI profesional (Tailwind v4 + shadcn/ui, harbour theme)
- [x] Step 10A — Worker registry API + portal UI (Overview, Catalog, Workers)
- [x] Step 10B — Agent binary (`apps/agent`: register + heartbeat loop)
- [x] Step 10C.1 — Deployments API (create/list/get + agent claim/update)
- [ ] Step 10C.2 — Agent poll job + Docker build/run + PATCH status
- [ ] Step 10C.3 — Portal UI tombol Deploy + list deployments

## Yang sudah jalan

```bash
cd deploy/compose && docker compose up -d
cd apps/api && go run .
cd apps/web && npm run dev

# terminal terpisah — agent (setelah API jalan)
cd apps/agent && go run .
```

| Endpoint / UI | Auth | Hasil |
|---------------|------|-------|
| `POST /api/v1/auth/register` | publik | buat user |
| `POST /api/v1/auth/login` | publik | JWT token |
| `GET /api/v1/auth/me` | Bearer | profil user |
| `GET/POST/PUT/DELETE /api/v1/services` | viewer+ / developer+ | catalog CRUD |
| `GET /api/v1/templates`, `POST /api/v1/scaffold` | viewer+ / developer+ | golden path |
| `POST /api/v1/workers/register` | publik | agent register |
| `POST /api/v1/workers/{id}/heartbeat` | publik | agent heartbeat |
| `GET /api/v1/workers` | viewer+ | list workers (portal) |
| `POST /api/v1/services/{id}/deployments` | developer+ | buat deploy (`pending`); service harus punya `workspace_path` |
| `GET /api/v1/services/{id}/deployments` | viewer+ | list per service |
| `GET /api/v1/deployments` | viewer+ | list semua |
| `GET /api/v1/deployments/{id}` | viewer+ | detail |
| `PATCH /api/v1/deployments/{id}` | developer+ | update status (portal/curl JWT) |
| `POST /api/v1/agent/jobs/next` | publik | claim 1 job pending → `claimed` (204 jika kosong) |
| `PATCH /api/v1/agent/deployments/{id}` | publik | agent update status tanpa JWT |
| Portal `/login`, `/register` | — | auth gate |
| Portal `/overview`, `/catalog`, `/worker` | JWT | app shell (sidebar + topbar) |

Env API: `AUTH_JWT_SECRET` (default dev-only-change-me)

Env agent:

| Variable | Default | Keterangan |
|----------|---------|------------|
| `SAILORPORT_API_URL` | `http://localhost:8080` | base URL API |
| `SAILORPORT_WORKER_NAME` | hostname mesin | nama worker di registry |
| `SAILORPORT_HEARTBEAT_INTERVAL` | `15s` | interval heartbeat |

Role: `admin`, `developer`, `viewer`

### Deployments (10C.1)

- Migrasi: `00005_create_deployments.sql`
- Lapisan: `model` → `store` → `service` → `handler`
- Status: `pending` → `claimed` → `building` → `running` | `failed` | `stopped`
- Claim atomik (`FOR UPDATE SKIP LOCKED`) + join service → `service_name`, `workspace_path`
- Create menolak service tanpa workspace (scaffold dulu)
- **Belum:** agent belum poll/docker; portal belum tombol Deploy

### Web UI (portal setelah login)

- **Layout:** `AppShell` — sidebar + topbar (Linear-style), harbour theme + dark mode
- **Routes (flat):** `/overview`, `/catalog`, `/worker` (bukan nested)
- **Catalog UX:** daftar saja; **Create service** = scaffold dari template (dialog); **Register existing** = metadata saja (dialog sekunder); edit/delete via dialog
- **Workers:** tabel dengan status badge, relative last seen
- **Overview:** metrik services/workers + panel recent
- **Shared:** `src/components/app/` (DataPanel, Toolbar, EmptyState, …)
- **Auth:** Inter Variable, `AuthLayout` centred (tetap terpisah dari app shell)

Template **masih di disk** (`templates/go-api/`), bukan DB. Service hasil scaffold punya `template_id` + `workspace_path` di Postgres.

Setelah `git pull` di mesin baru: `cd apps/web && npm install`

## Known debt (sengaja ditunda)

- **Delete service** hanya hapus row DB; folder workspace di disk **tidak** dihapus → recreate dengan nama sama bisa gagal (`workspace folder already exists`)
- **Agent stop** (Ctrl+C) belum otomatis set status `offline` di API
- **Template management** belum CRUD di DB/portal
- **Agent endpoints publik** (claim/update) — token agent menyusul saat harden
- **go-api template** belum punya Dockerfile (dibutuhkan di 10C.2)

## Next action (Step 10C.2)

1. Tambah `Dockerfile.tmpl` di `templates/go-api/`
2. Agent: poll `POST /api/v1/agent/jobs/next`
3. `docker build` + `docker run` dari `workspace_path`
4. `PATCH /api/v1/agent/deployments/{id}` → `building` / `running` / `failed`

Lalu 10C.3: tombol Deploy + list deployments di portal Catalog.

## Cara lanjut di mesin lain

Lihat `docs/CONTINUE.md` dan paste prompt dari `docs/RESUME-PROMPT.md`.
