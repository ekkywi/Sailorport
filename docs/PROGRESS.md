# Sailorport — Progress

> Update file ini setiap selesai 1 step. Ini sumber kebenaran saat pindah mesin.

## Status saat ini

- **Step selesai:** 13d — Environments portal UI (dropdown deploy + tampil slug)
- **Step berikutnya:** Opsional — runtime/stop-start per environment; logs; multi-agent targeting
- **Terakhir dikerjakan:** 2026-08-14 — deploy ke dev/staging/prod dari portal; container name `sailorport-{service}-{env}`
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
- [x] Step 12a — User management API (list users, patch role — admin only)
- [x] Step 12b — Portal users page + RBAC UI (hide actions for viewer)
- [x] Step 12c — Admin create user (temporary password / invite-style MVP)
- [x] Step 12d — Admin disable / enable user (blokir login)
- [x] Step 12e — Admin reset password (temporary password)
- [x] Step 12f — Soft-delete user (rename email + `deleted_at`)
- [x] Harden — agent token (register/heartbeat/claim/update)
- [x] R0 — latest deploy di catalog (API `latest_deployment` + kolom Deploy di portal)
- [x] R1 — agent docker helpers (`Stop` / `Start` / `Remove`)
- [x] R2 — runtime controls (stop/start via agent job + portal UI)
- [x] R3 Delete cleanup — stop/rm container saat hapus service
- [x] R4 — multi-port deploy (port unik per container di host agent)
- [x] Environments (dev/staging/prod) — 13a–13d

## Yang sudah jalan

### Mode development (disarankan sehari-hari)

Hanya Postgres di Docker; API/web/agent di host untuk hot-reload:

```bash
cd deploy/compose && docker compose up -d postgres
cd apps/api && go run .
cd apps/web && npm run dev

# terminal terpisah — agent (setelah API jalan)
cd apps/agent && SAILORPORT_AGENT_TOKEN=dev-agent-token go run .
```

Portal Vite: `http://localhost:5173` (proxy ke API `:8080`).

### Mode self-host / pack (Step 11)

Control plane penuh di Compose:

```bash
cd deploy/compose && docker compose up -d --build
# web http://localhost:5173  api http://localhost:8080  postgres :5433

# agent tetap di host (butuh Docker CLI)
cd apps/agent && SAILORPORT_API_URL=http://localhost:8080 SAILORPORT_AGENT_TOKEN=dev-agent-token go run .
```

| Endpoint / UI | Auth | Hasil |
|---------------|------|-------|
| `POST /api/v1/auth/register` | publik | buat user |
| `POST /api/v1/auth/login` | publik | JWT token |
| `GET /api/v1/auth/me` | Bearer | profil user |
| `GET /api/v1/users` | admin | list semua user |
| `POST /api/v1/users` | admin | create user (email, name, password, role) |
| `PATCH /api/v1/users/{id}` | admin | ubah `role` **atau** `disabled` (bukan keduanya sekaligus); tidak boleh ubah role/disable diri sendiri; id harus UUID |
| `POST /api/v1/users/{id}/password` | admin | set temporary password baru; tidak boleh reset password diri sendiri → **204** |
| `DELETE /api/v1/users/{id}` | admin | soft-delete (rename email + `deleted_at`); tidak boleh hapus diri sendiri → **204** |
| `GET/POST/PUT/DELETE /api/v1/services` | viewer+ / developer+ | catalog CRUD; `GET` list menyertakan `latest_deployment` per service |
| `GET /api/v1/templates`, `POST /api/v1/scaffold` | viewer+ / developer+ | golden path |
| `POST /api/v1/workers/register` | agent token | agent register |
| `POST /api/v1/workers/{id}/heartbeat` | agent token | agent heartbeat |
| `GET /api/v1/workers` | viewer+ | list workers (portal) |
| `GET /api/v1/environments` | viewer+ | list environment (dev/staging/prod) |
| `POST /api/v1/services/{id}/deployments` | developer+ | buat deploy (`pending`); body `{"environment":"staging"}` (default `dev`); service harus punya `workspace_path` |
| `GET /api/v1/services/{id}/deployments` | viewer+ | list per service |
| `GET /api/v1/deployments` | viewer+ | list semua |
| `GET /api/v1/deployments/{id}` | viewer+ | detail |
| `PATCH /api/v1/deployments/{id}` | developer+ | update status (portal/curl JWT) |
| `POST /api/v1/services/{id}/runtime/stop` | developer+ | enqueue stop container (deployment harus `running`) → **202** |
| `POST /api/v1/services/{id}/runtime/start` | developer+ | enqueue start container (deployment harus `stopped`) → **202** |
| `POST /api/v1/agent/jobs/next` | agent token | claim 1 deploy job pending → `claimed` (204 jika kosong) |
| `PATCH /api/v1/agent/deployments/{id}` | agent token | agent update deploy status |
| `POST /api/v1/agent/runtime/next` | agent token | claim 1 runtime job (`stop`/`start`) |
| `PATCH /api/v1/agent/runtime/{id}` | agent token | agent selesai runtime job; API update deployment → `stopped`/`running` |
| Portal `/login`, `/register` | — | auth gate |
| Portal `/overview`, `/catalog`, `/worker`, `/users` | JWT | app shell; `/users` admin-only (redirect non-admin) |

