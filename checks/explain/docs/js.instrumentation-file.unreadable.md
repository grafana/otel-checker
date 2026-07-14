---
id: js.instrumentation-file.unreadable
title: 'The instrumentation file could not be read'
severity: error
---

## Why this matters

When you run `otel-checker check sdk --language=js
--manual-instrumentation`, the checker inspects the file that wires up your
tracer and meter providers so it can verify which exporters are in use
(OTLP vs Console) and other manual-instrumentation details. If the file
can't be read — wrong path, wrong working directory, permission denied,
or the file doesn't exist yet — the exporter checks that depend on it are
skipped, and a real misconfiguration in that file will go unnoticed until
runtime.

## How to fix

1. Confirm the file path you passed to `--instrumentation-file`. The path
   is resolved from the directory where you invoke `otel-checker`.
2. Check the file actually exists at that path and that the current user
   can read it (`ls -l <path>`).
3. If the file lives elsewhere in the repo, pass an explicit path:

```bash
otel-checker check sdk --language=js --manual-instrumentation \
  --instrumentation-file=./src/telemetry/instrumentation.ts
```

The path can be absolute or relative to the current working directory.

## Example

Repo layout:

```text
my-app/
├── package.json
└── src/
    └── tracing.ts   ← instrumentation file
```

Correct invocation from `my-app/`:

```bash
otel-checker check sdk --language=js --manual-instrumentation \
  --instrumentation-file=src/tracing.ts
```

## Related

- `js.exporter.console-debug` — the check that inspects the
  instrumentation file for Console exporter usage.
