---
id: dotnet.version.empty
title: '.NET version output was empty'
severity: error
---

## Why this matters

`otel-checker` runs `dotnet --version` to verify the .NET SDK meets the
minimum supported version. This error fires when the command *ran
successfully* — no error exit code — but produced empty output. That's
unusual and almost always signals a broken or incomplete .NET
installation: the binary exists on `PATH` (otherwise you'd hit
`dotnet.version.unknown`), but the SDK it's meant to invoke isn't
reachable.

Common causes:

- Only the .NET *runtime* is installed, not the SDK. `dotnet` still
  exists but has no version to report.
- A `global.json` in the working directory (or a parent) pins an SDK
  version that isn't installed. `dotnet` refuses to select any SDK and
  prints nothing.
- The installation was interrupted and left `dotnet` present but the
  `sdk` directory empty.

## How to fix

1. Confirm what `dotnet` is actually reporting when run directly:

   ```bash
   dotnet --version
   dotnet --list-sdks
   ```

   If `--list-sdks` shows nothing, no SDK is installed for this runtime.

2. If a `global.json` is pinning an unavailable SDK, either install the
   pinned version or relax the pin:

   ```bash
   grep -R "version" global.json
   dotnet --list-sdks
   ```

   Adjust `global.json` to a version you have, or use `rollForward` to
   allow newer feature bands:

   ```json
   {
     "sdk": {
       "version": "8.0.100",
       "rollForward": "latestFeature"
     }
   }
   ```

3. Install a supported SDK from
   [dotnet.microsoft.com/download](https://dotnet.microsoft.com/download).

## Example

Healthy output after installing .NET 8:

```bash
$ dotnet --version
8.0.401
$ dotnet --list-sdks
8.0.401 [/usr/share/dotnet/sdk]
```

## Related

- `dotnet.version.unknown` — related failure when `dotnet` itself is
  missing or fails.
- `dotnet.version.too-old` — related check when an SDK is present but
  older than the supported minimum.
