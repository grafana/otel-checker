package grafana

import (
	"testing"

	"github.com/grafana/otel-checker/checks/utils"
)

func TestCheckEnvVarsGrafana(t *testing.T) {
	correct := correctWith(map[string]string{})
	tests := []utils.EnvVarTestCase{
		{
			Name:     "all required env vars set correctly",
			EnvVars:  correct,
			Language: "python",
			ExpectedChecks: []string{
				"Grafana Cloud: OTEL_EXPORTER_OTLP_PROTOCOL is set to 'http/protobuf'",
				"Grafana Cloud: OTEL_EXPORTER_OTLP_ENDPOINT set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
				"Grafana Cloud: OTEL_EXPORTER_OTLP_HEADERS is set correctly",
			},
		},
		{
			Name: "incorrect protocol",
			EnvVars: correctWith(map[string]string{
				"OTEL_EXPORTER_OTLP_PROTOCOL": "grpc",
			}),
			Language: "python",
			ExpectedErrors: []string{
				"Grafana Cloud: OTEL_EXPORTER_OTLP_PROTOCOL must be set to 'http/protobuf'",
			},
			IgnoreChecks: true,
		},
		{
			Name:       "nothing set",
			EnvVars:    map[string]string{},
			Language:   "python",
			Components: []string{"beyla"},
			ExpectedErrors: []string{
				"Grafana Cloud: OTEL_EXPORTER_OTLP_PROTOCOL must be set to 'http/protobuf'",
				"Grafana Cloud: OTEL_EXPORTER_OTLP_ENDPOINT is not set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
				"Grafana Cloud: OTEL_EXPORTER_OTLP_HEADERS is not set. Value should have 'Authorization=Basic%20...'",
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
		"OTEL_EXPORTER_OTLP_PROTOCOL": "http/protobuf",
		"OTEL_EXPORTER_OTLP_ENDPOINT": "https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
		"OTEL_EXPORTER_OTLP_HEADERS":  "Authorization=Basic%20dXNlcm5hbWU6cGFzc3dvcmQ=",
	}
	for k, v := range add {
		m[k] = v
	}
	return m
}