Env API: `AUTH_JWT_SECRET` (default `dev-only-change-me`), `SAILORPORT_AGENT_TOKEN` (default `dev-agent-token` — ganti di production)

Env agent:

| Variable | Default | Keterangan |
|----------|---------|------------|
| `SAILORPORT_API_URL` | `http://localhost:8080` | base URL API |
| `SAILORPORT_AGENT_TOKEN` | `dev-agent-token` | harus sama dengan API |
| `SAILORPORT_WORKER_NAME` | hostname mesin | nama worker di registry |
| `SAILORPORT_HEARTBEAT_INTERVAL` | `15s` | interval heartbeat |
| `SAILORPORT_POLL_INTERVAL` | `5s` | interval poll job deploy |
| `SAILORPORT_DEPLOY_PORT_BASE` | `18080` | awal rentang host port workload |
| `SAILORPORT_DEPLOY_PORT_COUNT` | `32` | jumlah port di rentang (`18080`–`18111`) |

Role: `admin`, `developer`, `viewer`

**Admin pertama:** register user biasa, lalu promote di Postgres:

```bash
docker exec -it sailorport-postgres psql -U sailorport -d sailorport \
  -c "UPDATE users SET role = 'admin' WHERE email = 'you@example.com';"
```

Login ulang agar JWT berisi role baru.

### Soft-delete user (12f)

Desain: soft delete (bukan `DELETE FROM`); rename email agar UNIQUE bebas; filter `deleted_at IS NULL`; no self-delete.

| Sub-step | Status | Isi |
|----------|--------|-----|
| 1 Migrasi | ✅ | `00009_add_users_deleted_at.sql` — kolom `deleted_at TIMESTAMPTZ` + index |
| 2 Model | ✅ | `User.DeletedAt *time.Time` |
| 3 Store | ✅ | filter aktif + `SoftDelete` rename `email \|\| '__deleted__' \|\| id` + `disabled=true` |
| 4 Service | ✅ | `Users.SoftDelete` — UUID + no self-delete |
| 5 Handler | ✅ | `DELETE /api/v1/users/{id}` → **204** (admin) |
| 6 Portal | ✅ | tombol Delete + konfirmasi; Enable juga pakai konfirmasi |

```bash
curl -X DELETE http://localhost:8080/api/v1/users/{id} \
  -H "Authorization: Bearer $TOKEN" -w '\n%{http_code}\n'
# sukses → 204; self-delete → 403; id tidak ada → 404
```

### User management API (12a + 12c + 12d + 12e + 12f)

- Lapisan: `store/user` → `service/users` → `handler/user`
- Migrasi `00008_add_users_disabled.sql` — kolom `users.disabled BOOLEAN NOT NULL DEFAULT false`
- Migrasi `00009_add_users_deleted_at.sql` — kolom `users.deleted_at`
- `GET /api/v1/users` — admin only (response termasuk `disabled`; hanya `deleted_at IS NULL`)
- `POST /api/v1/users` body `{email,name,password,role}` — admin only; role boleh `admin`/`developer`/`viewer`; email unik → **409**
- `PATCH /api/v1/users/{id}` — admin only; body `{"role":"..."}` **atau** `{"disabled":true|false}` (bukan keduanya); self-change → **403**; id non-UUID → **400**
- `POST /api/v1/users/{id}/password` body `{"password":"..."}` — admin only; min 8 chars; self-reset → **403**; sukses → **204**
- `DELETE /api/v1/users/{id}` — admin only; soft-delete; self-delete → **403**; sukses → **204**
- Login / `me`: akun `disabled` → **401**; soft-deleted tidak ketemu (store filter)
- Invite-style MVP: admin set / reset temporary password, bagikan manual (belum email SMTP)

### Portal users UI + catalog RBAC (12b + 12c + 12d + 12e + 12f)

- Feature: `apps/web/src/features/users/` (`type`, `api`, `UsersPage`)
- Helper RBAC: `apps/web/src/lib/rbac.ts` — `isAdmin()`, `canWriteCatalog()`
- Route `/users` + guard admin; non-admin → `/overview`
- Sidebar: sections Workspace / Platform / Administration; **Users** hanya admin
- `UsersPage`: tabel user, dropdown role, status active/disabled, **Create user**, **Reset password**, **Disable** / **Enable** (keduanya konfirmasi), **Delete** (soft-delete + konfirmasi); baris diri sendiri = badge + “You”
- Catalog: `canWriteCatalog` — viewer tanpa Create / Deploy / Edit / Delete (desktop + mobile)

