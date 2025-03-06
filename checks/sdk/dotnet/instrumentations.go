package dotnet

import (
	"fmt"
)

func ReadAvailableInstrumentations() map[string]string {
	return map[string]string{
		"Microsoft.AspNetCore.Hosting":  "OpenTelemetry.Instrumentation.AspNetCore",
		"Microsoft.EntityFrameworkCore": "OpenTelemetry.Instrumentation.EntityFrameworkCore",
		"Microsoft.Data.SqlClient":      "OpenTelemetry.Instrumentation.SqlClient",
		"StackExchange.Redis":           "OpenTelemetry.Instrumentation.StackExchangeRedis",
		"System.Net.Http":               "OpenTelemetry.Instrumentation.Http",
	}
}

func ImplicitPackagesForSdk(sdk string) ([]string, error) {
	base := []string{
		"System.Net.Http",
	}

	switch sdk {
	case "Microsoft.NET.Sdk.Web":
		// TODO what other libs do these SDKs bundle?
		// case "Microsoft.NET.Sdk.Worker":
		// case "Microsoft.NET.Sdk.Razor":
		// case "Microsoft.NET.Sdk.Function":
		return append(base, "Microsoft.AspNetCore.Hosting"), nil
	case "Microsoft.NET.Sdk":
		return base, nil
	}

	return nil, fmt.Errorf("Unrecognized SDK: %s", sdk)
}
