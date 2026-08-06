# Sailorport — Progress

> Update file ini setiap selesai 1 step. Ini sumber kebenaran saat pindah mesin.

## Status saat ini

- **Step selesai:** 6
- **Step berikutnya:** 7 (Update & delete service di portal web)
- **Terakhir dikerjakan:** 2026-08-06
- **Mesin terakhir:** rumah / lokal

## Checklist step belajar

- [x] Step 0 — Struktur repo + README + git init
- [x] Step 1 — Go install + endpoint GET `/healthz`
- [x] Step 2 — Refactor `internal/` + config env + POST `/api/v1/echo`
- [x] Step 3 — PostgreSQL via Docker Compose + koneksi DB
- [x] Step 4 — Migrasi tabel `services`
- [x] Step 5 — CRUD catalog API
- [x] Step 6 — Portal web menampilkan catalog (list + create)
- [ ] Step 7 — Update & delete service di portal
- [ ] Step 8 — Scaffold / golden path (1 template)
- [ ] Step 9 — Auth OIDC + RBAC dasar
- [ ] Step 10 — Worker registry + agent deploy

## Yang sudah jalan

```bash
# 1. Postgres
cd deploy/compose && docker compose up -d

# 2. API
cd apps/api && go run .

# 3. Portal web
cd apps/web && npm install && npm run dev
```

| Endpoint / UI | Method / aksi | Hasil |
|---------------|---------------|-------|
| `/healthz` | GET | `{"status":"ok","service":"sailorport-api","version":"0.1.0"}` |
| `/api/v1/echo` | POST | `{"reply":"Sailorport received: ..."}` |
| `/api/v1/services` | GET | list services (`[]` jika kosong) |
| `/api/v1/services` | POST | create service → `201` |
| `/api/v1/services/{id}` | GET | get satu service |
| `/api/v1/services/{id}` | PUT | update service |
| `/api/v1/services/{id}` | DELETE | hapus → `204` |
| `http://localhost:5173` | UI | list catalog + form create |

Env vars API: `PORT`, `APP_ENV`, `APP_VERSION`, `DATABASE_URL`

Default `DATABASE_URL`:
`postgres://sailorport:sailorport@localhost:5432/sailorport?sslmode=disable`

Portal: Vite proxy `/api` dan `/healthz` → `http://localhost:8080`  
API: middleware CORS mengizinkan `http://localhost:5173`

## Struktur file saat ini

```text
Sailorport/
├── apps/api/
│   ├── main.go
│   ├── go.mod
│   └── internal/
│       ├── config/
│       ├── db/
│       ├── handler/   (health, echo, services, cors)
│       ├── migrate/   (goose + SQL 00001_create_services)
│       ├── model/
│       └── store/
├── apps/web/          ← React + TS + Vite (Step 6)
│   ├── vite.config.ts (proxy ke API)
│   └── src/
│       ├── App.tsx
│       ├── api.ts
│       ├── types.ts
│       ├── App.css
│       └── main.tsx
├── apps/worker/       ← kosong
├── apps/agent/        ← kosong
├── packages/shared/   ← kosong
├── deploy/compose/    ← docker-compose.yml (Postgres)
└── docs/
```

## Catatan belajar

### Backend (Step 0–5)

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
- CORS middleware: izinkan browser Vite memanggil API lintas-port

### Frontend (Step 6)

- Vite = tooling frontend (dev server + hot reload)
- `useState` = state yang berubah di layar
- `useEffect(..., [])` = jalan sekali saat halaman dibuka
- `fetch` = HTTP dari browser (mirip curl)
- `import type { ... }` wajib untuk tipe saja (`FormEvent`) karena `verbatimModuleSyntax`
- Nama file di Linux case-sensitive: `App.css` ≠ `app.css`
- Nama export harus sama: `listServices` ≠ `listenServices` (typo bikin blank putih)
- Blank putih → cek Console browser (F12), bukan hanya terminal Vite

## Blocker / pertanyaan

- Cek nama file CORS: ada kemungkinan typo `cors..go` (titik ganda). Rename ke `cors.go` jika masih salah.
- Portal belum punya UI update/delete (API sudah siap).

## Next action (urutan Step 7)

1. Tombol/form edit service (panggil `PUT /api/v1/services/{id}`)
2. Tombol hapus service (panggil `DELETE /api/v1/services/{id}`)
3. Pastikan file CORS bernama `cors.go`
4. Test end-to-end: list → create → update → delete dari browser

## Cara lanjut di mesin lain

Lihat `docs/CONTINUE.md` dan paste prompt dari `docs/RESUME-PROMPT.md` ke chat Cursor baru.
