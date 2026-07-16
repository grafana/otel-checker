---
id: dotnet.coreclr-enable-profiling.value-mismatch
title: 'CORECLR_ENABLE_PROFILING is not set to 1'
severity: error
---

## Why this matters

`CORECLR_ENABLE_PROFILING` is the main switch for the .NET CLR profiling
API — the mechanism the OpenTelemetry .NET auto-instrumentation uses to
attach to your process, hook method entry/exit, and emit spans and metrics
without you touching the source. The Common Language Runtime (CLR) only reads the four `CORECLR_*`
profiler variables at process startup when this switch is set to `1`.

If the variable is missing, `0`, or any other value, the CLR skips profiler
attach entirely: none of `CORECLR_PROFILER`, `CORECLR_PROFILER_PATH`, or
`OTEL_DOTNET_AUTO_HOME` are consulted, the OpenTelemetry native profiler
never loads, and no telemetry is emitted — even though everything else
looks correctly configured.

## How to fix

Set `CORECLR_ENABLE_PROFILING` to exactly `1` in the environment that runs
the .NET process (not the shell that builds it).

- Local shell (Linux / macOS):

  ```bash
  export CORECLR_ENABLE_PROFILING=1
  ```

- PowerShell (Windows):

  ```powershell
  $env:CORECLR_ENABLE_PROFILING="1"
  ```

- Dockerfile:

  ```dockerfile
  ENV CORECLR_ENABLE_PROFILING=1
  ```

- Kubernetes:

  ```yaml
  env:
    - name: CORECLR_ENABLE_PROFILING
      value: "1"
  ```

Verify the value the process actually sees:

```bash
printenv CORECLR_ENABLE_PROFILING
```

Or in PowerShell:

```powershell
Write-Output $env:CORECLR_ENABLE_PROFILING
```

## Example

Complete set of env vars the OpenTelemetry .NET auto-instrumentation
requires (Linux / macOS):

```bash
export CORECLR_ENABLE_PROFILING=1
export CORECLR_PROFILER='{918728DD-259F-4A6A-AC2B-B85E1B658318}'
export CORECLR_PROFILER_PATH=/opt/opentelemetry/linux-x64/OpenTelemetry.AutoInstrumentation.Native.so
export OTEL_DOTNET_AUTO_HOME=/opt/opentelemetry
```

Same set in PowerShell (Windows):

```powershell
$env:CORECLR_ENABLE_PROFILING="1"
$env:CORECLR_PROFILER="{918728DD-259F-4A6A-AC2B-B85E1B658318}"
$env:CORECLR_PROFILER_PATH="$env:OTEL_DOTNET_AUTO_HOME\win-x64\OpenTelemetry.AutoInstrumentation.Native.dll"
$env:OTEL_DOTNET_AUTO_HOME="$env:PROGRAMFILES\OpenTelemetry\.NET AutoInstrumentation"
```

## Related

- [OpenTelemetry .NET auto-instrumentation configuration](https://opentelemetry.io/docs/zero-code/dotnet/configuration/)
- `dotnet.coreclr-profiler.value-mismatch` — companion check for the
  profiler GUID.
- `dotnet.coreclr-profiler-path.unset` — companion check for the native
  library path.
- `dotnet.otel-dotnet-auto-home.unset` — companion check for the
  auto-instrumentation home directory.
