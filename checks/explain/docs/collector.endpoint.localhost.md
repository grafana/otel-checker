---
id: collector.endpoint.localhost
title: 'Collector otlphttp exporter endpoint points at localhost'
severity: warning
---

## Why this matters

The Collector's `otlphttp` exporter is configured with a `localhost`
endpoint — that means the Collector will try to forward telemetry to
another OTLP receiver on the same machine, not to Grafana Cloud. That's
supported in multi-hop topologies (e.g. a load-balancing Collector in
front of another Collector), but in a simple single-Collector setup it
usually means someone left a development address in the config and no
data is reaching Grafana Cloud.

## How to fix

Replace the endpoint with the Grafana Cloud OTLP gateway URL for your
stack (found under **Connections → OpenTelemetry (OTLP) → Send data** in
the Grafana Cloud UI):

```yaml
exporters:
  otlphttp:
    endpoint: https://otlp-gateway-prod-us-east-0.grafana.net/otlp
    auth:
      authenticator: basicauth/grafana_cloud
```

Do not include a per-signal suffix (`/v1/traces`, `/v1/metrics`,
`/v1/logs`). The `otlphttp` exporter appends those automatically.

If the localhost endpoint is intentional (multi-hop), scope your
verification to the *upstream* Collector that actually egresses to
Grafana Cloud instead — this warning is a heads-up, not a bug, for that
topology.

## Example

Complete Grafana Cloud otlphttp exporter block, with basicauth:

```yaml
extensions:
  basicauth/grafana_cloud:
    client_auth:
      username: "<instance_id>"
      password: "<api_key>"

exporters:
  otlphttp:
    endpoint: https://otlp-gateway-prod-us-east-0.grafana.net/otlp
    auth:
      authenticator: basicauth/grafana_cloud

service:
  extensions: [basicauth/grafana_cloud]
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [otlphttp]
```

## Related

- [Grafana Cloud OTLP endpoint reference](https://grafana.com/docs/grafana-cloud/send-data/otlp/send-data-otlp/)
- `collector.endpoint.invalid-format` — companion check when the value
  is neither localhost nor a Grafana Cloud URL.
