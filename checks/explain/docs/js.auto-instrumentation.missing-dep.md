---
id: js.auto-instrumentation.missing-dep
title: 'Required auto-instrumentation dependency is missing from package.json'
severity: error
---

## Why this matters

The OpenTelemetry zero-code (auto-instrumentation) flow for Node.js requires
two packages: `@opentelemetry/auto-instrumentations-node`, which bundles the
per-library patches, and `@opentelemetry/api`, which the patches call into
to record spans, metrics, and logs. Both must be declared in `package.json`
so that `npm install` (or your CI's install step) puts them under
`node_modules/` where Node can load them at startup.

If either package is missing, the
`--require @opentelemetry/auto-instrumentations-node/register` hook fails
or the instrumented libraries have no API to call, and your service
produces no telemetry despite the rest of the setup looking correct.

## How to fix

Install the missing package(s) as regular dependencies (not
dev-dependencies, since they must be present in production):

```bash
npm install --save @opentelemetry/auto-instrumentations-node @opentelemetry/api
```

Confirm both entries appear under `"dependencies"` in `package.json` and
re-run `otel-checker check sdk --language=js`.

## Example

Minimal `package.json` fragment for zero-code Node.js instrumentation:

```json
{
  "dependencies": {
    "@opentelemetry/api": "^1.9.0",
    "@opentelemetry/auto-instrumentations-node": "^0.66.0"
  }
}
```

Then start your app with the register hook:

```bash
export NODE_OPTIONS="--require @opentelemetry/auto-instrumentations-node/register"
node ./server.js
```

## Related

- [OpenTelemetry Node.js zero-code setup](https://opentelemetry.io/docs/zero-code/js/)
