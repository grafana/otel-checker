---
id: grafana-cloud.credentials.unauthorized
title: 'Grafana Cloud rejected the configured credentials'
severity: error
---

## Why this matters

The credential test sent a real request to `<OTEL_EXPORTER_OTLP_ENDPOINT>
/v1/metrics` using the `Authorization` header from
`OTEL_EXPORTER_OTLP_HEADERS`, and the Grafana Cloud gateway responded
with `401 Unauthorized`. That's an authoritative rejection — the network
reached the gateway, but the token is missing, malformed, or revoked.
Your app running with the same environment variables will hit the same
`401` on every OTLP export, and no telemetry will land in Grafana Cloud.

Common causes:

- The token is missing the `Basic` scheme prefix (the value is the raw
  base64 payload, without a leading `Basic` followed by a space).
- The base64 payload isn't `instance_id:api_key` — for example, you
  encoded just the API key on its own.
- The API key was rotated or deleted in the Grafana Cloud UI.
- The token belongs to a different stack than the endpoint URL points
  at.

## How to fix

1. Re-generate the token from the Grafana Cloud stack that matches the
   endpoint URL. In the UI: **Connections → OpenTelemetry (OTLP) → Send
   data**. It gives you the exact `Authorization: Basic <token>` line
   to paste.

2. Or build it yourself.

   Linux / macOS:

   ```bash
   INSTANCE_ID="1234567"
   API_KEY="glc_eyJvIjoi..."           # from the same stack
   TOKEN=$(printf '%s:%s' "$INSTANCE_ID" "$API_KEY" | base64)
   export OTEL_EXPORTER_OTLP_HEADERS="Authorization=Basic $TOKEN"
   ```

   Windows (PowerShell):

   ```powershell
   $InstanceId = "1234567"
   $ApiKey = "glc_eyJvIjoi..."
   $Token = [Convert]::ToBase64String([Text.Encoding]::UTF8.GetBytes("$InstanceId`:$ApiKey"))
   $env:OTEL_EXPORTER_OTLP_HEADERS = "Authorization=Basic $Token"
   ```

   Check that `INSTANCE_ID` and the region embedded in the endpoint
   correspond to the same Grafana Cloud stack.

3. If you rotated the API key recently, update every environment (dev,
   staging, prod, CI secrets) — a stale value in any of them will hit
   this same error.

## Example

Full env-var set for a direct-to-cloud pipeline:

```bash
export OTEL_EXPORTER_OTLP_ENDPOINT="https://otlp-gateway-prod-us-east-0.grafana.net/otlp"
export OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf
export OTEL_EXPORTER_OTLP_HEADERS="Authorization=Basic $(printf '%s:%s' "$INSTANCE_ID" "$API_KEY" | base64)"
otel-checker check grafana-cloud --language=<lang>
```

## Related

- [Grafana Cloud OTLP authentication](https://grafana.com/docs/grafana-cloud/send-data/otlp/send-data-otlp/)
- `grafana-cloud.headers.missing-auth` — related check when the header
  is missing entirely instead of present-but-rejected.
- `grafana-cloud.credentials.network-error` — related error when the
  request never reached the gateway.
