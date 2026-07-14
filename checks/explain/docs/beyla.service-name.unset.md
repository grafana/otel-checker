---
id: beyla.service-name.unset
title: 'BEYLA_SERVICE_NAME is not set'
severity: warning
---

## Why this matters

`BEYLA_SERVICE_NAME` sets the `service.name` resource attribute on every
span and metric Beyla emits — the label your backend uses to identify
which service the telemetry came from. When it's unset, Beyla falls
back to auto-detecting the name from the executable of the process it
attaches to, which is often something generic like `python`, `node`, or
the container's entrypoint script. That makes it impossible to tell one
Beyla-instrumented service apart from another in Grafana Cloud
dashboards, Service Map, and alerts.

## How to fix

Set `BEYLA_SERVICE_NAME` in the environment that runs Beyla — the same
environment where you set `BEYLA_OPEN_PORT`. Use a stable, human-readable
name (kebab-case or dotted); pick something that describes the service,
not the process:

```bash
export BEYLA_SERVICE_NAME="checkout"
```

Docker:

```dockerfile
ENV BEYLA_SERVICE_NAME=checkout
```

Kubernetes:

```yaml
env:
  - name: BEYLA_SERVICE_NAME
    value: checkout
```

If you're running Beyla as a sidecar, match the value to the
application container's own `OTEL_SERVICE_NAME` so traces stitched
together across both instrumentations report the same service.

## Example

Minimal Beyla env-var set:

```bash
export BEYLA_SERVICE_NAME="checkout"
export BEYLA_OPEN_PORT=8080
export GRAFANA_CLOUD_SUBMIT=metrics,traces
export GRAFANA_CLOUD_INSTANCE_ID="1234567"
export GRAFANA_CLOUD_API_KEY="glc_eyJvIjoi..."
```

## Related

- [Beyla configuration: service naming](https://grafana.com/docs/beyla/latest/configure/options/)
- `env.otel-service-name.unset` — the equivalent check on the
  OpenTelemetry SDK side; both control the `service.name` attribute.
