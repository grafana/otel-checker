---
id: config.endpoint.missing
title: 'No OTLP HTTP endpoint declared for a signal in the config file'
severity: warning
---

## Why this matters

The declarative config file loaded fine, but one of the three signals
(`traces`, `metrics`, `logs`) does not declare an `otlp_http` exporter
endpoint anywhere under its provider. Without an endpoint the SDK
either drops that signal silently or falls back to a hard-coded
default that isn't Grafana Cloud, so telemetry for that signal will
never reach the customer's stack.

The checker looks for endpoints at these paths:

- Traces: `tracer_provider.processors[].batch.exporter.otlp_http.endpoint`
  and `.simple.exporter.otlp_http.endpoint`.
- Metrics: `meter_provider.readers[].periodic.exporter.otlp_http.endpoint`
  and `.pull.exporter.otlp_http.endpoint`.
- Logs: `logger_provider.processors[].batch.exporter.otlp_http.endpoint`
  and `.simple.exporter.otlp_http.endpoint`.

Only `otlp_http` is checked in this iteration.

## How to fix

Add an `otlp_http` exporter to the missing signal's provider:

```yaml
tracer_provider:
  processors:
    - batch:
        exporter:
          otlp_http:
            endpoint: https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/traces
            headers:
              - name: Authorization
                value: Basic <base64-encoded-credentials>
```

Repeat for `meter_provider` (using `readers` and `periodic`) and
`logger_provider` (using `processors` and `batch`).

## Related

- [OpenTelemetry Configuration spec — providers](https://opentelemetry.io/docs/specs/otel/configuration/)
- `grafana-cloud.signal-endpoint.invalid-format` — related failure
  when an endpoint IS declared but doesn't match the Grafana Cloud URL.
