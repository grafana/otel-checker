---
id: grafana-cloud.endpoint.invalid-format
title: 'OTEL_EXPORTER_OTLP_ENDPOINT is not a Grafana Cloud OTLP gateway URL'
severity: error
---

## Why this matters

The Grafana Cloud check expects `OTEL_EXPORTER_OTLP_ENDPOINT` to match
`https://*.grafana.net/otlp` — the address of the Grafana Cloud OTLP
gateway. Anything else means the SDK is either pointed at a different
backend, a local collector that will later forward to Grafana Cloud, or
a plain-HTTP endpoint that the gateway won't accept.

The exact pattern the checker enforces:

- Scheme must be `https` (the gateway rejects plaintext).
- Host must end with `.grafana.net`.
- Path must be exactly `/otlp` (no trailing signal path like `/v1/traces`;
  the SDK appends that itself).

## How to fix

Look up your stack's OTLP endpoint from the Grafana Cloud UI
(**Connections → OpenTelemetry (OTLP) → Send data**), then export the
full URL:

```bash
export OTEL_EXPORTER_OTLP_ENDPOINT="https://otlp-gateway-prod-us-east-0.grafana.net/otlp"
```

Do not include a per-signal suffix (`/v1/traces`, `/v1/metrics`,
`/v1/logs`). The OpenTelemetry SDK appends those automatically for the
HTTP protocol.

If you're sending through a local Collector first, that's an intentional
choice — but the Grafana Cloud check is only valid against the direct
endpoint. Either point at the gateway directly for this check, or scope
the check to just the Collector (`otel-checker check collector`).

## Example

```bash
# Before — path is missing /otlp
export OTEL_EXPORTER_OTLP_ENDPOINT="https://otlp-gateway-prod-us-east-0.grafana.net"

# Before — path has the signal suffix
export OTEL_EXPORTER_OTLP_ENDPOINT="https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/traces"

# After
export OTEL_EXPORTER_OTLP_ENDPOINT="https://otlp-gateway-prod-us-east-0.grafana.net/otlp"
```

## Related

- [Grafana Cloud OTLP endpoint reference](https://grafana.com/docs/grafana-cloud/send-data/otlp/send-data-otlp/)
- `grafana-cloud.endpoint.unset` — companion check when the variable is
  missing entirely.
- `grafana-cloud.endpoint.localhost` — companion check when the value
  points at a local Collector rather than Grafana Cloud.
