# Sailorport — Product Vision

> Dokumen positioning produk. Baca ini sebelum fitur besar (Git deploy, catalog apps, webhook).
> Detail implementasi harian: `docs/PROGRESS.md`. Peta step: `docs/ROADMAP.md`.

## Apa itu Sailorport?

**Self-hosted Internal Developer Platform (IDP)** — portal + API + agent untuk **mendaftar, deploy, dan mengoperasikan service** di infra sendiri.

Bukan hanya “tempat scaffold template”. **Catalog** (`/catalog`, API `services`) = **inventory pusat** semua yang di-deploy: custom app, infra app, atau hasil scaffold.

Tagline evolusi:

- Dulu: *catalog, **pave**, and ship* (golden path dominan)
- Arah: *catalog, deploy, and **ship*** (deploy dari Git + apps siap pakai dominan; pave opsional)

## Dua jalur deploy (target produk)

Semua jalur **berakhir di catalog yang sama** — deploy, env, worker policy, logs, stop/start tetap dari baris catalog.

### Jalur 1 — Custom app (Git + Dockerfile) — **primary path**

Developer sudah punya repo; IDP hanya urus deploy.

```text
Developer: repo + Dockerfile (syarat wajib)
  → Register / link service di catalog (repo_url, branch, dockerfile_path)
  → Deploy: agent git clone/pull → docker build → docker run
  → Update: webhook push atau redeploy manual
  → Rollback: redeploy commit/tag/image sebelumnya
```

| Aspek | Detail |
|-------|--------|
| Kontrak deploy | **Dockerfile** di repo (platform tidak tebak stack) |
| Source of truth | Git repo tim |
| Build | `docker build` di agent |
| Status implementasi | ✅ Step 19 (19a–19d) |

Ini selaras industri (Render, Railway, internal IDP + Git).

### Jalur 2 — Catalog app (platform-provided) — **secondary path**

App/infra siap pakai dari platform; tanpa repo developer.

```text
User: pilih "PostgreSQL" / "Redis" / "Gitea" / … dari catalog apps
  → Agent: docker pull image → docker run (+ env, volume, port)
  → Tidak ada git clone / docker build
```

| Aspek | Detail |
|-------|--------|
| Contoh | Postgres, Redis, Gitea, AdGuard, … |
| Kontrak | Manifest platform (image, env, volumes, ports) |
| Status implementasi | **Done (MVP)** — 22a–22f: create, deploy, portal, agent pull/run |

### Jalur 3 — Golden path scaffold (opsional) — **sudah ada (MVP)**

Starter code dari template; cocok demo / tim yang belum punya repo.

```text
Portal: Create service → template go-api
  → Generate folder data/workspaces/{name}/
  → Developer kembangkan kode di workspace (manual / git init sendiri)
  → Deploy: agent docker build workspace lokal
```

| Aspek | Detail |
|-------|--------|
| Template saat ini | `templates/go-api` (minimal /healthz + Dockerfile) |
| Status | **Ada** — bukan jalur utama jangka panjang |
| Register existing | Metadata saja di catalog, tanpa workspace — **belum** auto-deploy |

## Catalog tetap dipertahankan

| Pertanyaan | Jawaban |
|------------|---------|
| UI `/catalog` tetap? | **Ya** — pusat kontrol semua service |
| Backend `services` tetap? | **Ya** — perlu diperluas field (`source_type`, `repo_url`, …) |
| Scaffold dihapus? | **Tidak** — jadi opsi ketiga |
| Dua portal terpisah? | **Tidak** — satu catalog, banyak cara create |

## Worker & environment (sudah jalan — Step 18)

- Environment (dev/staging/prod) ≠ worker (node)
- Worker labels: `tier`, `environments` dari agent env
- Deploy policy: worker harus allow environment → **409** jika melanggar
- Portal: DeployDialog filter worker by env; Workers page kolom Tier / Environments

Berlaku untuk **semua jalur** deploy ke depan.

## Apa yang sudah jalan vs belum (ringkas)

| Area | Status |
|------|--------|
| Catalog CRUD + deploy UI | ✅ |
| Scaffold template `go-api` | ✅ |
| Agent docker build/run (workspace lokal) | ✅ |
| Environments, runtime stop/start, logs | ✅ |
| Multi-agent targeting + worker labels/policy | ✅ |
| Git clone/pull sebelum deploy | ✅ Step 19c |
| Service fields `repo_url`, `source_type` | ✅ Step 19a–19b |
| Portal Add from Git | ✅ Step 19d |
| Webhook auto-deploy | ✅ Step 20 (20a–20e) |
| Rollback / redeploy commit | ✅ Step 21 |
| Catalog apps (Postgres, Redis, …) | ✅ Step 22 (22a–22f) |
| Admin edit worker labels | ⬜ post-MVP |

## Keputusan produk (jangan dilanggar di MVP berikutnya)

1. **Jangan** dukung “repo sembarang tanpa Dockerfile” — kontrak = container buildable.
2. **Jangan** wajibkan scaffold dari IDP — developer nyata mulai dari Git.
3. **Jangan** campur “template go-api” dengan “catalog app Postgres” — beda jalur, beda agent flow.
4. **Do** pertahankan satu catalog sebagai inventory.
5. **Do** webhook/rollback setelah Git path — Step 19–21 ✅; catalog apps Step 22 ✅.

## Referensi industri

| Pola | Contoh |
|------|--------|
| IDP + catalog | Backstage, Port |
| Git + Dockerfile deploy | Render, Fly, Railway |
| One-click infra apps | CapRover, Coolify templates |
| Golden path (opsional) | Backstage Software Templates |

Sailorport = **self-hosted IDP ringan** (monolith + agent), bukan PaaS publik penuh.
