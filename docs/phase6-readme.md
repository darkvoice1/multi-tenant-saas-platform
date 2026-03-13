# 第六阶段：可观测性与测试体系

## 目标
- 结构化日志、指标与链路追踪可演示。
- 基础测试与夹具脚本可复现。
- 覆盖率指标可展示并具备最小门禁。

## 功能概览
- 结构化日志：每个请求输出 `trace_id/tenant_id/user_id` 等字段。
- 指标：Prometheus `/metrics` 暴露请求与业务指标。
- 链路追踪：OpenTelemetry（stdout 或 OTLP/Jaeger）。
- 测试：新增单测样例，提供夹具脚本。
- 覆盖率门禁：`coverage.ps1` 支持最小覆盖率阈值校验。

## 快速开始
### 1) 启动后端
```
$env:JWT_SECRET="dev_change_me"
$env:DATABASE_URL="postgres://postgres:postgres@localhost:5432/saas_platform?sslmode=disable"
$env:OTEL_EXPORTER="stdout"
$env:OTEL_SERVICE_NAME="saas-platform-api"
cd backend
go run ./cmd/api
```

### 2) 访问指标
- `http://localhost:8080/metrics`

### 3) Docker 观测组件
```
docker compose up -d prometheus grafana jaeger
```
- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3000` (admin/admin)
- Jaeger: `http://localhost:16686`

### 4) OTLP 追踪
```
$env:OTEL_EXPORTER="otlp"
$env:OTEL_EXPORTER_OTLP_ENDPOINT="localhost:4317"
$env:OTEL_SERVICE_NAME="saas-platform-api"
```

## 业务指标
- 登录成功率：`saas_login_total{result="success|failure"}`
- 审批耗时：`saas_task_approval_duration_seconds`

## 测试与夹具
```
go test ./... -cover
```
```
# 生成演示数据
$env:FIXTURE_TENANT_SLUG="demo"
$env:FIXTURE_PROJECTS=2
$env:FIXTURE_TASKS=5
go run ./cmd/fixtures
```

## 覆盖率门禁
```
# 默认阈值 20%
./scripts/coverage.ps1

# 自定义阈值
./scripts/coverage.ps1 -MinCoverage 30
```

