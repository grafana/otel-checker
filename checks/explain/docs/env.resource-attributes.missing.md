---
id: env.resource-attributes.missing
title: 'Recommended resource attribute is missing'
severity: warning
---

## Why this matters

`OTEL_RESOURCE_ATTRIBUTES` attaches a set of key-value pairs to every span,
metric, and log emitted by your service — Grafana Cloud uses them to filter,
group, and correlate telemetry across services, environments, versions, and
instances.

The checker flags these four attributes when they're missing:

- `service.namespace` — logical grouping for services (e.g. `shop`,
  `payments`). Lets you scope a Service Map to one product area.
- `deployment.environment.name` — `production`, `staging`, `dev`. Without
  it, production traces mingle with local ones on every dashboard.
- `service.instance.id` — unique per-process identifier (pod name, host,
  process id). Required to distinguish two replicas of the same service.
- `service.version` — the running build's version. Lets you see whether
  a spike in errors coincides with a deployment.

## How to fix

Set them via `OTEL_RESOURCE_ATTRIBUTES` in the environment that runs the
service. The value is comma-separated `key=value` pairs:

```bash
export OTEL_RESOURCE_ATTRIBUTES="service.namespace=shop,deployment.environment.name=production,service.instance.id=$HOSTNAME,service.version=1.4.2"
```

For attributes that vary per instance (`service.instance.id`), interpolate
a runtime value — `$HOSTNAME` on most Linux systems, `$POD_NAME` on
Kubernetes via the downward API, or a UUID generated at startup.

## Example

Kubernetes deployment fragment with all four attributes plus
`OTEL_SERVICE_NAME`:

```yaml
env:
  - name: OTEL_SERVICE_NAME
    value: checkout
  - name: POD_NAME
    valueFrom:
      fieldRef:
        fieldPath: metadata.name
  - name: OTEL_RESOURCE_ATTRIBUTES
    value: "service.namespace=shop,deployment.environment.name=production,service.instance.id=$(POD_NAME),service.version=1.4.2"
```

## Related

- [OpenTelemetry resource semantic conventions](https://opentelemetry.io/docs/specs/semconv/resource/)
- `env.otel-service-name.unset` — companion check for `service.name`, the
  most important resource attribute.
