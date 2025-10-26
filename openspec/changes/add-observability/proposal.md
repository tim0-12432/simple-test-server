# Proposal: Add OpenTelemetry Observability

## Why
Enable production-grade monitoring and debugging by integrating OpenTelemetry for distributed tracing and metrics. This will help developers and operators understand system behavior, diagnose issues, and track performance across protocol handlers and Docker operations.

## What
Add OpenTelemetry instrumentation with:
- Configuration model for traces, metrics, and exporters (OTLP, Prometheus)
- Observability initialization with graceful shutdown
- HTTP middleware for request tracing with request IDs
- Metrics endpoint for Prometheus scraping

## Impact
- **Positive**: Better visibility into application behavior, easier debugging, production-ready observability
- **Neutral**: Minimal performance overhead when enabled; zero overhead when disabled
- **Risk**: Additional dependencies (OpenTelemetry SDK); mitigated by feature flag

## Acceptance Criteria
1. Configuration via environment variables (OTEL_*) with sensible defaults
2. Observability can be enabled/disabled via `OTEL_ENABLED`
3. Request ID middleware captures or generates X-Request-ID header
4. Prometheus metrics endpoint exposed when configured
5. Graceful shutdown of telemetry providers
6. All tests pass; no behavioral changes when observability disabled
