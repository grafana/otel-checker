package env

import (
	"fmt"
	"strings"

	"github.com/grafana/otel-checker/checks/utils"
)

// Common environment variables used across the project
var (
	OtelServiceName = EnvVar{
		Name:        "OTEL_SERVICE_NAME",
		Recommended: true,
		Message:     "It's recommended the environment variable OTEL_SERVICE_NAME to be set to your service name, for easier identification",
	}

	OtelResourceAttributes = EnvVar{
		Name:        "OTEL_RESOURCE_ATTRIBUTES",
		Recommended: true,
		Message:     "It's recommended to set OTEL_RESOURCE_ATTRIBUTES with key-value pairs for resource attributes (e.g., \"key1=value1,key2=value2\")",
	}

	OtelMetricsExporter = exporterEnvVar("OTEL_METRICS_EXPORTER", "Metrics")
	OtelTracesExporter  = exporterEnvVar("OTEL_TRACES_EXPORTER", "Traces")
	OtelLogsExporter    = exporterEnvVar("OTEL_LOGS_EXPORTER", "Logs")
)

// ResourceAttribute represents a recommended OpenTelemetry resource attribute
type ResourceAttribute struct {
	Name         string
	Description  string
	ExampleValue string
}

// ParseResourceAttributes parses the OTEL_RESOURCE_ATTRIBUTES environment variable
// Format: "key1=value1,key2=value2"
func ParseResourceAttributes() map[string]string {
	attributes := make(map[string]string)

	// Get resource attributes from environment variable
	resourceAttrsEnv := GetValue(OtelResourceAttributes)
	if resourceAttrsEnv != "" {
		// Split by comma to get key-value pairs
		pairs := strings.Split(resourceAttrsEnv, ",")
		for _, pair := range pairs {
			// Split by = to get key and value
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) == 2 {
				key := strings.TrimSpace(kv[0])
				value := strings.TrimSpace(kv[1])
				attributes[key] = value
			}
		}
	}

	return attributes
}

// CheckResourceAttributes checks if recommended OpenTelemetry resource attributes are configured
func CheckResourceAttributes(reporter *utils.ComponentReporter) {
	// Define the recommended resource attributes based on Grafana documentation
	// See: https://grafana.com/docs/grafana-cloud/monitor-applications/application-observability/instrument/resource-attributes/
	recommendedAttributes := []ResourceAttribute{
		{
			Name:         "service.namespace",
			Description:  "An optional namespace for service.name",
			ExampleValue: "shop",
		},
		{
			Name:         "deployment.environment.name",
			Description:  "Name of the deployment environment (e.g. 'staging' or 'production')",
			ExampleValue: "production",
		},
		{
			Name:         "service.instance.id",
			Description:  "The unique instance, e.g. the pod name",
			ExampleValue: "checkout-123",
		},
		{
			Name:         "service.version",
			Description:  "The application version, to see if a new version has introduced a bug",
			ExampleValue: "1.2",
		},
	}

	attributes := ParseResourceAttributes()

	for _, attr := range recommendedAttributes {
		value, exists := attributes[attr.Name]

		if exists && value != "" {
			reporter.AddSuccessfulCheck(
				fmt.Sprintf("Resource attribute %s is set to '%s'", attr.Name, value))
		} else {
			reporter.AddWarningWithExplain("env.resource-attributes.missing",
				fmt.Sprintf("Set OTEL_RESOURCE_ATTRIBUTES=\"%s=%s\": %s", attr.Name, attr.ExampleValue, attr.Description))
		}
	}

	// Special handling for service.name which can be set via OTEL_SERVICE_NAME or as a resource attribute
	// Note: According to OpenTelemetry spec, if both are set, OTEL_SERVICE_NAME takes precedence
	serviceNameValue, serviceNameExists := attributes["service.name"]
	otelServiceNameValue := GetValue(OtelServiceName)

	if otelServiceNameValue != "" {
		reporter.AddSuccessfulCheck(fmt.Sprintf("Service name is set via OTEL_SERVICE_NAME to '%s'", otelServiceNameValue))
	} else if serviceNameExists && serviceNameValue != "" {
		reporter.AddSuccessfulCheck(fmt.Sprintf("Service name is set via OTEL_RESOURCE_ATTRIBUTES to '%s'", serviceNameValue))
	} else {
		reporter.AddWarningWithExplain("env.otel-service-name.unset",
			"Set OTEL_SERVICE_NAME=\"checkout\": The application name")
	}
}

func CheckCommon(r *utils.ComponentReporter, language string) {
	CheckExporterEnvVars(r, language)

	CheckResourceAttributes(r)
}

func CheckExporterEnvVars(r *utils.ComponentReporter, language string) {
	CheckEnvVars(r, language,
		OtelMetricsExporter,
		OtelTracesExporter,
		OtelLogsExporter)
}

func exporterEnvVar(key string, name string) EnvVar {
	return EnvVar{
		Name:         key,
		Required:     false,
		DefaultValue: "otlp",
		Validator: func(value string, language string, reporter *utils.ComponentReporter) {
			switch value {
			case "":
				reporter.AddSuccessfulCheck(fmt.Sprintf("%s is unset, with a default value of 'otlp'", key))
			case "otlp":
				reporter.AddSuccessfulCheck(fmt.Sprintf("The value of %s is set to 'otlp' (default value)", key))
			case "console":
				reporter.AddSuccessfulCheck(fmt.Sprintf("The value of %s is set to 'console'", key))
			case "none":
				reporter.AddErrorWithExplain("env.exporter.disabled",
					fmt.Sprintf("The value of %s cannot be 'none'. Change the value to 'otlp' or 'console', or leave it unset", key))
			default:
				reporter.AddErrorWithExplain("env.exporter.invalid-value",
					fmt.Sprintf("The value of %s must be 'otlp' or 'console' (or unset). Got '%s'", key, value))
			}
		},
		Description: name + " exporter configuration",
	}
}
