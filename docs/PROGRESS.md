# Sailorport — Progress

> Update file ini setiap selesai 1 step. Ini sumber kebenaran saat pindah mesin.

## Status saat ini

- **Step selesai:** 20b — `POST /api/v1/webhooks/github` (parse + acknowledge)
- **MVP core:** selesai (catalog, scaffold, deploy agent, env, runtime, logs, audit, multi-agent)
- **Step berikutnya:** 20c — validasi signature HMAC; lalu 20d match repo → deploy; 20e portal secret/toggle
- **Terakhir dikerjakan:** 2026-08-21 — Step 20a–20b (webhook data model + public endpoint)
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
- [x] Step 14 — runtime per environment (14a–14d)
- [x] Step 15 — logs end-to-end (15a–15e: API + agent + portal)
- [x] Step 16 — audit log (16a–16e: record + list API + portal)
- [x] Step 17 — Multi-agent targeting (17a–17e)
- [x] Step 18 — Worker capabilities (18a–18c: labels + deploy policy + portal filter)
- [x] Step 19a — Service Git fields (migration + model + store)
- [x] Step 19b — API create/update Git source + deploy gate
- [x] Step 19c — Agent git clone/pull + docker build
- [x] Step 19d — Portal Add from Git
- [x] Step 20a — Webhook fields on services (migration + catalog defaults)
- [x] Step 20b — Public `POST /api/v1/webhooks/github` (parse push, ack)
- [ ] Step 20c — Validate `X-Hub-Signature-256`
- [ ] Step 20d — Match repo/branch → create deployment
- [ ] Step 20e — Portal webhook secret + auto-deploy toggle

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
| `GET/POST/PUT/DELETE /api/v1/services` | viewer+ / developer+ | catalog CRUD; `GET` list: `latest_deployment` (global terbaru) + `env_deployments` (map slug → deploy terbaru per env) |
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
| `POST /api/v1/services/{id}/runtime/logs` | viewer+ | enqueue logs job (deployment harus `running`/`stopped`) → **202** |
| `GET /api/v1/runtime/{id}` | viewer+ | get runtime job by ID (poll status + output) |
| `GET /api/v1/audit` | admin | list audit events (`?limit=50`, max 200) |
| `POST /api/v1/agent/jobs/next` | agent token | claim 1 deploy job pending → `claimed` (204 jika kosong) |
| `PATCH /api/v1/agent/deployments/{id}` | agent token | agent update deploy status |
| `POST /api/v1/agent/runtime/next` | agent token | claim 1 runtime job (`stop`/`start`/`logs`) |
| `PATCH /api/v1/agent/runtime/{id}` | agent token | agent selesai runtime job; API update deployment → `stopped`/`running`; logs → output only |
| Portal `/login`, `/register` | — | auth gate |
| Portal `/overview`, `/catalog`, `/worker`, `/users`, `/audit` | JWT | app shell; `/users` dan `/audit` admin-only (redirect non-admin) |

Env API: `AUTH_JWT_SECRET` (default `dev-only-change-me`), `SAILORPORT_AGENT_TOKEN` (default `dev-agent-token` — ganti di production)

Env agent:

| Variable | Default | Keterangan |
|----------|---------|------------|
| `SAILORPORT_API_URL` | `http://localhost:8080` | base URL API |
| `SAILORPORT_AGENT_TOKEN` | `dev-agent-token` | harus sama dengan API |
| `SAILORPORT_WORKER_NAME` | hostname mesin | nama worker di registry |
| `SAILORPORT_WORKER_TIER` | — | label `tier` (contoh `nonprod`, `prod`) |
| `SAILORPORT_WORKER_ENVIRONMENTS` | — | label `environments` comma-separated (`dev,staging`) |
| `SAILORPORT_WORKER_LABELS` | — | JSON extra labels (merge; TIER/ENVIRONMENTS override) |
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
- **Stop / Start** (square/play): hanya jika `latest_deployment` = `running` / `stopped`; konfirmasi dialog; via runtime job + agent
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
- Portal: tombol Stop (■) / Start (▶) di baris catalog; **konfirmasi** sebelum enqueue; refresh otomatis setelah aksi

### Delete container cleanup (R3)

- Migrasi `00007_runtime_remove_action.sql` — `runtime_jobs.action` boleh `remove`
- Migrasi `00014_runtime_jobs_remove_survive.sql` — kolom `environment_slug`; `deployment_id` nullable + `ON DELETE SET NULL` (bukan CASCADE)
- `Catalog.Delete`: enqueue job `remove` (slug tersimpan di job) **sebelum** hapus row; hapus service tidak menghapus job
- `ClaimNext` / `Get` / `Update` tidak `JOIN deployments` — agent tetap bisa claim `remove` setelah catalog hilang
- Wiring: `catalog.SetCleanupEnqueue(runtimeSvc)` di `main.go` (hindari circular dependency)
- Agent: `handleRuntime` case `remove` → `docker.Remove` (`rm -f sailorport-{name}-{env}`, idempotent)
- API `UpdateFromAgent` untuk `remove`: **tidak** update deployment (row sudah terhapus)
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
- **Catalog UX:** daftar + kolom deploy terakhir; Create/Deploy/Stop/Start (konfirmasi)/Edit/Delete hanya `admin`/`developer`; viewer read-only + History
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

