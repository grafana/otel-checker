---
id: collector.pipelines.metrics-otlphttp-missing
title: 'Metrics pipeline does not export to otlphttp'
severity: warning
---

## Why this matters

A Collector `service.pipelines.metrics` block wires receivers →
processors → exporters for metrics. The exporters list here is the
definitive answer to "where do my metrics go?" — if `otlphttp` isn't in
that list, the Collector never hands metrics to the exporter that ships
them to Grafana Cloud, and dashboards backed by those metrics stop
updating.

The exporter can be defined perfectly elsewhere in the file; if it isn't
*referenced* in the metrics pipeline, it's dead code as far as metrics
are concerned.

## How to fix

Add `otlphttp` (or whichever named-instance form you're using, e.g.
`otlphttp/grafana_cloud`) to `service.pipelines.metrics.exporters`:

```yaml
service:
  pipelines:
    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlphttp]
```

Then confirm the exporter itself is defined in the top-level `exporters`
block with a valid Grafana Cloud endpoint (see
`collector.endpoint.invalid-format`).

If you also scrape a `prometheus` receiver or `hostmetrics`, add them to
the same pipeline — the exporter list determines destination, the
receiver list determines source:

```yaml
receivers:
  otlp:
    protocols:
      http:
        endpoint: 0.0.0.0:4318
  prometheus:
    config:
      scrape_configs:
        - job_name: 'app'
          static_configs:
            - targets: ['localhost:8080']

service:
  pipelines:
    metrics:
      receivers: [otlp, prometheus]
      processors: [batch]
      exporters: [otlphttp]
```

## Example

Minimal metrics pipeline that ships to Grafana Cloud:

```yaml
exporters:
  otlphttp:
    endpoint: https://otlp-gateway-prod-us-east-0.grafana.net/otlp
    auth:
      authenticator: basicauth/grafana_cloud

service:
  pipelines:
    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlphttp]
```

## Related

- `collector.pipelines.traces-otlphttp-missing` — companion check for
  the traces pipeline.
- `collector.pipelines.logs-otlphttp-missing` — companion check for the
  logs pipeline.
- [OpenTelemetry Collector pipelines](https://opentelemetry.io/docs/collector/configuration/#pipelines)
