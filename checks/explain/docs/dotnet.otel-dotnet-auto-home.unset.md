---
id: dotnet.otel-dotnet-auto-home.unset
title: 'OTEL_DOTNET_AUTO_HOME is unset'
severity: error
---

## Why this matters

`OTEL_DOTNET_AUTO_HOME` points at the root of the OpenTelemetry .NET
auto-instrumentation distribution — the directory that contains the
managed helper assemblies, the per-platform native profiler libraries,
and the startup hook DLL that the CLR loads before your app's own code
runs. The auto-instrumentation loader reads this variable at startup to
find every file it needs to attach; when it's unset, the loader has no
anchor point, and the profiler either fails to initialize or silently
falls back to a no-op path.

## How to fix

Install the OpenTelemetry .NET auto-instrumentation distribution and point
`OTEL_DOTNET_AUTO_HOME` at its root directory.

- Follow the [OpenTelemetry .NET auto-instrumentation installation
  guide](https://opentelemetry.io/docs/zero-code/dotnet/) to download
  and extract the distribution. Common install locations:

  - Linux: `/opt/opentelemetry`
  - macOS: `/opt/opentelemetry` or `$HOME/otel-dotnet-auto`
  - Windows: `"%PROGRAMFILES%\OpenTelemetry\.NET AutoInstrumentation\"`

- Set the variable in the environment that runs the .NET process:

  ```bash
  export OTEL_DOTNET_AUTO_HOME=/opt/opentelemetry
  ```

  Dockerfile:

  ```dockerfile
  ENV OTEL_DOTNET_AUTO_HOME=/opt/opentelemetry
  ```

  Kubernetes:

  ```yaml
  env:
    - name: OTEL_DOTNET_AUTO_HOME
      value: /opt/opentelemetry
  ```

Verify the directory exists and contains the expected sub-folders
(`net`, `linux-x64` / `osx-arm64` / `win-x64`, etc.):

```bash
ls "$OTEL_DOTNET_AUTO_HOME"
```

## Example

Complete set of env vars the OpenTelemetry .NET auto-instrumentation
requires:

```bash
export OTEL_DOTNET_AUTO_HOME=/opt/opentelemetry
export CORECLR_ENABLE_PROFILING=1
export CORECLR_PROFILER='{918728DD-259F-4A6A-AC2B-B85E1B658318}'
export CORECLR_PROFILER_PATH=$OTEL_DOTNET_AUTO_HOME/linux-x64/OpenTelemetry.AutoInstrumentation.Native.so
```

## Related

- [OpenTelemetry .NET auto-instrumentation installation](https://opentelemetry.io/docs/zero-code/dotnet/)
- `dotnet.coreclr-profiler-path.unset` — companion check for the native
  library path, which is normally built from this home directory.
- `dotnet.coreclr-enable-profiling.value-mismatch` — companion check for
  the profiling main switch.
- `dotnet.coreclr-profiler.value-mismatch` — companion check for the
  profiler GUID.
