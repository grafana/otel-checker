package env

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
				"Common Environment Variables: OTEL_SERVICE_NAME is set to 'test-service'",
				"Common Environment Variables: The value of OTEL_METRICS_EXPORTER is set to 'otlp' (default value)",
				"Common Environment Variables: The value of OTEL_TRACES_EXPORTER is set to 'otlp' (default value)",
				"Common Environment Variables: The value of OTEL_LOGS_EXPORTER is set to 'otlp' (default value)",
			},
		},
		{
			Name: "missing service name",
			EnvVars: correctWith(map[string]string{
				"OTEL_SERVICE_NAME": "",
			}),
			Language: "python",
			ExpectedWarnings: []string{
				"Common Environment Variables: It's recommended the environment variable OTEL_SERVICE_NAME to be set to your service name, for easier identification",
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
				"Common Environment Variables: The value of OTEL_METRICS_EXPORTER cannot be 'none'. Change the value to 'otlp' or leave it unset",
				"Common Environment Variables: The value of OTEL_TRACES_EXPORTER cannot be 'none'. Change the value to 'otlp' or leave it unset",
				"Common Environment Variables: The value of OTEL_LOGS_EXPORTER cannot be 'none'. Change the value to 'otlp' or leave it unset",
			},
			IgnoreChecks: true,
		},
		{
			Name:       "nothing set",
			EnvVars:    map[string]string{},
			Language:   "python",
			Components: []string{"beyla"},
			ExpectedWarnings: []string{
				"Common Environment Variables: It's recommended the environment variable OTEL_SERVICE_NAME to be set to your service name, for easier identification",
			},
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
					CheckCommonEnvVars(c, language)
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
