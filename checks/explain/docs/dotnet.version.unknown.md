---
id: dotnet.version.unknown
title: 'Could not determine the installed .NET version'
severity: error
---

## Why this matters

`otel-checker` runs `dotnet --version` and parses the output to verify
the .NET SDK meets the minimum supported version. This error means either
the `dotnet` CLI isn't on the current `PATH`, the command exited non-zero,
or its output couldn't be parsed as a version number. Without a version
to check, the tool can't tell you whether your runtime is compatible with
the OpenTelemetry .NET packages you're about to install — a version
mismatch is one of the most common reasons instrumentation appears to
install cleanly but never emits telemetry.

## How to fix

1. Verify the .NET SDK is installed and callable from the same shell that
   runs `otel-checker`:

   ```bash
   which dotnet   # `where dotnet` on Windows
   dotnet --version
   ```

2. If `dotnet` isn't found, install it. Grab an installer from
   [dotnet.microsoft.com/download](https://dotnet.microsoft.com/download)
   or use your package manager (`apt install dotnet-sdk-8.0`,
   `brew install --cask dotnet-sdk`, `winget install Microsoft.DotNet.SDK.8`).

3. If .NET is installed but only in another user's environment (common
   under CI-managed installs), make sure the same shell profile is
   sourced. On GitHub Actions use `actions/setup-dotnet` before invoking
   the checker.

## Example

GitHub Actions workflow that makes .NET visible to `otel-checker`:

```yaml
- uses: actions/setup-dotnet@v5
  with:
    dotnet-version: '8.0.x'
- run: otel-checker check sdk --language=dotnet
```

## Related

- `dotnet.version.too-old` — related check that fires when the SDK is
  found but is older than the supported minimum.
- `dotnet.version.empty` — related check when `dotnet --version` runs
  but returns an empty string.
- [.NET downloads](https://dotnet.microsoft.com/download)
