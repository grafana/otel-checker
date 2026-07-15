---
id: beyla.grafana-cloud-submit.unset
title: 'GRAFANA_CLOUD_SUBMIT is not set'
severity: error
---

## Why this matters

When Beyla ships telemetry directly to Grafana Cloud (rather than through
a local OpenTelemetry Collector), `GRAFANA_CLOUD_SUBMIT` tells it which
signals to send. It's a comma-separated list — typical values are
`metrics`, `traces`, or `metrics,traces`. Without it, Beyla doesn't know
which pipelines to enable and will not submit anything to Grafana
Cloud.

Note: this variable is part of Beyla's built-in Grafana Cloud submission
shortcut. If you're routing through a Collector instead (using OTLP env
vars), Beyla ignores the `GRAFANA_CLOUD_*` variables and you don't need
this one.

## How to fix

Set `GRAFANA_CLOUD_SUBMIT` to the comma-separated list of signals you
want Beyla to send:

```bash
export GRAFANA_CLOUD_SUBMIT=metrics,traces
```

In Docker:

```dockerfile
ENV GRAFANA_CLOUD_SUBMIT=metrics,traces
```

In Kubernetes:

```yaml
env:
  - name: GRAFANA_CLOUD_SUBMIT
    value: "metrics,traces"
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
- `beyla.grafana-cloud-instance-id.unset` / `beyla.grafana-cloud-api-key.unset`
  — companion checks for the credentials Beyla uses to authenticate.
