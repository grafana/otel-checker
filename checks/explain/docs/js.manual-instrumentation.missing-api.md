---
id: js.manual-instrumentation.missing-api
title: 'Manual instrumentation is missing the @opentelemetry/api dependency'
severity: error
---

## Why this matters

Every manual OpenTelemetry Node.js integration ultimately routes calls
through `@opentelemetry/api` — that's the package exposing `trace.getTracer`,
`metrics.getMeter`, context propagation, and the type definitions your code
imports to start spans and record metrics. Without it declared in
`package.json`, `npm install` won't bring it in as a top-level dependency,
imports of `@opentelemetry/api` will resolve to whatever transitive version
your SDK pulled in (if any), and version drift can cause silent no-ops
where spans are never registered against the active tracer provider.

## How to fix

Install the package as a regular dependency:

```bash
npm install --save @opentelemetry/api
```

Confirm it appears under `"dependencies"` in `package.json`, then re-run
`otel-checker check sdk --language=js --manual-instrumentation`.

## Example

Minimal `package.json` fragment for manual instrumentation:

```json
{
  "dependencies": {
    "@opentelemetry/api": "^1.9.0",
    "@opentelemetry/sdk-node": "^0.208.0",
    "@opentelemetry/exporter-trace-otlp-http": "^0.208.0"
  }
}
```

Usage in code:

```js
import { trace } from '@opentelemetry/api';

const tracer = trace.getTracer('my-service');
const span = tracer.startSpan('handle-request');
// ... work ...
span.end();
```

## Related

- [OpenTelemetry Node.js manual instrumentation guide](https://opentelemetry.io/docs/languages/js/instrumentation/)
