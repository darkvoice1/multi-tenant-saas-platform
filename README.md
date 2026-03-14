# 多租户协同 SaaS 平台

一个面向团队协作的多租户 SaaS 示例项目，覆盖租户隔离、RBAC、审计、可观测性、测试体系与 CI/CD 基线。

## 主要功能
- 多租户隔离与租户上下文注入
- 认证与权限（JWT + RBAC）
- 项目/任务/评论等协作流程
- 附件上传与对象存储（MinIO）
- 审计日志与安全基线
- 可观测性（Prometheus/Grafana/Jaeger）
- 测试体系（单测/集成/E2E）与 CI

## 技术栈
- 后端：Go + Gin + GORM + PostgreSQL
- 前端：Vue 3 + Vite + TypeScript + Pinia
- 基础设施：Docker Compose + Nginx + MinIO + Redis
- 可观测性：Prometheus + Grafana + Jaeger
- CI：GitHub Actions

## 目录结构
- `backend/` 后端服务
- `frontend/` 前端应用
- `e2e/` 端到端测试
- `docs/` 设计与阶段文档
- `scripts/` 演示脚本

## 一键本地部署（推荐）
前置：已安装 Docker Desktop 与 Docker Compose。

在项目根目录执行：

```bash
docker compose up -d --build
```

默认访问：
- 前端：`http://localhost:5173`
- 后端：`http://localhost:8080`
- MinIO：`http://localhost:9001`
- Grafana：`http://localhost:3000`
- Jaeger：`http://localhost:16686`

停止：
```bash
docker compose down
```

## 环境分层（dev/staging/prod）
- 后端：`backend/.env.dev`、`backend/.env.staging`、`backend/.env.prod`
- 前端：`frontend/.env.development`、`frontend/.env.staging`、`frontend/.env.production`
- Compose 会读取 `APP_ENV` 对应的后端环境文件（默认 `dev`）。

Windows PowerShell 示例：
```powershell
$env:APP_ENV="staging"
docker compose up -d --build
```

## 本地开发
### 后端
```bash
cd backend
cp .env.example .env
# 启动服务
# go run ./cmd/api
```

### 前端
```bash
cd frontend
npm install
npm run dev
```

## CI 与测试
- CI 工作流：`.github/workflows/ci.yml`
- 覆盖率门禁：`backend/scripts/coverage.ps1`
- 集成测试：`go test ./internal/db -tags=integration -count=1`
- E2E 测试：`e2e/` 目录

## 反向代理说明
前端容器使用 Nginx 托管静态文件，并将 `/api` 代理到后端。
配置文件：`frontend/nginx.conf`

## 灰度/回滚与发布策略
- 灰度/回滚说明：`docs/canary-rollback-phase7.md`
- 版本策略：`docs/release-strategy-phase7.md`
- 演示脚本：`scripts/rollout.ps1`、`scripts/rollback.ps1`

## 常见问题
- 镜像拉取慢：请配置 Docker Desktop 代理或使用镜像加速。
- 前端接口失败：确认后端是否已启动，或 Nginx 是否正确代理 `/api`。

## 许可证
MIT