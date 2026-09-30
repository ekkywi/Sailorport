# Sailorport — Roadmap

Peta besar dan backlog. Detail harian: `docs/PROGRESS.md`. Visi: `docs/PRODUCT.md`. Peta docs: `docs/README.md`.

## Visi produk

Sailorport = **self-hosted IDP** dengan **software catalog** sebagai pusat kontrol.

Dua jalur deploy (keduanya masuk catalog yang sama):

| Jalur | Siapa | Kontrak | Status |
|-------|-------|---------|--------|
| **Custom app** | Developer (repo sendiri) | Git + **Dockerfile** | ✅ Step 19 (19a–19d) |
| **Catalog app** | Platform (Postgres, Redis, …) | Image/manifest platform | Step 22 (22a–22b ✅) |
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
| 21 | Rollback / redeploy commit or tag | ✅ done |
| 22 | Catalog apps (Postgres, Redis, Gitea, …) | ✅ 22a–22f done |

Urutan berikutnya: opsional — backlog Tabel A/C atau Known debt Medium. Pass A/B/C Critical+High ✅ 2026-09-28. Step 33–35 hardening ✅.

## Saran pengembangan ke depan (backlog ide)

> Dicatat 2026-09-23 agar chat/mesin baru tidak kehilangan state.  
> **Step 33 ✅ (2026-09-24):** webhook delivery dedupe complete.  
> Item lain tetap kandidat; jangan kerjakan bersamaan. Detail Known debt: `docs/QC.md`.

### Cara pakai

1. Baca visi di `docs/PRODUCT.md` (jangan langgar keputusan produk).
2. Item **dikunci** → kerjakan lewat sub-step di `PROGRESS.md`.
3. Item lain: pilih setelah Step 33 selesai, lalu pecah jadi Step 34a…
4. Centang/ubah status di tabel saat selesai.

### A — Fitur produk (nilai user jelas)

| Ide | Kenapa | Estimasi | Status |
|-----|--------|----------|--------|
| Catalog app **Gitea** (atau MinIO / AdGuard) | Jalur sekunder siap; pola `command` + env sudah ada | ~1–2 step kecil | kandidat |
| **Private Git** (deploy key / token) | QC debt: sekarang hanya public clone | medium | ✅ Step 37 complete — HTTPS PAT; SSH deploy key menyusul |
| **Volume persist** catalog apps | Postgres/Redis hilang data saat recreate container | medium (agent + manifest) | ✅ Step 36 complete — generic `volumes` in manifest |
| **Health / open URL** di catalog | Port sudah ada; UX klik buka app | kecil (web) | ✅ Step 38 complete |
| **Notifikasi deploy gagal** (audit + badge/toast) | Audit ada; kurang sinyal ke user | kecil–medium | kandidat |
| **Service detail page** | Catalog padat; butuh halaman satu service | medium (web) | kandidat |
| **Invite user** (email / invite link) | Sekarang admin set password manual | medium | kandidat |

### B — Hardening (siap demo / TA / expose)

| Ide | Kenapa | Status |
|-----|--------|--------|
| **Webhook dedupe** (`X-GitHub-Delivery`) | Cegah double deploy dari replay | ✅ Step 33 complete |
| **Login rate limit** | Brute force murah | ✅ Step 34 complete |
| **Agent identity lebih ketat** | Shared token + `worker_id` self-reported on **claim** | **Partial:** PATCH must match claimer ✅ Pass B. **Remaining (candidate):** per-worker token / bind identity on claim |
| **Pass B** lalu **Pass C** QC | Agent + portal formal review | Pass B ✅ + Pass C ✅ 2026-09-28 |
| **CORS PATCH / origin** | Pecah jika portal tidak lewat proxy | ✅ Step 35 complete |

### C — Modul / polish teknis

| Ide | Yang dipelajari | Status |
|-----|-----------------|--------|
| Unit test agent `git.Sync` (temp repo) | Test tanpa DB penuh | kandidat |
| Pagination + search catalog / directory | Query + UI | kandidat |
| OpenAPI / typed client | Kontrak API | kandidat |
| Template kedua (mis. `node-api`) | Scaffold path | kandidat |
| Interceptor **401 logout** portal | Session hygiene | ✅ Pass C complete |

### D — Sengaja jangan dulu

Kecuali keputusan produk baru (diskusi dulu):

- Kubernetes / Helm penuh
- Multi-tenant SaaS
- Full Backstage plugin
- CI runner sendiri (bukan mengganti GitHub Actions)
- Repo tanpa Dockerfile / buildpack auto-detect (sudah di **Ditunda**)

### Panduan pilih cepat

| Prioritas | Sarankan mulai dari |
|-----------|---------------------|
| Fitur terlihat | Catalog app Gitea **atau** service detail page **atau** notifikasi deploy gagal |
| IDP lebih lengkap | Invite user **atau** service detail page |
| Keamanan sisa | Agent identity (per-worker token) **atau** webhook rate limit |

## Fase 3 — Ops & polish

- [x] Worker admin lite (edit labels, decommission) — Step 28
- Stuck deploy reconcile + notifikasi *(juga ada di tabel A di atas)*
- Docs-as-code, CI visibility, scorecards ringan
- Secrets management *(encrypt catalog env Step 25 ✅; OIDC / vault menyusul)*
- OIDC auth

## Ditunda (jangan di awal)

- Repo sembarang **tanpa Dockerfile**
- Auto-detect stack (buildpack) — opsional jauh nanti
- Infra provisioning (Terraform)
- GitOps executor (Argo/Flux) penuh
- Policy engine (OPA)
- Multi-cluster Kubernetes deep

## Early milestones (Step 0–18) — ✅ complete

Historical reference. Full checklist: `docs/PROGRESS.md`.

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
