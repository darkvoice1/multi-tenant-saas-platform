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
3. Set `DATABASE_URL` and `JWT_SECRET` in `.env`
4. Run migrations
   - `migrate -path migrations -database "$env:DATABASE_URL" up`
5. (Optional) Seed demo tenant + admin user
   - `go run ./cmd/seed`
6. Run API
   - `go run ./cmd/api`

## Auth
- `POST /auth/login` with `tenant_id`, `email`, `password`
- `POST /auth/refresh`
- `POST /auth/logout`
- `GET /auth/me` (requires `Authorization: Bearer <token>`)

## OAuth2/OIDC Mock
1. Get authorization code
   - `GET /auth/oidc/mock/authorize?tenant_id=<tenant_uuid>&email=<email>&state=demo`
2. Exchange code for tokens
   - `POST /auth/oidc/callback` with JSON `{ "code": "<code>", "state": "demo" }`

## Tenant Context
- Tenant ID is required in `X-Tenant-ID` for `/api/*`
- Example:
  - `curl -H "X-Tenant-ID: t_demo" http://localhost:8080/api/tenant/echo`
