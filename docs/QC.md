# Sailorport — QC Checkpoint

> Checklist kualitas ringan, dijalankan berkala (tiap 3–5 step atau sebelum step besar berikutnya). Bukan audit penuh — fokus: bug nyata, kontrak API, smoke path utama. Update tanggal + temuan setiap kali dijalankan.

## Cara pakai

1. Jalankan bagian **Automated checks** dulu (build/vet/test).
2. Jalankan **Smoke manual** kalau ada perubahan di path deploy/webhook/redeploy.
3. Catat temuan baru di **Known debt**; pindahkan ke **Fixed** kalau sudah beres.
4. Update `docs/PROGRESS.md` kalau ada bugfix yang perlu dicatat di step history.
5. Sebelum “siap produksi” / pakai **model mahal**: siapkan **Evidence pack**, lalu jalankan **Production review** (satu pass A/B/C per chat).

---

## Automated checks

```bash
# API
cd apps/api && go build ./... && go vet ./... && go test ./...

# Agent
cd apps/agent && go build ./... && go vet ./... && go test ./...

# Web
cd apps/web && npm run build   # tsc -b && vite build
```

Semua harus **exit 0** sebelum lanjut step baru. Kalau salah satu merah, itu blocker — jangan lanjut fitur baru dulu.

**Terakhir dijalankan:** 2026-08-28 — API ✅ (build+vet+test), Agent ✅ (build+vet+test), Web ✅ (tsc+vite build).

---

## Smoke manual (path utama)

Jalankan minimal setelah perubahan di `deployments`, `webhook`, atau `agent`:

1. Login portal → Catalog kosong/isi tampil
2. Deploy service **scaffold** → status `running`, healthz OK
3. Add service dari **Git** → Deploy → agent clone + build + run
4. History (`DeploymentsDialog`) → SHA tampil, badge status update live
5. **Redeploy** pada deployment lama yang punya `git_sha` → job baru `pending` → agent checkout SHA yang sama (cek log agent: `sha="..."`)
6. Webhook (opsional): push GitHub → signature valid → auto-deploy sesuai `auto_deploy_environment`
7. Stop / Start container dari portal → status ikut berubah
8. GET `/api/v1/deployments/{id}` → mengembalikan **satu** deployment (bukan list) — regression check untuk bug di bawah
9. Webhook tanpa/dengan signature salah, dan repo yang tidak terdaftar → **401** yang sama (tidak ada `no matching service`) — regression check A-H1
10. Disable user di portal, lalu pakai token lama user itu ke `GET /api/v1/services` → **401** tanpa perlu tunggu token kedaluwarsa — regression check A-H3
11. `PATCH /api/v1/deployments/{id}` dengan JWT developer → **405** (route portal sudah dihapus); `PATCH /api/v1/agent/deployments/{id}` dengan agent token tetap **200** — regression check A-H4
12. **Catalog app (Step 22):** Portal Add service → From catalog → PostgreSQL → Deploy dev → `running`; agent log `catalog_app pull` (bukan git/build). Stop / Start / Logs / History OK. Atau curl di `docs/PROGRESS.md` tes 22c–22d.
13. **Catalog app versions (Step 24):** Portal Add from catalog → pilih **15-alpine** → create → Deploy dev → agent log `pull image=postgres:15-alpine`; `docker inspect` container = `postgres:15-alpine`. Default tanpa pilih versi = `postgres:16-alpine`. Curl invalid image → **400** (`docs/PROGRESS.md` tes 24c).
14. **Encrypt catalog env (Step 25):** Set `SAILORPORT_SECRETS_KEY=$(openssl rand -hex 32)`, restart API → log `encrypted at rest`. Create catalog service dengan password → DB row `POSTGRES_PASSWORD` prefix `enc:1…`; GET API redact `*_set`; deploy → container env plaintext benar. Regression: unset key di dev → plaintext store OK (`docs/PROGRESS.md` tes 25d).

---

## Known debt (belum dikerjakan, sengaja ditunda)

