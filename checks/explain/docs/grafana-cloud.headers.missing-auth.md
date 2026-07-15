---
id: grafana-cloud.headers.missing-auth
title: 'OTEL_EXPORTER_OTLP_HEADERS is missing the Authorization header'
severity: error
---

## Why this matters

Grafana Cloud's OTLP gateway requires an `Authorization: Basic <token>`
header on every OTLP request — the token is base64(`instance_id:api_key`)
and identifies which Grafana Cloud stack the data belongs to. The
OpenTelemetry SDK reads that header from the `OTEL_EXPORTER_OTLP_HEADERS`
environment variable, formatted as comma-separated `key=value` pairs.

If `Authorization=Basic <token>` is missing from the value, the gateway
rejects the requests with `401 Unauthorized`, and no telemetry reaches
Grafana Cloud.

## How to fix

Generate the token from your Grafana Cloud stack's OTLP page, then
export it. Linux / macOS:

```bash
export OTEL_EXPORTER_OTLP_HEADERS="Authorization=Basic <base64-token>"
```

Windows (PowerShell):

```powershell
$env:OTEL_EXPORTER_OTLP_HEADERS = "Authorization=Basic <base64-token>"
```

To build the token yourself, base64-encode `instance_id:api_key`. See
the Example section below for a `printf | base64` (bash) and
`[Convert]::ToBase64String(...)` (PowerShell) recipe.

The header name is case-sensitive per HTTP; the checker looks for
`Authorization=Basic` followed by a space and the base64 token. Python
SDKs URL-encode the value, so use `Authorization=Basic%20<base64-token>`
when setting for a Python service.

Avoid checking the token into a shell profile or Dockerfile — pull it
from a secret store at runtime. In Kubernetes:

```yaml
env:
  - name: OTEL_EXPORTER_OTLP_HEADERS
    valueFrom:
      secretKeyRef:
        name: grafana-cloud-otlp
        key: authorization        # "Authorization=Basic <base64-token>"
```

## Example

Building the token manually on Linux / macOS:

```bash
INSTANCE_ID="1234567"                  # from Grafana Cloud UI
API_KEY="glc_eyJvIjoiO..."             # from Grafana Cloud UI
TOKEN=$(printf '%s:%s' "$INSTANCE_ID" "$API_KEY" | base64)
export OTEL_EXPORTER_OTLP_HEADERS="Authorization=Basic $TOKEN"
```

Same on Windows (PowerShell):

```powershell
$InstanceId = "1234567"
$ApiKey = "glc_eyJvIjoiO..."
$Token = [Convert]::ToBase64String([Text.Encoding]::UTF8.GetBytes("$InstanceId`:$ApiKey"))
$env:OTEL_EXPORTER_OTLP_HEADERS = "Authorization=Basic $Token"
```

## Related

- [Grafana Cloud OTLP authentication](https://grafana.com/docs/grafana-cloud/send-data/otlp/send-data-otlp/)
- `grafana-cloud.credentials.unauthorized` — companion check that
  actually calls the endpoint and reports what the gateway said.
- `grafana-cloud.credentials.skipped` — companion warning when the
  credential check can't run because prerequisites are missing.
