# 第四阶段说明
## 目标
- 数据库工程化（迁移规范、索引设计、慢 SQL 记录）
- 缓存一致性方案输出与复现说明

## 本阶段新增
- 迁移与索引
  - `backend/migrations/000006_indexes.up.sql`
  - `backend/migrations/000006_indexes.down.sql`
- 文档
  - `docs/db-cache-phase4.md`：数据库工程化与缓存一致性说明
  - `docs/slow-sql-report-phase4.md`：慢 SQL 优化记录模板
- 配置
  - `DB_LOG_LEVEL`、`SLOW_SQL_THRESHOLD`
  - `docker-compose.yml` 新增 Redis

## 说明
- 慢 SQL 记录通过 GORM Logger 输出，阈值默认 `200ms`。
- 缓存一致性方案以文档形式输出，后续阶段落地实现。
