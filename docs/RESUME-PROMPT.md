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

**Step terakhir selesai:** **Step 23c** — API create dengan `catalog_env`, validasi manifest, simpan ke DB, redact di list/get (`POSTGRES_PASSWORD_set`).

**Step berikutnya:** **23d** — agent claim job dapat env dari API; hapus hardcode `POSTGRES_PASSWORD=changeme`.

**Visi produk (ringkas):**
- Sailorport **tetap IDP**; **catalog** = inventory pusat
- **Jalur utama:** Git + Dockerfile → agent sync → build → run
- **Jalur sekunder:** catalog apps (Postgres, …) → pull image → run
- **Scaffold** = opsional; **Register only** = metadata tanpa deploy
- `catalog-apps/` ≠ `templates/`

**Yang sudah jalan (jangan ulang):**
- MVP core Step 0–18
- Step 19–21 Git + webhook + redeploy by SHA
- Step 22 catalog apps (API + agent + portal)
- Step 23a–23c catalog env schema + DB + API create/redact
- Migrasi `00017`–`00021`

**Yang belum / opsional:**
- Agent pakai `catalog_env` dari claim job (23d); portal form env (23e)
- Manifest app tambahan (Redis, …)
- Private Git credentials; polish edit Git fields
- Pass B/C production review sebelum expose publik

**Catatan teknis:**
- Tambah catalog app baru: folder `catalog-apps/<id>/manifest.json` (API auto-list; web ikut API)
- Env: `SAILORPORT_CATALOG_APPS`, `SAILORPORT_TEMPLATES`
- Agent token: `Authorization: Bearer $SAILORPORT_AGENT_TOKEN`
- Redact secret catalog env **hanya di handler** via `PublicView` (sama pola `webhook_secret_set`)

**Cara jalankan lokal:** `docs/SETUP.md` + `docs/PROGRESS.md` (mode development).

Tolong baca `docs/PROGRESS.md` **Next action** dan lanjutkan dari situ dengan gaya panduan detail seperti sebelumnya.

---

*Update bagian "Step terakhir selesai" setiap selesai step baru.*
