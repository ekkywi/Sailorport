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

**Step terakhir selesai:** Step 10C.2 — agent poll job (`SAILORPORT_POLL_INTERVAL`), `docker build`/`docker run` dari `workspace_path`, update status `building`/`running`/`failed`. Template `go-api` punya `Dockerfile.tmpl`. Client agent: `ClaimNext`, `UpdateDeployment`.

**Step berikutnya:** Step 10C.3 — portal UI tombol Deploy + list deployments di Catalog.

**Catatan produk:**
- Create service (default) = scaffold dari template + daftar catalog
- Register existing = metadata saja, tanpa generate folder
- Template masih di folder `templates/`, bukan database
- Delete catalog tidak hapus folder workspace (debt known)
- Agent claim/update masih endpoint publik (harden nanti)
- Deploy MVP: satu host port (`SAILORPORT_DEPLOY_PORT_BASE`, default 18080)

**Cara jalankan lokal:** lihat `docs/PROGRESS.md` (compose + api + web + agent).

Tolong lanjutkan Step 10C.3 dengan gaya panduan detail seperti sebelumnya.

---

*Catatan: update bagian "Step terakhir selesai" di prompt ini setiap kali Anda menyelesaikan step baru.*
