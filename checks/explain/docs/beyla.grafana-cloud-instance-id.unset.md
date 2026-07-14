---
id: beyla.grafana-cloud-instance-id.unset
title: 'GRAFANA_CLOUD_INSTANCE_ID is not set'
severity: error
---

## Why this matters

When Beyla ships telemetry directly to Grafana Cloud (rather than through
a local OpenTelemetry Collector), `GRAFANA_CLOUD_INSTANCE_ID` identifies
which Grafana Cloud stack the data belongs to. It's the numeric
"instance ID" you can find on your stack's OTLP page in the Grafana Cloud
UI (**Connections → OpenTelemetry (OTLP) → Send data**). Combined with
`GRAFANA_CLOUD_API_KEY`, it forms the basic-auth credentials Beyla uses
to authenticate every request.

Without it, Beyla can't authenticate and the Grafana Cloud endpoint
rejects every request with `401 Unauthorized`.

Note: this variable is part of Beyla's built-in Grafana Cloud submission
shortcut. If you're routing through a Collector instead (using OTLP env
vars like `OTEL_EXPORTER_OTLP_HEADERS`), Beyla ignores the
`GRAFANA_CLOUD_*` variables.

## How to fix

Copy the instance ID from your Grafana Cloud stack's OTLP page and set:

```bash
export GRAFANA_CLOUD_INSTANCE_ID="1234567"
```

Docker:

```dockerfile
ENV GRAFANA_CLOUD_INSTANCE_ID=1234567
```

Kubernetes — pull the value from a `Secret` rather than hardcoding it:

```yaml
env:
  - name: GRAFANA_CLOUD_INSTANCE_ID
    valueFrom:
      secretKeyRef:
        name: grafana-cloud
        key: instance-id
```

## Example

Full direct-to-Grafana-Cloud env-var set:

```bash
export BEYLA_OPEN_PORT=8080
export GRAFANA_CLOUD_SUBMIT=metrics,traces
export GRAFANA_CLOUD_INSTANCE_ID="1234567"
export GRAFANA_CLOUD_API_KEY="glc_eyJvIjoi..."
```

## Related

- [Beyla: Grafana Cloud configure](https://grafana.com/docs/beyla/latest/configure/)
- `beyla.grafana-cloud-api-key.unset` — companion check for the API key
  that pairs with this instance ID.
- `beyla.grafana-cloud-submit.unset` — companion check for the
  signal-list variable.