### Agent token (Harden)

- Env bersama: `SAILORPORT_AGENT_TOKEN` (API + agent; default dev `dev-agent-token`)
- Middleware `withAgentToken` — header `Authorization: Bearer <token>` (constant-time compare)
- Dilindungi: `POST /workers/register`, `POST /workers/{id}/heartbeat`, `POST /agent/jobs/next`, `PATCH /agent/deployments/{id}`, `POST /agent/runtime/next`, `PATCH /agent/runtime/{id}`
- Agent client mengirim token di setiap request
- Compose API: `SAILORPORT_AGENT_TOKEN` di service `api`
- Tested: tanpa token → **401**; token benar (no job) → **204**

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
- Host port: `AllocateHostPort` — reuse mapping container yang sama (redeploy), else port bebas di `[PortBase, PortBase+Count)`
- Hasil port disimpan di `deployments.port` (portal link `:port/healthz`)
- Template `go-api` punya `Dockerfile.tmpl` (scaffold baru otomatis dapat `Dockerfile`)
- Service scaffold **lama** (sebelum 10C.2) perlu `Dockerfile` manual di workspace

**Test deploy end-to-end:**

```bash
# API + agent jalan; scaffold service baru (dapat Dockerfile)
# Buat deployment pending (JWT developer+), agent akan pick up otomatis
curl http://localhost:18080/healthz   # service yang di-deploy
```

### Portal Deploy UI (10C.3 + R0)

- Feature: `apps/web/src/features/deployments/` (`types`, `api`, `DeployDialog`, `DeploymentsDialog`); `features/environments/` untuk list env
- Catalog kolom **Deploy**: badge environment + status deploy terakhir (`latest_deployment`), waktu relatif, link `:port/healthz` jika `running`
- **History** (icon jam) — buka dialog deployments tanpa membuat job baru (semua role)
- **Deploy** (icon rocket) — dialog pilih environment (dev/staging/prod) lalu create deployment; hanya `admin`/`developer` + service punya `workspace_path`
- Dialog deployments: badge environment + status; refresh + poll 3s saat job aktif; `onRefreshCatalog` sinkron kolom Deploy di tabel
- Catalog polling 5s (silent) saat ada service dengan status `pending`/`claimed`/`building`
- **Stop / Start** (square/play): hanya jika `latest_deployment` = `running` / `stopped`; via runtime job + agent
- **AppShell:** sidebar desktop collapse/expand (localStorage); **main content full width** (bukan `max-w-5xl`)
- Catalog table: kolom Service / Owner / Deploy / Origin / Actions; path workspace truncate CSS

### Template fix

- `templates/go-api/main.go.tmpl`: `ListenAndServe(":8080", mux)` (bukan `nil`) — tanpa ini `/healthz` di container selalu 404

### Runtime controls (R1 + R2)

- **R1** `apps/agent/internal/docker/container.go` — `ContainerName`, `Stop`, `Start`, `Remove` (idempotent jika container tidak ada)
- **R2** migrasi `00006_create_runtime_jobs.sql`; lapisan `model` → `store` → `service` → `handler`
- Portal enqueue: `POST /services/{id}/runtime/stop|start` → job `pending`
- Agent poll `POST /agent/runtime/next` → `docker stop|start` → `PATCH /agent/runtime/{id}` status `done`
- API saat job `done`: update deployment terkait ke `stopped` atau `running`
- Portal: tombol Stop (■) / Start (▶) di baris catalog; refresh otomatis setelah aksi

### Delete container cleanup (R3)

- Migrasi `00007_runtime_remove_action.sql` — `runtime_jobs.action` boleh `remove`
- `Catalog.Delete`: enqueue job `remove` (jika ada deployment) **sebelum** hapus row; lalu workspace cleanup
- Wiring: `catalog.SetCleanupEnqueuer(runtimeSvc)` di `main.go` (hindari circular dependency)
- Agent: `handleRuntime` case `remove` → `docker.Remove` (`rm -f`, idempotent)
- API `UpdateFromAgent` untuk `remove`: **tidak** update deployment (sudah CASCADE saat delete service)
- Portal: dialog delete menjelaskan cleanup container + workspace

### Latest deploy di catalog (R0)

- API: field `latest_deployment` di `model.Service` (omitempty)
- Store: `DeploymentsStore.LatestByServices()` — `DISTINCT ON (service_id)` order `created_at DESC`
- Service: `Catalog.List()` enrich via `DeploymentReader`; `NewCatalog(repo, deployments, workspaceDir)`
- Wiring: `deploymentsStore` dibuat **sebelum** `NewCatalog` di `main.go`

