# ByFood Assessment (Monorepo)

Full-stack implementation of the ByFood assessment.

- `server/` — Golang REST API + Postgres + Swagger + tests
- `client/` — Next.js app (App Router, TypeScript, Tailwind) with API route proxies, modals, and toast context

## Quick Start

### 1) Server

```bash
cd server
# Optional: spin Postgres in Docker and run API with migrations
make run/dev
# or run just the API if DB is already up
make run
```

### 2) Client

```bash
cd client
# (ensure API_BASE_URL=http://localhost:4748)
npm install
npm run dev
```
