# Sailorport — AI Context Brief (Gemini / NotebookLM)

> Upload atau tempel dokumen ini sebagai **sumber konteks tetap**.
> Lalu unggah jurnal PDF / bib sebagai sumber tambahan.
> Tujuan AI: membantu sitensis jurnal, outline bab, dan menyambungkan teori ke modul nyata — **bukan** mengarang fitur yang tidak ada.

**Bahasa kerja:** Indonesia (kecuali kutipan istilah teknis EN yang sudah baku).  
**Peran AI:** asisten riset & penulisan akademik untuk mahasiswa konsentrasi software / pengembangan perangkat lunak.

**Dokumen terkait:** `docs/ACADEMIC-GUIDE.md` (outline laporan→TA, tabel sitensis, checklist).

---

## 1. Apa yang harus AI pahami dulu

### Produk

**Sailorport** = self-hosted **Internal Developer Platform (IDP)** OSS.

Tagline: *Self-hosted developer port — catalog, deploy, and ship.*

Komponen utama:

- **Web portal** (`apps/web`) — UI catalog, workers, users, deploy
- **API** (`apps/api`) — Go, berlapis: handler → service → store → model
- **Agent** (`apps/agent`) — di node worker; git/build/run Docker; poll job dari API
- **PostgreSQL** — data persistensi
- **Catalog** — inventory pusat semua service yang dikelola platform

Prinsip arsitektur penting:

- **Control plane ≠ runtime** — API tidak menjalankan container; agent yang eksekusi
- Deploy Git: repo + Dockerfile → agent clone/pull → `docker build` → `docker run`
- Ada juga jalur catalog apps (pull image) dan scaffold template (opsional)

### Fokus akademik SAAT INI (laporan)

Bukan seluruh aplikasi secara mendalam.

**Judul laporan (tanpa nama produk):**  
Implementasi Fitur Auto-Deploy Berbasis Webhook GitHub pada Internal Developer Platform

**Penamaan:** jangan pakai “Sailorport” di judul. Di naskah boleh disebut sebagai studi kasus / platform yang dikembangkan penulis.

**Modul inti:** webhook GitHub → verifikasi signature → auto-deploy

**Endpoint:** `POST /api/v1/webhooks/github` (publik, tanpa JWT)

**Alur wajib dipahami:**

1. GitHub mengirim event HTTP (webhook)
2. API hanya memproses event `push` (ping/tag diabaikan sesuai implementasi)
3. Parse repo + branch dari payload
4. Cari service di catalog yang cocok `repo_url` / clone URL
5. Verifikasi **HMAC-SHA256** via header `X-Hub-Signature-256` + `webhook_secret`
6. Jika auto-deploy aktif → buat **deployment**
7. Agent (sudah ada) yang kemudian build/run — untuk laporan cukup disebut sebagai tahap lanjutan, bukan fokus utama

**File implementasi (referensi, jangan dikarang ulang seolah beda sistem):**

- `apps/api/internal/handler/webhook.go`
- `apps/api/internal/service/webhook.go`
- `apps/api/internal/service/webhook_signature.go`
- `apps/api/internal/model/webhook.go`

### Fokus nanti (TA — jangan dicampur ke laporan kecuali diminta)

- Jelaskan aplikasi **keseluruhan**
- Kontribusi utama tetap modul webhook auto-deploy
- Judul TA disarankan (tanpa nama produk):  
  Pengembangan Internal Developer Platform Self-hosted dengan Fitur Auto-Deploy Berbasis Webhook GitHub

---

## 2. Cara memakai jurnal yang diunggah user

Jika user mengunggah jurnal / paper:

1. **Klasifikasikan** ke salah satu kelompok:
   - **A** CI/CD / continuous deployment
   - **B** Webhook / event-driven integration
   - **C** HMAC / message authentication / webhook security
   - **D** DevOps / platform engineering / IDP
   - **E** Public API security / shared secret
   - **F** Software architecture / orchestration
2. Ringkas **temuan inti** (2–4 kalimat), bukan rewrite abstrak panjang.
3. Tulis **1 kalimat relevansi** ke Sailorport, contoh:
   - “Mendukung justifikasi auto-deploy setelah push ke repositori.”
   - “Mendukung pemilihan HMAC untuk otentikasi payload webhook publik.”
4. Jangan paksakan relevansi jika paper tidak nyambung; katakan “kurang relevan untuk modul ini”.
5. Utamakan **sintesis antar paper** (bandingkan/gabungkan), bukan daftar abstrak satu per satu.

### Format sitensis yang diminta (default)

| No | Penulis (tahun) | Kelompok A–F | Fokus | Temuan utama | Relevansi ke Sailorport | Saran bab |
|----|-----------------|--------------|-------|--------------|-------------------------|-----------|

---

## 3. Aturan ketat untuk AI (anti-halusinasi)