| Item | Dampak | Catatan |
|------|--------|---------|
| Tes unit untuk `Deployments.Redeploy` / `Create` dengan `git_sha` | Sedang | Belum ada `deployment_test.go`; behavior baru divalidasi manual |
| Tes untuk `git.Sync` + `checkoutSHA` (agent) | Sedang | Butuh repo git lokal di test; belum ada `sync_test.go` |
| Redeploy = rebuild dari SHA, bukan restore container instan | Rendah (by design) | Didokumentasikan di `PROGRESS.md` Step 21; jangan "perbaiki" tanpa diskusi |
| Private Git repo credentials | Rendah | Belum didukung; hanya public clone URL |
| Bundle web >500KB (vite warning) | Rendah | Belum perlu code-splitting di skala MVP ini |
| Worker admin lite (edit label/decommission) | Rendah | Post-MVP, sudah di roadmap |
| Webhook: tanpa dedupe `X-GitHub-Delivery` + tanpa rate limit | Sedang | Pass A. `service/webhook.go:113` — replay satu payload sah bisa terus membuat row `deployments` pending |
| Webhook deploy tidak mengirim `payload.After` sebagai `git_sha` | Sedang | Pass A. `service/webhook.go:113` — `ack.commit_sha` melaporkan SHA push, tapi deploy pakai tip branch (konsisten keputusan terkunci; bisa beda commit kalau ada push menyusul) |
| `webhook_secret` tidak bisa dikosongkan lewat API | Sedang | Pass A. `service/catalog.go:339` — empty = pertahankan yang lama; perlu sentinel eksplisit untuk revoke |
| Claim job: `worker_id` self-reported + satu shared agent token | Sedang | Pass A. `store/deployment.go:112` — agen mana pun bisa kirim `worker_id` node lain dan mencuri job bertarget. `FOR UPDATE SKIP LOCKED` sendiri sudah benar (tidak ada double-claim) |
| `DeploymentsStore.Update` pakai `COALESCE(NULLIF($n,''))` | Sedang | Pass A. `store/deployment.go:142` — `error_message` tidak bisa dikosongkan setelah retry sukses; transisi status tidak dijaga (`running` → `pending` diterima) |
| CORS: `Access-Control-Allow-Methods` tanpa `PATCH`, origin hardcode | Sedang | Pass A. `handler/cors.go:8` — belum terasa karena dev pakai vite proxy dan compose pakai nginx same-origin (`apps/web/nginx.conf:7`); pecah kalau portal diarahkan langsung ke `:8080` (PATCH users role/disable gagal preflight) |
| Login tanpa rate limit / lockout | Sedang | Pass A. `handler/auth.go:35` — brute force murah, terutama selama register masih terbuka |
| `handler/scaffold.go` mengembalikan `model.Service` tanpa `PublicService` | Rendah | Pass A. `handler/scaffold.go:47` — aman sekarang (scaffold selalu menulis secret `""`), tapi satu-satunya jalur serialisasi service yang tidak lewat redaksi |
| `webhook_secret` tanpa `omitempty` | Rendah | Pass A. `model/service.go:20` — response portal selalu memuat `"webhook_secret":""` (noise kontrak, bukan kebocoran) |
| Batas body webhook 1 MiB | Rendah | Pass A. `handler/webhook.go:13` — push dengan banyak commit bisa dibalas 413 dan deploy terlewat senyap |
| `UpdateDeploymentRequest.WorkerID` diabaikan store | Rendah | Pass A. `service/deployment.go:156` — di-trim lalu tidak pernah dipakai `store.Update` |
| JWT tanpa validasi `iss` / `aud` | Rendah | Pass A. `internal/auth/jwt.go:32` — aman karena `SigningMethodHS256` dipaksa dan secret tunggal |
| Compose: password Postgres default + port 5433 dipublish | Rendah | Pass A. `deploy/compose/docker-compose.yml:6` — dev convenience; jangan dipakai apa adanya di host publik |
| Portal masih menampilkan form `/register` padahal register sudah tertutup | Rendah | Efek fix A-C2. Setelah admin pertama ada, submit selalu **403** dengan pesan jelas; menyembunyikan link/form = Pass C |
| `model.RegisterRequest.Role` diabaikan | Rendah | Efek fix A-C2. Akun pertama selalu `admin`; field dibiarkan supaya kontrak web tidak pecah |
| Dua register serentak di instalasi kosong bisa jadi dua admin | Rendah | Efek fix A-C2. Gate-nya `COUNT(*)` lalu INSERT (bukan atomik); praktis tidak relevan karena siapa pun yang duluan register tetap dapat admin |
| Portal tidak auto-logout saat 401 di tengah sesi | Rendah | Efek fix A-H3. Token user yang di-disable langsung ditolak API, tapi portal baru menghapus token saat `me()` gagal (refresh / buka ulang); interceptor 401 global = Pass C |
| Satu query user tambahan per request ber-JWT | Rendah (by design) | Efek fix A-H3. Alternatif `token_version`/cache sengaja tidak dipakai supaya disable & ganti role langsung berlaku |

