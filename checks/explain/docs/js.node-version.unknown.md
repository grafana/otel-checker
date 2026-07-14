---
id: js.node-version.unknown
title: 'Could not determine the installed Node.js version'
severity: error
---

## Why this matters

`otel-checker` runs `node -v` and parses the output to verify the minimum
Node.js version. This error means either the `node` binary isn't on the
current `PATH`, the command exited non-zero, or the output didn't look
like a version string (`v22.x.y`). Without a version to check, the tool
can't tell you whether your runtime is compatible with the OpenTelemetry
Node.js packages you're about to install — a version mismatch is one of
the most common reasons instrumentation appears to install cleanly but
never emits telemetry.

## How to fix

1. Verify Node is installed and callable from the same shell that runs
   `otel-checker`:

   ```bash
   which node
   node -v
   ```

2. If `node` isn't found, install it. Use a version manager like `nvm` or
   `fnm` so you don't need `sudo`, or install a supported release from
   [nodejs.org](https://nodejs.org/).

3. If Node is installed but only in another user's environment (common
   under version managers on CI runners), ensure the same shell profile
   is sourced. On GitHub Actions use `actions/setup-node` before invoking
   the checker.

## Example

GitHub Actions workflow that makes Node visible to `otel-checker`:

```yaml
- uses: actions/setup-node@v7
  with:
    node-version: '22'
- run: otel-checker check sdk --language=js
```

## Related

- `js.node-version.too-old` — related check that fires when Node is found
  but is older than the supported minimum.
- [Node.js downloads](https://nodejs.org/en/download)
