---
id: config.env-var.unresolved
title: 'The declarative config references an environment variable that is unset and has no default'
severity: error
---

## Why this matters

The declarative config file uses `${VAR}` (or the equivalent
`${env:VAR}`) substitution to inject a value from the process
environment
([spec](https://opentelemetry.io/docs/specs/otel/configuration/data-model/#environment-variable-substitution)).
The checker evaluated the substitution using the environment it was
launched with, found the variable unset, and could not use a fallback
because none was provided after `:-`.

At application startup, the SDK would either refuse to load the file
or resolve the placeholder to an empty string — either way the field
is broken and no telemetry for the
affected signal reaches Grafana Cloud.

## How to fix

Pick one of:

1. Export the referenced variable in the environment the application
   runs in (Dockerfile `ENV`, Kubernetes pod spec `env:`, systemd
   unit, `.envrc`, …). Then re-run the checker with the same
   environment so it sees the value.

2. Add a fallback in the config file itself:

   ```yaml
   endpoint: ${OTEL_EXPORTER_OTLP_ENDPOINT:-http://localhost:4318}/v1/traces
   ```

   The fallback is used verbatim when the variable is unset, so it
   must be a valid URL (or whatever the field expects).

3. Replace the substitution with the literal value the customer
   intends to use:

   ```yaml
   endpoint: https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/traces
   ```

## How the checker resolves substitutions

- `${VAR}` and `${env:VAR}` — read from `os.Getenv("VAR")`. Unset or
  empty triggers this finding.
- `${VAR:-default}` — read from the environment; use `default`
  verbatim when unset or empty. The checker then validates the
  resolved value.

The checker resolves substitutions using the environment it was
launched with.

## Related

- [OpenTelemetry Configuration spec — environment variable substitution](https://opentelemetry.io/docs/specs/otel/configuration/data-model/#environment-variable-substitution)
- `grafana-cloud.endpoint.unset` — sibling failure when the
  configuration path is via environment variables instead of the
  declarative file.