**Delete protection (14d):** delete diblokir jika ada environment `running` (HTTP 409). Jika **prod** masih running → HTTP 403. Cleanup enqueue `remove` untuk **setiap** latest deploy per env.

### Step 14 — runtime per environment (Opsi A)

| Sub-step | Status | Isi |
|----------|--------|-----|
| 14a API | ✅ | `env_deployments` di `GET /services` — `LatestPerEnvByServices` (`DISTINCT ON service+env`) |
| 14b API | ✅ | `POST .../runtime/stop|start` + body `{"environment":"staging"}` |
| 14c Portal | ✅ | Kolom Deploy: 3 baris env + stop/start per env |
| 14d Delete | ✅ | Enqueue `remove` per env; blokir delete jika ada env running (prod → 403) |

**Tes 14a:**

```bash
curl -s http://localhost:8080/api/v1/services -H "Authorization: Bearer $TOKEN" \
  | jq '.[0] | {name, latest: .latest_deployment.environment_slug, env_deployments: (.env_deployments | keys)}'
```

**Tes 14d:**

```bash
# running env → 409; prod running → 403; semua stopped → 204 + job remove per env
curl -s -o /dev/null -w "%{http_code}\n" -X DELETE \
  "http://localhost:8080/api/v1/services/SERVICE_ID" \
  -H "Authorization: Bearer $TOKEN"

# job remove harus tetap ada (deployment_id NULL) sampai agent selesai
docker exec -it sailorport-postgres psql -U sailorport -d sailorport -c \
  "SELECT action, status, environment_slug, deployment_id IS NULL AS dep_gone FROM runtime_jobs WHERE action='remove' ORDER BY created_at DESC LIMIT 5;"

# setelah agent poll (~5s): container hilang
docker ps -a --filter name=sailorport-
```

### Runtime queue hardening (post-13d)

- Migrasi `00012_runtime_jobs_deployment_fk.sql` — FK `deployment_id` (semula CASCADE; di 00014 jadi SET NULL)
- `HasActiveJob` — tolak stop/start/remove ganda saat job masih `pending`/`claimed`
- Portal: konfirmasi AlertDialog sebelum Stop / Start
- **Fix 00014:** job `remove` tidak ikut terhapus saat service di-delete; slug disimpan di `runtime_jobs.environment_slug`

### Compose full stack (Step 11)

- `deploy/compose/docker-compose.yml` — `postgres` + `api` + `web`
- `apps/api/Dockerfile`, `apps/web/Dockerfile`, `apps/web/nginx.conf` (proxy `/api` + `/healthz`)
- Volume: `../../templates` → `/templates` (bind); **workspaces** → named volume `sailorport_workspaces` (prod tanpa `chown` host)
- API startup: `EnsureWorkspaceDir` — gagal cepat jika folder tidak writable
- Agent **tidak** di compose (perlu Docker daemon di host untuk deploy workload)
- Deploy E2E agent: lebih cocok API host (`go run`) agar path workspace di DB = path host

### Logs end-to-end (Step 15a–15e)

| Sub-step | Status | Isi |
|----------|--------|-----|
| 15a Migrasi | ✅ | `00013_runtime_jobs_logs.sql` — action `logs` + kolom `output TEXT` |
| 15b API | ✅ | `POST .../runtime/logs` enqueue; `GET /runtime/{id}` poll; agent update menyimpan output |
| 15c Agent | ✅ | `docker.Logs(container, 200)` → truncate 64KB → update job `done` + output |
| 15d Store | ✅ | `Update` SQL persist `output` via `COALESCE(NULLIF($4,''), output)` |
| 15e Portal | ✅ | `LogsDialog` + tombol `ScrollText` per env (running/stopped); semua role; POST → poll GET → `<pre>` |

**Tes logs:**

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@sailorport.com","password":"changeme"}' | jq -r .token)

JOB=$(curl -s -X POST "http://localhost:8080/api/v1/services/SERVICE_ID/runtime/logs" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"environment":"dev"}' | jq -r .id)

sleep 8
curl -s "http://localhost:8080/api/v1/runtime/$JOB" \
  -H "Authorization: Bearer $TOKEN" | jq '{status, output: (.output | length)}'
