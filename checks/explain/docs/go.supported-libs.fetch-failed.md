---
id: go.supported-libs.fetch-failed
title: 'Could not fetch the supported Go libraries list'
severity: error
---

## Why this matters

The supported-libraries check compares each dependency in your `go.mod`
against a curated catalog of Go libraries that have OpenTelemetry
instrumentation available. That catalog is embedded in the `otel-checker`
binary as `supported-libraries.yaml` under `checks/sdk/go/`. If parsing
the embedded file fails, the check can't tell you which of your
dependencies have instrumentation available — the main output of
`check sdk --language=go`.

The most common cause is running a locally-built binary against a stale
or malformed `supported-libraries.yaml` (e.g. after an interrupted
`mise run generate`).

## How to fix

If you're running an installed release, upgrade to the latest version —
the embedded catalog ships with the binary:

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

Verify the check works on a small module:

```bash
cd $(mktemp -d) && go mod init tmp && go get go.opentelemetry.io/otel@latest
otel-checker check sdk --language=go
```

## Related

- [opentelemetry-go-contrib](https://github.com/open-telemetry/opentelemetry-go-contrib) —
  upstream source of the supported-libraries list.
- `go.go-mod.unreadable` — related failure when your project's
  dependency file itself can't be read.
- `internal.sdk.version-range` — related internal error when a version
  range in the catalog can't be parsed.
