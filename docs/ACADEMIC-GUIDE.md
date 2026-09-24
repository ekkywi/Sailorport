# Sailorport — Panduan Akademik (Laporan → Tugas Akhir)

> Dokumen ini untuk penyusunan laporan semester dan perluasan ke TA.
> Bukan pengganti `PRODUCT.md` / `ARCHITECTURE.md` / `PROGRESS.md`.
> Sumber teknis produk tetap dokumen engineering; dokumen ini hanya kerangka akademik.

**Terakhir diisi:** 2026-09-15  
**Status laporan:** judul dikunci; sitensis jurnal  
**Status TA:** belum dimulai

---

## 0. Peta cepat: Laporan vs TA

| Aspek | Laporan (sekarang) | Tugas Akhir (semester depan) |
|-------|--------------------|------------------------------|
| Lingkup | **1 modul** | **Aplikasi keseluruhan** + 1 kontribusi utama |
| Judul fokus | Webhook GitHub → auto-deploy | IDP self-hosted + auto-deploy webhook |
| Kedalaman | Alur + uji modul | Arsitektur sistem + modul fokus + evaluasi |
| Jurnal | Minimal ~20 (sitensis) — untuk **laporan** | Biasanya lebih + teori lebih dalam (perluas D & F) |
| Deliverable kode | Cukup modul yang sudah ada | Platform utuh + mungkin perbaikan/pengembangan tambahan |

**Judul laporan (dikunci):**  
Implementasi Fitur Auto-Deploy Berbasis Webhook GitHub pada Internal Developer Platform

**Judul TA (disarankan, bisa diubah saat proposal):**  
Pengembangan Internal Developer Platform Self-hosted dengan Fitur Auto-Deploy Berbasis Webhook GitHub

**Catatan penamaan:** nama produk **Sailorport tidak dipakai di judul** (nama buatan penulis). Di naskah, sebut sebagai objek studi / studi kasus, mis. *“platform IDP self-hosted yang dikembangkan penulis (Sailorport)”*.

**Prinsip:** laporan = zoom-in satu modul; TA = zoom-out sistem + modul yang sama sebagai kontribusi.  
**Kuota 20 jurnal** di §2 = target **laporan** (mayoritas kelompok A/B/C); bukan kuota TA.

---

## 1. Objek teknis yang dibahas (anchor ke kode)

Gunakan ini agar jurnal dan bab selalu “nyambung” ke implementasi nyata.

### Modul fokus laporan

- Endpoint: `POST /api/v1/webhooks/github`
- File inti:
  - `apps/api/internal/handler/webhook.go`
  - `apps/api/internal/service/webhook.go`
  - `apps/api/internal/service/webhook_signature.go`
  - `apps/api/internal/model/webhook.go`
- Konsep kunci:
  1. Endpoint publik (tanpa JWT)
  2. Filter event `push`
  3. Parse repo + branch
  4. Match service di catalog (`repo_url` / clone URL)
  5. Verifikasi HMAC-SHA256 (`X-Hub-Signature-256`)
  6. Cek auto-deploy + environment
  7. Create deployment → agent yang menjalankan build/run (konteks saja)

### Konteks sistem (wajib singkat di laporan; wajib dalam di TA)

- Control plane (API) ≠ runtime (agent)
- Catalog sebagai inventory service
- Deploy dari Git + Dockerfile
- Portal: secret webhook + toggle auto-deploy (redact secret di response)

---

## 2. Guideline pencarian & sitensis jurnal

### 2.1 Kuota minimal 20 — pecah per kelompok

Isi kolom “Jumlah terkumpul” saat mencari.

| Kode | Kelompok tema | Target | Jumlah terkumpul | Dipakai di bab |
|------|---------------|--------|------------------|----------------|
| A | CI/CD, continuous deployment/delivery | 4–5 |  | Latar belakang, rumusan masalah |
| B | Webhook, event-driven integration | 3–4 |  | Landasan teori mekanisme |
| C | HMAC / message authentication / webhook security | 3–4 |  | Landasan teori keamanan |
| D | DevOps, platform engineering, IDP, self-hosted deploy | 3–4 |  | Konteks platform (lebih tebal di TA) |
| E | Public API security, shared-secret auth | 2–3 |  | Desain endpoint publik |
| F | Software architecture (layered, orchestration agent) | 2–3 |  | Perancangan (wajib lebih dalam di TA) |
| **Total** | | **≥ 20** | **0** | |

### 2.2 Keyword pencarian

**Utama (EN):**

- `continuous deployment` / `continuous delivery` / `CI/CD automation`
- `webhook integration` / `HTTP webhook` / `event-driven webhook`
- `HMAC SHA-256` / `webhook signature verification`
- `securing public webhook endpoints`
- `internal developer platform` / `platform engineering`
- `self-hosted continuous deployment` / `container deployment automation`
- `shared secret authentication HTTP` / `API request integrity`

