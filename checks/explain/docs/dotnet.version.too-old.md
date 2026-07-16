---
id: dotnet.version.too-old
title: '.NET version is below the recommended minimum'
severity: error
---

## Why this matters

`otel-checker` reads the local .NET SDK version via `dotnet --version` and
expects at least .NET 8. That's the oldest version the OpenTelemetry .NET
distribution officially supports: earlier majors (.NET 7 and older) are
past end-of-life
([endoflife.date/dotnet](https://endoflife.date/dotnet)) and the
OpenTelemetry .NET auto-instrumentation and manual SDK packages no longer
guarantee compatibility with them.

Installing the OpenTelemetry NuGet packages on .NET 7 or older may
succeed at restore time but throw at process startup with runtime errors
that don't obviously point at OpenTelemetry, or silently emit degraded
telemetry.

## How to fix

Upgrade the .NET SDK to a supported release — .NET 8 (LTS) or newer:

- Install a current SDK from
  [dotnet.microsoft.com/download](https://dotnet.microsoft.com/download)
  or your distro's package manager (`apt`, `brew`, `winget`, etc.).
- Bump the `TargetFramework` in every `.csproj` you want the checker to
  cover:

  ```xml
  <TargetFramework>net8.0</TargetFramework>
  ```

- In Docker, use a supported SDK base image:

  ```dockerfile
  FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
  ```

- On CI, pin an explicit SDK version with `actions/setup-dotnet`
  (or the equivalent):

  ```yaml
  - uses: actions/setup-dotnet@v5
    with:
      dotnet-version: '8.0.x'
  ```

Verify the version the current environment sees:

```bash
dotnet --version
# 8.0.x
```

## Example

`global.json` pinning the SDK for a repo:

```json
{
  "sdk": {
    "version": "8.0.100",
    "rollForward": "latestFeature"
  }
}
```

## Related

- [.NET release schedule](https://dotnet.microsoft.com/platform/support/policy/dotnet-core)
- [endoflife.date/dotnet](https://endoflife.date/dotnet) — current status
  of every .NET major.
- `dotnet.version.unknown` — related failure when `dotnet` isn't
  callable at all.
