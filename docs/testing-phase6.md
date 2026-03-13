# 测试说明（Phase 6）

## 单元测试
```
go test ./... -cover
```

## 夹具脚本
```
$env:FIXTURE_TENANT_SLUG="demo"
$env:FIXTURE_PROJECTS=2
$env:FIXTURE_TASKS=5
go run ./cmd/fixtures
```

