---
id: sdk.library.unsupported
title: 'A library is not auto-instrumented by OpenTelemetry'
severity: warning
---

## Why this matters

For each library the checker finds in your project's dependency
manifest (`package.json`, `requirements.txt`, `go.mod`, `.csproj`), it
cross-references the corresponding OpenTelemetry auto-instrumentation
catalog. When a library isn't in that catalog — no matching entry, or a
matching entry whose version range doesn't cover the version you have
installed — you'll only get telemetry for that library if you add
manual instrumentation yourself.

This warning only fires in debug mode (`--debug`). By default it's
suppressed, since projects typically have dozens of dependencies whose
instrumentation status doesn't matter for day-to-day observability.

## How to fix

Nothing to fix on the checker side — the warning is informational. What
to do depends on how important the library is to your observability:

- **Library you care about tracing** (a database client, an HTTP
  library, a message broker):
  - Check the specific *version* against the catalog. Version ranges
    are strict: bumping (or occasionally pinning down to) a supported
    range can bring the library into coverage.
  - Search the OpenTelemetry contrib repo for your language — new
    instrumentation modules are added regularly:
    [JS](https://github.com/open-telemetry/opentelemetry-js-contrib) ·
    [Python](https://github.com/open-telemetry/opentelemetry-python-contrib) ·
    [Go](https://github.com/open-telemetry/opentelemetry-go-contrib) ·
    [.NET](https://github.com/open-telemetry/opentelemetry-dotnet-contrib).
  - Add manual instrumentation with the language's OpenTelemetry API
    for the code paths that matter.

- **Library that doesn't need tracing** (a build-time tool, a logging
  helper, a test framework): the warning is safe to ignore.

- **Actively want it instrumented but no upstream support exists**:
  open an issue in the appropriate contrib repo — the maintainers add
  libraries based on community demand.

## Example

Manual instrumentation with the OpenTelemetry API for a library that
has no built-in coverage (Node.js):

```js
import { trace } from '@opentelemetry/api';

const tracer = trace.getTracer('my.custom.instr');

async function query(sql) {
  const span = tracer.startSpan('custom-db.query', {
    attributes: { 'db.statement': sql },
  });
  try {
    return await myLibrary.query(sql);
  } finally {
    span.end();
  }
}
```

## Related

- [OpenTelemetry manual instrumentation guides](https://opentelemetry.io/docs/languages/)
- `java.library.unsupported` — the Java-specific variant of this check.
