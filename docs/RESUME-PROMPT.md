# Prompt untuk Lanjut Chat Baru

> Maintainer workflow — paste ke chat Cursor baru. Bukan dokumen produk untuk visitor.  
> Peta docs: `docs/README.md`.

Copy teks di bawah ke chat baru.

---

Saya lanjut proyek **Sailorport** (self-hosted IDP: catalog, deploy, ship via agent).

**Mode kerja:** coding manual, panduan step-by-step dengan penjelasan jelas. Jangan refactor besar tanpa diminta.

**Baca dulu (urutan):**
- `docs/README.md` — peta dokumentasi
- `docs/PRODUCT.md` — visi produk & dua jalur deploy
- `docs/PROGRESS.md` — status + next (sumber kebenaran)
- `docs/ROADMAP.md` — backlog kandidat
- `docs/ARCHITECTURE.md` — aturan lapisan
- `docs/AGENTS.md` — konvensi proyek
- `docs/QC.md` — checks, Known debt, Pass A/B/C

**Stack:** Go (api/agent) + React/TS (web) + PostgreSQL + Docker Compose.

**Step terakhir selesai:** **Step 38** — Health / Open URL; smoke Open/Health di scaffold running (2026-09-30). Agent workspace default selaras API (`data/workspaces`).

**Step berikutnya:** **Step 39** — Notifikasi deploy gagal (lonceng topbar + badge unread). Dikunci: developer = fails on owned services; admin = all deploy fails; viewer = no bell (MVP); shell siap jenis notif lain nanti.

**Visi produk (ringkas):**
- Sailorport **tetap IDP**; **catalog** = inventory pusat
- **Jalur utama:** Git + Dockerfile → agent sync → build → run
- **Jalur sekunder:** catalog apps (Postgres, Redis, …) → pull image → run + env/command dari manifest
- **Scaffold** = opsional; **Register only** = metadata tanpa deploy
- `catalog-apps/` ≠ `templates/`

**Yang sudah jalan (jangan ulang):**
- MVP core Step 0–18
- Step 19–21 Git + webhook + redeploy by SHA
- Step 22–27 catalog apps (env, versions, encrypt, command, Redis)
- Step 28 worker admin lite
- Step 29 service ownership (ACL + transfer + directory picker)
- Step 30 first-run `/setup` (status, create admin, gate; bootstrap tidak lewat `/register`)
- Step 31 app settings: `registration_open`, admin Settings page, public registration-status, Register → developer, portal Sign up kondisional
- Step 32 global deployments list ACL (`ListByOwner` / admin sees all) + unit tests
- Step 33 webhook delivery dedupe (`webhook_deliveries` + unit tests replay)
- Step 34 login rate limit (in-memory per IP, 429)
- Step 35 CORS (`CORS_ORIGINS` allowlist + PATCH in Allow-Methods)
- **Pass B** (2026-09-28): agent path containment, git/docker arg harden, PATCH `worker_id` = claimer
- **Pass C** (2026-09-28): portal 401 → logout; Redeploy `canWrite`; client redact `webhook_secret`
- Migrasi melalui `00024_create_webhook_deliveries.sql` (app_settings tetap `00023`)

**Yang belum / opsional:**
- Backlog / Known debt Medium: `docs/ROADMAP.md`, `docs/QC.md`

**Auth / setup (ingat):**
- Instalasi kosong → paksa `/setup` → admin pertama
- Default `registration_open=false`; buka lewat **Administration → Settings**
- Self-register selalu role `developer`; admin buat user lewat Users

**Cara jalankan lokal:** `docs/SETUP.md` + `docs/PROGRESS.md` (mode development).

Tolong baca `docs/PROGRESS.md` **Next action** dan lanjutkan dari situ dengan gaya panduan detail seperti sebelumnya.

---

*Update bagian "Step terakhir selesai" setiap selesai step baru.*
