# Sailorport

Self-hosted internal developer platform — catalog, pave, and ship.

## Apa ini

Sailorport adalah platform internal yang bisa di-host sendiri, berisi:

- **Catalog** — daftar semua service beserta pemilik dan repo-nya
- **Golden path** — buat service baru dari template siap pakai
- **Deploy via agent** — deploy ke server sendiri lewat agent yang menjalankan Docker

## Struktur

```text
apps/web          portal — catalog CRUD + scaffold
apps/api          control plane — layered + scaffold
apps/worker       background jobs (Go) — belum
apps/agent        agent di node (Go) — belum
templates/        golden path templates (`go-api`)
packages/shared   kontrak & tipe bersama — belum
deploy/compose    Docker Compose (Postgres :5433)
docs/             progress, architecture, panduan lanjut
```

## Quick start (lokal)

```bash
# 1. Database
cd deploy/compose && docker compose up -d

# 2. API
cd apps/api && go run .

# 3. Portal (terminal lain)
cd apps/web && npm install && npm run dev
```

- API: `http://localhost:8080/healthz`
- Portal: `http://localhost:5173`

## Lanjut di mesin lain

Chat Cursor tidak ikut pindah antar device. Yang ikut pindah: **repo Git + folder `docs/`**.

| Situasi | Baca |
|---------|------|
| Pindah ke laptop/rumah | [`docs/CONTINUE.md`](docs/CONTINUE.md) |
| Setup tool di mesin baru | [`docs/SETUP.md`](docs/SETUP.md) |
| Lihat step mana yang sudah selesai | [`docs/PROGRESS.md`](docs/PROGRESS.md) |
| Buka chat Cursor baru | [`docs/RESUME-PROMPT.md`](docs/RESUME-PROMPT.md) |
| Peta besar proyek | [`docs/ROADMAP.md`](docs/ROADMAP.md) |

**Ritual singkat sebelum tutup:**

```bash
git add .
git commit -m "docs: update progress"
git push
```

## Status

Step 8 selesai (scaffold golden path `go-api`). Step berikutnya: auth OIDC + RBAC.

Arsitektur: lihat [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md).
