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

**Step terakhir selesai:** **19 (19a–19d)** — Git-backed deploy end-to-end:
- 19a: migrasi `source_type`, `repo_url`, `branch`, `dockerfile_path`
- 19b: API create/update Git + deploy gate
- 19c: agent `git clone`/`pull` lalu `docker build`
- 19d: portal Add service chooser + form From Git + Deploy now

**Step berikutnya:** **20 — webhook auto-deploy** (push → create deployment). Alternatif: rollback (21), polish edit Git fields di UI, private repo credentials.

**Visi produk (ringkas):**
- Sailorport **tetap IDP**; **catalog** = inventory pusat
- **Jalur utama:** Git + Dockerfile → agent sync → build → run
- **Scaffold** = opsional; **Register only** = metadata tanpa deploy
- **Catalog apps** (Postgres/Redis) = Step 22+

**Yang sudah jalan (jangan ulang):**
- MVP core Step 0–18 (auth, catalog, env, deploy, runtime, logs, audit, multi-agent, worker policy)
- Step 19 Git path: API fields, agent sync, portal Add from Git
- Agent workspace: `SAILORPORT_WORKSPACE` (clone ke `{workspace}/{serviceName}`)
- Migrasi `00017_service_git_source.sql` — pastikan kolom ada di DB

**Yang belum:**
- Webhook auto-deploy (20)
- Rollback / redeploy by commit (21)
- Catalog apps (22)
- Private Git credentials; edit Git fields di dialog Edit

**Catatan teknis:**
- Deploy body: `{"environment":"staging","worker_id":"<uuid>"}` opsional
- Agent token: `Authorization: Bearer $SAILORPORT_AGENT_TOKEN`
- Repo Git publik dulu (belum SSH/token)
- Agent env: lihat `apps/agent/.env.example`

**Cara jalankan lokal:** `docs/SETUP.md` + `docs/PROGRESS.md` (mode development).

Tolong lanjutkan **Step 20** (webhook auto-deploy) dengan gaya panduan detail seperti sebelumnya. Baca `docs/PRODUCT.md`.

---

*Update bagian "Step terakhir selesai" setiap selesai step baru.*
