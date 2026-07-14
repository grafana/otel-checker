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

- **No restore yet**: `dotnet list package` requires `project.assets.json`,
  which `dotnet restore` (or an implicit restore on `build`/`run`)
  produces. On a freshly-cloned repo the command errors until you
  restore.
- **Old SDK**: `--format json` was added in .NET SDK 7.0.200. Older SDKs
  emit human-readable text that the checker can't parse.
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

3. Check your SDK is new enough for `--format json`:

   ```bash
   dotnet --version
   # 7.0.200 or later
   ```

   If it's older, upgrade to a supported SDK — see
   `dotnet.version.too-old`.

4. If a corrupt `obj/` is the cause, wipe it and re-restore:

   ```bash
   rm -rf obj bin
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
- `dotnet.version.too-old` — related check when the SDK predates
  `--format json` support.