**Pelengkap (ID):**

- `deploy otomatis` / `continuous deployment`
- `integrasi webhook`
- `autentikasi HMAC`
- `otomatisasi DevOps`
- `platform pengembang`

**Contoh query siap tempel:**

1. `"continuous deployment" DevOps automation`
2. `"webhook" "HMAC" authentication`
3. `"webhook signature" verification security`
4. `"event-driven" webhook integration software`
5. `"internal developer platform" OR "platform engineering" deployment`
6. `"self-hosted" continuous deployment Docker`
7. `public webhook endpoint security shared secret`
8. `"message authentication code" HMAC API`

### 2.3 Kriteria lolos jurnal

Jurnal **layak** jika memenuhi minimal 2:

- [ ] Membahas konsep yang dipakai di modul (CD, webhook, HMAC, DevOps, API security, arsitektur)
- [ ] Bisa diringkas dalam 2–4 kalimat temuan
- [ ] Bisa dihubungkan ke keputusan desain Sailorport (1 kalimat relevansi)
- [ ] Sumber wajar (IEEE/ACM/Elsevier/Springer/Sinta; tahun ideal 2018–2025; teori klasik boleh lebih tua)

Jurnal **hindari** jika:

- Hanya judul mirip tapi isi ML/IoT/blockchain tanpa kaitan webhook/HMAC/CD
- Tidak ada temuan yang bisa disitasi (hanya tutorial blog tanpa peer-review — boleh pelengkap, jangan dihitung sebagai “jurnal” jika kampus melarang)

### 2.4 Tabel sitensis (isi sampai ≥ 20 baris)

| No | Penulis (tahun) | Kelompok (A–F) | Fokus singkat | Temuan / konsep utama | Relevansi ke Sailorport | Dipakai di bab |
|----|-----------------|----------------|---------------|------------------------|-------------------------|----------------|
| 1 | | | | | | |
| 2 | | | | | | |
| 3 | | | | | | |
| 4 | | | | | | |
| 5 | | | | | | |
| 6 | | | | | | |
| 7 | | | | | | |
| 8 | | | | | | |
| 9 | | | | | | |
| 10 | | | | | | |
| 11 | | | | | | |
| 12 | | | | | | |
| 13 | | | | | | |
| 14 | | | | | | |
| 15 | | | | | | |
| 16 | | | | | | |
| 17 | | | | | | |
| 18 | | | | | | |
| 19 | | | | | | |
| 20 | | | | | | |

**Cara menulis kolom Relevansi (contoh pola):**

- “Mendukung kebutuhan otomatisasi deploy setelah perubahan kode di repositori.”
- “Mendukung pemilihan webhook sebagai pemicu event dari GitHub ke API.”
- “Mendukung penggunaan HMAC untuk verifikasi keaslian payload endpoint publik.”

### 2.5 Checklist mencocokkan jurnal ↔ laporan/TA

Untuk setiap klaim di naskah, pastikan ada sitasi:

| Klaim di naskah | Kelompok jurnal | Sudah ada sitasi? |
|-----------------|-----------------|-------------------|
| Deploy manual lambat / perlu otomasi | A | |
| Webhook cocok untuk notifikasi event | B | |
| Payload harus diverifikasi (HMAC) | C | |
| Endpoint publik berisiko tanpa proteksi | E | |
| Platform/IDP memudahkan siklus deploy | D | |
| Pemisahan orkestrasi vs eksekusi (API vs agent) | F | |

---

## 3. Outline laporan (semester ini)

Sesuaikan nomor bab dengan template kampus; substansi mengikuti ini.

### BAB I Pendahuluan

- Latar belakang (masalah deploy manual → kebutuhan auto-deploy)
- Rumusan masalah (3–4 pertanyaan, fokus modul webhook)
- Tujuan & manfaat
- Batasan: hanya push GitHub; service Git; verifikasi signature; tidak bahas seluruh fitur platform secara mendalam; nama produk (Sailorport) hanya sebagai objek studi, tidak di judul

### BAB II Tinjauan Pustaka / Landasan Teori

Susun dari sitensis (bukan copy abstrak):

1. Continuous deployment / CI-CD
2. Webhook & integrasi event-driven
3. HMAC-SHA256 & keamanan webhook
4. (Singkat) DevOps / IDP sebagai konteks
5. Penelitian / praktik terkait + **celah** yang diisi modul ini

### BAB III Metodologi / Perancangan

- Objek: modul webhook Sailorport
- Alur sistem (flowchart): GitHub push → API → verify → create deployment
- Kebutuhan fungsional & non-fungsional (keamanan signature, ignore event non-push)
- Rancangan data singkat: `webhook_secret`, auto-deploy flag, env target

### BAB IV Implementasi & Pengujian

- Penjelasan implementasi mengikuti 7 langkah di §1
- Jangan dump seluruh source; kutip potongan penting + jelaskan
- Skenario uji minimal:

