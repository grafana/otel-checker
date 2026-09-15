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
3. **A vendor / distribution-specific extension**: the top-level
   `distribution:` block is the spec-defined place for these; anything
   else is unofficial.

In all three cases, the SDK will silently drop the field at runtime.

## How to fix

- Compare the field name against the
  [schema reference](https://github.com/open-telemetry/opentelemetry-configuration/tree/main/schema).
  Fix typos, remove obsolete fields, or move vendor-specific settings
  under `distribution:`.
- If the field is legitimate and the checker's schema is out of
  date, regenerate `checks/config/config.gen.go` against the latest
  upstream schema:

  ```bash
  mise run generate-config-schema
  ```

- If the field is intentionally custom and you want the checker to
  stop warning about it, use `distribution:` (which is
  `additionalProperties: true` per spec and is skipped by this check).

## Related

- [OpenTelemetry Configuration schema — schema/](https://github.com/open-telemetry/opentelemetry-configuration/tree/main/schema)
- `config.file.parse-error` — related failure when the YAML itself
  is malformed rather than schema-drift.
