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

**Step terakhir selesai:** Step 10C.3 — portal tombol Deploy di Catalog, dialog list deployments (status + poll + link healthz), sidebar AppShell collapsible, fix template `go-api` `ListenAndServe(mux)`. Flow scaffold → deploy → `curl :18080/healthz` sudah teruji dari portal.

**Step berikutnya:** Step 11 — Docker Compose full stack (packaging self-host).

**Catatan produk:**
- Create service (default) = scaffold dari template + daftar catalog
- Register existing = metadata saja, tanpa generate folder
- Template masih di folder `templates/`, bukan database
- Delete catalog menghapus workspace di bawah `data/workspaces` (path legacy `/tmp` di-skip)
- Agent claim/update masih endpoint publik (harden nanti)
- Deploy MVP: satu host port (`SAILORPORT_DEPLOY_PORT_BASE`, default 18080)
- Workspace default: `data/workspaces` (bukan `/tmp`)

**Cara jalankan lokal:** lihat `docs/PROGRESS.md` (compose + api + web + agent).

Tolong lanjutkan Step 11 dengan gaya panduan detail seperti sebelumnya.

---

*Catatan: update bagian "Step terakhir selesai" di prompt ini setiap kali Anda menyelesaikan step baru.*
