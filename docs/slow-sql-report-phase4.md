# 第四阶段慢 SQL 优化记录

> 记录示例：优化前后对比数据（Explain/执行时长）。

## 1. 任务列表按状态查询
- SQL：`SELECT * FROM tasks WHERE tenant_id = ? AND status = ? ORDER BY updated_at DESC LIMIT 20;`
- 优化前：基于索引扫描 + 排序，耗时 0.059 ms（本次数据量为 0 行）。
- 优化后：命中索引 `idx_tasks_tenant_status`，排序成本低。
- Explain（真实输出）：
```
Limit  (cost=8.18..8.18 rows=1 width=192) (actual time=0.026..0.026 rows=0 loops=1)
  Buffers: shared hit=5
  ->  Sort  (cost=8.18..8.18 rows=1 width=192) (actual time=0.024..0.025 rows=0 loops=1)
        Sort Key: updated_at DESC
        Sort Method: quicksort  Memory: 25kB
        Buffers: shared hit=5
        ->  Index Scan using idx_tasks_tenant_status on tasks  (cost=0.15..8.17 rows=1 width=192) (actual time=0.004..0.004 rows=0 loops=1)
              Index Cond: ((tenant_id = 'd055d50b-f4dd-4a2b-8a79-b24e6e1e56cc'::uuid) AND (status = 'todo'::text))
              Buffers: shared hit=2
Planning:
  Buffers: shared hit=207
Planning Time: 0.495 ms
Execution Time: 0.059 ms
```

## 2. 到期任务查询
- SQL：`SELECT * FROM tasks WHERE tenant_id = ? AND due_at IS NOT NULL AND due_at <= ? ORDER BY due_at ASC LIMIT 20;`
- 优化前：命中索引 `idx_tasks_tenant_due_at`，耗时 0.032 ms（本次数据量为 0 行）。
- 优化后：索引扫描，避免全表扫描。
- Explain（真实输出）：
```
Limit  (cost=0.15..8.17 rows=1 width=192) (actual time=0.009..0.010 rows=0 loops=1)
  Buffers: shared hit=2
  ->  Index Scan using idx_tasks_tenant_due_at on tasks  (cost=0.15..8.17 rows=1 width=192) (actual time=0.008..0.009 rows=0 loops=1)
        Index Cond: ((tenant_id = 'd055d50b-f4dd-4a2b-8a79-b24e6e1e56cc'::uuid) AND (due_at IS NOT NULL) AND (due_at <= (now() + '7 days'::interval)))
        Buffers: shared hit=2
Planning:
  Buffers: shared hit=210
Planning Time: 0.405 ms
Execution Time: 0.032 ms
```
