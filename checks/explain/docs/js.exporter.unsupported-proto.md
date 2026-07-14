---
id: js.exporter.unsupported-proto
title: '@opentelemetry/exporter-trace-otlp-proto is not supported by Grafana Cloud'
severity: error
---

## Why this matters

OTLP has two HTTP variants: `http/protobuf` and `http/json`. Grafana Cloud's
OTLP gateway (`https://otlp-gateway-<zone>.grafana.net/otlp`) accepts
`http/protobuf`, but the specific package `@opentelemetry/exporter-trace-
otlp-proto` targets a different combination that Grafana Cloud does not
accept. If it's the exporter your instrumentation file wires up, exports
will be rejected and no traces will reach Grafana Cloud even though
everything else — endpoint, credentials, resource attributes — is set up
correctly.

## How to fix

Replace the dependency with `@opentelemetry/exporter-trace-otlp-http`:

```bash
npm uninstall @opentelemetry/exporter-trace-otlp-proto
npm install --save @opentelemetry/exporter-trace-otlp-http
```

Then update any imports in your instrumentation file to point at the new
package and re-run `otel-checker check sdk --language=js`.

## Example

```js
// Before — rejected by Grafana Cloud.
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-proto';

// After — works with the Grafana Cloud OTLP gateway.
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http';

const exporter = new OTLPTraceExporter({
  url: `${process.env.OTEL_EXPORTER_OTLP_ENDPOINT}/v1/traces`,
});
```

## Related

- [Grafana Cloud OTLP endpoint reference](https://grafana.com/docs/grafana-cloud/send-data/otlp/send-data-otlp/)
