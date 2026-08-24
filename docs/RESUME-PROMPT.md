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

**Stack:** Go (api/agent) + React/TS (web) + PostgreSQL + Docker Compose.

**Step terakhir selesai:** **20 (20a–20e)** — webhook auto-deploy end-to-end:
- 20a: migrasi webhook fields + catalog/store
- 20b–20d: public endpoint, HMAC, match → create deployment
- 20e: portal edit git (secret generate, auto-deploy toggle/env); API redact secret (`webhook_secret_set`)

**Step berikutnya:** **21 — rollback / redeploy**. Alternatif: polish edit Git fields di UI, private repo credentials, catalog apps (22).

**Visi produk (ringkas):**
- Sailorport **tetap IDP**; **catalog** = inventory pusat
- **Jalur utama:** Git + Dockerfile → agent sync → build → run
- **Scaffold** = opsional; **Register only** = metadata tanpa deploy
- **Catalog apps** (Postgres/Redis) = Step 22+

**Yang sudah jalan (jangan ulang):**
- MVP core Step 0–18
- Step 19 Git path
- Step 20 webhook auto-deploy lengkap (API + portal)
- Migrasi `00017` (Git) + `00018` (webhook)

**Yang belum:**
- Rollback / redeploy by commit (21)
- Catalog apps (22)
- Private Git credentials; edit Git fields di dialog Edit

**Catatan teknis:**
- Webhook: `POST /api/v1/webhooks/github` (no JWT); `X-Hub-Signature-256`
- Portal: Edit git service → webhook section
- List/Get JSON: `webhook_secret` always empty; use `webhook_secret_set`
- Agent token: `Authorization: Bearer $SAILORPORT_AGENT_TOKEN`

**Cara jalankan lokal:** `docs/SETUP.md` + `docs/PROGRESS.md` (mode development).

Tolong lanjutkan **Step 21** (rollback / redeploy) dengan gaya panduan detail seperti sebelumnya. Baca `docs/PRODUCT.md`.

---

*Update bagian "Step terakhir selesai" setiap selesai step baru.*
