---
id: grafana-cloud.endpoint.localhost
title: 'OTEL_EXPORTER_OTLP_ENDPOINT points at localhost'
severity: warning
---

## Why this matters

`OTEL_EXPORTER_OTLP_ENDPOINT` is set to a localhost URL, which normally
means the SDK is shipping telemetry to a local OpenTelemetry Collector
that forwards it on to Grafana Cloud. That's a supported and common
architecture — but it also means the Grafana Cloud checks can't verify
that telemetry actually reaches Grafana Cloud from this environment.
Anything downstream of the local Collector (endpoint format,
authentication headers, credential validity) is invisible to
`otel-checker`.

The warning fires so you know the Grafana Cloud portion of the checker
is not exercising the real gateway; it does not mean the setup is
broken.

## How to fix

Two options depending on which topology you want to validate:

1. **You want to check the Collector-then-cloud topology.** Leave the
   endpoint pointed at localhost and run the Collector-side check
   instead:

   ```bash
   otel-checker check collector --collector-config-path=./otel/
   ```

   That inspects the Collector's `config.yaml` for the Grafana Cloud
   exporter, its endpoint format, and auth headers — the checks that
   move here once a Collector is in the middle.

2. **You want to check direct-to-cloud from this environment.**
   Temporarily point the SDK straight at the gateway:

   ```bash
   export OTEL_EXPORTER_OTLP_ENDPOINT="https://otlp-gateway-prod-us-east-0.grafana.net/otlp"
   export OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf
   export OTEL_EXPORTER_OTLP_HEADERS="Authorization=Basic <token>"
   otel-checker check grafana-cloud --language=<lang>
   ```

## Example

Typical layout with a local Collector:

```text
SDK  →  localhost:4318 (Collector)  →  https://otlp-gateway-*.grafana.net/otlp
```

The Grafana Cloud endpoint/protocol/auth checks apply to the Collector's
config, not the SDK's environment variables — run
`otel-checker check collector` to cover that leg.

## Related

- [Grafana Cloud OTLP endpoint reference](https://grafana.com/docs/grafana-cloud/send-data/otlp/send-data-otlp/)
- `grafana-cloud.credentials.skipped` — related warning that fires
  when a localhost endpoint prevents the credential test from running.
- `grafana-cloud.endpoint.invalid-format` — companion check for
  non-localhost values that don't match the Grafana Cloud pattern.
