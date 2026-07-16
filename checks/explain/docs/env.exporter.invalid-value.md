---
id: env.exporter.invalid-value
title: 'Exporter env var is set to an unsupported value'
severity: error
---

## Why this matters

`OTEL_TRACES_EXPORTER`, `OTEL_METRICS_EXPORTER`, and `OTEL_LOGS_EXPORTER`
accept a fixed set of values: `otlp` (the default, ships to any OTLP
endpoint including Grafana Cloud), `console` (prints to stdout for local
debugging), `none` (disables the signal), or unset (same as `otlp`). Any
other value is unrecognized by the SDK.

Behavior when the value is unrecognized varies by language: some SDKs
error out at startup and refuse to initialize, some silently fall back to
`none`, some log a warning that scrolls past before anyone notices.
Either way, telemetry for that signal doesn't reach Grafana Cloud.

## How to fix

Set the variable to one of the four supported values or remove it and let
the default apply:

```bash
export OTEL_TRACES_EXPORTER=otlp     # ship to your OTLP endpoint
export OTEL_METRICS_EXPORTER=otlp
export OTEL_LOGS_EXPORTER=otlp

# or, for local debugging:
export OTEL_TRACES_EXPORTER=console
```

If the value came from a template or config generator, check the substitution
— a common cause is an unresolved placeholder like `${EXPORTER}` being
stored literally.

Note that changing the exporter type usually means also setting the
matching `OTEL_EXPORTER_OTLP_ENDPOINT`, `OTEL_EXPORTER_OTLP_PROTOCOL`, and
`OTEL_EXPORTER_OTLP_HEADERS` — see the Grafana Cloud OTLP endpoint
reference for the right values.

## Example

Common misconfiguration and its fix:

```bash
# Before — 'otpl' is a typo.
export OTEL_TRACES_EXPORTER=otpl

# After.
export OTEL_TRACES_EXPORTER=otlp
```

## Related

- [OpenTelemetry exporter env-var specification](https://opentelemetry.io/docs/specs/otel/protocol/exporter/)
- [Grafana Cloud OTLP endpoint reference](https://grafana.com/docs/grafana-cloud/send-data/otlp/send-data-otlp/)
- `env.exporter.disabled` — related check for the specific value `none`.
