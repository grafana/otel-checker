package grafana

import (
	"testing"

	"otel-checker/checks/utils"
)

func TestCheckEnvVarsGrafana(t *testing.T) {
	correct := correctWith(map[string]string{})
	tests := []utils.EnvVarTestCase{
		{
			Name:     "all required env vars set correctly",
			EnvVars:  correct,
			Language: "python",
			ExpectedChecks: []string{
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
			Name: "missing service name",
			EnvVars: correctWith(map[string]string{
				"OTEL_SERVICE_NAME": "",
			}),
			Language: "python",
			ExpectedWarnings: []string{
				"Grafana Cloud: It's recommended the environment variable OTEL_SERVICE_NAME to be set to your service name, for easier identification",
			},
			IgnoreChecks: true,
		},
		{
			Name: "incorrect protocol",
			EnvVars: correctWith(map[string]string{
				"OTEL_EXPORTER_OTLP_PROTOCOL": "grpc",
			}),
			Language: "python",
			ExpectedErrors: []string{
				"Grafana Cloud: OTEL_EXPORTER_OTLP_PROTOCOL is not set to 'http/protobuf'",
			},
			IgnoreChecks: true,
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
				"Grafana Cloud: The value of OTEL_METRICS_EXPORTER cannot be 'none'. Change the value to 'otlp' or leave it unset",
				"Grafana Cloud: The value of OTEL_TRACES_EXPORTER cannot be 'none'. Change the value to 'otlp' or leave it unset",
				"Grafana Cloud: The value of OTEL_LOGS_EXPORTER cannot be 'none'. Change the value to 'otlp' or leave it unset",
			},
			IgnoreChecks: true,
		},
		{
			Name: "beyla component with required env vars",
			EnvVars: correctWith(map[string]string{
				"BEYLA_SERVICE_NAME":        "test-service",
				"BEYLA_OPEN_PORT":           "8080",
				"GRAFANA_CLOUD_SUBMIT":      "metrics,traces",
				"GRAFANA_CLOUD_INSTANCE_ID": "test-instance",
				"GRAFANA_CLOUD_API_KEY":     "test-key",
			}),
			Language:   "python",
			Components: []string{"beyla"},
			ExpectedChecks: []string{
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
			Name:       "nothing set",
			EnvVars:    map[string]string{},
			Language:   "python",
			Components: []string{"beyla"},
			ExpectedErrors: []string{
				"Grafana Cloud: OTEL_EXPORTER_OTLP_PROTOCOL is not set to 'http/protobuf'",
				"Grafana Cloud: OTEL_EXPORTER_OTLP_ENDPOINT is not set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
				"Grafana Cloud: OTEL_EXPORTER_OTLP_HEADERS is not set. Value should have 'Authorization=Basic%20...'"},
			ExpectedWarnings: []string{
				"Grafana Cloud: It's recommended the environment variable OTEL_SERVICE_NAME to be set to your service name, for easier identification"},
			ExpectedChecks: []string{
				"Grafana Cloud: OTEL_METRICS_EXPORTER is unset, with a default value of 'otlp'",
				"Grafana Cloud: OTEL_TRACES_EXPORTER is unset, with a default value of 'otlp'",
				"Grafana Cloud: OTEL_LOGS_EXPORTER is unset, with a default value of 'otlp'",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			utils.RunEnvVarComponentTest(t, tt, "Grafana Cloud", checkEnvVarsGrafana)
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
