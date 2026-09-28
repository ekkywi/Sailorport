# AI / contributor context — Sailorport

Context for assistants and maintainers continuing work in this repository.

## Project

**Sailorport** — self-hosted internal developer platform (OSS).

Tagline: *Self-hosted developer port — catalog, deploy, and ship.*

Core capabilities: **software catalog**, agent-based deploy, environments, worker health/policy, RBAC, audit. Scaffold golden path is **optional** (not the primary path).

**Product vision:** `docs/PRODUCT.md` — read before Git deploy / catalog-app features.  
**Doc map:** `docs/README.md`.

## Repo

- GitHub: `github.com/ekkywi/Sailorport`
- Go module API: `github.com/ekkywi/sailorport/apps/api`
- Go module Agent: `github.com/ekkywi/sailorport/apps/agent`

## Stack

| Component | Path | Status |
|-----------|------|--------|
| Portal | `apps/web` | auth + catalog deploy/runtime + workers + users + RBAC |
| API | `apps/api` | layered + scaffold + workers + deployments |
| Worker | `apps/worker` | not implemented (job queue / orchestrator) |
| Agent | `apps/agent` | register + heartbeat + poll deploy/runtime; host ports; Bearer agent token |
| Templates | `templates/` | `go-api` (on disk) |
| Shared contracts | `packages/shared` | not yet |
| Compose | `deploy/compose` | Postgres + API + web; workspaces named volume |

Infra: PostgreSQL (yes), Redis (not yet), local JWT (yes), OIDC (not yet).

## Contribution / learning guidance

Work is done in **small, testable steps** (this repo is also used as a structured learning path).

When guiding implementation:

1. One step = one runnable, testable outcome
2. Prefer complete file contents over fragmentary snippets for new concepts
3. Explain new ideas clearly; avoid drive-by refactors
4. Do not add features outside the requested step
5. End each step with: how to test + suggested commit message

## Architecture

```
Developer → Web Portal
              → API (Go)
                 → PostgreSQL (+ Redis later)
                 → Worker (jobs) [not yet]
                 → Agent on node (register, heartbeat, poll, docker build/run)
                 → status callback
```

Principle: the control plane does not run containers directly; the agent does.

## Workers (runtime nodes)

- Worker = node with Docker + agent; **not** the same as environment `dev`/`staging`/`prod`.
- One worker may run many environments (containers `sailorport-{service}-{env}`).
- Workers appear via **agent register + heartbeat** — no admin CRUD create/delete in MVP.
- `labels` (JSONB): sent on register (Step 18). Env `SAILORPORT_WORKER_TIER`, `SAILORPORT_WORKER_ENVIRONMENTS`, optional `SAILORPORT_WORKER_LABELS`. Deploy policy returns 409 if env not allowed.
- Deploy: optional `worker_id`; `target_worker_id` + claim filter route jobs.

**Auth note:** agents currently share one `SAILORPORT_AGENT_TOKEN`. Pass B requires PATCH updates to send matching claimer `worker_id`. Per-worker tokens remain a backlog item (`ROADMAP` / `QC` Known debt).

## Portal routes (after login)

| Path | Content |
|------|---------|
| `/overview` | services + workers summary |
| `/catalog` | services, deploy history, runtime, CRUD |
| `/worker` | workers + status (admin can edit labels / decommission) |
| `/users` | admin user management |
| `/settings` | admin: `registration_open` |
| `/audit` | admin audit trail |

Auth: `/login`, `/setup` (first-run), `/register` (when registration is open).

## Catalog — mental model

**Catalog = all platform-managed services** (not a “template store”).

| Entry path | Status | Deploy |
|------------|--------|--------|
| Scaffold (`go-api`) | ✅ | Build local `workspace_path` |
| Register metadata | ✅ | No auto workspace |
| Git + Dockerfile | ✅ | clone/pull → build → run |
| Catalog app (Postgres, Redis, …) | ✅ | pull image → run |

All paths share **`/catalog`** for deploy, env, logs, and runtime.

## Scaffold (optional golden path)

- Scaffold = template → `data/workspaces/{name}/` → catalog row
- Develop in the workspace after scaffold; template used once
- Register existing = metadata only
- Delete = enqueue `remove` → DB + workspace cleanup → `docker rm`

## Coding conventions

- Read `docs/ARCHITECTURE.md` before new features
- API: `/api/v1/...`; health: `GET /healthz`
- Thin `main.go` — wiring only
- Flow: `handler` → `service` → `store` (no domain SQL in handlers)
- API errors: `{"error":"..."}`
- Portal: `src/features/<domain>/`; thin `App.tsx`
- CORS: `CORS_ORIGINS` allowlist; methods include PATCH
- Commits: `feat(api):`, `feat(web):`, `feat(agent):`, `docs:`, `fix:`

## Resume workflow

1. Read `docs/PROGRESS.md` — “Step berikutnya”
2. If unset: pick from `docs/ROADMAP.md` backlog
3. Do not redo completed steps
4. After a step: update PROGRESS (+ ROADMAP/QC if needed) + commit
5. QC / production review: `docs/QC.md`

## MVP success criteria

`docker compose up` → agent → worker online → catalog service → deploy → status/logs.

**Done:** MVP core + Steps 18–35 + Pass A/B/C Critical/High. **Next:** optional backlog in `docs/ROADMAP.md`.
