---
id: env.node-options.recommended-unset
title: 'The recommended environment variable NODE_OPTIONS is unset'
severity: warning
---

## Why this matters

The OpenTelemetry zero-code (auto-instrumentation) flow for Node.js works
by hooking into Node's module loader at process startup. To make that
happen, Node has to load `@opentelemetry/auto-instrumentations-node/register`
*before* any application code runs. The `NODE_OPTIONS` environment variable
is the standard way to inject that `--require` flag without touching your
source.

If `NODE_OPTIONS` is unset — or set to something that doesn't include the
register hook — the auto-instrumentation never installs, and no telemetry
will be emitted from the popular libraries it would normally patch even 
though the SDK and instrumentation packages are installed.

## How to fix

Set `NODE_OPTIONS` so that the register hook runs at startup:

```bash
export NODE_OPTIONS="--require @opentelemetry/auto-instrumentations-node/register"
```

Or pass the `--require` (`-r`) flag directly when launching Node:

```bash
node --require @opentelemetry/auto-instrumentations-node/register app.js
```

The value must be visible to the process that *runs* your application —
not the shell that builds or packages it. In containers, set it in the
container's environment (e.g. Dockerfile `ENV` or Kubernetes `env:`), not
during the image build.

## Example

```bash
# One-time setup: install the packages.
npm install --save @opentelemetry/auto-instrumentations-node @opentelemetry/api

# Export NODE_OPTIONS in the same environment that starts the app.
export NODE_OPTIONS="--require @opentelemetry/auto-instrumentations-node/register"

# Start the app as usual.
node ./server.js
```

Verify the variable is visible to the running process with:

```bash
printenv NODE_OPTIONS
```

## Related

[OpenTelemetry Docs on NODE_OPTIONS for zero-code instrumentation](https://opentelemetry.io/docs/zero-code/js/#configuring-the-module)
