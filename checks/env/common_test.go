package env

import (
	"os"
	"reflect"
	"testing"

	"otel-checker/checks/utils"
)

func TestParseResourceAttributes(t *testing.T) {
	// Save original environment values and restore after test
	origValue := os.Getenv("OTEL_RESOURCE_ATTRIBUTES")
	defer os.Setenv("OTEL_RESOURCE_ATTRIBUTES", origValue)

	tests := []struct {
		name          string
		envValue      string
		expectedAttrs map[string]string
	}{
		{
			name:          "empty env var",
			envValue:      "",
			expectedAttrs: map[string]string{},
		},
		{
			name:          "single key-value pair",
			envValue:      "service.name=my-service",
			expectedAttrs: map[string]string{"service.name": "my-service"},
		},
		{
			name:     "multiple key-value pairs",
			envValue: "service.name=my-service,deployment.environment.name=production,service.version=1.0.0",
			expectedAttrs: map[string]string{
				"service.name":                "my-service",
				"deployment.environment.name": "production",
				"service.version":             "1.0.0",
			},
		},
		{
			name:     "pairs with spaces",
			envValue: "service.name = my-service , deployment.environment.name = production",
			expectedAttrs: map[string]string{
				"service.name":                "my-service",
				"deployment.environment.name": "production",
			},
		},
		{
			name:     "malformed pair is ignored",
			envValue: "service.name=my-service,malformed-no-equals,deployment.environment.name=production",
			expectedAttrs: map[string]string{
				"service.name":                "my-service",
				"deployment.environment.name": "production",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("OTEL_RESOURCE_ATTRIBUTES", tt.envValue)
			result := ParseResourceAttributes()

			if !reflect.DeepEqual(result, tt.expectedAttrs) {
				t.Errorf("ParseResourceAttributes() = %v, want %v", result, tt.expectedAttrs)
			}
		})
	}
}

func TestCheckResourceAttributes(t *testing.T) {
	tests := []utils.EnvVarTestCase{
		{
			Name: "all attributes present including service.name",
			EnvVars: map[string]string{
				"OTEL_SERVICE_NAME":        "my-service",
				"OTEL_RESOURCE_ATTRIBUTES": "service.namespace=my-namespace,deployment.environment.name=production,service.instance.id=instance-1,service.version=1.0.0",
			},
			Language: "test",
			ExpectedChecks: []string{
				"Resource Attributes: Service name is set via OTEL_SERVICE_NAME to 'my-service'",
				"Resource Attributes: Resource attribute service.namespace is set to 'my-namespace'",
				"Resource Attributes: Resource attribute deployment.environment.name is set to 'production'",
				"Resource Attributes: Resource attribute service.instance.id is set to 'instance-1'",
				"Resource Attributes: Resource attribute service.version is set to '1.0.0'",
			},
		},
		{
			Name: "service.name in resource attributes, no OTEL_SERVICE_NAME",
			EnvVars: map[string]string{
				"OTEL_RESOURCE_ATTRIBUTES": "service.name=my-service,deployment.environment.name=production",
			},
			Language: "test",
			ExpectedChecks: []string{
				"Resource Attributes: Service name is set via OTEL_RESOURCE_ATTRIBUTES to 'my-service'",
				"Resource Attributes: Resource attribute deployment.environment.name is set to 'production'",
			},
			ExpectedWarnings: []string{
				"Resource Attributes: Set OTEL_RESOURCE_ATTRIBUTES=\"service.namespace=shop\": An optional namespace for service.name",
				"Resource Attributes: Set OTEL_RESOURCE_ATTRIBUTES=\"service.instance.id=checkout-123\": The unique instance, e.g. the pod name",
				"Resource Attributes: Set OTEL_RESOURCE_ATTRIBUTES=\"service.version=1.2\": The application version, to see if a new version has introduced a bug",
			},
		},
		{
			Name: "no service.name in resource attributes, OTEL_SERVICE_NAME set",
			EnvVars: map[string]string{
				"OTEL_SERVICE_NAME":        "my-service",
				"OTEL_RESOURCE_ATTRIBUTES": "deployment.environment.name=production",
			},
			Language: "test",
			ExpectedChecks: []string{
				"Resource Attributes: Service name is set via OTEL_SERVICE_NAME to 'my-service'",
				"Resource Attributes: Resource attribute deployment.environment.name is set to 'production'",
			},
			ExpectedWarnings: []string{
				"Resource Attributes: Set OTEL_RESOURCE_ATTRIBUTES=\"service.namespace=shop\": An optional namespace for service.name",
				"Resource Attributes: Set OTEL_RESOURCE_ATTRIBUTES=\"service.instance.id=checkout-123\": The unique instance, e.g. the pod name",
				"Resource Attributes: Set OTEL_RESOURCE_ATTRIBUTES=\"service.version=1.2\": The application version, to see if a new version has introduced a bug",
			},
		},
		{
			Name: "both service.name and OTEL_SERVICE_NAME set",
			EnvVars: map[string]string{
				"OTEL_SERVICE_NAME":        "my-otel-service",
				"OTEL_RESOURCE_ATTRIBUTES": "service.name=my-service,deployment.environment.name=production",
			},
			Language: "test",
			ExpectedChecks: []string{
				"Resource Attributes: Service name is set via OTEL_SERVICE_NAME to 'my-otel-service'",
				"Resource Attributes: Resource attribute deployment.environment.name is set to 'production'",
			},
			ExpectedWarnings: []string{
				"Resource Attributes: Set OTEL_RESOURCE_ATTRIBUTES=\"service.namespace=shop\": An optional namespace for service.name",
				"Resource Attributes: Set OTEL_RESOURCE_ATTRIBUTES=\"service.instance.id=checkout-123\": The unique instance, e.g. the pod name",
				"Resource Attributes: Set OTEL_RESOURCE_ATTRIBUTES=\"service.version=1.2\": The application version, to see if a new version has introduced a bug",
			},
		},
		{
			Name:     "no attributes present, no OTEL_SERVICE_NAME",
			EnvVars:  map[string]string{},
			Language: "test",
			ExpectedWarnings: []string{
				"Resource Attributes: Set OTEL_RESOURCE_ATTRIBUTES=\"service.namespace=shop\": An optional namespace for service.name",
				"Resource Attributes: Set OTEL_RESOURCE_ATTRIBUTES=\"deployment.environment.name=production\": Name of the deployment environment (staging or production)",
				"Resource Attributes: Set OTEL_RESOURCE_ATTRIBUTES=\"service.instance.id=checkout-123\": The unique instance, e.g. the pod name",
				"Resource Attributes: Set OTEL_RESOURCE_ATTRIBUTES=\"service.version=1.2\": The application version, to see if a new version has introduced a bug",
				"Resource Attributes: Set OTEL_SERVICE_NAME=\"checkout\": The application name",
			},
		},
		{
			Name: "multiple resource attributes in OTEL_RESOURCE_ATTRIBUTES",
			EnvVars: map[string]string{
				"OTEL_RESOURCE_ATTRIBUTES": "service.name=my-service,service.namespace=my-namespace,service.version=1.0.0",
			},
			Language: "test",
			ExpectedChecks: []string{
				"Resource Attributes: Service name is set via OTEL_RESOURCE_ATTRIBUTES to 'my-service'",
				"Resource Attributes: Resource attribute service.namespace is set to 'my-namespace'",
				"Resource Attributes: Resource attribute service.version is set to '1.0.0'",
			},
			ExpectedWarnings: []string{
				"Resource Attributes: Set OTEL_RESOURCE_ATTRIBUTES=\"deployment.environment.name=production\": Name of the deployment environment (staging or production)",
				"Resource Attributes: Set OTEL_RESOURCE_ATTRIBUTES=\"service.instance.id=checkout-123\": The unique instance, e.g. the pod name",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			utils.RunEnvVarComponentTest(t, tt, "Resource Attributes",
				func(reporter utils.Reporter, c *utils.ComponentReporter, language string, components []string) {
					CheckResourceAttributes(c)
				})
		})
	}
}

