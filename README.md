# Sportz44

**Multi-Club Football Analytics & Community Platform**

A mobile-first, real-time football intelligence hub for African football fans — combining live match data, predictive analytics, vibrant community engagement, and betting insights. Built with a 100% Go backend and a React Native mobile app.

> Full product specification: [`docs/SPORT_PLATFORM_PRD.md`](docs/SPORT_PLATFORM_PRD.md)

---

## Features

| Feature | Description |
|---------|-------------|
| **Real-Time Match Hub** | Live scores, lineups, and minute-by-minute events (goals, cards, subs, injuries) updated every 30s |
| **Squad Analytics** | Player cards, season stats, form trends, injury tracking, and side-by-side comparisons |
| **Fixture Calendar** | 30-day fixture window with customizable push / WhatsApp / SMS notifications |
| **News Aggregation** | Club-specific news auto-pulled from RSS feeds (Marca, AS, ESPN, Sky Sports, BBC, etc.) |
| **Fan Community** | Per-match forum threads + live WebSocket chat with moderation |
| **Betting Insights** | Win/draw/loss probabilities, xG visualization, odds comparison, and prediction leaderboards |
| **WhatsApp Bot** | Twilio-powered bot for fixtures, results, top scorers, and live score subscriptions |
| **Admin Dashboard** | Analytics, content management, community moderation, and audit logs |
| **Multi-Club Support** | Follow multiple clubs (e.g. Real Madrid + Manchester City) simultaneously |

---

## Tech Stack

| Component | Technology |
|-----------|-----------|
| API Server | Go (Gin) |
| Background Jobs | Go (goroutines) |
| Real-Time | Go WebSocket + Redis Pub/Sub |
| Database | PostgreSQL (schema managed by **GORM AutoMigrate**) |
| Cache / Queue | Redis |
| Mobile App | React Native (Expo) |
| Admin Panel | React.js |
| Notifications | Firebase Cloud Messaging + Twilio WhatsApp |
| External Data | API-Sports, Odds API, RSS feeds |

---

## Architecture

```
┌─────────────────────────────────┐
│         Your VPS (Ubuntu)       │
├─────────────────────────────────┤
│ Nginx (Reverse Proxy, SSL)      │
│  └─ Go API Server (port 8080)   │
│       • /api/*   REST API       │
│       • /ws      WebSocket      │
│       • /webhook Twilio         │
│       • /admin   Admin Panel UI │
│ PostgreSQL (5432) │ Redis (6379)│
│ Go Background Services:         │
│  ├─ Match Listener              │
│  ├─ News Aggregator             │
│  ├─ Odds Sync                   │
│  └─ Injury Tracker              │
└─────────────────────────────────┘
```

**Data flow:** External APIs → Go services (poll + dedupe) → Redis Pub/Sub → WebSocket broadcast / Push / WhatsApp.

---

## Project Structure

```
sportz44/
├── backend/
│   ├── cmd/
│   │   └── api/              # Single API server (port 8080): REST + WebSocket + match listener goroutine
│   ├── internal/
│   │   ├── api/              # handlers, middleware, routes
│   │   ├── models/           # GORM models (AutoMigrate source of truth)
│   │   ├── services/         # Business logic (incl. in-process match listener)
│   │   ├── repository/       # Database access
│   │   ├── external/         # API-Sports, Football-Data, FlashScore, TheSportsDB clients
│   │   └── ws/               # WebSocket hub / client
│   ├── pkg/                  # config, database, cache, jwt, logger
│   ├── go.mod
│   └── .env.example
├── deploy/                   # PM2, build/deploy scripts, nginx, VPS setup
├── .github/workflows/        # CI + CD (separate files)
├── bruno/                    # API test collection (Bruno)
├── app/                      # React Native app
├── admin-panel/              # React.js admin dashboard
└── docs/                     # PRD, API docs, architecture, schema, deployment
```

---

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL & Redis
- API-Sports key (and optionally Odds API, Twilio, Firebase credentials)

### Local Setup

```bash
# 1. Clone
git clone https://github.com/icekidtech/Sportz44.git
cd Sportz44

# 2. Configure environment
cp backend/.env.example backend/.env
# Edit backend/.env with your API keys and DB credentials

# 3. Create database (schema is auto-migrated by GORM on first boot)
createdb sportz44
redis-cli FLUSHALL

# 4. Build all Go binaries
cd backend
go build -o bin/api ./cmd/api

# 5. Run the API server (port 8080): REST + WebSocket + in-process match listener
./bin/api
```

> The API server runs the live match listener as an in-process goroutine, so
> only **one** process is needed.

### Key Environment Variables

```bash
DATABASE_URL=postgresql://user:password@localhost:5432/sportz44
REDIS_URL=redis://localhost:6379
API_SPORTS_KEY=your_api_sports_key
JWT_SECRET=your_super_secret_key_min_32_chars
ADMIN_EMAIL=your@email.com
ADMIN_PASSWORD=secure_password
```

See [`docs/SPORT_PLATFORM_PRD.md` → Appendix B](docs/SPORT_PLATFORM_PRD.md) for the full list.

---

## Deployment

The backend is deployed to an Ubuntu VPS behind **Nginx + PM2**, with **GitHub
Actions** for CI/CD. The database schema is created automatically via GORM
`AutoMigrate` on first boot — no manual migrations required.

- **PM2 + CI/CD + VPS setup:** see [`deploy/README.md`](deploy/README.md)
- **Nginx config:** [`deploy/nginx.conf`](deploy/nginx.conf)
- **CI workflow:** [`.github/workflows/ci.yml`](.github/workflows/ci.yml)
- **CD workflow:** [`.github/workflows/cd.yml`](.github/workflows/cd.yml)

---

## Documentation

- **Product Requirements:** [`docs/SPORT_PLATFORM_PRD.md`](docs/SPORT_PLATFORM_PRD.md)
- API endpoints, database schema, caching strategy, and monitoring details are all covered in the PRD.

---

## License

See repository for license details.
