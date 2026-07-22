---
id: collector.pipelines.traces-otlp-missing
title: 'Traces pipeline does not accept OTLP input'
severity: warning
---

## Why this matters

A Collector `service.pipelines.traces` block wires receivers → processors
→ exporters for traces. The receivers list is the definitive answer to
"where do my traces come from?" — if no OTLP receiver (`otlp` or
`otlp/<name>`) is in that list, this pipeline won't accept traces emitted
via OTLP by any OpenTelemetry SDK.

That's not always a bug — pipelines fed exclusively by non-OTLP
receivers are legitimate. The warning is a heads-up so you can confirm
the trace source is what you intended.

## How to fix

If your services emit OTLP, add `otlp` (or a named `otlp/<name>` id) to
`service.pipelines.traces.receivers`:

```yaml
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlphttp]
```

Multiple receivers are supported — you can fan spans in from OTLP and
another source simultaneously:

```yaml
receivers: [otlp, zipkin]
```

Then confirm the receiver itself is defined in the top-level `receivers`
block (see `collector.receivers.http-protocol-missing`).

## Related

- `collector.pipelines.logs-otlp-missing` — companion check for the
  logs pipeline.
- `collector.pipelines.metrics-otlp-missing` — companion check for the
  metrics pipeline.
- `collector.receivers.http-protocol-missing` — related check on the
  OTLP receiver's `http` protocol.
- [OpenTelemetry Collector pipelines](https://opentelemetry.io/docs/collector/configuration/#pipelines)
