---
id: grafana-cloud.signal-endpoint.invalid-format
title: 'Signal-specific OTLP endpoint is not a Grafana Cloud OTLP gateway URL'
severity: error
---

## Why this matters

OpenTelemetry SDKs read three optional signal-specific endpoint
variables that override the base `OTEL_EXPORTER_OTLP_ENDPOINT` for one
signal each:

- `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT`
- `OTEL_EXPORTER_OTLP_METRICS_ENDPOINT`
- `OTEL_EXPORTER_OTLP_LOGS_ENDPOINT`

Unlike the base variable, these are *not* transformed by the SDK — the
value is used verbatim. That means they must include the signal path
themselves: `/v1/traces`, `/v1/metrics`, `/v1/logs`. The Grafana Cloud
check requires each set variable to match
`https://*.grafana.net/otlp/v1/<signal>`. Anything else fails: the
gateway rejects the request or the SDK falls through to a wrong URL.

## How to fix

Set the affected variable to the Grafana Cloud gateway URL with the
signal path appended. The finding message names which variable and
which signal is expected.

Linux / macOS:

```bash
export OTEL_EXPORTER_OTLP_TRACES_ENDPOINT="https://otlp-gateway-<zone>.grafana.net/otlp/v1/traces"
export OTEL_EXPORTER_OTLP_METRICS_ENDPOINT="https://otlp-gateway-<zone>.grafana.net/otlp/v1/metrics"
export OTEL_EXPORTER_OTLP_LOGS_ENDPOINT="https://otlp-gateway-<zone>.grafana.net/otlp/v1/logs"
```

Windows (PowerShell):

```powershell
$env:OTEL_EXPORTER_OTLP_TRACES_ENDPOINT = "https://otlp-gateway-<zone>.grafana.net/otlp/v1/traces"
$env:OTEL_EXPORTER_OTLP_METRICS_ENDPOINT = "https://otlp-gateway-<zone>.grafana.net/otlp/v1/metrics"
$env:OTEL_EXPORTER_OTLP_LOGS_ENDPOINT = "https://otlp-gateway-<zone>.grafana.net/otlp/v1/logs"
```

Common mistakes that trigger this error:

- Missing the signal path (`.../otlp` without `/v1/<signal>`).
- Wrong signal path (`.../otlp/v1/traces` set in
  `OTEL_EXPORTER_OTLP_METRICS_ENDPOINT`).
- Trailing extra path segments.
- `http://` instead of `https://`.
- Host that doesn't end with `.grafana.net`.

## Example

Before / after — the `_TRACES_ENDPOINT` value was left without the
signal path:

```bash
# Before
export OTEL_EXPORTER_OTLP_TRACES_ENDPOINT="https://otlp-gateway-prod-us-east-0.grafana.net/otlp"

# After
export OTEL_EXPORTER_OTLP_TRACES_ENDPOINT="https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/traces"
```

You can also unset the signal-specific variable and rely on
`OTEL_EXPORTER_OTLP_ENDPOINT` — the SDK will append the signal path
itself.

## Related

- [OpenTelemetry OTLP endpoint configuration](https://opentelemetry.io/docs/specs/otel/protocol/exporter/)
- [Grafana Cloud OTLP endpoint reference](https://grafana.com/docs/grafana-cloud/send-data/otlp/send-data-otlp/)
- `grafana-cloud.endpoint.invalid-format` — companion check for the
  base `OTEL_EXPORTER_OTLP_ENDPOINT` variable, which has the *opposite*
  rule (must NOT include the signal path).
- `grafana-cloud.endpoint.unset` — companion check that fires when the
  base variable is missing AND at least one signal-specific variable is
  also missing.
