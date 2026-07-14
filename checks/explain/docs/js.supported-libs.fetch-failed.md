---
id: js.supported-libs.fetch-failed
title: 'Could not fetch the supported JavaScript libraries list'
severity: error
---

## Why this matters

The supported-libraries check compares the dependencies in your project
against a curated list of Node.js libraries that have OpenTelemetry
instrumentation available. That list is embedded in the `otel-checker`
binary as `supported-libraries.yaml`. If parsing the embedded file fails,
the check can't tell you which of your dependencies have instrumentation
available — one of the most useful outputs of `check sdk --language=js`.

The most common cause is running a locally-built binary against a stale
or malformed `supported-libraries.yaml` (e.g. after an interrupted
`mise run generate`).

## How to fix

If you're running an installed release, upgrade to the latest version —
the embedded list ships with the binary:

```bash
go install github.com/grafana/otel-checker/cmd/otel-checker@latest
```

If you're developing locally, regenerate the supported-libraries YAML
from the upstream OTel contrib repos and rebuild:

```bash
mise run generate
mise run build
```

Then re-run the check.

## Example

Verify the embedded file parses by running the checker on a small
project and looking for the "Supported Libraries" section in the output:

```bash
otel-checker check sdk --language=js
```

## Related

- [opentelemetry-js-contrib](https://github.com/open-telemetry/opentelemetry-js-contrib) —
  upstream source of the supported-libraries list.
- `AGENTS.md` in this repo — describes `mise run generate`.