func TestCheckExporterEnvVars(t *testing.T) {
	correct := correctWith(map[string]string{})
	tests := []utils.EnvVarTestCase{
		{
			Name:     "all required env vars set correctly",
			EnvVars:  correct,
			Language: "python",
			ExpectedChecks: []string{
				"Common Environment Variables: The value of OTEL_METRICS_EXPORTER is set to 'otlp' (default value)",
				"Common Environment Variables: The value of OTEL_TRACES_EXPORTER is set to 'otlp' (default value)",
				"Common Environment Variables: The value of OTEL_LOGS_EXPORTER is set to 'otlp' (default value)",
			},
		},
		{
			Name: "exporters set to none",
			EnvVars: correctWith(map[string]string{
				"OTEL_METRICS_EXPORTER": "none",
				"OTEL_TRACES_EXPORTER":  "none",
				"OTEL_LOGS_EXPORTER":    "none",
			}),
			Language: "python",
			ExpectedErrors: []string{
				"Common Environment Variables: The value of OTEL_METRICS_EXPORTER cannot be 'none'. Change the value to 'otlp' or leave it unset",
				"Common Environment Variables: The value of OTEL_TRACES_EXPORTER cannot be 'none'. Change the value to 'otlp' or leave it unset",
				"Common Environment Variables: The value of OTEL_LOGS_EXPORTER cannot be 'none'. Change the value to 'otlp' or leave it unset",
			},
			IgnoreChecks: true,
		},
		{
			Name:             "nothing set",
			EnvVars:          map[string]string{},
			Language:         "python",
			Components:       []string{"beyla"},
			ExpectedWarnings: []string{},
			ExpectedChecks: []string{
				"Common Environment Variables: The value of OTEL_METRICS_EXPORTER is set to 'otlp' (default value)",
				"Common Environment Variables: The value of OTEL_TRACES_EXPORTER is set to 'otlp' (default value)",
				"Common Environment Variables: The value of OTEL_LOGS_EXPORTER is set to 'otlp' (default value)",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			utils.RunEnvVarComponentTest(t, tt, "Common Environment Variables",
				func(reporter utils.Reporter, c *utils.ComponentReporter, language string, components []string) {
					CheckExporterEnvVars(c, language)
				})
		})
	}
}

func correctWith(add map[string]string) map[string]string {
	m := map[string]string{
		"OTEL_SERVICE_NAME":     "test-service",
		"OTEL_METRICS_EXPORTER": "otlp",
		"OTEL_TRACES_EXPORTER":  "otlp",
		"OTEL_LOGS_EXPORTER":    "otlp",
	}
	for k, v := range add {
		m[k] = v
	}
	return m
}
