# Setup Mesin Baru (Sailorport)

Panduan ini dipakai saat pertama kali membuka proyek di laptop/komputer lain.

## Prasyarat

| Tool | Untuk apa | Cek versi |
|------|-----------|-----------|
| Git | sinkron kode antar mesin | `git --version` |
| Go 1.26+ | API, worker, agent | `go version` |
| Docker + Compose | Postgres, Redis, self-host | `docker --version` & `docker compose version` |
| Node.js 22+ | portal web | `node --version` |

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
`postgres://sailorport:sailorport@localhost:5433/sailorport?sslmode=disable`

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

## Jalankan portal web

```bash
cd apps/web
npm install
npm run dev
```

Buka `http://localhost:5173` — login/register, lalu portal dengan sidebar (bisa collapse di desktop): Overview, Catalog, Workers.

Vite mem-proxy `/api` dan `/healthz` ke API di `:8080`. API harus sudah jalan.

Portal mendukung: auth (JWT), catalog CRUD + scaffold + **Deploy** (dialog status), worker list, overview dashboard.

## Jalankan agent (Step 10B + 10C.2)

Terminal terpisah, API harus sudah jalan. Agent butuh **Docker CLI** di PATH untuk build/run deploy.

```bash
cd apps/agent
SAILORPORT_API_URL=http://localhost:8080 \
SAILORPORT_WORKER_NAME=local-dev \
SAILORPORT_HEARTBEAT_INTERVAL=15s \
SAILORPORT_POLL_INTERVAL=5s \
SAILORPORT_DEPLOY_PORT_BASE=18080 \
go run .
```

Cek portal `/worker` — worker muncul **online** dengan heartbeat berkala.

Register/heartbeat/claim/update deploy **tanpa JWT** (endpoint publik untuk agent).

**Test deploy (curl, setelah scaffold service baru):**

```bash
# 1. Login → dapat token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"you@example.com","password":"yourpass"}' | jq -r .token)

# 2. Buat deployment pending (ganti SERVICE_ID)
curl -X POST "http://localhost:8080/api/v1/services/SERVICE_ID/deployments" \
  -H "Authorization: Bearer $TOKEN"

# 3. Agent poll → build → run; cek health service
curl http://localhost:18080/healthz
```

## Environment variables

| Variable | Default | Keterangan |
|----------|---------|------------|
| `PORT` | `8080` | port HTTP API |
| `APP_ENV` | `development` | environment label |
| `APP_VERSION` | `0.1.0` | versi API di response health |
| `DATABASE_URL` | lihat di atas (port **5433**) | koneksi Postgres |
| `AUTH_JWT_SECRET` | `dev-only-change-me` | secret JWT (ganti di production) |

Agent (`apps/agent`):

| Variable | Default | Keterangan |
|----------|---------|------------|
| `SAILORPORT_API_URL` | `http://localhost:8080` | base URL API |
| `SAILORPORT_WORKER_NAME` | hostname | nama worker |
| `SAILORPORT_HEARTBEAT_INTERVAL` | `15s` | interval heartbeat |
| `SAILORPORT_POLL_INTERVAL` | `5s` | interval poll job deploy |
| `SAILORPORT_DEPLOY_PORT_BASE` | `18080` | host port container deploy |

Contoh:

```bash
PORT=9090 APP_ENV=production APP_VERSION=0.2.0 go run .
```

## Struktur kode saat ini

```text
apps/api/
├── main.go
└── internal/
    ├── config/, db/, migrate/
    ├── handler/   (health, services, scaffold, auth, workers, deployments, middleware)
    ├── service/   (catalog, scaffold, auth, workers, deployments)
    ├── store/     (service, user, worker, deployment)
    ├── model/     (service, user, worker, deployment)
    ├── template/  (registry + generate dari disk)
    └── auth/      (jwt, password)

apps/agent/
├── main.go
└── internal/
    ├── config/
    ├── client/    (HTTP ke API: workers + deployments)
    ├── docker/    (docker build + run)
    └── agent/     (register + heartbeat + poll/deploy loop)

apps/web/
├── vite.config.ts
├── components.json
└── src/
    ├── App.tsx
    ├── index.css          # Tailwind v4 + harbour theme
    ├── components/
    │   ├── app/           # DataPanel, Toolbar, EmptyState, …
    │   └── ui/            # shadcn primitives + dialog
    ├── layouts/           # AuthLayout, AppShell
    ├── features/
    │   ├── auth/
    │   ├── catalog/
    │   ├── deployments/   # Deploy API client + DeploymentsDialog
    │   ├── overview/
    │   ├── scaffold/      # CreateServiceForm
    │   └── workers/
    └── lib/               # http.ts, utils.ts, theme.ts
```

## Setelah setup

1. Baca `docs/PROGRESS.md` — lihat step terakhir yang selesai
2. Buka chat Cursor baru
3. Copy prompt dari `docs/RESUME-PROMPT.md`
4. Lanjut step berikutnya
