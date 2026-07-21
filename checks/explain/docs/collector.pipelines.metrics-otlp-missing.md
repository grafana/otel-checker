---
id: collector.pipelines.metrics-otlp-missing
title: 'Metrics pipeline does not accept OTLP input'
severity: warning
---

## Why this matters

A Collector `service.pipelines.metrics` block wires receivers →
processors → exporters for metrics. The receivers list is the
definitive answer to "where do my metrics come from?" — if no OTLP
receiver (`otlp` or `otlp/<name>`) is in that list, this pipeline won't
accept metrics emitted via OTLP by any OpenTelemetry SDK.

That's not always a bug — pipelines fed exclusively by non-OTLP
receivers (`prometheus` scraping, `hostmetrics`) are legitimate
for infrastructure-metrics setups where nothing is pushing OTLP.
The warning is a heads-up so you can confirm the metric source is
what you intended.

## How to fix

If your services emit OTLP metrics, add `otlp` (or a named `otlp/<name>`
id) to `service.pipelines.metrics.receivers`:

```yaml
service:
  pipelines:
    metrics:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlphttp]
```

Multiple receivers are supported — you can fan metrics in from OTLP and
a Prometheus scrape simultaneously:

```yaml
receivers: [otlp, prometheus]
```

Then confirm the receiver itself is defined in the top-level `receivers`
block (see `collector.receivers.http-protocol-missing`).

## Related

- `collector.pipelines.traces-otlp-missing` — companion check for the
  traces pipeline.
- `collector.pipelines.logs-otlp-missing` — companion check for the
  logs pipeline.
- `collector.receivers.http-protocol-missing` — related check on the
  OTLP receiver's `http` protocol.
- [OpenTelemetry Collector pipelines](https://opentelemetry.io/docs/collector/configuration/#pipelines)
