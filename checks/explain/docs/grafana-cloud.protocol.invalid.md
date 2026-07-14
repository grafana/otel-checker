---
id: grafana-cloud.protocol.invalid
title: 'OTEL_EXPORTER_OTLP_PROTOCOL must be http/protobuf for Grafana Cloud'
severity: error
---

## Why this matters

The Grafana Cloud OTLP gateway
(`https://otlp-gateway-<region>.grafana.net/otlp`) only accepts OTLP over
HTTP with a Protobuf body — the `http/protobuf` protocol in OpenTelemetry
terms. If `OTEL_EXPORTER_OTLP_PROTOCOL` is set to `grpc` or `http/json`,
the SDK either can't connect at all or the gateway rejects the payload,
and no telemetry reaches Grafana Cloud.

## How to fix

Set the protocol explicitly to `http/protobuf` in the environment that
runs your service:

```bash
export OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf
```

In Docker:

```dockerfile
ENV OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf
```

In Kubernetes:

```yaml
env:
  - name: OTEL_EXPORTER_OTLP_PROTOCOL
    value: http/protobuf
```

Per-signal overrides (`OTEL_EXPORTER_OTLP_TRACES_PROTOCOL`,
`OTEL_EXPORTER_OTLP_METRICS_PROTOCOL`,
`OTEL_EXPORTER_OTLP_LOGS_PROTOCOL`) take precedence over the general
variable when set — check that none of them override the value back to
`grpc` or `http/json`.

## Example

Complete Grafana Cloud OTLP env-var set:

```bash
export OTEL_EXPORTER_OTLP_ENDPOINT="https://otlp-gateway-prod-us-east-0.grafana.net/otlp"
export OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf
export OTEL_EXPORTER_OTLP_HEADERS="Authorization=Basic <base64-token>"
```

## Related

- [OpenTelemetry OTLP protocol specification](https://opentelemetry.io/docs/specs/otlp/)
- [Grafana Cloud OTLP endpoint reference](https://grafana.com/docs/grafana-cloud/send-data/otlp/send-data-otlp/)
- `grafana-cloud.endpoint.unset` / `grafana-cloud.endpoint.invalid-format`
  — companion checks for the endpoint URL.
