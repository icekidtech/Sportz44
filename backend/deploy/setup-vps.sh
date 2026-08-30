#!/usr/bin/env bash
# setup-vps.sh — one-time provisioning of a fresh Ubuntu VPS for Sportz44.
#
# Installs: Go toolchain, PostgreSQL, Redis, PM2, Nginx, and clones the repo.
# Run as a user with sudo (NOT as root).
#
# Usage:
#   ./deploy/setup-vps.sh
#
# After this completes, you must:
#   1. Create backend/.env with real secrets (see backend/.env.example)
#   2. Run: pm2 start deploy/ecosystem.config.js && pm2 save && pm2 startup
#   3. Configure Nginx (see deploy/nginx.conf) and point DNS at the server

set -euo pipefail

echo "==> Updating system packages..."
sudo apt-get update && sudo apt-get upgrade -y

echo "==> Installing build tools..."
sudo apt-get install -y build-essential git curl ca-certificates gnupg

echo "==> Installing Go..."
GO_VERSION="1.26"
if ! command -v go >/dev/null 2>&1; then
  curl -OL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
  sudo rm -rf /usr/local/go
  sudo tar -C /usr/local -xzf "go${GO_VERSION}.linux-amd64.tar.gz"
  rm "go${GO_VERSION}.linux-amd64.tar.gz"
  echo 'export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin' | sudo tee /etc/profile.d/go.sh
  export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
fi
go version

echo "==> Installing PostgreSQL..."
if ! command -v psql >/dev/null 2>&1; then
  sudo apt-get install -y postgresql postgresql-contrib
  sudo systemctl enable --now postgresql
fi

echo "==> Creating sportz44 database and user..."
sudo -u postgres psql -tc "SELECT 1 FROM pg_roles WHERE rolname='sportz44'" | grep -q 1 || \
  sudo -u postgres psql -c "CREATE USER sportz44 WITH PASSWORD 'change_me';"
sudo -u postgres psql -tc "SELECT 1 FROM pg_database WHERE datname='sportz44'" | grep -q 1 || \
  sudo -u postgres createdb -O sportz44 sportz44

echo "==> Installing Redis..."
if ! command -v redis-server >/dev/null 2>&1; then
  sudo apt-get install -y redis-server
  sudo systemctl enable --now redis-server
fi

echo "==> Installing Node.js + PM2..."
if ! command -v node >/dev/null 2>&1; then
  curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
  sudo apt-get install -y nodejs
fi
sudo npm install -g pm2

echo "==> Installing Nginx..."
if ! command -v nginx >/dev/null 2>&1; then
  sudo apt-get install -y nginx
  sudo systemctl enable --now nginx
fi

echo "==> Cloning Sportz44 repo..."
APP_DIR="${APP_DIR:-/opt/sportz44}"
if [ ! -d "$APP_DIR" ]; then
  sudo mkdir -p "$APP_DIR"
  sudo chown "$USER":"$USER" "$APP_DIR"
  git clone https://github.com/icekidtech/Sportz44.git "$APP_DIR"
fi

echo "==> Creating log directory..."
sudo mkdir -p /var/log/sportz44
sudo chown "$USER":"$USER" /var/log/sportz44

echo ""
echo "============================================================"
echo " VPS setup complete. Next steps:"
echo "   1. Create $APP_DIR/backend/.env with real secrets"
echo "      (see backend/.env.example)"
echo "   2. cd $APP_DIR && ./deploy/build.sh"
echo "   3. pm2 start deploy/ecosystem.config.js && pm2 save"
echo "   4. Run 'pm2 startup' and paste the printed command"
echo "   5. Configure Nginx with deploy/nginx.conf"
echo "============================================================"
