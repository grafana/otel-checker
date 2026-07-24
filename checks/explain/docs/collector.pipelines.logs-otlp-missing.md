---
id: collector.pipelines.logs-otlp-missing
title: 'Logs pipeline does not accept OTLP input'
severity: warning
---

## Why this matters

A Collector `service.pipelines.logs` block wires receivers → processors
→ exporters for log records. The receivers list is the definitive answer
to "where do my logs come from?" — if no OTLP receiver (`otlp` or
`otlp/<name>`) is in that list, this pipeline won't accept log records
emitted via OTLP by any OpenTelemetry SDK.

That's not always a bug — pipelines fed exclusively by non-OTLP
receivers are legitimate for services that log to stdout / files rather
than via the OTel SDK. The warning is a heads-up so you can confirm the
log source is what you intended.

## How to fix

If your services emit OTLP logs, add `otlp` (or a named `otlp/<name>`
id) to `service.pipelines.logs.receivers`:

```yaml
service:
  pipelines:
    logs:
      receivers: [otlp]
      processors: [batch]
      exporters: [otlphttp]
```

Multiple receivers are supported — you can fan logs in from OTLP and a
file source simultaneously:

```yaml
receivers: [otlp, filelog]
```

Then confirm the receiver itself is defined in the top-level `receivers`
block (see `collector.receivers.http-protocol-missing`).

## Related

- `collector.pipelines.traces-otlp-missing` — companion check for the
  traces pipeline.
- `collector.pipelines.metrics-otlp-missing` — companion check for the
  metrics pipeline.
- `collector.receivers.http-protocol-missing` — related check on the
  OTLP receiver's `http` protocol.
- [OpenTelemetry Collector pipelines](https://opentelemetry.io/docs/collector/configuration/#pipelines)
