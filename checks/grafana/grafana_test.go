package grafana

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	"otel-checker/checks/utils"
)

func TestCheckEnvVarsGrafana(t *testing.T) {
	correct := correctWith(map[string]string{})
	tests := []struct {
		name             string
		envVars          map[string]string
		language         string
		components       []string
		expectedErrors   []string
		expectedChecks   []string
		expectedWarnings []string
		ignoreWarnings   bool
		ignoreErrors     bool
		ignoreChecks     bool
	}{
		{
			name:     "all required env vars set correctly",
			envVars:  correct,
			language: "python",
			expectedChecks: []string{
				"Grafana Cloud: OTEL_SERVICE_NAME is set",
				"Grafana Cloud: OTEL_EXPORTER_OTLP_PROTOCOL set to 'http/protobuf'",
				"Grafana Cloud: The value of OTEL_METRICS_EXPORTER is set to 'otlp'",
				"Grafana Cloud: The value of OTEL_TRACES_EXPORTER is set to 'otlp'",
				"Grafana Cloud: The value of OTEL_LOGS_EXPORTER is set to 'otlp'",
				"Grafana Cloud: OTEL_EXPORTER_OTLP_ENDPOINT set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
				"Grafana Cloud: OTEL_EXPORTER_OTLP_HEADERS is set correctly",
			},
		},
		{
			name: "missing service name",
			envVars: correctWith(map[string]string{
				"OTEL_SERVICE_NAME": "",
			}),
			language: "python",
			expectedWarnings: []string{
				"Grafana Cloud: It's recommended the environment variable OTEL_SERVICE_NAME to be set to your service name, for easier identification",
			},
			ignoreChecks: true,
		},
		{
			name: "incorrect protocol",
			envVars: correctWith(map[string]string{
				"OTEL_EXPORTER_OTLP_PROTOCOL": "grpc",
			}),
			language: "python",
			expectedErrors: []string{
				"Grafana Cloud: OTEL_EXPORTER_OTLP_PROTOCOL is not set to 'http/protobuf'",
			},
			ignoreChecks: true,
		},
		{
			name: "exporters set to none",
			envVars: correctWith(map[string]string{
				"OTEL_METRICS_EXPORTER": "none",
				"OTEL_TRACES_EXPORTER":  "none",
				"OTEL_LOGS_EXPORTER":    "none",
			}),
			language: "python",
			expectedErrors: []string{
				"Grafana Cloud: The value of OTEL_METRICS_EXPORTER cannot be 'none'. Change the value to 'otlp' or leave it unset",
				"Grafana Cloud: The value of OTEL_TRACES_EXPORTER cannot be 'none'. Change the value to 'otlp' or leave it unset",
				"Grafana Cloud: The value of OTEL_LOGS_EXPORTER cannot be 'none'. Change the value to 'otlp' or leave it unset",
			},
			ignoreChecks: true,
		},
		{
			name: "beyla component with required env vars",
			envVars: correctWith(map[string]string{
				"BEYLA_SERVICE_NAME":        "test-service",
				"BEYLA_OPEN_PORT":           "8080",
				"GRAFANA_CLOUD_SUBMIT":      "metrics,traces",
				"GRAFANA_CLOUD_INSTANCE_ID": "test-instance",
				"GRAFANA_CLOUD_API_KEY":     "test-key",
			}),
			language:   "python",
			components: []string{"beyla"},
			expectedChecks: []string{
				"Grafana Cloud: OTEL_SERVICE_NAME is set",
				"Grafana Cloud: OTEL_EXPORTER_OTLP_PROTOCOL set to 'http/protobuf'",
				"Grafana Cloud: The value of OTEL_METRICS_EXPORTER is set to 'otlp'",
				"Grafana Cloud: The value of OTEL_TRACES_EXPORTER is set to 'otlp'",
				"Grafana Cloud: The value of OTEL_LOGS_EXPORTER is set to 'otlp'",
				"Grafana Cloud: OTEL_EXPORTER_OTLP_ENDPOINT set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
				"Grafana Cloud: OTEL_EXPORTER_OTLP_HEADERS is set correctly",
			},
		},
		{
			name:       "nothing set",
			envVars:    map[string]string{},
			language:   "python",
			components: []string{"beyla"},
			expectedErrors: []string{
				"Grafana Cloud: OTEL_EXPORTER_OTLP_PROTOCOL is not set to 'http/protobuf'",
				"Grafana Cloud: OTEL_EXPORTER_OTLP_ENDPOINT is not set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
				"Grafana Cloud: OTEL_EXPORTER_OTLP_HEADERS is not set. Value should have 'Authorization=Basic%20...'"},
			expectedWarnings: []string{
				"Grafana Cloud: It's recommended the environment variable OTEL_SERVICE_NAME to be set to your service name, for easier identification"},
			expectedChecks: []string{
				"Grafana Cloud: OTEL_METRICS_EXPORTER is unset, with a default value of 'otlp'",
				"Grafana Cloud: OTEL_TRACES_EXPORTER is unset, with a default value of 'otlp'", "Grafana Cloud: OTEL_LOGS_EXPORTER is unset, with a default value of 'otlp'",
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

			if !tt.ignoreErrors {
				require.Equal(t, tt.expectedErrors, component.Errors, "errors mismatch")
			}

			if !tt.ignoreChecks {
				require.Equal(t, tt.expectedChecks, component.Checks, "checks mismatch")
			}

			if !tt.ignoreWarnings {
				require.Equal(t, tt.expectedWarnings, component.Warnings, "warnings mismatch")
			}
		})
	}
}

func correctWith(add map[string]string) map[string]string {
	m := map[string]string{
		"OTEL_SERVICE_NAME":           "test-service",
		"OTEL_EXPORTER_OTLP_PROTOCOL": "http/protobuf",
		"OTEL_METRICS_EXPORTER":       "otlp",
		"OTEL_TRACES_EXPORTER":        "otlp",
		"OTEL_LOGS_EXPORTER":          "otlp",
		"OTEL_EXPORTER_OTLP_ENDPOINT": "https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
		"OTEL_EXPORTER_OTLP_HEADERS":  "Authorization=Basic%20dXNlcm5hbWU6cGFzc3dvcmQ=",
	}
	for k, v := range add {
		m[k] = v
	}
	return m
}
