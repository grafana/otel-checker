---
id: collector.endpoint.invalid-format
title: 'Collector otlphttp exporter endpoint is not a Grafana Cloud OTLP gateway URL'
severity: error
---

## Why this matters

The Collector's `otlphttp` exporter is configured with an endpoint that
doesn't match the Grafana Cloud OTLP gateway pattern
(`https://*.grafana.net/otlp`). That usually means one of three things:

- A placeholder like `https://your-instance.example.com/otlp` was left
  in the config and never replaced.
- The URL points at a *different* backend entirely (e.g. an internal
  system). Legitimate, but the Grafana Cloud check will never pass.
- A typo, wrong region, or extra path segment (`/v1/traces`) got baked
  into the endpoint.

Either way, the exporter as configured will not deliver telemetry to
Grafana Cloud.

## How to fix

Look up the correct OTLP endpoint from the Grafana Cloud UI (**Connections
→ OpenTelemetry (OTLP) → Send data**) and paste it in verbatim:

```yaml
exporters:
  otlphttp:
    endpoint: https://otlp-gateway-prod-us-east-0.grafana.net/otlp
    auth:
      authenticator: basicauth/grafana_cloud
```

Endpoint requirements the checker enforces:

- Scheme must be `https`.
- Host must end with `.grafana.net`.
- Path must be exactly `/otlp` (no `/v1/traces`, `/v1/metrics`, or
  `/v1/logs` suffix — the exporter appends those itself).

If the current endpoint intentionally points at a different backend, the
Grafana Cloud check isn't the right check to run against this Collector —
scope your validation accordingly.

## Example

```yaml
# Before — wrong path
exporters:
  otlphttp:
    endpoint: https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/traces

# After
exporters:
  otlphttp:
    endpoint: https://otlp-gateway-prod-us-east-0.grafana.net/otlp
```

## Related

- [Grafana Cloud OTLP endpoint reference](https://grafana.com/docs/grafana-cloud/send-data/otlp/send-data-otlp/)
- `collector.endpoint.localhost` — companion check when the endpoint
  still points at a local Collector.
