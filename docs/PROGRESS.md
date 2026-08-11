# Sailorport — Progress

> Update file ini setiap selesai 1 step. Ini sumber kebenaran saat pindah mesin.

## Status saat ini

- **Step selesai:** 11 (Docker Compose full stack)
- **Step berikutnya:** Harden (agent token) + multi-port deploy (opsional); atau Environments
- **Terakhir dikerjakan:** 2026-08-11 — compose api+web+postgres; docs dual-mode (dev vs self-host)
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
- [x] Step 10C.2 — Agent poll job + Docker build/run + PATCH status
- [x] Step 10C.3 — Portal UI tombol Deploy + list deployments
- [x] Step 11 — Docker Compose full stack
- [ ] Harden — agent token + (opsional) multi-port deploy
- [ ] Environments (dev/staging/prod)

## Yang sudah jalan

### Mode development (disarankan sehari-hari)

Hanya Postgres di Docker; API/web/agent di host untuk hot-reload:

```bash
cd deploy/compose && docker compose up -d postgres
cd apps/api && go run .
cd apps/web && npm run dev

# terminal terpisah — agent (setelah API jalan)
cd apps/agent && go run .
```

Portal Vite: `http://localhost:5173` (proxy ke API `:8080`).

### Mode self-host / pack (Step 11)

Control plane penuh di Compose:

```bash
cd deploy/compose && docker compose up -d --build
# web http://localhost:5173  api http://localhost:8080  postgres :5433

# agent tetap di host (butuh Docker CLI)
cd apps/agent && SAILORPORT_API_URL=http://localhost:8080 go run .
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
| `SAILORPORT_POLL_INTERVAL` | `5s` | interval poll job deploy |
| `SAILORPORT_DEPLOY_PORT_BASE` | `18080` | host port untuk container (MVP: satu port) |

Role: `admin`, `developer`, `viewer`

### Deployments (10C.1)

- Migrasi: `00005_create_deployments.sql`
- Lapisan: `model` → `store` → `service` → `handler`
- Status: `pending` → `claimed` → `building` → `running` | `failed` | `stopped`
- Claim atomik (`FOR UPDATE SKIP LOCKED`) + join service → `service_name`, `workspace_path`
- Create menolak service tanpa workspace (scaffold dulu)

### Agent deploy (10C.2)

- Poll `POST /api/v1/agent/jobs/next` setiap `SAILORPORT_POLL_INTERVAL`
- Flow: claim → `building` → `docker build` + `docker run` → `running` | `failed`
- Image tag: `sailorport/{service_name}:{deployment_id[:8]}`
- Container name: `sailorport-{service_name}`; map host port → container `:8080`
- Template `go-api` punya `Dockerfile.tmpl` (scaffold baru otomatis dapat `Dockerfile`)
- Service scaffold **lama** (sebelum 10C.2) perlu `Dockerfile` manual di workspace

**Test deploy end-to-end:**

```bash
# API + agent jalan; scaffold service baru (dapat Dockerfile)
# Buat deployment pending (JWT developer+), agent akan pick up otomatis
curl http://localhost:18080/healthz   # service yang di-deploy
```

### Portal Deploy UI (10C.3)

- Feature: `apps/web/src/features/deployments/` (`types`, `api`, `DeploymentsDialog`)
- Catalog: tombol **Deploy** (rocket) hanya jika service punya `workspace_path`
- Setelah create → dialog list deployments (status badge, refresh, poll 3s saat job aktif)
- Link `http://localhost:{port}/healthz` saat status `running`
- **AppShell:** sidebar desktop bisa collapse/expand (persisted di `localStorage`)
- Catalog table: layout actions lebih rapi (`table-fixed`, overflow)

### Template fix

- `templates/go-api/main.go.tmpl`: `ListenAndServe(":8080", mux)` (bukan `nil`) — tanpa ini `/healthz` di container selalu 404

### Web UI (portal setelah login)

- **Layout:** `AppShell` — sidebar collapsible + topbar (Linear-style), harbour theme + dark mode
- **Routes (flat):** `/overview`, `/catalog`, `/worker` (bukan nested)
- **Catalog UX:** daftar saja; **Create service** = scaffold; **Register existing** = metadata; edit/delete/deploy via actions; dialog deployments
- **Workers:** tabel dengan status badge, relative last seen
- **Overview:** metrik services/workers + panel recent
- **Shared:** `src/components/app/` (DataPanel, Toolbar, EmptyState, …)
- **Auth:** Inter Variable, `AuthLayout` centred (tetap terpisah dari app shell)

Template **masih di disk** (`templates/go-api/`), bukan DB. Service hasil scaffold punya `template_id` + `workspace_path` di Postgres.

Setelah `git pull` di mesin baru: `cd apps/web && npm install`

## Known debt (sengaja ditunda)

- **Template management** belum CRUD di DB/portal
- **Agent endpoints publik** (claim/update) — token agent menyusul saat harden
- **Deploy port** MVP pakai satu `PortBase` (18080); multi-service collision belum di-handle
- **Workspace lama** (path `/tmp/...`) tidak ikut terhapus saat delete (di luar root baru); scaffold ulang ke `data/workspaces`

### Debt yang sudah diperbaiki (2026-08-11)

- Workspace default → `data/workspaces` (bukan `/tmp`), override `SAILORPORT_WORKSPACE`
- Agent Ctrl+C mengirim heartbeat `offline`
- Delete service menghapus folder workspace jika path di bawah workspace root

### Compose full stack (Step 11)

- `deploy/compose/docker-compose.yml` — `postgres` + `api` + `web`
- `apps/api/Dockerfile`, `apps/web/Dockerfile`, `apps/web/nginx.conf` (proxy `/api` + `/healthz`)
- Volume: `../../templates` → `/templates`, `../../data/workspaces` → `/data/workspaces`
- Agent **tidak** di compose (perlu Docker daemon di host untuk deploy workload)

## Next action

1. Harden: token untuk endpoint agent (claim/update tidak publik)
2. Multi-port deploy (hindari collision di `18080`)
3. Atau Environments (dev/staging/prod)

## Cara lanjut di mesin lain

Lihat `docs/CONTINUE.md` dan paste prompt dari `docs/RESUME-PROMPT.md`.
