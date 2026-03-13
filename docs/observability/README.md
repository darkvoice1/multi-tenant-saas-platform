# 可观测性说明（Phase 6）

## 结构化日志
- 日志以 JSON 形式输出，包含：`trace_id/tenant_id/user_id/method/path/status/latency_ms`。
- 每个请求响应头会返回 `X-Trace-ID`。

## 指标
- `/metrics` 暴露 Prometheus 指标。
- 主要指标：
  - `saas_api_requests_total` / `saas_api_request_duration_seconds`
  - `saas_login_total`
  - `saas_task_approval_duration_seconds`

## 链路追踪
- 基于 OpenTelemetry 的 Gin 中间件。
- 支持 `stdout` 或 `otlp` 导出：
  - `OTEL_EXPORTER=stdout`
  - `OTEL_EXPORTER=otlp` + `OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317`

