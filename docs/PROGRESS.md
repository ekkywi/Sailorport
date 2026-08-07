# Sailorport — Progress

> Update file ini setiap selesai 1 step. Ini sumber kebenaran saat pindah mesin.

## Status saat ini

- **Step selesai:** 9 (Auth lokal + JWT + RBAC + UI login shadcn)
- **Step berikutnya:** 10 (Worker registry + agent deploy)
- **Terakhir dikerjakan:** 2026-08-06 — polish auth UI (Tailwind v4 + shadcn)
- **Mesin terakhir:** rumah / lokal

## Checklist step belajar

- [x] Step 0 — Struktur repo + README + git init
- [x] Step 1 — Go install + endpoint GET `/healthz`
- [x] Step 2 — Refactor `internal/` + config env + POST `/api/v1/echo`
- [x] Step 3 — PostgreSQL via Docker Compose + koneksi DB
- [x] Step 4 — Migrasi tabel `services`
- [x] Step 5 — CRUD catalog API
- [x] Step 6 — Portal web menampilkan catalog (list + create)
- [x] Step 7 — Update & delete service di portal (CRUD UI lengkap)
- [x] Step 7.5 — Architecture foundation (lapisan service, router, pecah web)
- [x] Step 8 — Scaffold / golden path (1 template `go-api`)
- [x] Step 9 — Auth lokal + JWT + RBAC
- [x] Step 9 polish — Auth UI profesional (Tailwind v4 + shadcn/ui)
- [ ] Step 10 — Worker registry + agent deploy

## Yang sudah jalan

```bash
cd deploy/compose && docker compose up -d
cd apps/api && go run .
cd apps/web && npm run dev
```

| Endpoint / UI | Auth | Hasil |
|---------------|------|-------|
| `POST /api/v1/auth/register` | publik | buat user |
| `POST /api/v1/auth/login` | publik | JWT token |
| `GET /api/v1/auth/me` | Bearer | profil user |
| `GET /api/v1/services` | viewer+ | list |
| `POST/PUT/DELETE services`, `POST scaffold` | developer/admin | mutasi |
| Portal login | — | split-screen AuthLayout, tabs sign in/register, token di localStorage |

Env: `AUTH_JWT_SECRET` (default dev-only-change-me)

Role: `admin`, `developer`, `viewer`

### Web UI stack (auth)

- Routes: `/login`, `/register` (react-router); unauthenticated → redirect `/login`
- Theme: ocean biru–putih (light) + harbour malam (dark); toggle di auth & dashboard (`localStorage`)
- Font: Geist Variable (sama seperti Vercel)
- Tailwind CSS v4 + shadcn primitives; layout `AuthLayout` full-bleed laut + gelombang
- Login/Register: halaman terpisah di `features/auth/`
- Dashboard catalog/scaffold masih pakai `src/styles/app.css` (belum dimigrasi)

Setelah `git pull` di mesin baru: `cd apps/web && npm install`

## Next action (Step 10)

Worker registry + health + agent binary skeleton.

## Cara lanjut di mesin lain

Lihat `docs/CONTINUE.md` dan paste prompt dari `docs/RESUME-PROMPT.md`.