---

## Fixed (riwayat bugfix dari QC)

| Tanggal | Bug | Fix |
|---------|-----|-----|
| 2026-08-26 | `GET /api/v1/deployments/{id}` memanggil `deployments.List` (semua deployment) bukan `deployments.Get(id)` | `DeploymentsHandler.Get` sekarang panggil `h.deployments.Get(ctx, r.PathValue("id"))` + `writeDeploymentError` |
| 2026-08-26 | **A-C1** Secret bawaan repo dipakai apa adanya: `AUTH_JWT_SECRET=dev-only-change-me`, `SAILORPORT_AGENT_TOKEN=dev-agent-token`, termasuk di compose | `config.Config.Validate()` menolak start kalau `APP_ENV != development` dan secret masih kosong/nilai dev (dipanggil dari `main.go`); compose wajib `${AUTH_JWT_SECRET:?…}` + `${SAILORPORT_AGENT_TOKEN:?…}` dari `deploy/compose/.env` (`.env.example` baru; `APP_ENV` bisa di-override) |
| 2026-08-26 | **A-C2** `POST /api/v1/auth/register` publik tanpa gate dan memberi role `developer` (boleh create service + deploy) | Register jadi jalur bootstrap saja: `UsersStore.Count` (termasuk soft-deleted) → kalau sudah ada user, **403** `registration is closed`; akun pertama otomatis role `admin` dan `role` dari request diabaikan. User berikutnya lewat `POST /api/v1/users` (admin) |
| 2026-08-26 | **A-H1** Webhook membalas beda-beda sebelum HMAC diverifikasi (`no matching service` / `secret not configured` / `no auto-deploy`) → oracle enumerasi repo bagi penyerang tanpa signature | `HandleGitHub` memverifikasi signature dulu; semua kegagalan auth jadi satu `ErrUnauthorized` (**401**) dan `ack` baru diisi setelah terverifikasi. Tes lama yang mengunci `no matching service` diganti `TestHandleGitHub_UnknownRepoIsUnauthorized` |
| 2026-08-26 | **A-H2** Secret diambil dari service pertama yang punya secret, tapi yang di-deploy `eligible[0]` → secret service A bisa memicu deploy service B pada repo yang sama | `filterVerifiedServices` menyisakan service yang **secret-nya sendiri** cocok dengan body; auto-deploy dan target dipilih hanya dari himpunan itu. Tes baru: `TestHandleGitHub_OtherServiceSecretCannotDeploy`, `TestHandleGitHub_DeploysServiceOwningTheSecret` |
| 2026-08-26 | **A-H4** `PATCH /api/v1/deployments/{id}` terbuka untuk role writer, padahal itu endpoint laporan status agent → developer bisa mengarang `status`/`git_sha`/`container_id` | Route portal dihapus dari `router.go`; laporan status hanya lewat `PATCH /api/v1/agent/deployments/{id}` (`withAgentToken`). Portal tidak pernah memakainya (`features/deployments/api.ts` hanya `POST …/redeploy`) |
| 2026-08-26 | **A-H3** `RequireRole` percaya `role`/`disabled` dari klaim JWT (TTL 24 jam) → user yang di-disable, di-soft-delete, atau diturunkan role tetap punya akses sampai token kedaluwarsa | `RequireAuth(secret, currentUserLookup)` memuat user dari DB (`service.Auth.Me` → tolak disabled / `deleted_at`), lalu menimpa `claims.Role` + `claims.Email` dengan nilai DB sebelum `RequireRole` berjalan. Biaya: satu query user per request ber-JWT |
| 2026-09-01 | Catalog app: `POSTGRES_PASSWORD=changeme` hardcoded di agent | Step 23d–23f: `catalog_env` dari claim job + portal form; agent `docker run -e` dari map env user |

