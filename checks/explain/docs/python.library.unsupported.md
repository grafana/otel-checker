---
id: python.library.unsupported
title: 'A library in requirements.txt is not auto-instrumented by OpenTelemetry'
severity: warning
---

## Why this matters

For each dependency in your `requirements.txt`, the checker cross-references
it against the OpenTelemetry Python contrib catalog (parsed live from
the [`opentelemetry-python-contrib`](https://github.com/open-telemetry/opentelemetry-python-contrib)
README on GitHub). When a library isn't in that catalog — no matching
entry, or a matching entry whose version range doesn't cover the version
you have pinned — you'll only get telemetry for that library if you add
manual instrumentation yourself.

This warning only fires in debug mode (`--debug`). By default it's
suppressed, since projects typically have dozens of dependencies whose
instrumentation status doesn't matter for day-to-day observability.

## How to fix

Nothing to fix on the checker side — the warning is informational. What
to do depends on how important the library is to your observability:

- **Library you care about tracing** (a database driver, an HTTP client,
  a web framework):
  - Check the specific *pinned version* against the catalog. Version
    ranges are strict: bumping (or occasionally pinning down to) a
    supported range can bring the library into coverage.
  - Search the
    [OpenTelemetry Python contrib repo](https://github.com/open-telemetry/opentelemetry-python-contrib/tree/main/instrumentation)
    — new instrumentation packages are added regularly.
  - Add manual instrumentation with the OpenTelemetry API for the code
    paths that matter.

- **Library that doesn't need tracing** (a linter, a test tool, a build
  helper): the warning is safe to ignore.

- **Actively want it instrumented but no upstream support exists**:
  open an issue on the
  [opentelemetry-python-contrib](https://github.com/open-telemetry/opentelemetry-python-contrib)
  repo — the maintainers add libraries based on community demand.

## Example

Manual instrumentation with the OpenTelemetry API for a library that
has no built-in coverage:

```python
from opentelemetry import trace

tracer = trace.get_tracer("my.custom.instr")

def query(sql: str):
    with tracer.start_as_current_span("custom-db.query") as span:
        span.set_attribute("db.statement", sql)
        return my_library.query(sql)
```

## Related

- [OpenTelemetry Python manual instrumentation](https://opentelemetry.io/docs/languages/python/instrumentation/)
- [opentelemetry-python-contrib instrumentation list](https://github.com/open-telemetry/opentelemetry-python-contrib/tree/main/instrumentation)
- `sdk.library.unsupported` — the cross-language variant of this check.