- Jangan mengarang endpoint, tabel DB, atau fitur yang tidak disebut di brief / sumber user.
- Jangan mengklaim paper mengatakan sesuatu yang tidak ada di teks yang diberikan.
- Jika konteks teknis kurang: tanya user atau tandai `[PERLU KONFIRMASI]`.
- Jangan perluas laporan menjadi pembahasan seluruh IDP kecuali user bilang “untuk TA”.
- Istilah konsisten: webhook, HMAC-SHA256, auto-deploy, deployment, agent, catalog, control plane.
- Keamanan: sebutkan bahwa secret di-redact di response portal; verifikasi signature penting karena endpoint publik.

### Yang boleh disebut sebagai keterbatasan (jika relevan, jangan dilebih-lebihkan)

- Endpoint webhook publik perlu proteksi signature
- Ide perbaikan masa depan (bukan klaim sudah ada): dedupe `X-GitHub-Delivery`, rate limit, memakai `git_sha` dari payload push

---

## 4. Prompt siap pakai (copy ke Gemini / chat NotebookLM)

### A) Sitensis satu jurnal

```
Berdasarkan AI Context Brief Sailorport dan jurnal yang saya unggah:
1) Klasifikasikan ke kelompok A–F
2) Ringkas temuan (2–4 kalimat)
3) Tuliskan relevansi ke modul webhook auto-deploy Sailorport
4) Sarankan di bab mana sitasi ini dipakai (laporan)
Jangan mengarang fitur di luar brief.
```

### B) Sitensis batch (≥20)

```
Susun tabel sitensis untuk semua jurnal yang diunggah.
Kolom: No | Penulis (tahun) | Kelompok | Fokus | Temuan | Relevansi Sailorport | Bab.
Pastikan sebaran kelompok A–F seimbang jika memungkinkan.
Tandai jurnal yang lemah/tidak relevan.
```

### C) Draft landasan teori

```
Tulis draf BAB II (landasan teori) untuk laporan modul webhook auto-deploy Sailorport.
Gunakan hanya brief + jurnal terunggah.
Struktur: CI/CD → webhook → HMAC/keamanan → konteks DevOps singkat → celah/posisi penelitian.
Akhiri setiap subbab dengan 1 paragraf penerapan pada Sailorport.
Gaya akademik Indonesia, sintesis (bukan tumpukan abstrak).
```

### D) Cocokkan klaim ↔ sitasi

```
Cek klaim berikut dan sarankan sitasi dari jurnal terunggah:
- Deploy manual perlu otomasi
- Webhook cocok untuk event push
- Payload harus diverifikasi HMAC
- Endpoint publik berisiko tanpa autentikasi pesan
Jika tidak ada jurnal yang cocok, bilang kekurangan sitasi apa yang harus dicari.
```

### E) Mode TA (hanya jika diminta)

```
Bantu outline TA: sistem Sailorport keseluruhan + kontribusi utama webhook auto-deploy.
Bedakan jelas mana konteks sistem dan mana fokus penelitian.
```

---

## 5. Outline laporan (agar AI tidak melebar)

1. **Pendahuluan** — masalah deploy manual; tujuan modul webhook
2. **Landasan teori** — dari jurnal kelompok A–E (F opsional singkat)
3. **Perancangan** — flowchart push → verify → deployment
4. **Implementasi & uji** — 7 langkah modul + skenario U1–U4
5. **Penutup** — kesimpulan + saran menuju TA

Skenario uji minimal:

- U1 signature valid + auto-deploy on → deployment dibuat
- U2 signature invalid → 401
- U3 event non-push → diabaikan
- U4 branch tidak cocok / auto-deploy off → tidak deploy

---

## 6. Keyword pencarian jurnal (jika AI diminta menyarankan query)

- `"continuous deployment" DevOps automation`
- `"webhook" "HMAC" authentication`
- `"webhook signature" verification security`
- `"event-driven" webhook integration software`
- `"internal developer platform" OR "platform engineering"`
- `public webhook endpoint security shared secret`

---

## 7. Paket sumber yang disarankan di NotebookLM

Unggah bersama:

1. Dokumen ini (`ACADEMIC-AI-BRIEF.md`)
2. `ACADEMIC-GUIDE.md`
3. Cuplikan singkat PRODUCT / arsitektur (opsional)
4. PDF jurnal (target ≥ 20)

Lalu instruksikan: **“Utamakan brief ini sebagai kebenaran produk; jurnal hanya untuk teori & sitasi.”**

---

## 8. Identitas singkat untuk jawaban AI

Jika diminta menjelaskan objek penelitian dalam 3 kalimat:

> Objek laporan adalah modul auto-deploy berbasis webhook GitHub pada sebuah internal developer platform (IDP) self-hosted (studi kasus: Sailorport). Modul menerima event push di endpoint publik, memverifikasi keaslian payload dengan HMAC-SHA256, lalu membuat deployment bila auto-deploy aktif. Eksekusi build/run dilakukan agent terpisah; laporan berfokus pada jalur webhook hingga terciptanya deployment.

Jika diminta menulis **judul**, gunakan versi tanpa nama produk di atas — jangan mengembalikan judul yang memuat “Sailorport”.
