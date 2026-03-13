# 第四阶段数据库工程化与缓存一致性说明

本文用于记录第四阶段的数据库工程化与缓存一致性方案，可直接作为复现与验收材料。

## 1. 迁移脚本规范
- 采用 golang-migrate 版本化脚本：`00000X_xxx.up.sql` / `00000X_xxx.down.sql`。
- 所有建表/索引使用 `IF NOT EXISTS`，便于多环境幂等执行。
- 新增索引统一放在单独迁移：`000006_indexes.*.sql`，可回滚。

## 2. 索引设计说明
本阶段新增索引在 `backend/migrations/000006_indexes.up.sql`：
- `tasks(tenant_id, status)`：支持状态筛选与看板统计。
- `tasks(tenant_id, due_at)`：支持到期时间排序与范围过滤。
- `tasks(tenant_id, assignee_id)`：支持“我的任务”查询。
- `tasks(tenant_id, updated_at)`：支持最近更新排序。
- `task_comments(tenant_id, created_at)`：支持最近评论列表。
- `notifications(tenant_id, read_at)`：支持未读通知统计。

## 3. 事务与隔离级别选择说明（含示例）
- 业务默认使用 PostgreSQL 的 `READ COMMITTED`（避免脏读，适合大多数场景）。
- 对于“审批/状态流转”类操作，建议在应用层使用事务包裹多表写入，确保一致性。
- 如需更严格一致性（如强一致报表），可在特定场景提升至 `REPEATABLE READ`。

示例（伪代码）：
```
BEGIN;
UPDATE tasks SET status = 'review' WHERE id = $1 AND tenant_id = $2;
INSERT INTO task_approvals(...);
COMMIT;
```

## 4. 慢 SQL 记录
- 在 `backend/internal/db/db.go` 中启用 GORM Logger。
- 通过 `SLOW_SQL_THRESHOLD`（默认 200ms）记录慢 SQL。
- 通过 `DB_LOG_LEVEL` 控制日志级别（默认 warn）。

## 5. 缓存策略与一致性方案
本阶段输出方案与建议，后续阶段落地：
- 穿透：布隆过滤器/空值缓存。
- 击穿：热点 key 加互斥锁或逻辑过期。
- 雪崩：过期时间随机化 + 预热。
- 双写一致性：优先延迟双删策略；必要时用消息队列补偿。
- 热点降级与限流：热点接口加本地缓存 + 限流（令牌桶/漏桶）。

## 6. 复现与验证
- 迁移执行：`migrate up` / `migrate down`。
- 索引检查：`\d+ <table>` 或 `SELECT * FROM pg_indexes WHERE tablename='tasks';`
- 慢 SQL 观察：将 `SLOW_SQL_THRESHOLD` 调低到 `20ms` 进行压测。
