---
id: go.go-mod.unreadable
title: 'go.mod could not be read'
severity: error
---

## Why this matters

`otel-checker` reads `go.mod` in the current directory to build the list
of your project's dependencies, then cross-references each one against
the OpenTelemetry Go contrib catalog. If the file can't be read — no
`go.mod` in the working directory, permission denied, or the file is
corrupt — the supported-libraries check has nothing to inspect and is
skipped, which is the main output of `check sdk --language=go`.

Common causes:

- Running the checker from the wrong directory (e.g. the repo root of a
  monorepo, while the Go module lives in a sub-directory).
- The Go module hasn't been initialized yet — `go.mod` truly does not
  exist.
- Permission bits prevent the invoking user from reading the file.

## How to fix

1. Verify `go.mod` exists and is readable from the current directory:

   ```bash
   ls -l go.mod
   ```

2. If your module lives in a sub-directory (a monorepo layout, a
   `cmd/` package with its own module, etc.), `cd` into that directory
   before running the checker:

   ```bash
   cd services/api
   otel-checker check sdk --language=go
   ```

3. If the module hasn't been initialized yet, create one:

   ```bash
   go mod init github.com/you/your-service
   go mod tidy
   ```

4. If the file exists but the read failed, fix permissions:

   ```bash
   chmod +r go.mod
   ```

## Example

Typical layout:

```text
my-service/
├── go.mod
├── go.sum
└── main.go
```

Invocation from `my-service/`:

```bash
otel-checker check sdk --language=go
```

## Related

- `go.supported-libs.fetch-failed` — related failure when the
  supported-libraries YAML can't be loaded.