# expect: status "done", output > 0
```

Portal: klik icon 📜 (ScrollText) di baris env → dialog logs → "Waiting for agent…" → output muncul.

### Audit log (Step 16a–16e)

| Sub-step | Status | Isi |
|----------|--------|-----|
| 16a Migrasi | ✅ | `00015_create_audit_events.sql` — tabel append-only + index |
| 16b Catalog delete | ✅ | Snapshot service sebelum hard-delete; actor dari JWT |
| 16c Events lain | ✅ | `service.create`/`update`, `user.create`/`role`/`disable`/`enable`/`password_reset`/`delete` |
| 16d API list | ✅ | `GET /api/v1/audit?limit=50` (admin) |
| 16e Portal | ✅ | `/audit` — tabel when/actor/action/resource/details; sidebar admin |

**Actions yang dicatat:** catalog create/update/delete (scaffold → create), admin user CRUD-ish. **Tidak** dicatat: login, list, agent poll.

**Tes audit:**

```bash
curl -s "http://localhost:8080/api/v1/audit?limit=10" \
  -H "Authorization: Bearer $TOKEN" | jq '.[0] | {action, resource_name, actor_email}'
```

Portal: login admin → sidebar **Audit** → lihat jejak aksi.

### Step 17 — Multi-agent targeting ✅

**Tujuan:** deploy dan runtime job hanya dijalankan oleh worker yang benar (bukan race).

| Sub-step | Status | Isi |
|----------|--------|-----|
| 17a Migrasi + model + store | ✅ | `00016_multi_agent_targeting.sql` (`target_worker_id` di deployments + runtime_jobs); model struct + scanner updated |
| 17b Deploy with worker | ✅ | `CreateDeploymentRequest.WorkerID`; `WorkersStore.Get` + `Workers.Get`; `resolveTargetWorker` (explicit/redeploy affinity/null); validasi online → 409; `writeDeploymentError` handle `ErrConflict` |
| 17c Claim deploy filter | ✅ | `ClaimNext` WHERE `target_worker_id IS NULL OR = $1` |
| 17d Runtime affinity | ✅ | Enqueue set target dari deployment.worker_id; `ClaimNext` runtime filter |
| 17e Portal worker picker | ✅ | DeployDialog — Any available + online workers; kirim `worker_id` ke API |

**Tes deploy dengan worker:**

```bash
# deploy tanpa worker (target = null atau affinity)
curl -s -X POST "http://localhost:8080/api/v1/services/$SVC/deployments" \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"environment":"dev"}' | jq '{target_worker_id, status}'

# deploy dengan worker explicit
curl -s -X POST "http://localhost:8080/api/v1/services/$SVC/deployments" \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d "{\"environment\":\"dev\",\"worker_id\":\"$WORKER_ID\"}" | jq '{target_worker_id, status}'
```

Portal: Deploy dialog → pilih environment + worker (Any available atau worker online). UI worker: search (≥4 online), scroll `max-h-52`, fallback jika worker prior offline.

### Checkpoint — Worker model (keputusan produk)

- Worker **self-register** via agent (`POST /workers/register`); **tidak** perlu admin CRUD tambah/hapus worker di MVP.
- **Environment** (dev/staging/prod) ≠ **Worker** (node Docker); satu worker boleh menampung banyak env (container terpisah).
- Pola infra umum: nonprod VM (dev+staging) + prod VM terpisah — didukung dengan **labels** (Step 18), bukan 1 worker = 1 env.
- Portal `/worker` = **monitoring** (read-only); admin edit labels / decommission = post-MVP.

### Step 18 — Worker capabilities ✅

| Sub-step | Status | Isi |
|----------|--------|-----|
| 18a Agent labels | ✅ | `parseWorkerLabels()`: `role=agent` + JSON `SAILORPORT_WORKER_LABELS` + override `TIER`/`ENVIRONMENTS`; register kirim `cfg.Labels` |
| 18b Deploy policy API | ✅ | `worker_env.go` + `validateWorkerForDeploy` di deploy; 409 jika env tidak diizinkan labels |
| 18c Portal filter | ✅ | `DeployDialog` filter worker by environment; Workers page kolom Tier + Environments |

**Contoh labels:**

```json
{ "role": "agent", "tier": "nonprod", "environments": "dev,staging" }
{ "role": "agent", "tier": "prod", "environments": "prod" }
```

**Tes 18a:**

```bash
cd apps/agent
SAILORPORT_WORKER_NAME=nonprod-01 \
SAILORPORT_WORKER_TIER=nonprod \
SAILORPORT_WORKER_ENVIRONMENTS=dev,staging \
go run .
```

Portal `/worker` atau `GET /api/v1/workers` harus menampilkan `tier` + `environments`. Restart agent dengan env berbeda → labels ter-update (upsert by name).

**Tip env file (opsional):** agent tidak auto-load `.env`; buat file lokal (gitignored) lalu `source`:

```bash
# apps/agent/.env.nonprod — export SAILORPORT_* …
source .env.nonprod && go run .
```

**Tes 18b (API policy):**

```bash
# staging + worker nonprod → 201
curl -s -X POST "http://localhost:8080/api/v1/services/$SVC/deployments" \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d "{\"environment\":\"staging\",\"worker_id\":\"$WORKER_ID\"}" | jq '{status, target_worker_id}'

