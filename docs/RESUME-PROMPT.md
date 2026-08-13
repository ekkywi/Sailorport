# Prompt untuk Lanjut Chat Baru

Copy semua teks di bawah ini ke chat Cursor baru di mesin lain.

---

Saya lanjut proyek **Sailorport** (self-hosted IDP: catalog, scaffold, deploy via agent).

**Mode belajar:** saya coding manual, Anda pandu step-by-step dengan penjelasan detail baris per baris. Jangan refactor besar tanpa diminta.

**Baca dulu file ini di repo:**
- `docs/PROGRESS.md` — step terakhir yang selesai + checklist 12f WIP
- `docs/ARCHITECTURE.md` — aturan lapisan (wajib diikuti)
- `docs/AGENTS.md` — konvensi & konteks proyek

**Stack:** Go (api/agent) + React/TS (web) + PostgreSQL + Docker Compose.

**Step terakhir selesai penuh:** 12e — admin reset password (`POST /api/v1/users/{id}/password` → 204; dialog Reset password di `/users`).

**Sedang dikerjakan:** **12f — soft-delete user** (Steps 1–4 sudah di-commit; lanjut dari Step 5).

**Checklist 12f:**
- [x] Migrasi `00009_add_users_deleted_at.sql`
- [x] Model `User.DeletedAt`
- [x] Store: filter `deleted_at IS NULL`, `SoftDelete` (rename email + disable)
- [x] Service: `Users.SoftDelete` (no self-delete)
- [ ] **Step 5:** Handler `DELETE /api/v1/users/{id}` + router (admin) → 204
- [ ] **Step 6:** Portal tombol Delete + konfirmasi di `/users`

**Desain soft-delete:**
- Bukan hard delete; set `deleted_at = NOW()`, `disabled = true`
- Email di-rename: `{email}__deleted__{uuid}` agar email lama bisa dipakai lagi
- List/Get/Login hanya user dengan `deleted_at IS NULL`
- Auth extra check (4c) **tidak perlu** — store filter sudah cukup

**Step berikutnya setelah 12f:** Environments (dev/staging/prod); opsional logs / multi-port.

**Catatan produk:**
- Create service (default) = scaffold dari template + daftar catalog
- Register existing = metadata saja, tanpa generate folder
- Template masih di folder `templates/`, bukan database
- Delete catalog: hapus workspace di bawah `data/workspaces` + enqueue container cleanup (`remove`) jika pernah deploy
- Admin Users: list, change role, create user, disable/enable, reset password; **soft-delete handler + UI belum**
- Agent endpoints butuh `Authorization: Bearer $SAILORPORT_AGENT_TOKEN` (bukan JWT user)
- Deploy MVP: satu host port (`SAILORPORT_DEPLOY_PORT_BASE`, default 18080)
- Workspace default: `data/workspaces` (bukan `/tmp`); Compose workspaces = named volume
- Portal RBAC UI: viewer read-only di Catalog (boleh History); Users page admin-only
- Catalog: kolom Deploy menampilkan `latest_deployment`; rocket = deploy baru, jam = history; square/play = stop/start runtime

**Cara jalankan lokal:** lihat `docs/PROGRESS.md` / `docs/SETUP.md` (dua mode).

Tolong lanjutkan **Step 12f Step 5** (HTTP delete handler + router) dengan gaya panduan detail seperti sebelumnya.

---

*Catatan: update bagian "Step terakhir selesai" di prompt ini setiap kali Anda menyelesaikan step baru.*