---

## Hasil Production review — Pass A (2026-08-26)

Scope: API auth, webhook, deploy, secret. Diff `ce88e2b^..HEAD` + uncommitted (Step 19–21). Evidence: build/vet/test API hijau.
Status: **Critical clear, High clear** — A-C1, A-C2, A-H1, A-H2, A-H4, A-H3 semua difix 2026-08-26 (lihat tabel **Fixed**). Sisa Pass A hanya Medium/Low di **Known debt**. Pass B (agent) dan Pass C (web) belum pernah dijalankan.

### Blocker sebelum expose publik (Critical)

| # | Lokasi | Masalah | Status |
|---|--------|---------|--------|
| A-C1 | `internal/config/config.go`, `deploy/compose/docker-compose.yml` | Secret bawaan repo (`dev-only-change-me` / `dev-agent-token`) bisa dipakai memalsukan JWT admin dan agent token | ✅ Fixed 2026-08-26 — `Config.Validate()` + compose wajib `.env` |
| A-C2 | `internal/handler/router.go:49`, `internal/service/auth.go` | Register publik memberi role `developer` → orang luar bisa deploy container arbitrer di worker node | ✅ Fixed 2026-08-26 — register = bootstrap admin pertama saja, sesudahnya 403 |

Catatan pilihan desain A-C2 (jangan diubah tanpa diskusi): akun pertama **otomatis `admin`** (menghapus langkah promote lewat SQL), gate-nya `COUNT(*)` semua baris `users` termasuk soft-deleted, dan tidak ada env `ALLOW_OPEN_REGISTRATION` — jalur satu-satunya untuk user tambahan adalah `POST /api/v1/users`.

### High

| # | Lokasi | Masalah | Dampak | Status / fix disarankan |
|---|--------|---------|--------|------------------------|
| A-H1 | `internal/service/webhook.go` | HMAC diverifikasi **setelah** payload dipakai query catalog dan setelah handler memutuskan balasan | Oracle enumerasi repo terdaftar tanpa signature | ✅ Fixed 2026-08-26 — verifikasi dulu, satu `401` generic. **Sisa:** tiap request tetap memicu `Catalog.List` tanpa auth (rate limit ada di Known debt) |
| A-H2 | `internal/service/webhook.go` | Secret diambil dari service pertama yang punya secret, tapi yang di-deploy `eligible[0]` (bisa service lain) | Pemilik secret service A bisa memicu deploy service B pada repo sama, termasuk ke `prod` | ✅ Fixed 2026-08-26 — target hanya dari service yang secret-nya memverifikasi body |
| A-H3 | `internal/handler/middleware.go` | `RequireRole` percaya `role` dari klaim JWT; TTL 24 jam; tanpa cek DB / token version | User yang di-disable, di-soft-delete, atau diturunkan ke `viewer` tetap bisa deploy/hapus service sampai 24 jam | ✅ Fixed 2026-08-26 — `RequireAuth` memuat user dari DB dan memakai role DB |
| A-H4 | `internal/handler/router.go` | `PATCH /api/v1/deployments/{id}` dibuka untuk role writer, padahal endpoint laporan status agent sudah ada (`withAgentToken`) | Developer bisa mengarang `status`/`image_tag`/`git_sha`/`container_id` → riwayat dan `git_sha` (dasar Redeploy Step 21) tidak bisa dipercaya | ✅ Fixed 2026-08-26 — route portal dihapus |

Catatan pilihan desain A-H1/A-H2 (jangan diubah tanpa diskusi): push ke repo terdaftar yang auto-deploy-nya **off** tetap dijawab **200 ignored** (bukan 401) supaya delivery GitHub tidak merah untuk pemilik secret yang sah; kalau beberapa service satu repo sama-sama eligible, tetap hanya **satu** yang di-deploy (`eligible[0]`) seperti perilaku sebelumnya; event non-`push` dan ref non-branch masih dijawab tanpa verifikasi karena isinya cuma memantulkan input pengirim.

