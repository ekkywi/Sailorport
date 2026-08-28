# Cara Lanjut di Mesin Lain

Chat Cursor **tidak tersinkron** antar device. Yang tersinkron adalah **repo Git** + file di folder `docs/`.

## Sebelum tutup / pindah mesin (di kantor)

```bash
cd ~/Projects/Sailorport

# 1. Cek perubahan
git status

# 2. Simpan semua
git add .
git commit -m "docs: update progress"

# 3. Push ke GitHub
git push
```

Kalau `git push` gagal karena belum ada remote, lihat bagian **Setup remote pertama kali** di bawah.

## Di rumah (atau mesin lain)

```bash
# pertama kali di mesin ini
git clone git@github.com:ekkywi/Sailorport.git
cd Sailorport

# atau kalau sudah pernah clone
cd ~/Projects/Sailorport
git pull
```

Lalu:

1. Buka folder `Sailorport` di Cursor
2. Baca `docs/PRODUCT.md` (visi) lalu `docs/PROGRESS.md` (step)
3. Buka chat baru
4. Paste isi `docs/RESUME-PROMPT.md`
5. Lanjut step berikutnya (lihat **Next action** di `docs/PROGRESS.md`)

## Setup remote pertama kali (sekali saja)

Jika repo belum ada di GitHub:

1. Buat repo kosong di GitHub: `ekkywi/Sailorport` (tanpa README)
2. Jalankan di mesin kantor:

```bash
cd ~/Projects/Sailorport
git branch -M main
git remote add origin git@github.com:ekkywi/Sailorport.git
git push -u origin main
```

HTTPS (kalau belum pakai SSH):

```bash
git remote add origin https://github.com/ekkywi/Sailorport.git
git push -u origin main
```

## Ritual setiap selesai 1 step

1. Update checklist di `docs/PROGRESS.md`
2. Isi "Terakhir dikerjakan" dan "Mesin terakhir"
3. Tulis 1–3 poin di "Catatan belajar pribadi"
4. `git commit` + `git push`

Dengan ritual ini, Anda tidak perlu mengandalkan ingatan chat.

## File penting

| File | Fungsi |
|------|--------|
| `docs/PRODUCT.md` | **visi produk** — dua jalur deploy, positioning IDP |
| `docs/PROGRESS.md` | step mana yang sudah/belum |
| `docs/SETUP.md` | install tool di mesin baru |
| `docs/AGENTS.md` | konteks untuk AI |
| `docs/RESUME-PROMPT.md` | teks siap paste ke chat baru |
| `docs/ROADMAP.md` | peta besar proyek |
