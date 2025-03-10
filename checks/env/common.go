package env

import (
	"fmt"
	"otel-checker/checks/utils"
)

// Common environment variables used across the project
var (
	// OpenTelemetry common variables
	OtelServiceName = EnvVar{
		Name:        "OTEL_SERVICE_NAME",
		Recommended: true,
		Message:     "It's recommended the environment variable OTEL_SERVICE_NAME to be set to your service name, for easier identification",
	}

	OtelMetricsExporter = exporterEnvVar("OTEL_METRICS_EXPORTER", "Metrics")
	OtelTracesExporter  = exporterEnvVar("OTEL_TRACES_EXPORTER", "Traces")
	OtelLogsExporter    = exporterEnvVar("OTEL_LOGS_EXPORTER", "Logs")
)

func CheckCommonEnvVars(r *utils.ComponentReporter, language string) {
	// Check common OpenTelemetry variables
	commonVars := []EnvVar{
		OtelServiceName,
		OtelMetricsExporter,
		OtelTracesExporter,
		OtelLogsExporter,
	}

	CheckEnvVars(r, language, commonVars...)
}

func exporterEnvVar(key string, name string) EnvVar {
	return EnvVar{
		Name:         key,
		Required:     false,
		DefaultValue: "otlp",
		Validator: func(value string, language string, reporter *utils.ComponentReporter) {
			if value == "none" {
				reporter.AddError(fmt.Sprintf("The value of %s cannot be 'none'. Change the value to 'otlp' or leave it unset", key))
			} else {
				if value == "" {
					reporter.AddSuccessfulCheck(fmt.Sprintf("%s is unset, with a default value of 'otlp'", key))
				} else {
					reporter.AddSuccessfulCheck(fmt.Sprintf("The value of %s is set to '%s' (default value)", key, value))
				}
			}
		},
		Description: name + " exporter configuration",
	}
}
