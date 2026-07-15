---
id: dotnet.dependencies.unreadable
title: 'Could not read .NET project dependencies'
severity: error
---

## Why this matters

`otel-checker` runs `dotnet list package --format json --include-transitive`
to get the full dependency graph — including transitive NuGet references
— so it can match every package against the list of known OpenTelemetry
instrumentations. This error fires when that command fails: either
`dotnet list package` returned a non-zero exit code, or its output wasn't
parseable JSON.

Common causes:

- **`dotnet` isn't installed** or isn't on the current `PATH`, so the
  `dotnet list package` invocation itself errors before it can read the
  project.
- **No restore yet**: `dotnet list package` requires `project.assets.json`,
  which `dotnet restore` (or an implicit restore on `build`/`run`)
  produces. On a freshly-cloned repo the command errors until NuGet
  packages have been restored.
- **Unsupported (old) SDK**: `--format json` is available in every
  currently-supported .NET SDK, so this is only relevant if you're
  running an SDK that has already reached end-of-life.
- **Broken project state**: a partial build, a corrupt
  `obj/project.assets.json`, or a `.csproj` that references missing
  targets can make `dotnet list package` fail.

Without a dependency list, the checker can't identify which of your
packages have OpenTelemetry instrumentation available.

## How to fix

1. Restore the project first:

   ```bash
   dotnet restore
   ```

2. Confirm the CLI is happy on its own:

   ```bash
   dotnet list package --format json --include-transitive
   ```

   If it errors, the error message points at the underlying problem
   (missing target, broken NuGet source, unauthenticated feed, etc.).

3. Confirm you're on a still-supported .NET SDK:

   ```bash
   dotnet --version
   ```

   If it's older than a currently-supported release, upgrade — see
   `dotnet.version.too-old`.

4. If a corrupt build output is the cause, clean and re-restore:

   ```bash
   dotnet clean
   dotnet restore
   ```

## Example

Clean restore + check flow:

```bash
cd MyApp
dotnet restore
otel-checker check sdk --language=dotnet
```

## Related

- `dotnet.project.no-dependencies` — related check when the command
  succeeds but returns an empty list.
- `dotnet.version.too-old` — related check when the SDK is older than
  the currently-supported minimum.
