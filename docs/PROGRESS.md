# Sailorport — Progress

> Update file ini setiap selesai 1 step. Ini sumber kebenaran saat pindah mesin.

## Status saat ini

- **Step selesai:** 5
- **Step berikutnya:** 6 (Portal web menampilkan catalog)
- **Terakhir dikerjakan:** 2026-08-05
- **Mesin terakhir:** kantor

## Checklist step belajar

- [x] Step 0 — Struktur repo + README + git init
- [x] Step 1 — Go install + endpoint GET `/healthz`
- [x] Step 2 — Refactor `internal/` + config env + POST `/api/v1/echo`
- [x] Step 3 — PostgreSQL via Docker Compose + koneksi DB
- [x] Step 4 — Migrasi tabel `services`
- [x] Step 5 — CRUD catalog API
- [ ] Step 6 — Portal web menampilkan catalog

## Yang sudah jalan

```bash
# 1. Postgres
cd deploy/compose && docker compose up -d

# 2. API
cd apps/api && go run .
```

| Endpoint | Method | Hasil |
|----------|--------|-------|
| `/healthz` | GET | `{"status":"ok","service":"sailorport-api","version":"0.1.0"}` |
| `/api/v1/echo` | POST | `{"reply":"Sailorport received: ..."}` |
| `/api/v1/services` | GET | list services (`[]` jika kosong) |
| `/api/v1/services` | POST | create service → `201` |
| `/api/v1/services/{id}` | GET | get satu service |
| `/api/v1/services/{id}` | PUT | update service |
| `/api/v1/services/{id}` | DELETE | hapus → `204` |

Env vars: `PORT`, `APP_ENV`, `APP_VERSION`, `DATABASE_URL`

Default `DATABASE_URL`:
`postgres://sailorport:sailorport@localhost:5432/sailorport?sslmode=disable`

## Struktur file saat ini

```text
Sailorport/
├── apps/api/
│   ├── main.go
│   ├── go.mod
│   └── internal/
│       ├── config/
│       ├── db/
│       ├── handler/   (health, echo, services)
│       ├── migrate/   (goose + SQL 00001_create_services)
│       ├── model/
│       └── store/
├── apps/web/          ← kosong (Step 6)
├── apps/worker/       ← kosong
├── apps/agent/        ← kosong
├── packages/shared/   ← kosong
├── deploy/compose/    ← docker-compose.yml (Postgres)
└── docs/
```

## Catatan belajar

- `package main` = program yang bisa dijalankan
- `:=` mendeklarasikan variabel baru; `!=` membandingkan
- Handler HTTP: `ServeHTTP` atau method + `HandleFunc("GET /path", ...)`
- `json.NewDecoder(r.Body).Decode(&req)` baca JSON dari client
- Status: 200 OK, 201 Created, 204 No Content, 400 Bad Request, 404 Not Found, 409 Conflict
- `internal/` = kode privat proyek
- `os.Getenv` / `DATABASE_URL` untuk config
- `database/sql` + driver `pgx`; `SELECT 1` untuk ping
- Goose + `embed` untuk migrasi SQL berversi
- Lapisan: handler → store (SQL) → model
- `$1`, `$2` = parameter query (hindari SQL injection)
- `r.PathValue("id")` ambil `{id}` dari route Go 1.22+

## Blocker / pertanyaan

- (kosong)

## Next action (urutan Step 6)

1. Init portal React + TypeScript + Vite di `apps/web`
2. Halaman list catalog (panggil `GET /api/v1/services`)
3. Form create service sederhana
4. Test end-to-end: web → API → Postgres

## Cara lanjut di mesin lain

Lihat `docs/CONTINUE.md` dan paste prompt dari `docs/RESUME-PROMPT.md` ke chat Cursor baru.
