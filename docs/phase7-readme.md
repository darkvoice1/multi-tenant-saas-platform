# 第七阶段：CI/CD与部署交付

## 目标
- 形成可持续交付的最小闭环：CI 通过、CD 本地一键部署、配置分层、反向代理。

## 已完成事项
- GitHub Actions CI（lint/test/build）
- Docker Compose 本地一键构建与启动
- 环境变量分层（dev/staging/prod）
- Nginx 反向代理与静态资源托管

## 本地部署（推荐）
1. 启动
   - `docker compose up -d --build`
2. 访问
   - 前端：http://localhost:5173
   - 后端：http://localhost:8080

## 环境切换
- 默认：dev（如果未设置 APP_ENV，会自动使用 dev）
- 切换示例：
  - `set APP_ENV=staging`（PowerShell: `$env:APP_ENV="staging"`）
  - `docker compose up -d --build`

## 反向代理说明
- 前端 Nginx 负责托管静态资源，并将 `/api` 代理到后端服务。
- 生产/预发环境前端默认使用 `/api` 作为 API Base。

## 灰度/回滚
- 参考文档：`docs/canary-rollback-phase7.md`
- 演示脚本：`scripts/rollout.ps1`、`scripts/rollback.ps1`

## 发布说明与版本策略
- 参考文档：`docs/release-strategy-phase7.md`