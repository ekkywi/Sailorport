# Agent Context — Sailorport

Dokumen ini memberi konteks ke AI saat user membuka chat baru di mesin lain.

## Project

**Sailorport** — self-hosted internal developer platform OSS.

Tagline: *Self-hosted developer port — catalog, pave, and ship.*

Fitur inti: software catalog, golden path scaffold, environments, deploy via agent, worker health, secrets, RBAC, audit.

## Repo

- GitHub: `github.com/ekkywi/Sailorport`
- Go module API: `github.com/ekkywi/sailorport/apps/api`

## Stack

| Komponen | Path | Status |
|----------|------|--------|
| Portal | `apps/web` | Step 6 selesai (list + create catalog) |
| API | `apps/api` | Step 5 selesai (CRUD services) + CORS |
| Worker | `apps/worker` | belum |
| Agent | `apps/agent` | belum |
| Shared contracts | `packages/shared` | belum |
| Compose | `deploy/compose` | Postgres jalan |

Infra: PostgreSQL (ada), Redis (belum), auth OIDC + RBAC (belum).

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
                 → Worker (jobs)
                 → Agent di node (Docker build/run)
                 → callback status
```

Prinsip: control plane tidak menjalankan container langsung; agent yang eksekusi.

## Coding conventions

- API routes: `/api/v1/...`
- Health check: `GET /healthz`
- `main.go` tipis — hanya wiring
- Logic di `internal/*`
- Portal: Vite + React + TypeScript; proxy `/api` ke `:8080`
- CORS: izinkan `http://localhost:5173` di development
- Commit: `feat(api):`, `feat(web):`, `docs:`, `chore:`, `fix(api):`, `fix(web):`

## Resume workflow

1. Baca `docs/PROGRESS.md` — cek "Step berikutnya"
2. Lanjutkan dari step itu, jangan ulang step selesai
3. Setelah step selesai, minta user update `docs/PROGRESS.md` + commit

## MVP v1 success criteria

`docker compose up` → install agent → register worker → scaffold service → deploy → lihat status/logs.
