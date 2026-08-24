# Sailorport — Roadmap

Peta besar proyek. Detail step harian ada di `docs/PROGRESS.md`.  
**Visi produk & dua jalur deploy:** `docs/PRODUCT.md`.

## Visi produk

Sailorport = **self-hosted IDP** dengan **software catalog** sebagai pusat kontrol.

Dua jalur deploy (keduanya masuk catalog yang sama):

| Jalur | Siapa | Kontrak | Status |
|-------|-------|---------|--------|
| **Custom app** | Developer (repo sendiri) | Git + **Dockerfile** | ✅ Step 19 (19a–19d) |
| **Catalog app** | Platform (Postgres, Redis, …) | Image/manifest platform | Rencana Step 22+ |
| **Scaffold** (opsional) | Golden path demo | Template `go-api` → workspace | ✅ Ada |

## Visi v1 (MVP — selesai)

Flow yang sudah jalan:

1. `docker compose up` control plane (atau dev mode: Postgres + go run)
2. Agent di node Docker → register worker
3. Create service (scaffold, **From Git**, atau register metadata)
4. Deploy ke environment (worker policy Step 18; Git → agent clone/build)
5. Runtime stop/start, logs, audit

## Stack

| Lapisan | Teknologi |
|---------|-----------|
| Portal | React + TypeScript + Vite + shadcn/ui |
| API | Go |
| Worker | Go (proses terpisah) — belum |
| Agent | Go (binary di node) |
| Database | PostgreSQL |
| Queue/cache | Redis — belum |
| Auth | JWT lokal (OIDC menyusul) |
| Runtime v1 | Docker di worker node |
| Packaging | Docker Compose |

## MVP (wajib) — ✅ selesai

- [x] Auth + RBAC dasar (JWT lokal)
- [x] Software catalog (API + portal CRUD)
- [x] Scaffolder / golden path (`go-api`) — **opsional ke depan**
- [x] Environments (dev/staging/prod)
- [x] Deploy via agent (build workspace lokal)
- [x] Worker registry + multi-agent + labels/policy (Step 18)
- [x] Runtime controls + logs + audit
- [x] Self-hosting pack (compose + docs)

## Fase 2 — Deploy seperti industri (prioritas)

| Step | Topik | Status |
|------|-------|--------|
| 19 | Git-backed service (19a–19d: model, API, agent sync, portal) | ✅ selesai |
| 20 | Webhook auto-deploy (20a–20e) | ✅ selesai |
| 21 | Rollback / redeploy commit or tag | ⬜ planned |
| 22 | Catalog apps (Postgres, Redis, Gitea, …) | ⬜ planned |

Urutan: **21 → 22**.

## Fase 3 — Ops & polish

- Worker admin lite (edit labels, decommission)
- Stuck deploy reconcile + notifikasi
- Docs-as-code, CI visibility, scorecards ringan
- Secrets management
- OIDC auth

## Ditunda (jangan di awal)

- Repo sembarang **tanpa Dockerfile**
- Auto-detect stack (buildpack) — opsional jauh nanti
- Infra provisioning (Terraform)
- GitOps executor (Argo/Flux) penuh
- Policy engine (OPA)
- Multi-cluster Kubernetes deep

## Fase pembelajaran (Step 0–18) — ✅ selesai

| Step | Topik |
|------|-------|
| 0–9 | Foundation, catalog, auth |
| 10 | Worker + agent + deploy |
| 11–12 | Compose, users admin |
| 13–16 | Environments, runtime, logs, audit |
| 17 | Multi-agent targeting |
| 18 | Worker labels + deploy policy + portal filter |
| 19 | Git-backed deploy (19a–19d) |
| 20 | Webhook auto-deploy (20a–20e) |

Lihat checklist lengkap di `docs/PROGRESS.md`.

## Post-MVP worker ops (referensi)

| Fase | Fitur |
|------|--------|
| Ops | Admin lite (edit labels, decommission) |
| Ops | Draining; auto-pick worker; worker detail |
| Scale | Workspace sync multi-node; capacity/port pool |
