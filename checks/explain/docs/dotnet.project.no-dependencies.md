---
id: dotnet.project.no-dependencies
title: 'The .NET project has no dependencies declared'
severity: error
---

## Why this matters

After running `dotnet list package --format json --include-transitive`,
`otel-checker` found zero projects with dependencies attached. This is
almost always one of two situations:

- The current directory contains a project file but the packages haven't
  been restored yet — `dotnet list package` reports an empty result set
  until `dotnet restore` has run at least once.
- The command ran against the wrong directory or a stale
  `project.assets.json`, so it looked at a project that genuinely has no
  NuGet references. In a real service that's a red flag by itself.

Without a dependency list, the checker can't tell you which of your
NuGet packages have OpenTelemetry instrumentation available, which is
the main output of `check sdk --language=dotnet`.

## How to fix

1. Restore packages first, so the dependency graph is populated:

   ```bash
   dotnet restore
   dotnet list package --include-transitive
   ```

   The list command should print at least the implicit framework
   packages. If it still prints nothing, the current directory doesn't
   contain a runnable project.

2. Confirm you're running the checker in the project directory (see
   `dotnet.project.not-found` for layout notes), not the repo root of a
   multi-project solution.

3. If the project really has no dependencies (a brand-new empty class
   library, for example), install at least the base packages your app
   actually uses, then re-run the checker.

## Example

Clean restore + check flow:

```bash
cd MyApp
dotnet restore
otel-checker check sdk --language=dotnet
```

## Related

- `dotnet.dependencies.unreadable` — related failure when the CLI can't
  produce a dependency list at all.
- `dotnet.project.not-found` — related failure when there's no project
  to check.
