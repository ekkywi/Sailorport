# Prompt untuk Lanjut Chat Baru

Copy semua teks di bawah ini ke chat Cursor baru di mesin lain.

---

Saya lanjut proyek **Sailorport** (self-hosted IDP: catalog, deploy, ship via agent).

**Mode belajar:** saya coding manual, Anda pandu step-by-step dengan penjelasan detail baris per baris. Jangan refactor besar tanpa diminta.

**Baca dulu file ini di repo (urutan penting):**
- `docs/PRODUCT.md` — **visi produk & dua jalur deploy** (wajib baca)
- `docs/PROGRESS.md` — step terakhir yang selesai + rencana Step 19+
- `docs/ARCHITECTURE.md` — aturan lapisan (wajib diikuti)
- `docs/AGENTS.md` — konvensi & konteks proyek

**Stack:** Go (api/agent) + React/TS (web) + PostgreSQL + Docker Compose.

**Step terakhir selesai:** 18c — worker capabilities (labels, deploy policy API, portal filter by env). **MVP core selesai.**

**Step berikutnya:** **19 — Git-backed custom app deploy** (`repo_url`, agent git clone/pull, Dockerfile path, build & run). Baca `docs/PRODUCT.md` Jalur 1.

**Visi produk (ringkas):**
- Sailorport **tetap IDP**; **catalog** (`/catalog`, API services) = inventory pusat semua service.
- **Jalur utama:** developer punya repo + Dockerfile → IDP deploy (pull → build → run). Scaffold **tidak wajib**.
- **Jalur sekunder (nanti):** catalog apps siap pakai (Postgres, Redis, Gitea) — pull image, tanpa Git.
- **Scaffold `go-api`:** opsional / golden path demo; deploy saat ini dari `data/workspaces/` lokal.
- Webhook & rollback **setelah** Step 19, bukan sebelum Git path.

**Yang sudah jalan (jangan ulang):**
- Catalog CRUD, deploy, history, env dev/staging/prod, worker picker, stop/start, logs, audit, RBAC
- Agent: register, heartbeat, claim deploy/runtime, docker build/run workspace lokal
- Worker labels (`SAILORPORT_WORKER_TIER`, `SAILORPORT_WORKER_ENVIRONMENTS`); deploy policy 409; portal filter worker by env
- Multi-agent: `target_worker_id`, claim filter, redeploy affinity

**Yang belum (Step 19+):**
- `repo_url`, `branch`, `source_type` di model service
- Agent `git clone` / `git pull` sebelum `docker build`
- Webhook, rollback, catalog apps (Postgres/Redis)

**Catatan teknis:**
- Create service (scaffold) = template → `data/workspaces/{name}/` → deploy build folder itu
- Register existing = metadata saja, tanpa workspace — belum auto-deploy
- Deploy body: `{"environment":"staging","worker_id":"<uuid>"}` opsional
- Agent token: `Authorization: Bearer $SAILORPORT_AGENT_TOKEN` (bukan JWT user)
- Workspace dev: `data/workspaces`; Compose: named volume
- Agent env file: `apps/agent/.env.example` → copy ke `.env.nonprod` lokal, `source .env.nonprod && go run .`

**Cara jalankan lokal:** `docs/SETUP.md` + `docs/PROGRESS.md` (mode development).

Tolong lanjutkan **Step 19** (Git-backed deploy) dengan gaya panduan detail seperti sebelumnya. Ikuti `docs/PRODUCT.md` — kontrak deploy = Dockerfile di repo.

---

*Update bagian "Step terakhir selesai" setiap selesai step baru.*
