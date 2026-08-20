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

**Step terakhir selesai:** 18c — portal filter worker by environment; Workers page tier/environments; Step 18 (worker capabilities) selesai.

**Step berikutnya:** Opsional — webhook auto-deploy; post-MVP admin edit worker labels.

**Catatan produk:**
- Create service (default) = scaffold dari template + daftar catalog
- Register existing = metadata saja, tanpa generate folder
- Template masih di folder `templates/`, bukan database
- Delete catalog: hapus workspace di bawah `data/workspaces` + enqueue `remove` untuk **setiap** env; ditolak jika ada env `running` (prod running → 403)
- Admin Users: list, change role, create, disable/enable (konfirmasi), reset password, **soft-delete** (rename email + `deleted_at`); belum email SMTP
- Agent endpoints butuh `Authorization: Bearer $SAILORPORT_AGENT_TOKEN` (bukan JWT user)
- Deploy: body `{"environment":"staging","worker_id":"<uuid>"}` opsional; portal Deploy dialog pilih worker (search + scroll jika banyak); redeploy affinity otomatis di API
- Workers: self-register via agent (bukan admin CRUD); `/worker` monitoring read-only; labels dari env agent (`SAILORPORT_WORKER_TIER`, `SAILORPORT_WORKER_ENVIRONMENTS`, optional JSON `SAILORPORT_WORKER_LABELS`); deploy policy API 409 jika env tidak diizinkan; portal DeployDialog filter worker by env
- Environments: `GET /api/v1/environments`; portal Deploy dialog; `GET /services` juga return `env_deployments` (map slug → latest deploy per env)
- Workspace default: `data/workspaces` (bukan `/tmp`); Compose workspaces = named volume
- Portal RBAC UI: viewer read-only di Catalog (boleh History); Users page admin-only
- Catalog: kolom Deploy menampilkan `latest_deployment` + badge environment; rocket = dialog deploy (pilih env), jam = history; square/play = stop/start runtime (**konfirmasi**); 📜 = logs (semua role)
- Logs: POST `/services/{id}/runtime/logs` → agent `docker logs --tail 200` → output di `runtime_jobs.output`; portal `LogsDialog` poll `GET /runtime/{job_id}` tiap 2s; viewer+ boleh akses
- Audit: tabel `audit_events` append-only; catalog hard-delete tetap tapi snapshot di audit; `GET /api/v1/audit` admin; portal `/audit` admin-only
- Delete container fix (00014): job `remove` survive SET NULL + `environment_slug`; agent `docker rm` setelah catalog hilang

**Cara jalankan lokal:** lihat `docs/PROGRESS.md` / `docs/SETUP.md` (dua mode).

Tolong lanjutkan fitur berikutnya (webhook auto-deploy atau post-MVP worker admin) dengan gaya panduan detail seperti sebelumnya.

---

*Catatan: update bagian "Step terakhir selesai" di prompt ini setiap kali Anda menyelesaikan step baru.*
