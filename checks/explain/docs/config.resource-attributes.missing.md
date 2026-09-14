---
id: config.resource-attributes.missing
title: 'The declarative config does not declare a recommended resource attribute'
severity: warning
---

## Why this matters

Resource attributes describe the service emitting telemetry — which
namespace it belongs to, which environment it runs in, which instance
it is, and which version is deployed. The Grafana Cloud UI groups,
filters, and correlates telemetry using these fields; missing values
degrade the experience (spans and metrics can't be traced back to a
specific deployment or version, alerts can't scope to an environment,
service maps show unnamed nodes).

When the declarative config file is the source of truth (which it is
when `--config-path` is passed or a default `otel-config.yaml` is
present), attributes belong under `resource.attributes` in that file
— NOT in the `OTEL_RESOURCE_ATTRIBUTES` environment variable.

## How to fix

Add each missing attribute under `resource.attributes` in the config
file:

```yaml
resource:
  attributes:
    - name: service.namespace
      value: shop
    - name: deployment.environment.name
      value: production
    - name: service.instance.id
      value: ${HOSTNAME}          # or ${POD_NAME} on Kubernetes
    - name: service.version
      value: "1.4.2"
```

Env-var substitution (`${VAR}` or `${VAR:-default}`) is supported, so
values that vary per-deployment can still be sourced from the runtime
environment without abandoning the declarative model:

```yaml
    - name: service.version
      value: ${SERVICE_VERSION:-dev}
```

Alternatively, `resource.attributes_list` accepts the same
comma-separated format as `OTEL_RESOURCE_ATTRIBUTES`, useful when a
platform injects attributes centrally:

```yaml
resource:
  attributes_list: ${OTEL_RESOURCE_ATTRIBUTES}
```

Per the [Resource schema](https://raw.githubusercontent.com/open-telemetry/opentelemetry-configuration/main/schema/resource.yaml),
`attributes` entries have HIGHER priority than `attributes_list`
entries when both declare the same attribute name.

## Related

- `env.resource-attributes.missing` — same finding when the checker
  is running against environment variables instead of a config file.
- `config.service-name.unset` — the service.name-specific variant.
