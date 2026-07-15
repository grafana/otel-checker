---
id: beyla.grafana-cloud-api-key.unset
title: 'GRAFANA_CLOUD_API_KEY is not set'
severity: error
---

## Why this matters

When Beyla ships telemetry directly to Grafana Cloud (rather than through
a local OpenTelemetry Collector), `GRAFANA_CLOUD_API_KEY` is the token it
uses to authenticate. Combined with `GRAFANA_CLOUD_INSTANCE_ID`, it forms
the basic-auth credentials for the Grafana Cloud OTLP gateway. Without
it, the gateway rejects every request with `401 Unauthorized` and no
telemetry reaches your stack.

Note: this variable is part of Beyla's built-in Grafana Cloud submission
shortcut. If you're routing through a Collector instead (using OTLP env
vars like `OTEL_EXPORTER_OTLP_HEADERS`), Beyla ignores the
`GRAFANA_CLOUD_*` variables.

## How to fix

Generate an API token with `metrics:write` and `traces:write` scopes in
the Grafana Cloud UI (**Home → Configuration → API keys** or **Access
Policies**), then set:

```bash
export GRAFANA_CLOUD_API_KEY="glc_eyJvIjoi..."
```

Avoid checking the key into source, shell profiles, or Dockerfiles —
pull it from a secret store at runtime.

Docker with build-time secret injection (or an env-file):

```dockerfile
# Prefer runtime env or docker secrets:
ENV GRAFANA_CLOUD_API_KEY=${GRAFANA_CLOUD_API_KEY}
```

Kubernetes:

```yaml
env:
  - name: GRAFANA_CLOUD_API_KEY
    valueFrom:
      secretKeyRef:
        name: grafana-cloud
        key: api-key
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
- `beyla.grafana-cloud-instance-id.unset` — companion check for the
  instance ID that pairs with this key.
- `beyla.grafana-cloud-submit.unset` — companion check for the signal
  list.
