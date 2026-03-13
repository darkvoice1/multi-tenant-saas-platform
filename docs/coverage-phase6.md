# 覆盖率门禁（Phase 6）

## 目的
- 输出可展示的覆盖率指标。
- 在本地进行最小质量门禁（可配置阈值）。

## 使用方式
在 `backend` 目录执行：
```
# 默认阈值 20%
./scripts/coverage.ps1

# 自定义阈值
./scripts/coverage.ps1 -MinCoverage 30
```

## 说明
- 生成 `coverage.out` 文件。
- 使用 `go tool cover -func` 计算总覆盖率。

