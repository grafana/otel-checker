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
				"OTEL_EXPORTER_OTLP_PROTOCOL is set to 'http/protobuf'",
				"OTEL_EXPORTER_OTLP_ENDPOINT set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
				"OTEL_EXPORTER_OTLP_HEADERS is set correctly",
			},
		},
		{
			Name: "incorrect protocol",
			EnvVars: correctWith(map[string]string{
				"OTEL_EXPORTER_OTLP_PROTOCOL": "grpc",
			}),
			Language: "python",
			ExpectedErrors: []string{
				"OTEL_EXPORTER_OTLP_PROTOCOL must be set to 'http/protobuf'",
			},
			IgnoreChecks: true,
		},
		{
			Name:       "nothing set",
			EnvVars:    map[string]string{},
			Language:   "python",
			Components: []string{"beyla"},
			ExpectedErrors: []string{
				"OTEL_EXPORTER_OTLP_PROTOCOL must be set to 'http/protobuf'",
				"OTEL_EXPORTER_OTLP_ENDPOINT is not set — required because signal-specific endpoint(s) missing for: traces, metrics, logs",
				"OTEL_EXPORTER_OTLP_HEADERS is not set. Value should have 'Authorization=Basic%20...'",
			},
		},
		{
			Name: "base endpoint has a signal suffix",
			EnvVars: correctWith(map[string]string{
				"OTEL_EXPORTER_OTLP_ENDPOINT": "https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/traces",
			}),
			Language: "python",
			ExpectedErrors: []string{
				"OTEL_EXPORTER_OTLP_ENDPOINT is not set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp (no signal suffix like /v1/traces)",
			},
			IgnoreChecks: true,
		},
		{
			Name: "all three signal-specific endpoints set, base unset",
			EnvVars: map[string]string{
				"OTEL_EXPORTER_OTLP_PROTOCOL":         "http/protobuf",
				"OTEL_EXPORTER_OTLP_HEADERS":          "Authorization=Basic%20dXNlcm5hbWU6cGFzc3dvcmQ=",
				"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT":  "https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/traces",
				"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT": "https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/metrics",
				"OTEL_EXPORTER_OTLP_LOGS_ENDPOINT":    "https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/logs",
			},
			Language: "python",
			ExpectedChecks: []string{
				"OTEL_EXPORTER_OTLP_PROTOCOL is set to 'http/protobuf'",
				"OTEL_EXPORTER_OTLP_HEADERS is set correctly",
				"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/traces",
				"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/metrics",
				"OTEL_EXPORTER_OTLP_LOGS_ENDPOINT set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/logs",
				"OTEL_EXPORTER_OTLP_ENDPOINT is unset — all signals are covered by signal-specific endpoints",
			},
		},
		{
			Name: "one signal-specific missing, base unset",
			EnvVars: map[string]string{
				"OTEL_EXPORTER_OTLP_PROTOCOL":         "http/protobuf",
				"OTEL_EXPORTER_OTLP_HEADERS":          "Authorization=Basic%20dXNlcm5hbWU6cGFzc3dvcmQ=",
				"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT":  "https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/traces",
				"OTEL_EXPORTER_OTLP_METRICS_ENDPOINT": "https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/metrics",
			},
			Language: "python",
			ExpectedErrors: []string{
				"OTEL_EXPORTER_OTLP_ENDPOINT is not set — required because signal-specific endpoint(s) missing for: logs",
			},
			IgnoreChecks: true,
		},
		{
			Name: "signal-specific endpoint missing /v1/traces path",
			EnvVars: correctWith(map[string]string{
				"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT": "https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
			}),
			Language: "python",
			ExpectedErrors: []string{
				"OTEL_EXPORTER_OTLP_TRACES_ENDPOINT is not set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/traces",
			},
			IgnoreChecks: true,
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
