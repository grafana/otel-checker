---
id: grafana-cloud.endpoint.not-grafana
title: 'OTLP endpoint is a valid URL but does not point at Grafana Cloud'
severity: warning
---

## Why this matters

The OTLP endpoint parses as a valid URL and is not `localhost`, but
the host does not match the Grafana Cloud pattern
(`*.grafana.net/otlp`). Telemetry from this exporter will be shipped
to whatever the URL resolves to — a third-party observability backend,
an in-cluster forwarder, or a mistyped hostname — not to the
customer's Grafana Cloud stack.

This is a warning, not an error, because a non-Grafana endpoint is a
legitimate configuration for many deployments (customer runs a
Collector or Alloy in front of Grafana Cloud, ships to multiple
backends, staging environment). The check surfaces it so the reviewer
can confirm the customer's intent.

## How to fix

- If Grafana Cloud is the intended backend, replace the endpoint with
  the Grafana Cloud OTLP URL for the customer's stack:

  ```yaml
  endpoint: https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/traces
  ```

  Look up the exact host in **Grafana Cloud → OpenTelemetry →
  Configure**.

- If a forwarder (Collector, Alloy, Beyla) is intentional, verify
  independently that the forwarder itself is pointed at Grafana Cloud
  (run `otel-checker check collector --collector-config-path=<path>`
  against that forwarder's config).

- If the endpoint is a typo — for example a customer moved from
  prod-us-east-0 to prod-eu-west-2 and updated the SDK config but not
  the exporter — correct the host and re-run the checker.

## Related

- `grafana-cloud.endpoint.localhost` — sibling warning when the
  endpoint points at `localhost`.
- `grafana-cloud.signal-endpoint.invalid-format` — related error when
  the endpoint isn't a valid URL at all.
- `collector.endpoint.not-grafana` — same finding on the Collector
  side.
