package grafana

import (
	"os"
	"testing"

	"otel-checker/checks/utils"
)

func TestCheckEnvVarsGrafana(t *testing.T) {
	tests := []struct {
		name             string
		envVars          map[string]string
		language         string
		components       []string
		expectedErrors   []string
		expectedChecks   []string
		expectedWarnings []string
	}{
		{
			name: "all required env vars set correctly",
			envVars: map[string]string{
				"OTEL_SERVICE_NAME":           "test-service",
				"OTEL_EXPORTER_OTLP_PROTOCOL": "http/protobuf",
				"OTEL_METRICS_EXPORTER":       "otlp",
				"OTEL_TRACES_EXPORTER":        "otlp",
				"OTEL_LOGS_EXPORTER":          "otlp",
				"OTEL_EXPORTER_OTLP_ENDPOINT": "https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
				"OTEL_EXPORTER_OTLP_HEADERS":  "Authorization=Basic dXNlcm5hbWU6cGFzc3dvcmQ=",
			},
			language:       "python",
			components:     []string{},
			expectedErrors: []string{},
			expectedChecks: []string{
				"OTEL_SERVICE_NAME is set",
				"OTEL_EXPORTER_OTLP_PROTOCOL set to 'http/protobuf'",
				"The value of OTEL_METRICS_EXPORTER is set to 'otlp'",
				"The value of OTEL_TRACES_EXPORTER is set to 'otlp'",
				"The value of OTEL_LOGS_EXPORTER is set to 'otlp'",
				"OTEL_EXPORTER_OTLP_ENDPOINT set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
				"OTEL_EXPORTER_OTLP_HEADERS is set correctly",
			},
			expectedWarnings: []string{},
		},
		{
			name: "missing service name",
			envVars: map[string]string{
				"OTEL_EXPORTER_OTLP_PROTOCOL": "http/protobuf",
			},
			language:       "python",
			components:     []string{},
			expectedErrors: []string{},
			expectedChecks: []string{},
			expectedWarnings: []string{
				"It's recommended the environment variable OTEL_SERVICE_NAME to be set to your service name, for easier identification",
			},
		},
		{
			name: "incorrect protocol",
			envVars: map[string]string{
				"OTEL_EXPORTER_OTLP_PROTOCOL": "grpc",
			},
			language:   "python",
			components: []string{},
			expectedErrors: []string{
				"OTEL_EXPORTER_OTLP_PROTOCOL is not set to 'http/protobuf'",
			},
			expectedChecks:   []string{},
			expectedWarnings: []string{},
		},
		{
			name: "exporters set to none",
			envVars: map[string]string{
				"OTEL_METRICS_EXPORTER": "none",
				"OTEL_TRACES_EXPORTER":  "none",
				"OTEL_LOGS_EXPORTER":    "none",
			},
			language:   "python",
			components: []string{},
			expectedErrors: []string{
				"The value of OTEL_METRICS_EXPORTER cannot be 'none'. Change the value to 'otlp' or leave it unset",
				"The value of OTEL_TRACES_EXPORTER cannot be 'none'. Change the value to 'otlp' or leave it unset",
				"The value of OTEL_LOGS_EXPORTER cannot be 'none'. Change the value to 'otlp' or leave it unset",
			},
			expectedChecks:   []string{},
			expectedWarnings: []string{},
		},
		{
			name: "beyla component with required env vars",
			envVars: map[string]string{
				"BEYLA_SERVICE_NAME":        "test-service",
				"BEYLA_OPEN_PORT":           "8080",
				"GRAFANA_CLOUD_SUBMIT":      "metrics,traces",
				"GRAFANA_CLOUD_INSTANCE_ID": "test-instance",
				"GRAFANA_CLOUD_API_KEY":     "test-key",
			},
			language:       "python",
			components:     []string{"beyla"},
			expectedErrors: []string{},
			expectedChecks: []string{
				"BEYLA_SERVICE_NAME is set",
				"BEYLA_SERVICE_NAME is set",
				"GRAFANA_CLOUD_SUBMIT is set correctly",
				"GRAFANA_CLOUD_INSTANCE_ID is set",
				"GRAFANA_CLOUD_API_KEY is set",
			},
			expectedWarnings: []string{},
		},
		{
			name:       "beyla component with missing env vars",
			envVars:    map[string]string{},
			language:   "python",
			components: []string{"beyla"},
			expectedErrors: []string{
				"BEYLA_OPEN_PORT must be set",
				"GRAFANA_CLOUD_SUBMIT must be set to 'metrics' and/or 'traces'",
				"GRAFANA_CLOUD_INSTANCE_ID must be set",
				"GRAFANA_CLOUD_API_KEY must be set",
			},
			expectedChecks: []string{},
			expectedWarnings: []string{
				"It's recommended the environment variable BEYLA_SERVICE_NAME to be set to your service name",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables for test
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				// Clean up environment variables after test
				for k := range tt.envVars {
					os.Unsetenv(k)
				}
			}()

			reporter := utils.Reporter{}
			component := reporter.Component("Grafana Cloud")
			checkEnvVarsGrafana(reporter, component, tt.language, tt.components)

			// Check errors
			if len(component.Errors) != len(tt.expectedErrors) {
				t.Errorf("expected %d errors, got %d", len(tt.expectedErrors), len(component.Errors))
			}
			for i, err := range component.Errors {
				if err != tt.expectedErrors[i] {
					t.Errorf("error %d: expected %q, got %q", i, tt.expectedErrors[i], err)
				}
			}

			// Check successful checks
			if len(component.Checks) != len(tt.expectedChecks) {
				t.Errorf("expected %d checks, got %d", len(tt.expectedChecks), len(component.Checks))
			}
			for i, check := range component.Checks {
				if check != tt.expectedChecks[i] {
					t.Errorf("check %d: expected %q, got %q", i, tt.expectedChecks[i], check)
				}
			}

			// Check warnings
			if len(component.Warnings) != len(tt.expectedWarnings) {
				t.Errorf("expected %d warnings, got %d", len(tt.expectedWarnings), len(component.Warnings))
			}
			for i, warning := range component.Warnings {
				if warning != tt.expectedWarnings[i] {
					t.Errorf("warning %d: expected %q, got %q", i, tt.expectedWarnings[i], warning)
				}
			}
		})
	}
}
