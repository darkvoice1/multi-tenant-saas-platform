# Backend (Gin + PostgreSQL)

## Requirements
- Go 1.22+
- PostgreSQL 14+
- golang-migrate (CLI)

## Quick Start
1. Copy env
   - `Copy-Item .env.example .env`
2. Set `DATABASE_URL` in `.env`
3. Run migrations
   - `migrate -path migrations -database "$env:DATABASE_URL" up`
4. Run API
   - `go run ./cmd/api`

## Tenant Context
- Tenant ID is required in `X-Tenant-ID` for `/api/*`
- Example:
  - `curl -H "X-Tenant-ID: t_demo" http://localhost:8080/api/tenant/echo`
