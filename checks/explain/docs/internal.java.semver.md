---
id: internal.java.semver
title: 'Internal: invalid Maven coordinate in the supported libraries data'
severity: internal
---

## Why this matters

This is an internal `otel-checker` error, not a problem with your
project. When walking the OpenTelemetry Java auto-instrumentation
catalog (`instrumentation-list.yaml`, fetched from GitHub at check
time), the checker expects each version entry to be a Maven coordinate
of the form `groupId:artifactId:version-range` — three colon-separated
parts, e.g. `com.amazonaws:aws-lambda-java-core:[1.0.0,)`. This error
fires when a catalog entry doesn't have exactly three parts, so the
checker can't split it into group/artifact/range and skips that entry.

The affected library will simply be missing from the supported-libraries
report. The rest of the checks still run.

Common causes:

- The upstream YAML format changed and now uses a different shape for
  version strings — a shape `otel-checker` doesn't yet know about.
- The entry itself is genuinely malformed upstream (missing a colon, a
  stray leading whitespace).

## How to fix

You don't need to fix anything in your project; this is a signal to file
a bug so `otel-checker` can handle the new format.

1. Note the version string quoted in the "Internal Error" message.
2. Grep the upstream file for it to confirm it's still there:

   ```bash
   curl -s https://raw.githubusercontent.com/open-telemetry/opentelemetry-java-instrumentation/refs/heads/main/docs/instrumentation-list.yaml \
     | grep -F "<the offending string>"
   ```

3. Open an issue on
   [grafana/otel-checker](https://github.com/grafana/otel-checker/issues)
   including the module name, the raw version string, and the checker
   version (`otel-checker version`).
4. In the meantime, upgrade to the latest release — the mapping may
   already have been added:

   ```bash
   go install github.com/grafana/otel-checker/cmd/otel-checker@latest
   ```

## Example

The parser expects entries like:

```text
com.amazonaws:aws-lambda-java-core:[1.0.0,)
org.springframework:spring-web:[5.0.0,)
```

Anything else (four parts, two parts, wrapped in extra syntax) triggers
this error.

## Related

- [instrumentation-list.yaml (upstream source of truth)](https://github.com/open-telemetry/opentelemetry-java-instrumentation/blob/main/docs/instrumentation-list.yaml)
- `java.supported-libs.fetch-failed` — related failure when the catalog
  can't be downloaded at all.
