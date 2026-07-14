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

For the Grafana Cloud `check` to pass this variable must be set to the
OTLP gateway URL of your Grafana Cloud stack.

## How to fix

Copy the OTLP endpoint from your Grafana Cloud stack (found under
**Connections → OpenTelemetry (OTLP) → Send data**), and export it in the
environment that runs your service:

```bash
export OTEL_EXPORTER_OTLP_ENDPOINT="https://otlp-gateway-prod-<region>.grafana.net/otlp"
```

Substitute `<region>` with the region your stack lives in
(`us-east-0`, `eu-west-2`, etc.). Set it in every environment that runs
the service — Dockerfile `ENV`, Kubernetes pod spec, systemd unit,
CI environment.

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
  value is set but doesn't match the Grafana Cloud pattern.
- `grafana-cloud.endpoint.localhost` — companion check when the value
  still points at a local Collector.
