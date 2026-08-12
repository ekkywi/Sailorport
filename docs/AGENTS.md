# Agent Context — Sailorport

Dokumen ini memberi konteks ke AI saat user membuka chat baru di mesin lain.

## Project

**Sailorport** — self-hosted internal developer platform OSS.

Tagline: *Self-hosted developer port — catalog, pave, and ship.*

Fitur inti: software catalog, golden path scaffold, environments, deploy via agent, worker health, secrets, RBAC, audit.

## Repo

- GitHub: `github.com/ekkywi/Sailorport`
- Go module API: `github.com/ekkywi/sailorport/apps/api`
- Go module Agent: `github.com/ekkywi/sailorport/apps/agent`

## Stack

| Komponen | Path | Status |
|----------|------|--------|
| Portal | `apps/web` | auth + app shell (full width) + overview/catalog/workers/users + Deploy/History UI + RBAC hide |
| API | `apps/api` | layered + scaffold/templates + workers + deployments |
| Worker | `apps/worker` | belum (job queue / orchestrator) |
| Agent | `apps/agent` | register + heartbeat + poll/deploy; Bearer agent token |
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

## Portal routes (setelah login)

| Path | Isi |
|------|-----|
| `/overview` | ringkasan services + workers |
| `/catalog` | daftar services + kolom deploy terakhir; create/edit/delete; History (dialog) vs Deploy (rocket) |
| `/worker` | daftar workers + status |

Auth: `/login`, `/register` — layout terpisah (`AuthLayout`).

## Catalog vs scaffold (mental model)

- **Create service** (default) = pilih template → generate workspace → daftar ke catalog
- **Register existing** = metadata saja di catalog, tanpa folder
- **Delete service** = hapus row DB + folder workspace jika path di bawah `SAILORPORT_WORKSPACE`

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

`docker compose up` → install agent → register worker → scaffold service → deploy → lihat status/logs.

(Saat ini: control plane `compose up --build`; agent token harden OK. Next: delete cleanup container / runtime controls / environments.)
