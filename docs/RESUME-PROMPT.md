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

**Step terakhir selesai:** **20c** — GitHub webhook signature validation:
- 20a: migrasi `webhook_secret`, `auto_deploy_enabled`, `auto_deploy_environment` + catalog/store/types
- 20b: `POST /api/v1/webhooks/github` parse push + ack
- 20c: HMAC-SHA256 `X-Hub-Signature-256`; match service by `clone_url`; 401 jika invalid

**Step berikutnya:** **20d** — match branch + `auto_deploy_enabled` → create deployment. Lalu 20e portal secret/toggle.

**Visi produk (ringkas):**
- Sailorport **tetap IDP**; **catalog** = inventory pusat
- **Jalur utama:** Git + Dockerfile → agent sync → build → run
- **Scaffold** = opsional; **Register only** = metadata tanpa deploy
- **Catalog apps** (Postgres/Redis) = Step 22+

**Yang sudah jalan (jangan ulang):**
- MVP core Step 0–18 (auth, catalog, env, deploy, runtime, logs, audit, multi-agent, worker policy)
- Step 19 Git path: API fields, agent sync, portal Add from Git
- Step 20a–20c webhook fields + public endpoint + signature validation
- Agent workspace: `SAILORPORT_WORKSPACE` (clone ke `{workspace}/{serviceName}`)
- Migrasi `00017` (Git) + `00018` (webhook) — pastikan kolom ada di DB

**Yang belum:**
- Webhook auto-deploy create deployment + portal UI (20d–20e)
- Rollback / redeploy by commit (21)
- Catalog apps (22)
- Private Git credentials; edit Git fields di dialog Edit

**Catatan teknis:**
- Deploy body: `{"environment":"staging","worker_id":"<uuid>"}` opsional
- Agent token: `Authorization: Bearer $SAILORPORT_AGENT_TOKEN`
- Webhook route **tanpa JWT** (GitHub tidak punya portal token)
- Repo Git publik dulu (belum SSH/token)
- Agent env: lihat `apps/agent/.env.example`

**Cara jalankan lokal:** `docs/SETUP.md` + `docs/PROGRESS.md` (mode development).

Tolong lanjutkan **Step 20d** (match repo/branch → create deployment) dengan gaya panduan detail seperti sebelumnya. Baca `docs/PRODUCT.md`.

---

*Update bagian "Step terakhir selesai" setiap selesai step baru.*
