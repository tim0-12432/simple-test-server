# Observability Specification

## Overview
OpenTelemetry integration for distributed tracing and metrics in Simple Test Server.

## Features
- **Configuration-driven**: Enable/disable via `OTEL_ENABLED` environment variable
- **Multiple exporters**: OTLP (gRPC/HTTP) and Prometheus
- **Request tracing**: Automatic request ID generation and propagation
- **Metrics endpoint**: Prometheus-compatible scraping endpoint
- **Graceful shutdown**: Clean telemetry provider shutdown during app termination

## Configuration

| Environment Variable      | Type    | Default          | Description                          |
|---------------------------|---------|------------------|--------------------------------------|
| `OTEL_ENABLED`            | bool    | `false`          | Enable observability features        |
| `OTEL_SERVICE_NAME`       | string  | `simple-test-server` | Service name for traces/metrics |
| `OTEL_ENV`                | string  | `production`     | Environment (dev, staging, prod)     |
| `OTEL_EXPORTER`           | string  | `prometheus`     | Exporter type (otlp, prometheus)     |
| `OTEL_OTLP_ENDPOINT`      | string  | `localhost:4317` | OTLP gRPC endpoint                   |
| `OTEL_PROMETHEUS_PORT`    | int     | `9090`           | Prometheus metrics server port       |
| `OTEL_TRACES_ENABLED`     | bool    | `true`           | Enable distributed tracing           |
| `OTEL_METRICS_ENABLED`    | bool    | `true`           | Enable metrics collection            |

## APIs

### Initialization
```go
func Init(ctx context.Context, cfg config.ObservabilityConfig) (shutdown func(context.Context) error, err error)
```
Sets up OpenTelemetry providers. Returns shutdown function for cleanup.

### Metrics Handler
```go
func MetricsHandler() http.Handler
```
Returns HTTP handler for Prometheus `/metrics` endpoint.

### Request ID Middleware
```go
func RequestIDMiddleware() gin.HandlerFunc
```
Gin middleware that reads or generates `X-Request-ID` header and injects into context.

## Integration Points
1. `main.go`: Initialize observability after config load; defer shutdown
2. `controllers/routes.go`: Install RequestIDMiddleware early in chain
3. Future: Instrument Docker operations, database calls, protocol handlers