Catatan pilihan desain A-H3 (jangan diubah tanpa diskusi): otorisasi dibaca dari DB tiap request (bukan token version / cache), jadi disable dan ganti role langsung berlaku dengan biaya satu query per request ber-JWT; `service.Auth.Me` dipakai sebagai satu-satunya sumber (`deleted_at IS NULL` + tolak `disabled`), dan `GET /api/v1/auth/me` memang jadi dua query — dibiarkan supaya lapisan handler tetap tipis.

### Yang sudah solid (jangan diutak-atik saat fix)

- HMAC benar: `crypto/hmac` + `hmac.Equal`, prefix `sha256=` divalidasi, decode hex dicek (`service/webhook_signature.go`).
- `withAgentToken` pakai `subtle.ConstantTimeCompare` (`handler/middleware.go:82`).
- Redaksi `webhook_secret` konsisten hanya di handler (`handler/service.go:29,50,59,80`); `Catalog.List` tetap membawa secret asli untuk HMAC — keputusan terkunci dipatuhi dan dikunci tes.
- Semua SQL pakai placeholder `$1…`, tidak ada concat input user; `ClaimNext` bebas double-claim.
- Bug `GET /deployments/{id}` di tabel **Fixed** terverifikasi beres di working tree.

---

## Prinsip QC di project ini

- **Jangan** refactor besar/DRY agresif saat QC — ini project belajar, boilerplate kecil (handler/store) OK.
- **Jangan** tambah abstraksi/framework baru hanya karena "lebih bersih" — ikuti `docs/ARCHITECTURE.md`.
- Fokus QC: bug fungsional, kontrak HTTP salah, kebocoran data (secret), regresi step sebelumnya.
- Kalau nemu bug saat QC → fix kecil boleh langsung; kalau butuh keputusan desain → catat di **Known debt**, jangan dadakan diputuskan sendiri.

---

## Evidence pack (sebelum review model mahal)

Jalankan dan simpan output singkat (atau paste ringkas ke chat):

```bash
git log --oneline -20
git status -sb
# Ganti range sesuai kebutuhan, contoh sejak Step 19:
git log --oneline ce88e2b^..HEAD
git diff --stat ce88e2b^..HEAD

cd apps/api && go build ./... && go vet ./... && go test ./...
cd apps/agent && go build ./... && go vet ./... && go test ./...
cd apps/web && npm run build
```

Lampirkan juga:

- Isi **Known debt** + **Fixed** di file ini
- Keputusan terkunci dari `.cursor/rules/sailorport.mdc` / `docs/ARCHITECTURE.md` (redact secret hanya di handler; redeploy = rebuild by SHA)

**Jangan** minta model mahal me-review seluruh repo dari Step 0 dalam satu chat.

---

## Production review — alur hemat biaya

```text
1. Model biasa / coding step     → implementasi + smoke (bagian atas)
2. Automated green               → build / vet / test / web build
3. Satu chat model mahal         → pilih Pass A ATAU B ATAU C (di bawah)
4. Fix Critical + High saja      → commit
5. Medium / Low                  → masuk Known debt (jangan habiskan budget polish)
```

---

## Security checklist (wajib di Pass A; skim di B/C)

- [ ] Route sensitif pakai `withRole` / `withAgentToken` yang tepat (user JWT ≠ agent token)
- [ ] Webhook publik: HMAC timing-safe, body bounded, event filter (ignore ping/tag yang tidak relevan)
- [ ] `webhook_secret` **tidak** muncul di JSON list/get portal (`webhook_secret_set` saja); secret tetap utuh di `Catalog.List` untuk HMAC
- [ ] Tidak ada SQL string-concat dari input user; placeholder `$1…` selaras dengan argumen
- [ ] Agent: argumen `git` / `docker` tidak menyisipkan shell user input mentah
- [ ] Workspace path / service name tidak membuka path traversal ke luar workspace
- [ ] Error 500 tidak membocorkan stack/internal detail ke client
- [ ] Claim job: `FOR UPDATE SKIP LOCKED` / filter `target_worker_id` tidak race double-claim

