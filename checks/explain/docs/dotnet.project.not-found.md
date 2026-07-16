---
id: dotnet.project.not-found
title: 'No .NET project file was found'
severity: error
---

## Why this matters

`otel-checker` looks for exactly one `.csproj` file in the current
directory to identify the .NET project it should check. That project
file drives every downstream check — the .NET SDK type (used to pick
implicit instrumentations), the NuGet dependency list (used to match
against known OpenTelemetry instrumentations), and the target framework.

This error fires when either:

- No `.csproj` file exists in the current directory. The checker doesn't
  recurse, so a solution laid out with each project in its own
  sub-directory (the common .NET layout) will hit this if run from the
  solution root.
- The directory exists but the checker can't read it (permission denied,
  the path doesn't exist, or it isn't a directory).
- The `.csproj` file exists but is unreadable or malformed XML — the
  loader wraps that under the same ID.

Without a project to inspect, the remaining .NET checks can't proceed.

## How to fix

1. Confirm you're running the checker from the project directory:

   ```bash
   ls *.csproj
   ```

   If it's empty, `cd` into the sub-directory that contains the
   `.csproj`.

2. If the project file exists but the load failed, the error message
   quotes the reason from the XML parser. Common cases:

   - **Malformed XML**: fix the syntax and re-run.
   - **Wrong root element**: the checker expects an SDK-style project
     (`<Project Sdk="…">`), not a legacy .NET Framework project.

3. If you have multiple `.csproj` files in the same directory (rare, but
   supported by MSBuild), the checker can't disambiguate. Move each
   project into its own sub-directory, or run `otel-checker` inside the
   sub-directory of the one you want checked.

## Example

Typical layout for a single-project repo:

```text
MyApp/
├── MyApp.csproj
├── Program.cs
└── ...
```

Run from `MyApp/`:

```bash
cd MyApp
otel-checker check sdk --language=dotnet
```

## Related

- `dotnet.dependencies.unreadable` — related failure when the project
  loads but its dependencies can't be listed.
- `dotnet.sdk.unrecognized` — related failure when the project loads
  but the SDK attribute isn't one the checker knows about.
