---
id: config.service-name.unset
title: 'The declarative config does not declare service.name'
severity: warning
---

## Why this matters

`service.name` is the single most important resource attribute: every
Grafana Cloud view — Application Observability, Traces Drilldown,
Service Map, per-service dashboards — is keyed on it. Without it, the
service shows up as `unknown_service` (the SDK default when nothing is
declared) and gets grouped with every other unnamed service in the
account.

When a declarative config file is loaded, `service.name` must be
declared under `resource.attributes` in that file. Per spec, SDKs
ignore `OTEL_SERVICE_NAME` and `OTEL_RESOURCE_ATTRIBUTES` when a
declarative config is loaded — setting those env vars alongside a
config file is a no-op.

## How to fix

Add `service.name` under `resource.attributes`:

```yaml
resource:
  attributes:
    - name: service.name
      value: checkout
```

To keep the same "env var if present, default otherwise" behavior the
`OTEL_SERVICE_NAME` env var used to provide, use substitution:

```yaml
resource:
  attributes:
    - name: service.name
      value: ${OTEL_SERVICE_NAME:-checkout}
```

The checker resolves the substitution against the current environment
and validates the result — so a `:-default` clause is what determines
whether the checker sees `service.name` as "set."

## Related

- `env.otel-service-name.unset` — same finding when the checker is
  running against environment variables instead of a config file.
- `config.resource-attributes.missing` — the general recommended
  resource-attribute check.
