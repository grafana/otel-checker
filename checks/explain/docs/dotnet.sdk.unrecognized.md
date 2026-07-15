---
id: dotnet.sdk.unrecognized
title: 'The .NET SDK reported in the project file is not recognized'
severity: error
---

## Why this matters

The top of every SDK-style .NET project file declares the SDK it uses:
`<Project Sdk="Microsoft.NET.Sdk">`, `<Project Sdk="Microsoft.NET.Sdk.Web">`,
and so on. `otel-checker` uses that SDK identifier to pick which
*implicit* framework packages your project has (ASP.NET Core hosting on
`Microsoft.NET.Sdk.Web`, `System.Net.Http` on the base SDK, etc.), so it
can report the OpenTelemetry instrumentations that apply even before you
add any explicit NuGet references.

This error means the SDK attribute in your `.csproj` is not one the
checker knows about yet. The most common cause is a specialized project
type (worker service, Razor library, Azure Functions, Blazor) whose
implicit-packages mapping hasn't been added, or a typo / non-standard
casing in the `Sdk` attribute.

## How to fix

1. Open the `.csproj` and confirm the `Sdk` attribute is one of the
   standard values:

   ```xml
   <Project Sdk="Microsoft.NET.Sdk">          <!-- library, console app -->
   <Project Sdk="Microsoft.NET.Sdk.Web">      <!-- ASP.NET Core -->
   <Project Sdk="Microsoft.NET.Sdk.Worker">   <!-- background worker -->
   <Project Sdk="Microsoft.NET.Sdk.Razor">    <!-- Razor class library -->
   ```

2. Fix any typo or case mismatch. The value is case-sensitive.

3. If your project uses a legitimately different SDK (Blazor, MAUI,
   Azure Functions, etc.), the checker can still analyze your explicit
   NuGet dependencies — the "implicit packages" branch is what's
   skipped. Open an issue in the `otel-checker` repo mentioning the SDK
   so a mapping can be added.

## Example

Standard ASP.NET Core project header:

```xml
<Project Sdk="Microsoft.NET.Sdk.Web">
  <PropertyGroup>
    <TargetFramework>net8.0</TargetFramework>
  </PropertyGroup>
  <!-- ... -->
</Project>
```

## Related

- [.NET SDK-style project reference](https://learn.microsoft.com/dotnet/core/project-sdk/overview)
- `dotnet.sdk.no-implicit-packages` — related warning when the SDK is
  recognized but has no implicit packages mapped.
