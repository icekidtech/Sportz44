# Sportz44 Deployment

This directory contains everything needed to deploy the Sportz44 backend to a
VPS, managed by **PM2**, with **GitHub Actions** for CI/CD.

## Architecture

The backend runs as a **single process** (`sportz44-api`). The live match
listener runs as an in-process goroutine inside the API server — no separate
binary or process is needed, saving memory and simplifying operations.

```
┌─────────────────────────────────────────────┐
│              VPS (Ubuntu)                   │
├─────────────────────────────────────────────┤
│  Nginx (reverse proxy + TLS)                │
│   └─ /api/*  → 127.0.0.1:8080               │
│   └─ /ws     → 127.0.0.1:8080 (WebSocket)   │
│                                             │
│  PM2 → sportz44-api (backend/bin/api)       │
│       ├─ REST API (Gin)                     │
│       ├─ WebSocket hub (Redis pub/sub)      │
│       └─ Match listener goroutine (30s)     │
│                                             │
│  PostgreSQL (5432) │ Redis (6379)           │
└─────────────────────────────────────────────┘
```

## Files

| File | Purpose |
|------|---------|
| `ecosystem.config.js` | PM2 process config (single `sportz44-api` app) |
| `build.sh` | Compiles `backend/bin/api` (production flags) |
| `deploy.sh` | Pulls latest code, builds, reloads PM2 (used by CD) |
| `setup-vps.sh` | One-time provisioning of a fresh Ubuntu VPS |
| `nginx.conf` | Reverse proxy for REST + WebSocket + TLS |

## CI/CD (GitHub Actions)

CI and CD are **separate workflow files** under `.github/workflows/`:

- **`ci.yml`** — runs on every push/PR to `main` and `development`:
  gofmt check → `go vet` → `go build` → `go test -race`.
- **`cd.yml`** — deploys to the VPS on every push to `main` (after CI passes),
  or manually via the Actions tab.

### Required GitHub secrets (for CD)

| Secret | Description |
|--------|-------------|
| `VPS_HOST` | Server IP or hostname |
| `VPS_USER` | SSH user (e.g. `deploy`) |
| `VPS_SSH_KEY` | Private SSH key for the deploy user |
| `VPS_APP_DIR` | Repo path on the server (default `/opt/sportz44`) |

## One-time VPS setup

```bash
# On the VPS, as a sudo user:
git clone https://github.com/icekidtech/Sportz44.git /opt/sportz44
cd /opt/sportz44
./deploy/setup-vps.sh        # installs Go, Postgres, Redis, Node, PM2, Nginx

# Create the env file with real secrets:
cp backend/.env.example backend/.env
# ... edit backend/.env ...

# Build and start:
./deploy/build.sh
pm2 start deploy/ecosystem.config.js
pm2 save
pm2 startup                 # run the printed command to enable boot persistence

# Configure Nginx:
sudo cp deploy/nginx.conf /etc/nginx/sites-available/sportz44
sudo ln -s /etc/nginx/sites-available/sportz44 /etc/nginx/sites-enabled/
sudo nginx -t && sudo systemctl reload nginx

# TLS (optional but recommended):
sudo apt-get install -y certbot python3-certbot-nginx
sudo certbot --nginx -d api.sportz44.app
```

## Deploying manually

```bash
./deploy/deploy.sh          # pulls latest, builds, reloads PM2
```

## Useful PM2 commands

```bash
pm2 status                  # list processes
pm2 logs sportz44-api       # tail logs
pm2 restart sportz44-api    # restart
pm2 reload sportz44-api     # zero-downtime reload
pm2 monit                   # live resource monitor
```
