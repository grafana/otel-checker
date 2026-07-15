---
id: grafana-cloud.credentials.skipped
title: 'Credential test was skipped because prerequisites are not met'
severity: warning
---

## Why this matters

The Grafana Cloud check normally sends a real HTTP request to the OTLP
gateway with the configured `Authorization` header and reads back the
status code — that's the only way to distinguish a good token from a
plausibly-shaped bad one. The credential test is skipped when either:

- `OTEL_EXPORTER_OTLP_ENDPOINT` points at `localhost` (there's no cloud
  endpoint to hit; telemetry is going through a local Collector first),
  or
- `OTEL_EXPORTER_OTLP_ENDPOINT` or `OTEL_EXPORTER_OTLP_HEADERS` is empty
  (nothing to try).

The warning is informational: the token itself might be perfectly valid,
`otel-checker` just didn't verify it in this run.

## How to fix

Depending on which prerequisite is missing:

- **Localhost endpoint on purpose**: run
  `otel-checker check collector --collector-config-path=<path>` instead
  — the Collector-side check inspects its exporter's endpoint and auth
  headers.

- **Endpoint unset**: set `OTEL_EXPORTER_OTLP_ENDPOINT` to the Grafana
  Cloud OTLP gateway URL (see `grafana-cloud.endpoint.unset`).

- **Headers unset**: set `OTEL_EXPORTER_OTLP_HEADERS` with a valid
  `Authorization=Basic <token>` value (see
  `grafana-cloud.headers.missing-auth`).

Re-run once both variables are populated with cloud-facing values:

```bash
otel-checker check grafana-cloud --language=<lang>
```

## Example

Direct-to-cloud env vars needed for the credential check to run:

```bash
export OTEL_EXPORTER_OTLP_ENDPOINT="https://otlp-gateway-<zone>.grafana.net/otlp"
export OTEL_EXPORTER_OTLP_HEADERS="Authorization=Basic <base64-token>"
export OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf
```

## Related

- `grafana-cloud.endpoint.localhost` — related warning explaining why
  Cloud-side checks aren't exercised when the endpoint is localhost.
- `grafana-cloud.endpoint.unset` / `grafana-cloud.headers.missing-auth`
  — the missing prerequisites that most commonly trigger this warning.
