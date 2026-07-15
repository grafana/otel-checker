---
id: java.library.unsupported
title: 'A Java library is not auto-instrumented by OpenTelemetry'
severity: warning
---

## Why this matters

For each library the checker finds in your Maven or Gradle dependencies,
it cross-references the OpenTelemetry Java auto-instrumentation catalog
([instrumentation-list.yaml](https://github.com/open-telemetry/opentelemetry-java-instrumentation/blob/main/docs/instrumentation-list.yaml)).
When a library is *not* covered — no matching entry, or a matching entry
that requires a Java version above the one you're running — you'll only
get telemetry for that library if you add manual instrumentation
yourself.

This warning only fires in debug mode (`--debug`). By default it's
suppressed, since projects typically have dozens of dependencies whose
instrumentation status doesn't concern the operator today.

## How to fix

Nothing to fix on the checker side — the warning is informational. What
to do depends on how important the library is to your observability:

- **Library that you actually care about tracing** (a database driver
  your app hits on every request, an HTTP client, a message broker):
  - Check whether the specific *version* is supported. The check
    matches on version ranges from the upstream YAML; upgrading the
    library sometimes brings it into a supported range.
  - Search the [OpenTelemetry Java contrib repo](https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/instrumentation)
    for the library — new modules are added regularly.
  - Consider adding manual instrumentation with the OpenTelemetry
    API for the code paths that matter.

- **Library that doesn't need tracing** (a build-time tool, a
  logging framework, an assertion library): the warning is safe to
  ignore.

- **Actively want it instrumented but no upstream support exists**:
  open an issue in the
  [opentelemetry-java-instrumentation](https://github.com/open-telemetry/opentelemetry-java-instrumentation)
  repo — the contributors add libraries based on community demand.

## Example

Manual instrumentation with the OpenTelemetry API for a library that has
no built-in coverage:

```java
import io.opentelemetry.api.GlobalOpenTelemetry;
import io.opentelemetry.api.trace.Span;
import io.opentelemetry.api.trace.Tracer;

Tracer tracer = GlobalOpenTelemetry.getTracer("my.custom.instr");

Span span = tracer.spanBuilder("db.query").startSpan();
try {
    return myLibrary.doWork();
} finally {
    span.end();
}
```

## Related

- [OpenTelemetry Java auto-instrumentation library list](https://github.com/open-telemetry/opentelemetry-java-instrumentation/blob/main/docs/supported-libraries.md)
- [OpenTelemetry Java manual instrumentation](https://opentelemetry.io/docs/languages/java/instrumentation/)
