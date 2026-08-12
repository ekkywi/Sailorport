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

**Step terakhir selesai:** Harden — `SAILORPORT_AGENT_TOKEN` melindungi register/heartbeat/claim/update agent. Tested: 401 tanpa token, 204 dengan Bearer `dev-agent-token`.

**Step berikutnya:** Delete service membersihkan container Docker (stop/rm); lalu runtime controls (stop/start/restart + logs). Opsional multi-port; Environments.

**Catatan produk:**
- Create service (default) = scaffold dari template + daftar catalog
- Register existing = metadata saja, tanpa generate folder
- Template masih di folder `templates/`, bukan database
- Delete catalog menghapus workspace di bawah `data/workspaces` (path legacy `/tmp` di-skip); **belum** stop/rm container
- Agent endpoints butuh `Authorization: Bearer $SAILORPORT_AGENT_TOKEN` (bukan JWT user)
- Deploy MVP: satu host port (`SAILORPORT_DEPLOY_PORT_BASE`, default 18080)
- Workspace default: `data/workspaces` (bukan `/tmp`); Compose workspaces = named volume
- Portal RBAC UI: viewer read-only di Catalog; Users page admin-only

**Cara jalankan lokal:** lihat `docs/PROGRESS.md` / `docs/SETUP.md` (dua mode).

Tolong lanjutkan cleanup delete → docker stop/rm (atau runtime controls) dengan gaya panduan detail seperti sebelumnya.

---

*Catatan: update bagian "Step terakhir selesai" di prompt ini setiap kali Anda menyelesaikan step baru.*
