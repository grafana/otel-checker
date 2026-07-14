---
id: js.manual-instrumentation.node-options-conflict
title: 'NODE_OPTIONS conflicts with --manual-instrumentation'
severity: error
---

## Why this matters

You invoked `otel-checker check sdk --language=js --manual-instrumentation`,
which asserts that your application configures the OpenTelemetry SDK
explicitly in code — creating tracer providers, registering exporters, and
starting spans yourself. But `NODE_OPTIONS` is currently set to
`--require @opentelemetry/auto-instrumentations-node/register`, which tells
Node to load the *auto*-instrumentation register hook at process startup.

Running both at once produces duplicate SDK initialization: two tracer
providers, two sets of exporters, doubled spans, or, more commonly, one
side silently winning while the other's spans disappear. Neither the
manual nor the auto behavior is guaranteed to be the one you observe.

## How to fix

Pick one instrumentation style and remove the other. If you meant to use
manual instrumentation, unset the auto-instrumentation require hook in
every environment that runs the app:

```bash
unset NODE_OPTIONS
```

In containers, remove the `NODE_OPTIONS` `ENV` in the Dockerfile or the
Kubernetes pod spec. In systemd or init scripts, remove the export.

If you meant to use auto-instrumentation, drop the `--manual-instrumentation`
flag from your `otel-checker` invocation instead.

## Example

Verify that `NODE_OPTIONS` no longer sets the register hook:

```bash
printenv NODE_OPTIONS
# (empty output means unset)
```

Then re-run:

```bash
otel-checker check sdk --language=js --manual-instrumentation
```

## Related

- `env.node-options.recommended-unset` — the counterpart warning when
  you *are* using auto-instrumentation and `NODE_OPTIONS` is unset.
- [OpenTelemetry Node.js instrumentation styles](https://opentelemetry.io/docs/languages/js/instrumentation/)
