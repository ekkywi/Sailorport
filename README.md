# Sailorport

Self-hosted internal developer platform — catalog, pave, and ship.

## Apa ini

Sailorport adalah platform internal yang bisa di-host sendiri, berisi:

- **Catalog** — daftar semua service beserta pemilik dan repo-nya
- **Golden path** — buat service baru dari template siap pakai
- **Deploy via agent** — deploy ke server sendiri lewat agent yang menjalankan Docker

## Struktur

```text
apps/web          portal (React + TypeScript) — belum dimulai
apps/api          control plane API (Go) — aktif
apps/worker       background jobs (Go) — belum
apps/agent        agent di node (Go) — belum
packages/shared   kontrak & tipe bersama — belum
deploy/compose    Docker Compose self-hosting — belum
docs/             progress, setup, panduan lanjut
```

## Quick start (API lokal)

```bash
cd apps/api
go run .
```

```bash
curl http://localhost:8080/healthz
```

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

Dalam pengembangan awal. Step 2 selesai (API health + echo). Step berikutnya: PostgreSQL.
