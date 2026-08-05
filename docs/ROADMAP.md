# Sailorport — Roadmap

Peta besar proyek. Detail step harian ada di `docs/PROGRESS.md`.

## Visi v1 (definisi suukses)

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

- [ ] Auth + RBAC dasar
- [ ] Software catalog
- [ ] Scaffolder / golden path
- [ ] Environments (dev/staging/prod)
- [ ] Deploy via agent
- [ ] Worker registry + health
- [ ] Audit log
- [ ] Self-hosting pack (compose + docs)

## Setelah MVP stabil

- Git webhook auto-deploy
- Rollback / redeploy
- Runtime controls (start/stop/restart + logs)
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
| 6 | Portal web catalog | berikutnya |
| 7 | Auth OIDC + RBAC | - |
| 8 | Worker registry | - |
| 9 | Agent deploy end-to-end | - |
| 10 | Docker Compose full stack | - |
