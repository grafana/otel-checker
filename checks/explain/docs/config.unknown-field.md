---
id: config.unknown-field
title: 'The declarative config declares a field not in the OpenTelemetry Configuration schema'
severity: warning
---

## Why this matters

The checker parses `--config-path` in strict mode against the model
generated from the upstream
[OpenTelemetry Configuration schema](https://github.com/open-telemetry/opentelemetry-configuration).
A field the file declares but the schema doesn't know about is one of
three things:

1. **A typo**: `schedulers_delay` instead of `schedule_delay`,
   `resources` (plural) instead of `resource`, etc.
2. **A deprecated or renamed field**: the schema evolved and the
   file is on an older shape.
3. **A checker schema that's out of date**: the OpenTelemetry
   Configuration schema added a field the checker hasn't been
   regenerated against yet.

The first two are typically real problems. The third is a false
positive; the checker doesn't have the context.

## How to fix

- If the field looks like a typo, compare against the
  [schema reference](https://github.com/open-telemetry/opentelemetry-configuration/tree/main/schema)
  and correct it in the config.
- If the field is legitimate and the checker's schema is out of
  date, please
  [open an issue](https://github.com/grafana/otel-checker/issues/new)
  so we can regenerate the model.

## Related

- [OpenTelemetry Configuration schema — schema/](https://github.com/open-telemetry/opentelemetry-configuration/tree/main/schema)
- `config.file.parse-error` — related failure when the YAML itself
  is malformed rather than schema-drift.
