#!/usr/bin/env bash
# deploy.sh — run on the VPS to pull the latest release, build, and restart.
#
# This is invoked by the GitHub Actions CD workflow (or manually) after a
# successful CI build. It assumes:
#   - The repo is cloned at $APP_DIR (default: /opt/sportz44)
#   - PM2 is installed and the ecosystem file is registered
#   - backend/.env exists on the server (secrets are NOT in the repo)
#
# Usage:
#   ./deploy/deploy.sh [app_dir]
#
# Env vars:
#   APP_DIR   — repo location on the server (default /opt/sportz44)

set -euo pipefail

APP_DIR="${APP_DIR:-/opt/sportz44}"
cd "$APP_DIR"

echo "==> Pulling latest code..."
git fetch --all --prune
git checkout "$(git describe --tags --abbrev=0 2>/dev/null || echo main)"
git pull --ff-only

echo "==> Building binary..."
./deploy/build.sh

echo "==> Restarting PM2 process..."
pm2 reload deploy/ecosystem.config.js --update-env || pm2 start deploy/ecosystem.config.js
pm2 save

echo "==> Deploy complete."
pm2 status
