---
id: dotnet.coreclr-profiler-path.unset
title: 'CORECLR_PROFILER_PATH is unset'
severity: error
---

## Why this matters

`CORECLR_PROFILER_PATH` tells the .NET Common Language Runtime (CLR)
where on disk to find the OpenTelemetry native profiler shared library
that corresponds to the GUID you selected via `CORECLR_PROFILER`. The
CLR loads that library at process startup and hands it profiling
callbacks — without a path, the CLR has nothing to load and the profiler
is silently skipped, even if `CORECLR_ENABLE_PROFILING=1` and
`CORECLR_PROFILER` are set correctly.

The correct value depends on your OS and CPU architecture: the
OpenTelemetry .NET distribution ships one native library per platform.

## How to fix

Point `CORECLR_PROFILER_PATH` at the platform-appropriate native library
inside your OpenTelemetry .NET auto-instrumentation distribution
(`$OTEL_DOTNET_AUTO_HOME/<platform>/OpenTelemetry.AutoInstrumentation.Native.<ext>`).

- Linux x64:

  ```bash
  export CORECLR_PROFILER_PATH=$OTEL_DOTNET_AUTO_HOME/linux-x64/OpenTelemetry.AutoInstrumentation.Native.so
  ```

- Linux ARM64:

  ```bash
  export CORECLR_PROFILER_PATH=$OTEL_DOTNET_AUTO_HOME/linux-arm64/OpenTelemetry.AutoInstrumentation.Native.so
  ```

- macOS (Apple silicon):

  ```bash
  export CORECLR_PROFILER_PATH=$OTEL_DOTNET_AUTO_HOME/osx-arm64/OpenTelemetry.AutoInstrumentation.Native.dylib
  ```

- Windows x64 (PowerShell):

  ```powershell
  $env:CORECLR_PROFILER_PATH = "$env:OTEL_DOTNET_AUTO_HOME\win-x64\OpenTelemetry.AutoInstrumentation.Native.dll"
  ```

  On Windows only x64 is supported by the OpenTelemetry .NET distribution —
  there is no x86 or ARM64 build.

Verify the file exists and is readable by the process user:

```bash
ls -l "$CORECLR_PROFILER_PATH"
```

## Example

Complete set of env vars the OpenTelemetry .NET auto-instrumentation
requires (Linux x64):

```bash
export OTEL_DOTNET_AUTO_HOME=/opt/opentelemetry
export CORECLR_ENABLE_PROFILING=1
export CORECLR_PROFILER='{918728DD-259F-4A6A-AC2B-B85E1B658318}'
export CORECLR_PROFILER_PATH=$OTEL_DOTNET_AUTO_HOME/linux-x64/OpenTelemetry.AutoInstrumentation.Native.so
```

## Related

- [OpenTelemetry .NET auto-instrumentation configuration](https://opentelemetry.io/docs/zero-code/dotnet/configuration/)
- `dotnet.coreclr-enable-profiling.value-mismatch` — companion check for
  the profiling main switch.
- `dotnet.coreclr-profiler.value-mismatch` — companion check for the
  profiler GUID.
- `dotnet.otel-dotnet-auto-home.unset` — companion check for the root
  directory this path is usually built from.
