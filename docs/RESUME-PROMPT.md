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

**Step terakhir selesai:** **22a–22b** — fondasi catalog apps:
- 22a: migrasi `catalog_app_id` / `image` / `container_port` + model/store/web types
- 22b: `catalog-apps/postgres/manifest.json` + `GET /api/v1/catalog-apps` (+ `/{id}`)

**Step berikutnya:** **22c** — API create/deploy `source_type=catalog_app` (isi image dari manifest). Lanjut 22d agent pull/run, 22e portal, 22f docs.

**Visi produk (ringkas):**
- Sailorport **tetap IDP**; **catalog** = inventory pusat
- **Jalur utama:** Git + Dockerfile → agent sync → build → run
- **Scaffold** = opsional; **Register only** = metadata tanpa deploy
- **Catalog apps** (Postgres/Redis) = Step 22 (sedang dikerjakan; jangan campur dengan `templates/`)

**Yang sudah jalan (jangan ulang):**
- MVP core Step 0–18
- Step 19–21 Git + webhook + redeploy by SHA
- Step 22a–22b catalog app fields + list manifests
- Migrasi `00017`–`00020`

**Yang belum:**
- 22c create `catalog_app`; 22d agent `docker pull`/`run`; 22e portal; 22f docs
- Private Git credentials; polish edit Git fields di dialog Edit
- Pass B/C production review (opsional sebelum expose publik; Pass A sudah clear)

**Catatan teknis:**
- `catalog-apps/` ≠ `templates/` (manifest image vs scaffold kode)
- `source_type=catalog_app` belum diizinkan di `validateSourceFields` sampai 22c
- Env: `SAILORPORT_CATALOG_APPS` (default cari folder `catalog-apps`)
- Agent token: `Authorization: Bearer $SAILORPORT_AGENT_TOKEN`

**Cara jalankan lokal:** `docs/SETUP.md` + `docs/PROGRESS.md` (mode development).

Tolong lanjutkan **Step 22c** (API create catalog_app dari manifest) dengan gaya panduan detail seperti sebelumnya. Baca `docs/PRODUCT.md`.

---

*Update bagian "Step terakhir selesai" setiap selesai step baru.*
