# 灰度发布与回滚说明（Phase 7）

## 灰度思路（适配本地构建）
- 使用 `staging` 作为灰度环境，验证通过后再切到 `prod`。
- 流程：staging 验证 -> prod 发布 -> 监控 -> 必要时回滚。

## 演示流程
1. 灰度发布到 staging
   - `set APP_ENV=staging`（PowerShell: `$env:APP_ENV="staging"`）
   - `docker compose up -d --build`
2. 验证关键功能
   - 登录、项目/任务、文件上传
3. 正式发布到 prod
   - `set APP_ENV=prod`（PowerShell: `$env:APP_ENV="prod"`）
   - `docker compose up -d --build`

## 回滚策略
- 使用上一个稳定 tag 的代码重新构建部署。
- 回滚示例：
  - `git checkout v1.0.0`
  - `set APP_ENV=prod`
  - `docker compose up -d --build`

## 演示脚本
- 灰度发布：`scripts/rollout.ps1`
- 回滚：`scripts/rollback.ps1`