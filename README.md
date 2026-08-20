# Sailorport

Self-hosted internal developer platform — **catalog, deploy, and ship**.

## Apa ini

Sailorport adalah **IDP self-hosted** untuk mendaftar, deploy, dan mengoperasikan service di infra sendiri:

- **Catalog** — inventory pusat semua service (custom app, infra app, hasil scaffold)
- **Deploy via agent** — build/run container di worker node (Docker)
- **Environments** — dev / staging / prod, worker policy, logs, audit, RBAC

Dua jalur deploy (rencana — detail di [`docs/PRODUCT.md`](docs/PRODUCT.md)):

1. **Custom app** — repo developer + Dockerfile → git pull → build → deploy *(Step 19+)*
2. **Catalog apps** — Postgres, Redis, Gitea, … one-click *(Step 22+)*
3. **Scaffold** (opsional) — starter dari template `go-api` *(sudah ada)*

## Struktur

```text
apps/web          portal — catalog, deploy, workers, users, audit
apps/api          control plane — layered API
apps/agent        agent di node — register, heartbeat, docker build/run
apps/worker       background jobs — belum
templates/        golden path go-api (opsional)
deploy/compose    Docker Compose (Postgres :5433)
docs/             progress, product vision, architecture, panduan lanjut
```

## Quick start (lokal)

```bash
# 1. Database
cd deploy/compose && docker compose up -d postgres

# 2. API
cd apps/api && go run .

# 3. Portal (terminal lain)
cd apps/web && npm install && npm run dev

# 4. Agent (terminal lain, setelah API jalan)
cd apps/agent
cp .env.example .env.nonprod   # sesuaikan, lalu:
source .env.nonprod && go run .
```

- API: `http://localhost:8080/healthz`
- Portal: `http://localhost:5173`

## Lanjut di mesin lain

Chat Cursor tidak ikut pindah antar device. Yang ikut: **repo Git + folder `docs/`**.

| Situasi | Baca |
|---------|------|
| **Visi produk & arah fitur** | [`docs/PRODUCT.md`](docs/PRODUCT.md) |
| Pindah ke laptop/rumah | [`docs/CONTINUE.md`](docs/CONTINUE.md) |
| Setup tool di mesin baru | [`docs/SETUP.md`](docs/SETUP.md) |
| Step mana yang sudah selesai | [`docs/PROGRESS.md`](docs/PROGRESS.md) |
| Buka chat Cursor baru | [`docs/RESUME-PROMPT.md`](docs/RESUME-PROMPT.md) |
| Peta besar proyek | [`docs/ROADMAP.md`](docs/ROADMAP.md) |

**Ritual sebelum tutup:** update `docs/PROGRESS.md` → `git commit` → `git push`

## Status

**MVP core selesai** (Step 0–18): catalog, scaffold, deploy agent, environments, runtime, logs, audit, multi-agent, worker labels/policy.

**Berikutnya:** Step 19 — Git-backed deploy (`repo_url` + Dockerfile). Lihat [`docs/ROADMAP.md`](docs/ROADMAP.md).

Arsitektur: [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md)
