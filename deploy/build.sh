#!/usr/bin/env bash
# build.sh — compile the Sportz44 backend binary for production.
#
# Produces:
#   backend/bin/api — the REST + WebSocket API server (also runs the live
#                     match listener as an in-process goroutine)
#
# Usage:
#   ./deploy/build.sh
#
# Set GOOS/GOARCH to cross-compile (e.g. GOOS=linux GOARCH=amd64 for a Linux VPS).

set -euo pipefail

cd "$(dirname "$0")/../backend"

echo "==> Building sportz44-api..."
CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/api ./cmd/api

echo "==> Done. Binary in backend/bin/:"
ls -lh bin/
