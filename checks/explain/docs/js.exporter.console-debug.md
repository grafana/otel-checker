---
id: js.exporter.console-debug
title: 'Instrumentation file uses a Console exporter'
severity: warning
---

## Why this matters

`ConsoleSpanExporter` and `ConsoleMetricExporter` write telemetry to stdout
as human-readable text. They are useful during local development to verify
that spans and metrics are being produced, but they do not send anything
over the network — the data never reaches Grafana Cloud, and nothing shows
up in Explore, dashboards, or alerts. Leaving a Console exporter wired into
a production instrumentation file usually means the OTLP exporter that was
supposed to replace it was never plugged in.

## How to fix

Swap the Console exporter for the OTLP HTTP exporter that ships telemetry
to Grafana Cloud:

- Traces: replace `ConsoleSpanExporter` with `OTLPTraceExporter` from
  `@opentelemetry/exporter-trace-otlp-http`.
- Metrics: replace `ConsoleMetricExporter` with `OTLPMetricExporter` from
  `@opentelemetry/exporter-metrics-otlp-http`.

If you want to keep console output for debugging, add the Console exporter
as an *additional* processor guarded by a flag (e.g. `if
(process.env.OTEL_DEBUG)`) rather than as the primary exporter.

## Example

```js
// Before — telemetry never leaves the process.
import { ConsoleSpanExporter } from '@opentelemetry/sdk-trace-node';
const exporter = new ConsoleSpanExporter();

// After — spans are shipped to Grafana Cloud via OTLP/HTTP.
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http';
const exporter = new OTLPTraceExporter({
  url: process.env.OTEL_EXPORTER_OTLP_ENDPOINT + '/v1/traces',
});
```

Endpoint and auth are picked up from `OTEL_EXPORTER_OTLP_ENDPOINT` and
`OTEL_EXPORTER_OTLP_HEADERS`.

## Related

- [OTLP HTTP exporter for Node.js](https://opentelemetry.io/docs/languages/js/exporters/)
