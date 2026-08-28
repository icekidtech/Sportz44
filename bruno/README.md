# Sportz44 API — Bruno Collection

This folder contains a [Bruno](https://www.usebruno.com/) collection for testing the
Sportz44 backend API. Bruno is a lightweight, offline-first API client that stores
collections as plain files (`.bru`) so they can be versioned in the repo.

## Setup

1. Install [Bruno](https://www.usebruno.com/downloads) (desktop app).
2. Open Bruno → **Open Collection** → select this `bruno/` folder.
3. Select the **Local** environment (top-right dropdown) — it points at
   `http://localhost:8080`.

## Prerequisites

- Postgres + Redis running locally (native macOS services).
- Backend running: `cd backend && go run ./cmd/api`
- (Optional) Match-listener running for live updates:
  `cd backend && go run ./cmd/match-listener`

## Suggested flow

1. **System → Health Check** — confirm the server is up.
2. **Auth → Register** — create a test user (uses `username`/`email`/`password`
   from the environment). If the user already exists, skip to Login.
3. **Auth → Login** — sets the HttpOnly access/refresh cookies in Bruno's cookie
   jar automatically. Subsequent protected requests will carry them.
4. **Auth → Me** — verify the session works.
5. **Matches → Sync Competition** — ingest fixtures for a competition/season
   (defaults: `competition=PD`, `season=2026`).
6. **Matches → List Matches** (and the filtered/paginated variants) — query the
   ingested data.
7. **Matches → Get Match / Events / Lineup** — inspect a single match
   (`match_id` from the environment, default `1`).
8. **System → Live Match Updates (WebSocket)** — connect to `ws://localhost:8080/ws`
   to observe live match events broadcast by the match-listener.

## Notes

- The match endpoints are protected by cookie auth. Bruno stores cookies
  automatically after Login, so no manual header setup is needed.
- The `/ws` WebSocket endpoint is public.
- Environment variables live in `environments/Local.bru`. Edit values there
  (e.g. `base_url`, credentials, `match_id`) rather than hardcoding them in
  requests.
