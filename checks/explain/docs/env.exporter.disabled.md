---
id: env.exporter.disabled
title: 'Exporter env var is set to "none"'
severity: error
---

## Why this matters

`OTEL_TRACES_EXPORTER`, `OTEL_METRICS_EXPORTER`, and `OTEL_LOGS_EXPORTER`
pick which exporter each signal flows through. Setting one to `none` is a
supported OTel value, but its explicit meaning is *disable this signal*
— the SDK still initializes, still runs your instrumentation, still
creates spans and records metrics, but drops all of them instead of
shipping them anywhere.

That configuration is intentional in some environments (a batch job that
should only emit logs, a test suite that shouldn't page anyone), but in a
service that's supposed to be observed it looks identical to a completely
misconfigured setup: no data reaches Grafana Cloud, dashboards stay
empty, and no error surfaces to point at the cause.

## How to fix

Either remove the variable entirely (the default is `otlp`), or set it to
one of the supported values:

```bash
unset OTEL_TRACES_EXPORTER          # falls back to the default 'otlp'
# or
export OTEL_TRACES_EXPORTER=otlp    # explicit
export OTEL_TRACES_EXPORTER=console # print to stdout for debugging
```

In containers, remove the `ENV` line from the Dockerfile or drop the
entry from the Kubernetes pod spec. In systemd or init scripts, remove
the export.

If you *do* mean to disable a signal, prefer scoping the change to the
one variable that needs it (rather than all three) and add a comment near
where it's set so a future reader knows it's intentional.

## Example

Enable OTLP for every signal:

```bash
export OTEL_EXPORTER_OTLP_ENDPOINT="https://otlp-gateway-prod-us-east-0.grafana.net/otlp"
export OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf
export OTEL_EXPORTER_OTLP_HEADERS="Authorization=Basic <token>"

# Leave the exporter env vars unset — the default is otlp for all three.
```

## Related

- [OpenTelemetry exporter env-var specification](https://opentelemetry.io/docs/specs/otel/protocol/exporter/)
