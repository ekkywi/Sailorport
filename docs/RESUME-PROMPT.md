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

**Step terakhir selesai:** **Step 26 (26a–26f)** — update `catalog_env` pada service existing (API merge + portal edit + dark-mode select).

**Step berikutnya:** belum ditetapkan — lihat **Next action** di `docs/PROGRESS.md` (opsional: Redis manifest, Pass B/C QC, worker admin lite).

**Visi produk (ringkas):**
- Sailorport **tetap IDP**; **catalog** = inventory pusat
- **Jalur utama:** Git + Dockerfile → agent sync → build → run
- **Jalur sekunder:** catalog apps (Postgres, …) → pull image → run + env dari portal
- **Scaffold** = opsional; **Register only** = metadata tanpa deploy
- `catalog-apps/` ≠ `templates/`

**Yang sudah jalan (jangan ulang):**
- MVP core Step 0–18
- Step 19–21 Git + webhook + redeploy by SHA
- Step 22 catalog apps (API + agent + portal)
- Step 23 catalog env (schema, DB, API, agent, portal, smoke)
- Step 24 catalog app versions (manifest, API image pick, portal dropdown, smoke)
- Step 25 encrypt catalog env at-rest (`EncryptedStore`, `SAILORPORT_SECRETS_KEY`)
- Step 26 update `catalog_env` on existing services (PUT merge + portal edit)
- Migrasi `00017`–`00021`

**Yang belum / opsional:**
- Manifest app tambahan (Redis, …)
- Pass B/C production review sebelum expose publik
- Worker admin lite (edit labels / decommission)

**Catatan teknis:**
- Tambah catalog app: folder `catalog-apps/<id>/manifest.json` + optional `env[]` + optional `versions[]`
- Portal **From catalog** baca schema env dari API; password = field `secret: true`; versi = dropdown `versions[]` → kirim `image`
- Edit catalog_app: `catalog_env` di PUT; secret kosong = keep; setelah ubah env → **redeploy**
- Agent dapat env penuh lewat `POST /api/v1/agent/jobs/next` → `catalog_env` map
- Redact secret di GET user: `POSTGRES_PASSWORD_set: true` (bukan nilai asli)
- Encrypt at-rest: `SAILORPORT_SECRETS_KEY` (hex 64 chars); hanya baris `secret: true`; dev kosong = plaintext OK
- Env: `SAILORPORT_CATALOG_APPS`, `SAILORPORT_TEMPLATES`, `SAILORPORT_AGENT_TOKEN`, `SAILORPORT_SECRETS_KEY`

**Cara jalankan lokal:** `docs/SETUP.md` + `docs/PROGRESS.md` (mode development).

Tolong baca `docs/PROGRESS.md` **Next action** dan lanjutkan dari situ dengan gaya panduan detail seperti sebelumnya.

---

*Update bagian "Step terakhir selesai" setiap selesai step baru.*
