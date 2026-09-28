# Sailorport

Self-hosted **internal developer platform (IDP)** — **catalog, deploy, and ship**.

Sailorport helps teams register services, deploy them to their own infrastructure via a node agent, and operate them across environments — with RBAC, audit, logs, and worker policy.

## Features

- **Software catalog** — central inventory (custom apps, catalog apps, optional scaffold)
- **Agent-based deploy** — clone/build/run (or pull image) on Docker worker nodes
- **Environments** — `dev` / `staging` / `prod`, runtime controls, logs, audit
- **Git + webhook** — Git-backed services, auto-deploy on push, redeploy by commit SHA
- **Catalog apps** — platform images (e.g. Postgres, Redis) with managed env/command

Primary product decisions: [`docs/PRODUCT.md`](docs/PRODUCT.md).

## Repository layout

```text
apps/web          Portal (React + TypeScript)
apps/api          Control plane API (Go)
apps/agent        Node agent — register, heartbeat, git sync, Docker
apps/worker       Background jobs (not implemented yet)
templates/        Optional golden-path template (`go-api`)
deploy/compose    Docker Compose pack
docs/             Product, architecture, progress, QC, setup
```

## Quick start (local)

```bash
# 1. Database
cd deploy/compose && docker compose up -d postgres

# 2. API
cd apps/api && go run .

# 3. Portal (another terminal)
cd apps/web && npm install && npm run dev

# 4. Agent (another terminal, after API is up)
cd apps/agent
cp .env.example .env.nonprod   # edit as needed
source .env.nonprod && go run .
```

- API health: `http://localhost:8080/healthz`
- Portal: `http://localhost:5173`

Full setup: [`docs/SETUP.md`](docs/SETUP.md). Compose self-host pack: `deploy/compose`.

## Documentation

Start with the **[documentation map](docs/README.md)**.

| Document | Purpose |
|----------|---------|
| [`docs/PRODUCT.md`](docs/PRODUCT.md) | Product vision & deploy paths |
| [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) | Layering rules |
| [`docs/PROGRESS.md`](docs/PROGRESS.md) | What is done / next |
| [`docs/ROADMAP.md`](docs/ROADMAP.md) | Backlog candidates |
| [`docs/QC.md`](docs/QC.md) | Checks, known debt, security passes |
| [`docs/SETUP.md`](docs/SETUP.md) | Tooling & local run |

Maintainer workflow (another machine / new chat): [`docs/CONTINUE.md`](docs/CONTINUE.md).

## Status

**MVP complete**, including Git-backed deploy, webhook auto-deploy, catalog apps, ownership/ACL, first-run setup, and production-review Pass A/B/C (Critical/High clear).

**Next:** optional backlog items in [`docs/ROADMAP.md`](docs/ROADMAP.md) (product features or remaining Medium known debt in `docs/QC.md`).

## License

See repository license file if present; otherwise treat as the author’s published source until a license is added.
