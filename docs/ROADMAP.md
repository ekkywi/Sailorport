# Sailorport — Roadmap

Peta besar proyek. Detail step harian ada di `docs/PROGRESS.md`.

## Visi v1 (definisi sukses)

Satu flow ini harus jalan di mesin orang lain:

1. `docker compose up` control plane
2. Pasang agent di 1 VM ber-Docker
3. Register worker
4. Scaffold service dari template
5. Deploy ke environment
6. Lihat logs + status sukses

## Stack

| Lapisan | Teknologi |
|---------|-----------|
| Portal | React + TypeScript + Vite + shadcn/ui |
| API | Go |
| Worker | Go (proses terpisah) |
| Agent | Go (binary di node) |
| Database | PostgreSQL |
| Queue/cache | Redis |
| Auth | OIDC + RBAC |
| Runtime v1 | Docker di worker node |
| Packaging | Docker Compose dulu |

## Urutan belajar & build

1. Catalog
2. Scaffold (1 template)
3. Agent deploy end-to-end
4. Secrets + webhook + harden
5. Docs + scorecard + CI links
6. Executor tambahan (GitOps/K8s)

## MVP (wajib)

- [x] Auth + RBAC dasar (JWT lokal; OIDC menyusul)
- [x] Software catalog (API + portal CRUD lengkap)
- [x] Scaffolder / golden path (template `go-api`)
- [x] Environments (dev/staging/prod)
- [x] Deploy via agent (API + agent docker + portal Deploy UI)
- [x] Worker registry + health (API + portal + agent heartbeat)
- [x] Audit log
- [x] Multi-agent targeting (deploy + runtime affinity; portal worker picker)
- [x] Self-hosting pack (compose + docs)

## Setelah MVP stabil

- Worker capabilities (labels + deploy policy by environment) — Step 18
- Git webhook auto-deploy
- Rollback / redeploy
- ~~Runtime controls (start/stop/restart + logs)~~ ✅
- Docs-as-code
- CI visibility (GitHub/GitLab)
- Scorecards ringan
- Stuck deploy reconcile + notifikasi

## Ditunda (jangan di awal)

- Infra provisioning (Terraform)
- GitOps executor (Argo/Flux)
- Policy engine (OPA)
- Multi-cluster Kubernetes deep
- Plugin marketplace / AI assistant

## Fase pembelajaran (tutorial step-by-step)

| Step | Topik | Status |
|------|-------|--------|
| 0 | Struktur repo + Git | selesai |
| 1 | Go HTTP server + `/healthz` | selesai |
| 2 | Refactor + config + POST echo | selesai |
| 3 | PostgreSQL via Compose + koneksi DB | selesai |
| 4 | Migrasi tabel `services` | selesai |
| 5 | CRUD catalog API | selesai |
| 6 | Portal web catalog (list + create) | selesai |
| 7 | Portal: update + delete service | selesai |
| 7.5 | Architecture foundation | selesai |
| 8 | Scaffold / golden path (1 template) | selesai |
| 9 | Auth lokal + JWT + RBAC | selesai |
| 10A | Worker registry API + portal UI | selesai |
| 10B | Agent binary (register + heartbeat) | selesai |
| 10C.1 | Deployments API + agent claim | selesai |
| 10C.2 | Agent docker build/run | selesai |
| 10C.3 | Portal UI Deploy | selesai |
| 11 | Docker Compose full stack | selesai |
| 12 | Harden (agent token) | selesai |
| 12a | User management API | selesai |
| 12b | Portal users + RBAC UI | selesai |
| 12c | Admin create user | selesai |
| 12d | Admin disable / enable user | selesai |
| 12e | Admin reset password | selesai |
| 12f | Soft-delete user (rename email + `deleted_at`) | selesai |
| R0 | Latest deploy di catalog + layout full width | selesai |
| R1 | Agent docker stop/start/remove helpers | selesai |
| R2 | Runtime controls (API + portal stop/start) | selesai |
| R3 | Delete cleanup (stop/rm container) | selesai |
| R4 | Multi-port deploy (host port unik per container) | selesai |
| 13 | Environments (dev/staging/prod) | selesai |
| 14 | Runtime per environment (stop/start per env) | selesai |
| 15 | Logs end-to-end (API + agent + portal) | selesai |
| 16 | Audit log (record + list API + portal) | selesai |
| 17 | Multi-agent targeting (17a–17e) | selesai |
| 18 | Worker capabilities (labels + deploy policy) | ✅ selesai (18a–18c) |

## Post-MVP worker roadmap (referensi)

| Fase | Fitur |
|------|--------|
| 18 | Agent labels; API validasi deploy by env; portal filter |
| Ops | Admin lite (edit labels, decommission stale worker) |
| Ops | Draining; auto-pick worker; worker detail (deployments on node) |
| Scale | Workspace sync multi-node; capacity/port pool; region labels |
