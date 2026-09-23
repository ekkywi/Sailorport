# Prompt untuk Lanjut Chat Baru

Copy semua teks di bawah ini ke chat Cursor baru di mesin lain.

---

Saya lanjut proyek **Sailorport** (self-hosted IDP: catalog, deploy, ship via agent).

**Mode belajar:** saya coding manual, Anda pandu step-by-step dengan penjelasan detail baris per baris. Jangan refactor besar tanpa diminta.

**Baca dulu file ini di repo (urutan penting):**
- `docs/PRODUCT.md` — **visi produk & dua jalur deploy** (wajib baca)
- `docs/PROGRESS.md` — step terakhir yang selesai + rencana berikutnya
- `docs/ARCHITECTURE.md` — aturan lapisan (wajib diikuti)
- `docs/AGENTS.md` — konvensi & konteks proyek
- `docs/QC.md` — automated/smoke + **prompt Production review** (Pass A/B/C) untuk model mahal

**Stack:** Go (api/agent) + React/TS (web) + PostgreSQL + Docker Compose.

**Step terakhir selesai:** **Step 32d** — unit tests `Deployments.List` ACL; Step 32 (global deployments list by owner) complete.

**Step berikutnya:** opsional — Pass B/C QC, catalog app lain (Gitea, …) (lihat `docs/PROGRESS.md` Next action).

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
- Migrasi melalui `00023_create_app_settings.sql`

**Yang belum / opsional:**
- Pass B/C production review sebelum expose publik
- Catalog app lain (Gitea, …)

**Auth / setup (ingat):**
- Instalasi kosong → paksa `/setup` → admin pertama
- Default `registration_open=false`; buka lewat **Administration → Settings**
- Self-register selalu role `developer`; admin buat user lewat Users

**Cara jalankan lokal:** `docs/SETUP.md` + `docs/PROGRESS.md` (mode development).

Tolong baca `docs/PROGRESS.md` **Next action** dan lanjutkan dari situ dengan gaya panduan detail seperti sebelumnya.

---

*Update bagian "Step terakhir selesai" setiap selesai step baru.*
