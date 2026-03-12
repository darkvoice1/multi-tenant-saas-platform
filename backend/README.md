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
4. Run API (loads `.env`, dev mode auto-applies SQL migrations)
   - `go run ./cmd/api`
5. Initialize demo tenant + admin user (UI)
   - Open `http://localhost:5173/bootstrap` and click initialize

## Auth
- `POST /auth/bootstrap` (dev only, first-time init)
- `POST /auth/login` with `tenant_id`, `email`, `password`
- `POST /auth/refresh`
- `POST /auth/logout`
- `GET /auth/me` (requires `Authorization: Bearer <token>`)

## OAuth2/OIDC Mock
1. Get authorization code
   - `GET /auth/oidc/mock/authorize?tenant_id=<tenant_uuid>&email=<email>&state=demo`
2. Exchange code for tokens
   - `POST /auth/oidc/callback` with JSON `{ "code": "<code>", "state": "demo" }`

## Core APIs (Phase 3)
- Projects: `GET/POST /api/projects`, `GET/PUT/DELETE /api/projects/:id`
- Tasks: `GET/POST /api/projects/:id/tasks`, `GET/PUT/DELETE /api/tasks/:id`
- Task status: `POST /api/tasks/:id/status`
- Task approval: `POST /api/tasks/:id/approve` (admin/manager)
- Comments: `GET/POST /api/tasks/:id/comments`
- Attachments: `GET/POST /api/tasks/:id/attachments`
- Attachment preview: `GET /api/attachments/:id/preview`
- Attachment download: `GET /api/attachments/:id/download`
- Notifications: `GET /api/notifications`, `POST /api/notifications/:id/read`

## Tenant Context
- Tenant ID is required in `X-Tenant-ID` for `/api/*`
- Example:
  - `curl -H "X-Tenant-ID: t_demo" http://localhost:8080/api/tenant/echo`

## Storage (Local / S3)
- Default is local disk (`STORAGE_BACKEND=local`, `STORAGE_DIR=storage`)
- MinIO (S3 compatible) via Docker Compose:
  - Start: `docker compose up -d minio`
  - Console: `http://localhost:9001` (user/pass: `minioadmin`)
  - Use:
    - `STORAGE_BACKEND=s3`
    - `S3_ENDPOINT=http://localhost:9000`
    - `S3_BUCKET=saas-platform`
    - `S3_ACCESS_KEY=minioadmin`
    - `S3_SECRET_KEY=minioadmin`
