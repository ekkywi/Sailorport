# Sailorport — Progress

> Update file ini setiap selesai 1 step. Ini sumber kebenaran saat pindah mesin.

## Status saat ini

- **Step selesai:** 2
- **Step berikutnya:** 3 (PostgreSQL via Docker Compose + koneksi DB)
- **Terakhir dikerjakan:** 2026-08-05
- **Mesin terakhir:** kantor

## Checklist step belajar

- [x] Step 0 — Struktur repo + README + git init
- [x] Step 1 — Go install + endpoint GET `/healthz`
- [x] Step 2 — Refactor `internal/` + config env + POST `/api/v1/echo`
- [x] Step 3 — PostgreSQL via Docker Compose + koneksi DB
- [x] Step 4 — Migrasi tabel `services`
- [ ] Step 5 — CRUD catalog API
- [ ] Step 6 — Portal web menampilkan catalog

## Yang sudah jalan

```bash
cd apps/api && go run .
```

| Endpoint | Method | Hasil |
|----------|--------|-------|
| `/healthz` | GET | `{"status":"ok","service":"sailorport-api","version":"0.1.0"}` |
| `/api/v1/echo` | POST | `{"reply":"Sailorport received: ..."}` |

Env vars: `PORT`, `APP_ENV`, `APP_VERSION`

## Struktur file saat ini

```text
Sailorport/
├── apps/api/          ← API Go (aktif)
├── apps/web/          ← kosong
├── apps/worker/       ← kosong
├── apps/agent/        ← kosong
├── packages/shared/   ← kosong
├── deploy/compose/    ← kosong
└── docs/              ← progress & panduan lanjut
```

## Catatan belajar

- `package main` = program yang bisa dijalankan
- `:=` mendeklarasikan variabel baru
- Handler HTTP implement `ServeHTTP(w, r)`
- `json.NewDecoder(r.Body).Decode(&req)` baca JSON dari client
- Status code: 200 OK, 201 Created, 400 Bad Request, 405 Method Not Allowed
- `internal/` = kode privat, hanya untuk proyek ini
- `os.Getenv` baca environment variable

## Blocker / pertanyaan

- (kosong)

## Next action (urutan Step 3)

1. Buat `deploy/compose/docker-compose.yml` — service Postgres
2. Tambah `DATABASE_URL` di config API
3. Koneksi DB dari Go (`database/sql` + driver `pgx`)
4. Test: API bisa `SELECT 1` ke Postgres

## Cara lanjut di mesin lain

Lihat `docs/CONTINUE.md` dan paste prompt dari `docs/RESUME-PROMPT.md` ke chat Cursor baru.
