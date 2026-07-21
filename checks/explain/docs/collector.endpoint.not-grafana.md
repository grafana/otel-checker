---
id: collector.endpoint.not-grafana
title: 'Collector exporter endpoint is a valid URL but not a Grafana Cloud endpoint'
severity: warning
---

## Why this matters

An `otlphttp` / `otlp_http` exporter in your Collector config is
configured with a well-formed URL that doesn't match the Grafana Cloud
OTLP gateway pattern (`https://*.grafana.net/otlp`). The Collector will
happily start with this configuration — the exporter is valid — but any
telemetry routed through it goes to a different backend, not to
Grafana Cloud.

This is a common intentional setup: the Collector fans out to multiple
backends, and only *some* exporters point at Grafana Cloud. The warning
is a heads-up, not a bug — it lists which exporters are shipping data
elsewhere so you can confirm each is on purpose.

## How to fix

If the endpoint is intentional (e.g. `otlp_http/prometheus` routing to a
Prometheus server), no action is required — the warning is informational
and safe to ignore. Verify by checking which exporter each pipeline uses
in `service.pipelines.*.exporters`.

If the endpoint was meant to be Grafana Cloud, replace it with your
stack's OTLP gateway URL (found under **Connections → OpenTelemetry
(OTLP) → Send data** in the Grafana Cloud UI):

```yaml
exporters:
  otlphttp/grafana_cloud:
    endpoint: https://otlp-gateway-<zone>.grafana.net/otlp
    auth:
      authenticator: basicauth/grafana_cloud
```

Do not include a per-signal suffix (`/v1/traces`, `/v1/metrics`,
`/v1/logs`). The `otlphttp` exporter appends those itself.

## Example

Common multi-backend fan-out where the warning is expected on the
non-Grafana exporters:

```yaml
exporters:
  otlp_http/grafana_cloud:
    endpoint: https://otlp-gateway-prod-us-east-0.grafana.net/otlp   # → success
  otlp_http/prometheus:
    endpoint: http://prometheus:9090/api/v1/otlp                    # → this warning
  otlp_http/loki:
    endpoint: http://loki:3100/otlp                                 # → this warning
```

## Related

- [Grafana Cloud OTLP endpoint reference](https://grafana.com/docs/grafana-cloud/send-data/otlp/send-data-otlp/)
- `collector.endpoint.localhost` — companion warning when a Collector
  exporter still points at `localhost`.
- `collector.endpoint.invalid-format` — companion error when the
  endpoint isn't a well-formed URL at all.
