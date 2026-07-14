---
id: internal.java.version-range
title: 'Internal: invalid Java version range in the supported libraries data'
severity: internal
---

## Why this matters

This is an internal `otel-checker` error, not a problem with your
project. When walking the OpenTelemetry Java auto-instrumentation
catalog (`instrumentation-list.yaml`, fetched from GitHub at check
time), the checker splits each Maven-style entry into
`groupId:artifactId:version-range` and then parses the range with the
same range grammar it uses for the SDK check
(`[1.0.0,4.0.0)`, `[5.0,)`, and so on). This error fires when the range
portion itself can't be parsed — the delimiters are wrong, an endpoint
isn't valid semver, or the string contains an unexpected character.

The affected library will simply be missing from the supported-libraries
report. The rest of the checks still run.

## How to fix

You don't need to fix anything in your project; this is a signal to file
a bug so `otel-checker` can handle the new format.

1. Note the module name and version range quoted in the "Internal Error"
   message.
2. Confirm it's still upstream:

   ```bash
   curl -s https://raw.githubusercontent.com/open-telemetry/opentelemetry-java-instrumentation/refs/heads/main/docs/instrumentation-list.yaml \
     | grep -F "<the offending range>"
   ```

3. Open an issue on
   [grafana/otel-checker](https://github.com/grafana/otel-checker/issues)
   with the module name, the raw range string, and the checker version.
4. In the meantime, upgrade to the latest release in case the range
   grammar has already been extended:

   ```bash
   go install github.com/grafana/otel-checker/cmd/otel-checker@latest
   ```

## Example

The parser accepts single versions or half/full ranges:

```text
5.0.0                # single version
[5.0.0,)             # 5.0.0 and up
[1.0.0,4.0.0)        # 1.0.0 up to (but not including) 4.0.0
(1.0.0,4.0.0]        # exclusive on the lower end, inclusive on the upper
```

Anything else — missing bracket, non-semver endpoint, extra commas —
triggers this error.

## Related

- [instrumentation-list.yaml (upstream source of truth)](https://github.com/open-telemetry/opentelemetry-java-instrumentation/blob/main/docs/instrumentation-list.yaml)
- `internal.java.semver` — related internal error when the whole Maven
  coordinate has the wrong shape.
- `internal.sdk.version-range` — same class of failure in the
  cross-language SDK check.
