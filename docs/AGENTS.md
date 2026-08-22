# Agent Context — Sailorport

Dokumen ini memberi konteks ke AI saat user membuka chat baru di mesin lain.

## Project

**Sailorport** — self-hosted internal developer platform OSS.

Tagline: *Self-hosted developer port — catalog, deploy, and ship.*

Fitur inti: **software catalog** (inventory pusat), deploy via agent, environments, worker health/policy, RBAC, audit. Golden path scaffold = **opsional** (bukan jalur utama).

**Visi produk lengkap:** `docs/PRODUCT.md` — baca sebelum fitur Git deploy / catalog apps.

## Repo

- GitHub: `github.com/ekkywi/Sailorport`
- Go module API: `github.com/ekkywi/sailorport/apps/api`
- Go module Agent: `github.com/ekkywi/sailorport/apps/agent`

## Stack

| Komponen | Path | Status |
|----------|------|--------|
| Portal | `apps/web` | auth + app shell + catalog deploy/runtime + workers + users (create/role/disable/reset/soft-delete) + RBAC |
| API | `apps/api` | layered + scaffold/templates + workers + deployments |
| Worker | `apps/worker` | belum (job queue / orchestrator) |
| Agent | `apps/agent` | register + heartbeat + poll deploy/runtime (stop/start/remove); host port unik; Bearer agent token |
| Templates | `templates/` | `go-api` (disk, bukan DB) |
| Shared contracts | `packages/shared` | belum |
| Compose | `deploy/compose` | Postgres + API + web; workspaces = named volume (no host chown) |

Infra: PostgreSQL (ada), Redis (belum), auth JWT lokal (ada), OIDC (belum).

## Learning mode (PENTING)

User **belum bisa coding** — belajar sambil mengetik manual.

Aturan panduan:

1. Satu step = satu hasil yang bisa dijalankan dan ditest
2. Beri kode lengkap per file, bukan potongan acak
3. Jelaskan baris per baris untuk konsep baru
4. Jangan refactor besar atau tambah fitur di luar step
5. Akhiri setiap step dengan: cara test + commit message

## Arsitektur

```
Developer → Web Portal
              → API (Go)
                 → PostgreSQL + Redis
                 → Worker (jobs) [belum]
                 → Agent di node (register, heartbeat, poll job, docker build/run)
                 → callback status
```

Prinsip: control plane tidak menjalankan container langsung; agent yang eksekusi.

## Workers (runtime nodes)

- Worker = node dengan Docker + agent; **bukan** sama dengan environment dev/staging/prod.
- Satu worker boleh menjalankan banyak environment (container `sailorport-{service}-{env}` terpisah).
- Data worker dari **agent register + heartbeat** — tidak ada admin CRUD create/delete worker di MVP.
- Kolom `labels` (JSONB): agent kirim saat register (Step 18a). Env `SAILORPORT_WORKER_TIER`, `SAILORPORT_WORKER_ENVIRONMENTS`, optional `SAILORPORT_WORKER_LABELS` JSON. Deploy policy (18b): API 409 jika env tidak diizinkan labels. Portal (18c): DeployDialog filter worker; Workers page kolom Tier/Environments.
- Deploy: optional `worker_id`; `target_worker_id` + claim filter memastikan job ke node yang benar.

## Portal routes (setelah login)

| Path | Isi |
|------|-----|
| `/overview` | ringkasan services + workers |
| `/catalog` | daftar services + deploy terakhir; History; Deploy; **Stop/Start** (runtime); create/edit/delete |
| `/worker` | daftar workers + status (read-only; self-register via agent) |
| `/users` | admin: list, create, role, disable/enable (confirm), reset password, soft-delete |
| `/audit` | admin: jejak aksi (catalog + user admin) |

Auth: `/login`, `/register` — layout terpisah (`AuthLayout`).

## Catalog — mental model (penting)

**Catalog = daftar semua service yang dikelola platform** (bukan “template store”).

| Cara masuk catalog | Status | Deploy |
|--------------------|--------|--------|
| Scaffold template (`go-api`) | ✅ ada | Build `workspace_path` lokal |
| Register existing (metadata) | ✅ ada | Belum auto-deploy |
| Git repo + Dockerfile | ✅ Step 19 | clone/pull → build → run |
| Catalog app (Postgres, Redis, …) | ⬜ Step 22+ | pull image → run |

Semua jalur berakhir di **satu UI `/catalog`** — deploy, env, logs, runtime sama.

## Scaffold (golden path — opsional)

- **Create service** (scaffold) = pilih template → generate `data/workspaces/{name}/` → daftar catalog
- Developer **mengembangkan kode di workspace** setelah scaffold; template dipakai **sekali** di awal
- **Register existing** = metadata saja, tanpa folder
- **Delete service** = enqueue job `remove` → hapus row DB + folder workspace → agent `docker rm`

## Coding conventions

- Baca `docs/ARCHITECTURE.md` sebelum menambah fitur
- API routes: `/api/v1/...`
- Health check: `GET /healthz`
- `main.go` tipis — hanya wiring
- Alur API: `handler` → `service` → `store` (jangan bypass service untuk domain logic)
- Error API JSON: `{"error":"..."}`
- Portal: fitur di `src/features/<domain>/`; `App.tsx` hanya shell + routing
- CORS: izinkan `http://localhost:5173` di development
- Commit: `feat(api):`, `feat(web):`, `feat(agent):`, `docs:`, `fix:`

## Resume workflow

1. Baca `docs/PROGRESS.md` — cek "Step berikutnya"
2. Lanjutkan dari step itu, jangan ulang step selesai
3. Setelah step selesai, minta user update `docs/PROGRESS.md` + commit

## MVP v1 success criteria

`docker compose up` → agent → worker online → service di catalog → deploy → status/logs.

**Selesai:** MVP core + Step 18–19 Git + Step 20a–20d webhook auto-deploy. **Next:** Step 20e portal UI (`docs/PRODUCT.md`, `docs/ROADMAP.md`).
