# Sailorport documentation

Map of project docs. If two files disagree, **`PROGRESS.md` wins** for “what is done / what is next.”

## Language

| Surface | Language |
|---------|----------|
| Repository `README.md` | English (public entry) |
| Product, architecture, setup, roadmap, QC | Indonesian (formal) + English technical terms |
| Progress journal, resume prompt, AI context | Indonesian (maintainer workflow) |

Technical terms stay in English (`deploy`, `webhook`, `claim`, `RBAC`, …).

## Read by goal

| Goal | Document |
|------|----------|
| Product vision & deploy paths | [`PRODUCT.md`](PRODUCT.md) |
| Current status & next work | [`PROGRESS.md`](PROGRESS.md) |
| Backlog ideas (candidates) | [`ROADMAP.md`](ROADMAP.md) |
| Quality checks, known debt, Pass A/B/C | [`QC.md`](QC.md) |
| Local / compose setup | [`SETUP.md`](SETUP.md) |
| Layering rules (handler → service → store) | [`ARCHITECTURE.md`](ARCHITECTURE.md) |
| Context for AI assistants | [`AGENTS.md`](AGENTS.md) |
| Continue work on another machine | [`CONTINUE.md`](CONTINUE.md) → paste [`RESUME-PROMPT.md`](RESUME-PROMPT.md) |
| Coursework / thesis notes (not core product) | [`ACADEMIC-GUIDE.md`](ACADEMIC-GUIDE.md), [`ACADEMIC-AI-BRIEF.md`](ACADEMIC-AI-BRIEF.md) |

## Audience

| Audience | Start here |
|----------|------------|
| Visitors / evaluators | Root `README.md` → `PRODUCT.md` → `ARCHITECTURE.md` |
| Operators / developers | `SETUP.md` → `PROGRESS.md` (status) → `QC.md` |
| Maintainers continuing a step | `PROGRESS.md` → `RESUME-PROMPT.md` |
| Academic writing | `ACADEMIC-*` only; engineering truth stays in PRODUCT / ARCHITECTURE / PROGRESS |

## Short glossary

| Term | Meaning |
|------|---------|
| **Catalog** | Central inventory of services in the portal/API |
| **Agent** | Process on a Docker node: register, heartbeat, claim jobs, build/run |
| **Claim** | Agent takes the next pending deploy/runtime job |
| **Worker** | Registered node (agent); not the same as env `dev`/`staging`/`prod` |
| **Redeploy** | New job that rebuilds from a stored `git_sha` (not instant container restore) |
| **Pass A/B/C** | Scoped production reviews (API / agent / web) recorded in `QC.md` |

## What not to treat as product docs

- `ACADEMIC-GUIDE.md` / `ACADEMIC-AI-BRIEF.md` — laporan / TA scaffolding
- Long step history inside `PROGRESS.md` — implementation journal; use the **status block at the top** for a quick view
