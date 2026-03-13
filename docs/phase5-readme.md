# 第五阶段：安全、审计与隔离增强

## 变更摘要
- 审计日志：新增 `audit_logs` 表与审计中间件，记录关键操作（人/时间/资源/状态）。
- 敏感操作二次确认：删除项目/任务需要 `?confirm=true`。
- 越权扫描脚本：`backend/cmd/scan-permissions` 输出跨租户关联异常。
- 租户限额：项目数、成员数、存储、请求速率限制（按租户配置）。
- 参数校验与错误码：关键输入长度校验；统一错误码中间件补齐 `code` 字段。
- 安全基线：弱口令校验、JWT 强度校验（非 dev 环境），CORS 白名单。

## 审计日志
- 接口：`GET /api/admin/audit-logs?limit=100&before=2026-03-13T00:00:00Z`
- 字段：`action/resource/resource_id/status_code/ip/user_agent/created_at` 等。

## 敏感操作确认
- 删除项目：`DELETE /api/projects/:id?confirm=true`
- 删除任务：`DELETE /api/tasks/:id?confirm=true`

## 越权扫描
在 `backend` 目录执行：
```
go run ./cmd/scan-permissions
```
输出为各类跨租户异常统计与样例 ID。

## 租户限额配置
`tenants` 表新增字段：
- `max_projects`
- `max_members`
- `max_storage_bytes`
- `max_requests_per_minute`

## 错误码规则
当响应状态码 >= 400 且为 JSON 时，自动补齐 `code` 字段，便于前端统一处理。

