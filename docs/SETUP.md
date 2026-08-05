# Setup Mesin Baru (Sailorport)

Panduan ini dipakai saat pertama kali membuka proyek di laptop/komputer lain.

## Prasyarat

| Tool | Untuk apa | Cek versi |
|------|-----------|-----------|
| Git | sinkron kode antar mesin | `git --version` |
| Go 1.26+ | API, worker, agent | `go version` |
| Docker + Compose | Postgres, Redis, self-host | `docker --version` & `docker compose version` |
| Node.js (nanti Step 6) | portal web | `node --version` |

## Install cepat (Ubuntu / Debian)

```bash
sudo apt update
sudo apt install -y git golang-go docker.io docker-compose-plugin
```

Aktifkan Docker tanpa sudo (opsional, logout-login setelah ini):

```bash
sudo usermod -aG docker $USER
```

## Clone project

```bash
git clone git@github.com:ekkywi/Sailorport.git
cd Sailorport
```

Kalau belum pakai SSH key, pakai HTTPS:

```bash
git clone https://github.com/ekkywi/Sailorport.git
cd Sailorport
```

## Jalankan Postgres

```bash
cd deploy/compose
docker compose up -d
docker compose ps
```

Connection string default:
`postgres://sailorport:sailorport@localhost:5432/sailorport?sslmode=disable`

## Jalankan API lokal

```bash
cd apps/api
go run .
```

Startup yang diharapkan:

```text
Database OK (SELECT 1)
Database migrations OK
Sailorport API (development) running on http://localhost:8080
```

Server di `http://localhost:8080` (atau port dari env `PORT`).

## Test cepat

Terminal baru (biarkan server tetap jalan):

```bash
# health check
curl http://localhost:8080/healthz

# echo endpoint
curl -X POST http://localhost:8080/api/v1/echo \
  -H "Content-Type: application/json" \
  -d '{"message":"hello sailorport"}'

# list services
curl http://localhost:8080/api/v1/services

# create service
curl -X POST http://localhost:8080/api/v1/services \
  -H "Content-Type: application/json" \
  -d '{"name":"payments-api","description":"Payment service","owner":"platform"}'
```

## Environment variables

| Variable | Default | Keterangan |
|----------|---------|------------|
| `PORT` | `8080` | port HTTP API |
| `APP_ENV` | `development` | environment label |
| `APP_VERSION` | `0.1.0` | versi API di response health |
| `DATABASE_URL` | lihat di atas | koneksi Postgres |

Contoh:

```bash
PORT=9090 APP_ENV=production APP_VERSION=0.2.0 go run .
```

## Struktur kode API saat ini

```text
apps/api/
├── main.go
├── go.mod
└── internal/
    ├── config/config.go
    ├── db/db.go
    ├── handler/
    │   ├── health.go
    │   ├── echo.go
    │   └── service.go
    ├── migrate/
    │   ├── migrate.go
    │   └── migrations/00001_create_services.sql
    ├── model/service.go
    └── store/service.go
```

## Setelah setup

1. Baca `docs/PROGRESS.md` — lihat step terakhir yang selesai
2. Buka chat Cursor baru
3. Copy prompt dari `docs/RESUME-PROMPT.md`
4. Lanjut step berikutnya