# prod + worker nonprod → 409
curl -s -o /tmp/out -w "%{http_code}\n" \
  -X POST "http://localhost:8080/api/v1/services/$SVC/deployments" \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d "{\"environment\":\"prod\",\"worker_id\":\"$WORKER_ID\"}"
```

**Tes 18c (portal):** Deploy dialog → pilih env → hanya worker yang allow env itu; Workers page → kolom Tier + Environments.

### Step 20 — Webhook auto-deploy (in progress)

| Sub-step | Status | Isi |
|----------|--------|-----|
| 20a Data model | ✅ | Migrasi `00018_service_webhook_autodeploy.sql`; model/store; catalog default `staging` + merge secret; web types; scaffold defaults |
| 20b Public endpoint | ✅ | `POST /api/v1/webhooks/github` — parse push (`repo`, `branch`, `commit`); ignore ping/tags; **no JWT** |
| 20c Signature | ⬜ | `X-Hub-Signature-256` HMAC-SHA256 |
| 20d Match → deploy | ⬜ | Cari service by repo/branch; create deployment jika `auto_deploy_enabled` |
| 20e Portal UI | ⬜ | Set secret + toggle + target environment |

**Tes 20b:**

```bash
curl -sS -X POST http://localhost:8080/api/v1/webhooks/github \
  -H 'Content-Type: application/json' \
  -H 'X-GitHub-Event: push' \
  -d '{"ref":"refs/heads/main","after":"deadbeef","repository":{"full_name":"acme/demo","clone_url":"https://github.com/acme/demo.git"},"pusher":{"name":"you"}}'
```

### Checkpoint — Product vision (2026-08-20)

Diskusi positioning produk (detail: **`docs/PRODUCT.md`**):

1. **Sailorport tetap IDP** — catalog = inventory pusat; bukan bergeser ke non-IDP.
2. **Jalur utama ke depan:** custom app dari **Git + Dockerfile** (developer setup mandiri; IDP untuk deploy).
3. **Jalur sekunder:** **catalog apps** siap pakai (Postgres, Redis, Gitea, …) — pull image, tanpa build.
4. **Scaffold `go-api`:** tetap ada sebagai **golden path opsional**, bukan syarat deploy.
5. **Webhook / rollback:** masuk **setelah** Git deploy (Step 19), bukan sebelum kontrak repo jelas.

**Yang belum di kode:** webhook signature + auto-deploy (20c–20e), rollback, catalog apps, private Git credentials.

## Rencana step berikutnya (belum dikerjakan)

| Step | Topik | Isi singkat |
|------|-------|-------------|
| 19a | Model + migrasi Git fields | ✅ `source_type`, `repo_url`, `branch`, `dockerfile_path` |
| 19b | API create/update Git | ✅ validasi `git`+`repo_url`; Update merge; deploy allow git atau workspace |
| 19c | Agent git sync | ✅ clone/pull → build → run |
| 19d | Portal Add from Git | ✅ `GitServiceForm`; Origin git di list; deploy button untuk git |
| 20a | Webhook data model | ✅ `webhook_secret`, `auto_deploy_enabled`, `auto_deploy_environment` |
| 20b | Webhook HTTP endpoint | ✅ `POST /api/v1/webhooks/github` parse + ack (belum signature/deploy) |
| 20c | Signature validation | HMAC-SHA256 `X-Hub-Signature-256` |
| 20d | Match → deploy | Match `clone_url`/`branch` → create deployment jika auto-deploy on |
| 20e | Portal webhook UI | Set secret + toggle auto-deploy + target env |
| 21 | Rollback / redeploy | Pin commit/tag; redeploy versi sebelumnya dari UI/API |
| 22 | Catalog apps | Manifest app (image, env, volume); deploy tanpa Git |
| — | Worker admin lite | Edit labels, decommission stale worker (post-MVP) |

## Next action

1. **Step 20c** — validasi signature webhook
2. **Step 20d–20e** — match service + create deployment + portal UI
3. **Step 21** — rollback / redeploy
4. Catalog apps (22) — setelah custom path matang
5. Opsional: edit Git fields di portal; private repo credentials

## Cara lanjut di mesin lain

1. `git pull`
2. Baca `docs/PRODUCT.md` (visi) + `docs/PROGRESS.md` (step)
3. Paste `docs/RESUME-PROMPT.md` ke chat Cursor baru
4. Lihat `docs/CONTINUE.md` untuk ritual Git
