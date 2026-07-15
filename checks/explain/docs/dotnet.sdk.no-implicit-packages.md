---
id: dotnet.sdk.no-implicit-packages
title: 'No implicit instrumentation packages found for this SDK'
severity: warning
---

## Why this matters

Some .NET project SDKs bundle framework packages implicitly — you don't
have to add a `<PackageReference>` for them, they're pulled in by the
SDK itself. `otel-checker` maps each supported SDK to that implicit list
so it can tell you which OpenTelemetry instrumentations apply based on
the SDK alone, before it even looks at your explicit NuGet dependencies.

This warning means the checker recognized your project's SDK but
currently has an empty implicit-packages list for it. That's not
necessarily wrong — some SDKs simply don't bundle instrumentable
packages — but it means the "implicit instrumentation" branch of the
report will be empty. Explicit NuGet references still get checked; only
the framework-implicit part is skipped.

## How to fix

Nothing to *fix* on your side — the warning is a hint that the checker's
implicit-package mapping for this SDK isn't populated yet. Two paths:

1. **If the SDK really has no notable implicit instrumentations**
   (e.g. a plain library that doesn't pull in HTTP or hosting): this
   warning is harmless. Ignore it, and rely on the explicit
   `<PackageReference>` matches in the report.

2. **If you believe the SDK *should* bundle instrumentable packages**
   (worker services and Razor libraries do bundle
   `Microsoft.Extensions.Hosting`, for example): open an issue on the
   [`grafana/otel-checker`](https://github.com/grafana/otel-checker/issues)
   repo pointing at the SDK. Adding the mapping is a small change to
   `checks/sdk/dotnet/instrumentations.go`.

## Example

Currently mapped SDKs (see `checks/sdk/dotnet/instrumentations.go`):

- `Microsoft.NET.Sdk` → `System.Net.Http`
- `Microsoft.NET.Sdk.Web` → `System.Net.Http`, `Microsoft.AspNetCore.Hosting`

Other SDKs are recognized but have no mapping yet.

## Related

- `dotnet.sdk.unrecognized` — related error when the SDK itself isn't
  recognized at all.
- [.NET SDK-style project reference](https://learn.microsoft.com/dotnet/core/project-sdk/overview)
