---
id: collector.pipelines.logs-otlphttp-missing
title: 'Logs pipeline does not export to otlphttp'
severity: warning
---

## Why this matters

A Collector `service.pipelines.logs` block wires receivers → processors
→ exporters for log records. The exporters list here is the definitive
answer to "where do my logs go?" — if `otlphttp` isn't in that list, the
Collector never hands log records to the exporter that ships them to
Grafana Cloud, and logs stop arriving in your stack.

The exporter can be defined perfectly elsewhere in the file; if it isn't
*referenced* in the logs pipeline, it's dead code as far as logs are
concerned.

## How to fix

Add `otlphttp` (or whichever named-instance form you're using, e.g.
`otlphttp/grafana_cloud`) to `service.pipelines.logs.exporters`:

```yaml
service:
  pipelines:
    logs:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlphttp]
```

Then confirm the exporter itself is defined in the top-level `exporters`
block with a valid Grafana Cloud endpoint (see
`collector.endpoint.invalid-format`).

If your service also uses a non-OTLP log source (filelog, syslog, etc.),
add those receivers to the same pipeline — the exporter list is what
determines where the output goes, receivers determine where it comes
from:

```yaml
receivers:
  otlp:
    protocols:
      http:
        endpoint: 0.0.0.0:4318
  filelog:
    include: [/var/log/app/*.log]

service:
  pipelines:
    logs:
      receivers: [otlp, filelog]
      processors: [batch]
      exporters: [otlphttp]
```

## Example

Minimal logs pipeline that ships to Grafana Cloud:

```yaml
exporters:
  otlphttp:
    endpoint: https://otlp-gateway-prod-us-east-0.grafana.net/otlp
    auth:
      authenticator: basicauth/grafana_cloud

service:
  pipelines:
    logs:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlphttp]
```

## Related

- `collector.pipelines.traces-otlphttp-missing` — companion check for
  the traces pipeline.
- `collector.pipelines.metrics-otlphttp-missing` — companion check for
  the metrics pipeline.
- [OpenTelemetry Collector pipelines](https://opentelemetry.io/docs/collector/configuration/#pipelines)
