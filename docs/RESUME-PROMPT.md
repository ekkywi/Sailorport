# Prompt untuk Lanjut Chat Baru

Copy semua teks di bawah ini ke chat Cursor baru di mesin lain.

---

Saya lanjut proyek **Sailorport** (self-hosted IDP: catalog, scaffold, deploy via agent).

**Mode belajar:** saya coding manual, Anda pandu step-by-step dengan penjelasan detail baris per baris. Jangan refactor besar tanpa diminta.

**Baca dulu file ini di repo:**
- `docs/PROGRESS.md` — step terakhir yang selesai
- `docs/ARCHITECTURE.md` — aturan lapisan (wajib diikuti)
- `docs/AGENTS.md` — konvensi & konteks proyek

**Stack:** Go (api/agent) + React/TS (web) + PostgreSQL + Docker Compose.

**Step terakhir selesai:** R4 — multi-port deploy (agent alokasi host port unik di rentang `SAILORPORT_DEPLOY_PORT_BASE` + count; reuse saat redeploy).

**Step berikutnya:** Environments (dev/staging/prod); opsional logs / multi-agent targeting.

**Catatan produk:**
- Create service (default) = scaffold dari template + daftar catalog
- Register existing = metadata saja, tanpa generate folder
- Template masih di folder `templates/`, bukan database
- Delete catalog: hapus workspace di bawah `data/workspaces` + enqueue container cleanup (`remove`) jika pernah deploy
- Admin Users: list, change role, create, disable/enable (konfirmasi), reset password, **soft-delete** (rename email + `deleted_at`); belum email SMTP
- Agent endpoints butuh `Authorization: Bearer $SAILORPORT_AGENT_TOKEN` (bukan JWT user)
- Deploy: host port unik per container di mesin agent (`SAILORPORT_DEPLOY_PORT_BASE` + `SAILORPORT_DEPLOY_PORT_COUNT`); hasil di `deployments.port`
- Workspace default: `data/workspaces` (bukan `/tmp`); Compose workspaces = named volume
- Portal RBAC UI: viewer read-only di Catalog (boleh History); Users page admin-only
- Catalog: kolom Deploy menampilkan `latest_deployment`; rocket = deploy baru, jam = history; square/play = stop/start runtime

**Cara jalankan lokal:** lihat `docs/PROGRESS.md` / `docs/SETUP.md` (dua mode).

Tolong lanjutkan Environments (dev/staging/prod) dengan gaya panduan detail seperti sebelumnya.

---

*Catatan: update bagian "Step terakhir selesai" di prompt ini setiap kali Anda menyelesaikan step baru.*
