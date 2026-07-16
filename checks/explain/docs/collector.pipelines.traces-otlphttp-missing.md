---
id: collector.pipelines.traces-otlphttp-missing
title: 'Traces pipeline does not export to otlphttp'
severity: warning
---

## Why this matters

A Collector `service.pipelines.traces` block wires receivers → processors
→ exporters for spans. The exporters list here is the definitive answer
to "where do my traces go?" — if `otlphttp` isn't in that list, the
Collector never hands spans to the exporter that ships them to Grafana
Cloud, and traces stop arriving in your stack.

The exporter can be defined perfectly elsewhere in the file; if it isn't
*referenced* in the traces pipeline, it's dead code as far as spans are
concerned.

## How to fix

Add `otlphttp` (or whichever named-instance form you're using, e.g.
`otlphttp/grafana_cloud`) to `service.pipelines.traces.exporters`:

```yaml
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlphttp]
```

Then confirm the exporter itself is defined in the top-level `exporters`
block with a valid Grafana Cloud endpoint (see
`collector.endpoint.invalid-format`).

Multiple exporters are supported — you can fan traces out to both
Grafana Cloud and, say, `debug`:

```yaml
exporters: [otlphttp, debug]
```

## Example

Minimal traces pipeline that ships to Grafana Cloud:

```yaml
receivers:
  otlp:
    protocols:
      http:
        endpoint: 0.0.0.0:4318
      grpc:
        endpoint: 0.0.0.0:4317

exporters:
  otlphttp:
    endpoint: https://otlp-gateway-prod-us-east-0.grafana.net/otlp
    auth:
      authenticator: basicauth/grafana_cloud

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlphttp]
```

## Related

- `collector.pipelines.logs-otlphttp-missing` — companion check for the
  logs pipeline.
- `collector.pipelines.metrics-otlphttp-missing` — companion check for
  the metrics pipeline.
- [OpenTelemetry Collector pipelines](https://opentelemetry.io/docs/collector/configuration/#pipelines)
