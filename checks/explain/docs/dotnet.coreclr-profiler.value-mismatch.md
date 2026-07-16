---
id: dotnet.coreclr-profiler.value-mismatch
title: 'CORECLR_PROFILER is not set to the OpenTelemetry profiler GUID'
severity: error
---

## Why this matters

`CORECLR_PROFILER` selects *which* profiler the .NET CLR loads when
profiling is enabled. It's a single GUID — the CLR only attaches one
profiler per process, so this variable is what determines whether the
OpenTelemetry native profiler or some other vendor's profiler is in
control.

For the OpenTelemetry .NET auto-instrumentation the value must be exactly
`{918728DD-259F-4A6A-AC2B-B85E1B658318}` (braces and hex digits included).
If it's set to another GUID — often left over from a previous APM,
profiler, or debugger — the OpenTelemetry profiler is never loaded, and
no traces or metrics reach Grafana Cloud even though everything else looks
correctly configured.

## How to fix

Set `CORECLR_PROFILER` to the OpenTelemetry GUID verbatim.

- Local shell (single-quote to keep the braces literal):

  ```bash
  export CORECLR_PROFILER='{918728DD-259F-4A6A-AC2B-B85E1B658318}'
  ```

- Dockerfile:

  ```dockerfile
  ENV CORECLR_PROFILER={918728DD-259F-4A6A-AC2B-B85E1B658318}
  ```

- Kubernetes:

  ```yaml
  env:
    - name: CORECLR_PROFILER
      value: "{918728DD-259F-4A6A-AC2B-B85E1B658318}"
  ```

If the GUID currently belongs to another APM's profiler, understand that
the CLR only accepts *one* profiler at a time. Choose which profiler you
want attached and unset the other vendor's configuration — you cannot
run both simultaneously.

Verify:

```bash
printenv CORECLR_PROFILER
# {918728DD-259F-4A6A-AC2B-B85E1B658318}
```

## Example

Complete set of env vars the OpenTelemetry .NET auto-instrumentation
requires:

```bash
export CORECLR_ENABLE_PROFILING=1
export CORECLR_PROFILER='{918728DD-259F-4A6A-AC2B-B85E1B658318}'
export CORECLR_PROFILER_PATH=/opt/opentelemetry/linux-x64/OpenTelemetry.AutoInstrumentation.Native.so
export OTEL_DOTNET_AUTO_HOME=/opt/opentelemetry
```

## Related

- [OpenTelemetry .NET auto-instrumentation configuration](https://opentelemetry.io/docs/zero-code/dotnet/configuration/)
- `dotnet.coreclr-enable-profiling.value-mismatch` — companion check for
  the profiling main switch.
- `dotnet.coreclr-profiler-path.unset` — companion check for the native
  library path the CLR loads under this GUID.
