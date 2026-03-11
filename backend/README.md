# Backend (Gin + PostgreSQL)

## Requirements
- Go 1.22+
- PostgreSQL 14+ (or Docker)
- golang-migrate (CLI)

## Quick Start
1. Start database (Docker)
   - `docker compose up -d`
2. Copy env
   - `Copy-Item .env.example .env`
3. Set `DATABASE_URL` in `.env`
4. Run migrations
   - `migrate -path migrations -database "$env:DATABASE_URL" up`
5. Run API
   - `go run ./cmd/api`

## Tenant Context
- Tenant ID is required in `X-Tenant-ID` for `/api/*`
- Example:
  - `curl -H "X-Tenant-ID: t_demo" http://localhost:8080/api/tenant/echo`
