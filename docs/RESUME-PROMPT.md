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

**Step terakhir selesai:** **21 (21a–21f)** — rollback / redeploy by `git_sha`:
- 21a: migrasi `deployments.git_sha`
- 21b–21c: agent catat HEAD + checkout SHA dari job
- 21d: `Create` optional `git_sha` + `POST /api/v1/deployments/{id}/redeploy`
- 21e: portal History → short SHA + tombol Redeploy
- 21f: docs

**Step berikutnya:** **22 — catalog apps**. Alternatif: polish edit Git fields di UI, private repo credentials.

**Visi produk (ringkas):**
- Sailorport **tetap IDP**; **catalog** = inventory pusat
- **Jalur utama:** Git + Dockerfile → agent sync → build → run
- **Scaffold** = opsional; **Register only** = metadata tanpa deploy
- **Catalog apps** (Postgres/Redis) = Step 22+

**Yang sudah jalan (jangan ulang):**
- MVP core Step 0–18
- Step 19 Git path
- Step 20 webhook auto-deploy
- Step 21 redeploy by SHA
- Migrasi `00017`–`00019`

**Yang belum:**
- Catalog apps (22)
- Private Git credentials; polish edit Git fields di dialog Edit

**Catatan teknis:**
- Redeploy = rebuild at old SHA (bukan restore container)
- Webhook tip branch: `git_sha` kosong; Redeploy pin SHA
- Agent: `Sync(repo, branch, dir, sha)` + `checkout --detach`
- Portal: Redeploy hanya jika history punya `git_sha`
- Agent token: `Authorization: Bearer $SAILORPORT_AGENT_TOKEN`

**Cara jalankan lokal:** `docs/SETUP.md` + `docs/PROGRESS.md` (mode development).

Tolong lanjutkan **Step 22** (catalog apps) dengan gaya panduan detail seperti sebelumnya. Baca `docs/PRODUCT.md`.

---

*Update bagian "Step terakhir selesai" setiap selesai step baru.*