| ID | Skenario | Diharapkan |
|----|----------|------------|
| U1 | Signature valid + auto-deploy on + branch cocok | Deployment dibuat |
| U2 | Signature invalid / hilang | 401 Unauthorized |
| U3 | Event bukan `push` (mis. ping) | Diabaikan |
| U4 | Branch tidak cocok / auto-deploy off | Tidak membuat deploy (atau ignore sesuai perilaku sistem) |

### BAB V Penutup

- Kesimpulan menjawab rumusan masalah
- Saran → menjadi benih proposal TA (perluas ke sistem utuh, rate limit webhook, dsb.; dedupe delivery ✅ Step 33)

### Lampiran

- Contoh curl uji signature
- Screenshot portal (secret + toggle auto-deploy)
- Tabel sitensis lengkap

---

## 4. Outline TA (semester depan)

### Perluasan wajib

- Jelaskan **keseluruhan aplikasi**: portal, API berlapis, catalog, agent, deploy Git, environments, auth/RBAC (ringkas tapi utuh)
- Tetap tetapkan **kontribusi utama**: modul webhook auto-deploy (+ integrasinya ke pipeline deploy)
- Evaluasi lebih lengkap: uji fungsional sistem + uji modul fokus + pembahasan keterbatasan

### BAB khas TA (mapping)

1. Pendahuluan — masalah platform deploy self-hosted + otomasi
2. Tinjauan pustaka — pakai ulang sitensis laporan, tambah D & F
3. Metodologi — R&D / waterfall / prototyping (ikuti aturan kampus)
4. Analisis & perancangan sistem — arsitektur global + sequence webhook
5. Implementasi — stack, modul utama, fokus webhook
6. Pengujian & pembahasan
7. Kesimpulan & saran

### Pengembangan tambahan (opsional, jika dosen minta “ada pengembangan baru”)

Pilih 1 yang realistis (sudah tercatat sebagai debt/ide di proyek):

- Dedupe `X-GitHub-Delivery` (anti-replay)
- Rate limit endpoint webhook
- Kirim `git_sha` dari payload push ke deployment
- Audit log khusus event webhook

Catat pilihan di sini: **________________**

---

## 5. Guideline penulisan agar konsisten dengan jurnal

1. **Setiap bagian teori** diakhiri 1 paragraf: “Dalam penelitian/laporan ini, konsep X diterapkan sebagai Y pada Sailorport.”
2. **Jangan** sitasi jurnal yang tidak muncul di pembahasan.
3. **Utamakan sintesis**: bandingkan 2–3 sumber (“A menekankan …, B menambahkan …, pada modul ini digabung menjadi …”).
4. Istilah kunci konsisten: *webhook*, *HMAC-SHA256*, *auto-deploy*, *deployment*, *agent*, *catalog*.
5. Screenshot & uji = bukti; jurnal = justifikasi ilmiah.

---

## 6. Checklist progres

### Laporan

- [ ] Judul & batasan disepakati
- [ ] ≥ 20 jurnal masuk tabel sitensis
- [ ] Setiap kelompok A–F terisi
- [ ] Outline bab diisi substansi
- [ ] 4 skenario uji dijalankan + bukti
- [ ] Draf BAB II mensitensis (bukan menumpuk abstrak)
- [ ] Revisi akhir + daftar pustaka format kampus

### Menuju TA

- [ ] Judul TA / proposal disepakati pembimbing
- [ ] Bab arsitektur sistem global ditulis ulang lebih dalam
- [ ] Kontribusi utama tetap webhook (atau diganti tertulis di §4)
- [ ] Pengembangan tambahan (jika ada) diimplementasi + diuji
- [ ] QC teknis produk: lihat `docs/QC.md`

---

## 7. Referensi dokumen engineering (baca saat menulis teknis)

| Butuh | Baca |
|-------|------|
| Visi produk / jalur deploy | `docs/PRODUCT.md` |
| Lapisan API & web | `docs/ARCHITECTURE.md` |
| Status fitur & cara smoke webhook | `docs/PROGRESS.md` (Step 20) |
| Debt keamanan webhook | `docs/QC.md` |
| Konteks agent / katalog | `docs/AGENTS.md` |

---

## 8. Catatan keputusan akademik (isi sendiri)

| Tanggal | Keputusan | Alasan |
|---------|-----------|--------|
| 2026-09-15 | Fokus laporan = webhook auto-deploy | Satu modul jelas, cocok sitensis & uji |
| 2026-09-15 | Judul laporan tanpa nama produk | Nama Sailorport buatan sendiri; judul pakai istilah IDP |
| 2026-09-15 | Judul final laporan = Implementasi Fitur Auto-Deploy Berbasis Webhook GitHub pada Internal Developer Platform | Disepakati penulis |
| | Judul final TA = … | |
| | Pengembangan tambahan TA = … | |
