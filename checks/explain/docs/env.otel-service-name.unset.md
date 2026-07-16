---
id: env.otel-service-name.unset
title: 'OTEL_SERVICE_NAME is not set'
severity: warning
---

## Why this matters

`OTEL_SERVICE_NAME` sets the `service.name` resource attribute — the label
that identifies the emitting service on every span, metric, and log line
your app sends. It is the single most important resource attribute in
OpenTelemetry: Grafana Cloud (and any other backend) uses it to group,
filter, and route telemetry into Application Observability views,
per-service dashboards, and Service Map nodes.

When it's unset (and `service.name` is not present in
`OTEL_RESOURCE_ATTRIBUTES` either), the SDK falls back to a generic value
like `unknown_service:node`. Telemetry from every service in your account
that also skipped setting the name collapses under that same label,
making it impossible to tell which app emitted what.

## How to fix

Set `OTEL_SERVICE_NAME` in the environment that runs your service — not
the shell that builds or packages it. Use kebab-case or dotted names and
keep it stable across deploys of the same service:

```bash
export OTEL_SERVICE_NAME="checkout"
```

In Docker:

```dockerfile
ENV OTEL_SERVICE_NAME=checkout
```

In Kubernetes:

```yaml
env:
  - name: OTEL_SERVICE_NAME
    value: checkout
```

Alternatively, include `service.name=<value>` in `OTEL_RESOURCE_ATTRIBUTES`.
Per the OpenTelemetry specification, when both are set, `OTEL_SERVICE_NAME`
wins.

## Example

Minimal setup for a single service:

```bash
export OTEL_SERVICE_NAME="checkout"
export OTEL_RESOURCE_ATTRIBUTES="service.namespace=shop,deployment.environment.name=production,service.version=1.4.2"
```

## Related

- [OpenTelemetry service.name specification](https://opentelemetry.io/docs/specs/semconv/resource/#service)
- `env.resource-attributes.missing` — companion check for the other
  recommended resource attributes (`service.namespace`,
  `deployment.environment.name`, etc.).
