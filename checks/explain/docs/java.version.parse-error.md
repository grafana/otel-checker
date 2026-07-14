---
id: java.version.parse-error
title: 'Could not parse a Java version constraint from the upstream catalog'
severity: error
---

## Why this matters

Despite what its name suggests, this error is *not* about your local Java
installation. It fires when the checker is walking the OpenTelemetry
Java auto-instrumentation catalog and finds a version constraint of the
form `Java N+` (e.g. `Java 8+`, `Java 17+`) where `N` couldn't be parsed
as an integer.

That indicates a shape change or corruption in the upstream
`instrumentation-list.yaml`: the catalog is downloaded from GitHub at
check time, and if the file the checker retrieved doesn't match the
version-string format it expects, the supported-libraries report for
that entry is skipped.

In practice this is a defensive check that shouldn't fire — the regex
already restricts the capture to a digit before we try to parse it.
Seeing it usually means either:

- The upstream YAML format changed in a way `otel-checker` doesn't
  handle yet.
- A cached / partially-downloaded catalog is corrupt.

## How to fix

1. Retry — a transient partial download will resolve on the next run:

   ```bash
   otel-checker check sdk --language=java
   ```

2. If it keeps firing, upgrade to the latest `otel-checker` release —
   the fetch and parse logic ships with the binary:

   ```bash
   go install github.com/grafana/otel-checker/cmd/otel-checker@latest
   ```

3. If it still fires after the upgrade, open an issue on the
   `grafana/otel-checker` repo. Include the offending version string
   from the error message so the mapping can be extended.

## Example

Compare the local expectation with what the upstream file actually
contains at the moment:

```bash
curl -s https://raw.githubusercontent.com/open-telemetry/opentelemetry-java-instrumentation/refs/heads/main/docs/instrumentation-list.yaml \
  | grep -o 'Java [0-9]*+' | sort -u
```

The expected shape is `Java <digit>+`. Any other form triggers this
error.

## Related

- [instrumentation-list.yaml (upstream source of truth)](https://github.com/open-telemetry/opentelemetry-java-instrumentation/blob/main/docs/instrumentation-list.yaml)
- `java.supported-libs.fetch-failed` — related error when the file
  itself can't be downloaded.
- `internal.java.semver` / `internal.java.version-range` — related
  internal errors when other version-string forms in the same catalog
  can't be parsed.
