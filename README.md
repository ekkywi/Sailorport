# Sailorport

Self-hosted internal developer platform — catalog, pave, and ship.

## Apa ini

Sailorport adalah platform internal yang bisa di-host sendiri, berisi:

- **Catalog** — daftar semua service beserta pemilik dan repo-nya
- **Golden path** — buat service baru dari template siap pakai
- **Deploy via agent** — deploy ke server sendiri lewat agent yang menjalankan Docker

## Struktur

- `apps/web` — portal (React + TypeScript)
- `apps/api` — control plane API
- `apps/worker` — background jobs
- `apps/agent` — agent yang berjalan di node
- `packages/shared` — kontrak & tipe bersama
- `deploy/compose` — Docker Compose untuk self-hosting
- `docs` — dokumentasi & catatan keputusan

## Status

Dalam pengembangan awal. Belum siap dipakai.