### Web UI (portal setelah login)

- **Layout:** `AppShell` — sidebar collapsible + topbar; main **full width** (`lg:px-8`, tanpa max-width container)
- **Sidebar sections:** Workspace (Overview), Platform (Catalog, Workers), Administration (Users — admin)
- **Routes (flat):** `/overview`, `/catalog`, `/worker`, `/users` (admin)
- **Catalog UX:** daftar + kolom deploy terakhir; Create/Deploy/Stop/Start/Edit/Delete hanya `admin`/`developer`; viewer read-only + History
- **Users:** admin create, ubah role, disable/enable (konfirmasi), reset password, soft-delete
- **Workers:** tabel dengan status badge, relative last seen
- **Overview:** metrik services/workers + panel recent
- **Shared:** `src/components/app/` (DataPanel, Toolbar, EmptyState, …)
- **Auth:** Inter Variable, `AuthLayout` centred (tetap terpisah dari app shell)

Template **masih di disk** (`templates/go-api/`), bukan DB. Service hasil scaffold punya `template_id` + `workspace_path` di Postgres.

Setelah `git pull` di mesin baru: `cd apps/web && npm install`

## Known debt (sengaja ditunda)

- **Template management** belum CRUD di DB/portal
- **Workspace lama** (path `/tmp/...`) tidak ikut terhapus saat delete (di luar root baru); scaffold ulang ke `data/workspaces`
- **Self-host API + agent host:** path workspace di DB adalah path container; agent host perlu API lokal untuk E2E deploy (atau solusi path-mapping nanti)

### Debt yang sudah diperbaiki

- Workspace default → `data/workspaces` (bukan `/tmp`), override `SAILORPORT_WORKSPACE`
- Agent Ctrl+C mengirim heartbeat `offline`
- Delete service menghapus folder workspace jika path di bawah workspace root
- **R3:** Delete service enqueue job `remove` → agent `docker rm -f sailorport-{name}`
- **R4:** Agent alokasi host port unik per container (bukan selalu 18080)

### Multi-port deploy (R4)

- Port unik **per Docker host** (bukan global fleet) — agent yang memilih
- Redeploy service **environment yang sama**: reuse port container `sailorport-{name}-{env}` jika masih ada
- Service/environment baru: port berikutnya yang tidak dipakai mapping `sailorport-*` dan tidak listen di host
- Pool habis → job `failed` dengan pesan rentang port
- Catalog sudah menampilkan `latest_deployment.port`

### Environments (13a–13d)

| Sub-step | Status | Isi |
|----------|--------|-----|
| 13a Migrasi + API | ✅ | `00010_create_environments.sql` seed dev/staging/prod; `GET /api/v1/environments` |
| 13b Deploy bound | ✅ | `00011_deployments_environment.sql`; create deploy dengan `environment` slug (default dev) |
| 13c Agent container | ✅ | `ContainerName(service, env)` → `sailorport-{service}-{env}`; runtime job bawa `environment_slug` |
| 13d Portal UI | ✅ | dialog Deploy pilih environment; badge slug di catalog + history deployments |

**Test deploy multi-env:**

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@sailorport.com","password":"changeme"}' | jq -r .token)

curl -s -X POST "http://localhost:8080/api/v1/services/SERVICE_ID/deployments" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"environment":"staging"}' | jq .

docker ps | grep sailorport-
# expect sailorport-{name}-dev and sailorport-{name}-staging on different ports
```

**Known limitation:** `latest_deployment` di catalog masih satu per service (bukan per env); stop/start runtime mengikuti deploy terbaru global — perbaikan opsional berikutnya.

**Debt runtime store:** perbaikan typo SQL di `store/runtime.go` (`deployments`, `COALESCE`, alias CTE).

### Compose full stack (Step 11)

- `deploy/compose/docker-compose.yml` — `postgres` + `api` + `web`
- `apps/api/Dockerfile`, `apps/web/Dockerfile`, `apps/web/nginx.conf` (proxy `/api` + `/healthz`)
- Volume: `../../templates` → `/templates` (bind); **workspaces** → named volume `sailorport_workspaces` (prod tanpa `chown` host)
- API startup: `EnsureWorkspaceDir` — gagal cepat jika folder tidak writable
- Agent **tidak** di compose (perlu Docker daemon di host untuk deploy workload)
- Deploy E2E agent: lebih cocok API host (`go run`) agar path workspace di DB = path host

## Next action

1. Opsional: runtime stop/start per environment (bukan latest global)
2. Opsional: container logs; multi-agent targeting (job prod → worker prod)

## Cara lanjut di mesin lain

Lihat `docs/CONTINUE.md` dan paste prompt dari `docs/RESUME-PROMPT.md`.