---

## File kritis per pass (lampirkan / @-mention di chat)

### Pass A — API: auth, webhook, deploy, secret

| Path | Kenapa |
|------|--------|
| `apps/api/internal/handler/router.go` | AuthZ route map |
| `apps/api/internal/handler/webhook.go` | Endpoint publik |
| `apps/api/internal/handler/deployment.go` | Create / Get / Redeploy |
| `apps/api/internal/service/webhook.go` | Signature + match + create deploy |
| `apps/api/internal/service/webhook_signature.go` | HMAC |
| `apps/api/internal/service/deployment.go` | Create `git_sha`, Redeploy, claim |
| `apps/api/internal/service/service_public.go` | Redact secret |
| `apps/api/internal/store/deployment.go` | ClaimNext SQL, INSERT `git_sha` |
| `apps/api/internal/handler/auth.go` + `internal/auth/` | JWT login/register |

### Pass B — Agent: git sync, docker, claim

| Path | Kenapa |
|------|--------|
| `apps/agent/internal/agent/agent.go` | Poll, resolveWorkDir, PATCH status |
| `apps/agent/internal/git/sync.go` | Sync / checkoutSHA / HeadSHA |
| `apps/agent/internal/client/api.go` | Job payload `git_sha` |
| `apps/agent/internal/docker/` | Build / run / port / stop |
| `apps/agent/internal/config/` | Labels, env, token |

### Pass C — Web: auth, secret UI, redeploy

| Path | Kenapa |
|------|--------|
| `apps/web/src/lib/http.ts` | Bearer token |
| `apps/web/src/features/auth/` | Login/session |
| `apps/web/src/features/catalog/WebhookSettingsFields.tsx` | Secret UI |
| `apps/web/src/features/catalog/types.ts` | `webhook_secret_set` |
| `apps/web/src/features/deployments/DeploymentsDialog.tsx` | Redeploy + SHA |
| `apps/web/src/features/deployments/api.ts` | `redeployDeployment` |

---

## Prompt template — Production review (copy ke chat model mahal)

Ganti bagian dalam `<...>`. Pilih **satu** pass per chat.

```text
Kamu review kode Sailorport untuk production readiness.

Baca dulu: docs/QC.md (Known debt + Fixed), docs/ARCHITECTURE.md, docs/PRODUCT.md.
Ikuti keputusan terkunci: redact webhook_secret hanya di handler PublicService;
redeploy = rebuild by git_sha (bukan restore container).

SCOPE: Pass <A|B|C> saja — file kritis ada di docs/QC.md bagian "File kritis per pass".
DIFF: <git range atau "uncommitted + last N commits">
EVIDENCE: automated checks terakhir <tanggal / hijau>.

TUJUAN: temukan bug, hole keamanan, kontrak HTTP/SQL salah, race, secret leak.
BUKAN tujuan: rewrite arsitektur, DRY besar, rename massal, fitur baru.

OUTPUT (wajib format ini):
## Critical
- file:line — masalah — dampak — fix disarankan (singkat)
## High
- …
## Medium
- …
## Low / debt OK
- …
## Ringkas
- 3 bullet: apa yang sudah solid; apa yang harus difix sebelum expose publik

Setelah daftar temuan, TANYA dulu sebelum mengedit. Fix hanya Critical/High kalau saya setuju.
```

Contoh isi SCOPE cepat:

- Pass A: “auth + webhook HMAC + deployment Get/Redeploy + secret redact”
- Pass B: “agent git checkout SHA + docker run + claim job”
- Pass C: “portal Bearer auth + webhook UI tidak expose secret + Redeploy button”

---

## Setelah review model mahal

1. Fix **Critical** + **High** → commit `fix:` terpisah kalau bisa.
2. **Medium/Low** → tambah baris di **Known debt** (jangan silent ignore).
3. Update **Terakhir dijalankan** di Automated checks bila perlu.
4. Jangan langsung Pass berikutnya sampai Critical Pass ini clear.
