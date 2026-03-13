# 审计日志说明（Phase 5）

## 表结构
`audit_logs` 记录关键操作：
- `tenant_id`、`user_id`：谁在操作
- `action`、`resource`、`resource_id`：做了什么
- `method`、`path`、`status_code`：请求与结果
- `ip`、`user_agent`：访问来源
- `created_at`：发生时间

## 记录范围
- `/api` 路径下所有已认证请求。
- 失败请求同样记录（便于追踪越权和误操作）。

## 查询接口
`GET /api/admin/audit-logs`
- `limit`：最多 200
- `before`：RFC3339 时间，用于翻页

