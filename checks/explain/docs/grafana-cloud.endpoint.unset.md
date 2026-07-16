---
id: grafana-cloud.endpoint.unset
title: 'OTEL_EXPORTER_OTLP_ENDPOINT is not set'
severity: error
---

## Why this matters

`OTEL_EXPORTER_OTLP_ENDPOINT` tells your application's OpenTelemetry SDK
where to ship its telemetry. When the variable is unset, most SDKs fall
back to `http://localhost:4318` (the default local Collector address),
so traces, metrics, and logs never leave the machine — they go to a
Collector that isn't running, and no data reaches Grafana Cloud.

The Grafana Cloud check treats the base variable as **required only for
signals that don't have a signal-specific override**. Each of the three
per-signal variables (`OTEL_EXPORTER_OTLP_TRACES_ENDPOINT`,
`OTEL_EXPORTER_OTLP_METRICS_ENDPOINT`,
`OTEL_EXPORTER_OTLP_LOGS_ENDPOINT`) fully substitutes for the base for
its signal. This error fires when the base is unset AND at least one
signal-specific variable is also unset — that signal has no endpoint to
ship to. The finding message names the missing signal(s).

## How to fix

Pick one of these two approaches:

**A. Set the base endpoint** — one variable covers every signal. Copy
the OTLP endpoint from your Grafana Cloud stack (found under
**Connections → OpenTelemetry (OTLP) → Send data**) and export it:

```bash
export OTEL_EXPORTER_OTLP_ENDPOINT="https://otlp-gateway-prod-<region>.grafana.net/otlp"
```

Substitute `<region>` with the region your stack lives in
(`us-east-0`, `eu-west-2`, etc.). Set it in every environment that runs
the service — Dockerfile `ENV`, Kubernetes pod spec, systemd unit,
CI environment.

**B. Set the missing signal-specific endpoint(s)** — appropriate if you
already use per-signal endpoints for other reasons. Note the signal
suffix is required here (see `grafana-cloud.signal-endpoint.invalid-format`):

```bash
export OTEL_EXPORTER_OTLP_TRACES_ENDPOINT="https://otlp-gateway-prod-<region>.grafana.net/otlp/v1/traces"
export OTEL_EXPORTER_OTLP_METRICS_ENDPOINT="https://otlp-gateway-prod-<region>.grafana.net/otlp/v1/metrics"
export OTEL_EXPORTER_OTLP_LOGS_ENDPOINT="https://otlp-gateway-prod-<region>.grafana.net/otlp/v1/logs"
```

## Example

Kubernetes fragment with the three variables the Grafana Cloud check
looks at:

```yaml
env:
  - name: OTEL_EXPORTER_OTLP_ENDPOINT
    value: "https://otlp-gateway-prod-us-east-0.grafana.net/otlp"
  - name: OTEL_EXPORTER_OTLP_PROTOCOL
    value: "http/protobuf"
  - name: OTEL_EXPORTER_OTLP_HEADERS
    valueFrom:
      secretKeyRef:
        name: grafana-cloud-otlp
        key: authorization
```

## Related

- [Grafana Cloud OTLP endpoint reference](https://grafana.com/docs/grafana-cloud/send-data/otlp/send-data-otlp/)
- `grafana-cloud.endpoint.invalid-format` — companion check when the
  base value is set but doesn't match the Grafana Cloud pattern.
- `grafana-cloud.signal-endpoint.invalid-format` — companion check for
  the per-signal variants.
- `grafana-cloud.endpoint.localhost` — companion check when the value
  still points at a local Collector.
