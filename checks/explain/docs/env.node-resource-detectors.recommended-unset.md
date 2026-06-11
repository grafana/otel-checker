---
id: env.node-resource-detectors.recommended-unset
title: 'The recommended environment variable OTEL_NODE_RESOURCE_DETECTORS is unset'
severity: warning
---

## Why this matters

`OTEL_NODE_RESOURCE_DETECTORS` controls which built-in resource detectors the
OpenTelemetry Node.js SDK runs when it starts up. Each detector attaches a
set of attributes to every emitted span, metric, and log — host name, OS
type, process identity, container metadata, cloud account, and so on — so
Grafana can correlate your telemetry with the
infrastructure it came from. When the variable is unset or doesn't include
the common detectors, your service telemetry shows up without identifying
attributes and becomes harder to filter, group, or alert on.

Setting the value to `all` will enable every detector that ships
with `@opentelemetry/auto-instrumentations-node`. If you'd rather opt in to
specific ones, the minimum recommended set is `env`, `host`, `os`, and
`serviceinstance`.

## How to fix

Set `OTEL_NODE_RESOURCE_DETECTORS` to one of:

- `all` — enables every built-in detector.
- A comma-separated list — pick specific detectors. Include at minimum
  `env,host,os,serviceinstance` to keep the standard correlation
  attributes.

The value must be set in the environment of the process that runs the
Node.js application — not the shell that *builds* it.

## Example

```bash
# Recommended: let every built-in detector run.
export OTEL_NODE_RESOURCE_DETECTORS=all

# Or, opt in to specific detectors:
export OTEL_NODE_RESOURCE_DETECTORS=env,host,os,serviceinstance

```

You can verify the variable is visible to the running process with
`printenv OTEL_NODE_RESOURCE_DETECTORS` inside the same shell or container.

## Related

[OpenTelemetry Docs on OTEL_NODE_RESOURCE_DETECTORS](https://opentelemetry.io/docs/zero-code/js/configuration/#sdk-resource-detector-configuration)
