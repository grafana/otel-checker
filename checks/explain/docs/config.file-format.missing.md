---
id: config.file-format.missing
title: 'The declarative config file does not declare a file_format'
severity: error
---

## Why this matters

The top-level `file_format` key tells the SDK which schema version to
apply when parsing the rest of the document
([spec](https://opentelemetry.io/docs/specs/otel/configuration/data-model/#yaml-file-format)).
Without it, different SDKs will pick different defaults — most refuse
to load the file at all, which surfaces at process startup as an
opaque configuration error rather than at review time. Pinning the
version also protects the file from silently changing meaning when the
SDK is upgraded and the default schema shifts.

## How to fix

Add `file_format` at the top of the file. Use the version your SDK
supports — `1.1` is the current stable release; `0.4` is the last
0.x line still accepted by some SDKs:

```yaml
file_format: "1.1"

resource:
  # ...
```

Quote the value. YAML parses an unquoted `1.10` as the float `1.1`
and drops the trailing zero — quote it to keep the exact string.

## Related

- [OpenTelemetry Configuration spec — file_format](https://opentelemetry.io/docs/specs/otel/configuration/data-model/#yaml-file-format)
- [Schema changelog](https://github.com/open-telemetry/opentelemetry-configuration/blob/main/CHANGELOG.md)
