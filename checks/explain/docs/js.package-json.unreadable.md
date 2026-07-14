---
id: js.package-json.unreadable
title: 'package.json could not be read'
severity: error
---

## Why this matters

`otel-checker` reads `package.json` to verify that the OpenTelemetry
packages your instrumentation depends on are declared as dependencies —
`@opentelemetry/api`, `@opentelemetry/auto-instrumentations-node`, and so
on — and to run the supported-libraries check against your dependency
list. If the file can't be read (missing, wrong directory, or permission
denied) all of those checks are skipped, and a project misconfiguration
will be missed.

## How to fix

1. Run `otel-checker` from the directory that contains your `package.json`
   (usually the repo root):

   ```bash
   ls package.json
   otel-checker check sdk --language=js
   ```

2. If your Node project lives in a subdirectory (e.g. a monorepo), `cd`
   into that package before running the checker, or pass an explicit
   path via `--package-json-path`:

   ```bash
   otel-checker check sdk --language=js --package-json-path=./services/api/
   ```

3. If the file exists but the process can't read it, check permissions:

   ```bash
   ls -l package.json
   ```

## Example

Monorepo layout:

```text
my-repo/
├── package.json      ← workspace root
└── services/
    └── api/
        └── package.json   ← service manifest to check
```

Invocation from `my-repo/`:

```bash
otel-checker check sdk --language=js --package-json-path=./services/api/
```

## Related

- `js.package-lock-json.unreadable` — related failure on the lock file.
- `js.auto-instrumentation.missing-dep` — a check that depends on
  `package.json` being readable.
