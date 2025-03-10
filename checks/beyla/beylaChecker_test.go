package beyla

import (
	"testing"

	"otel-checker/checks/utils"
)

func TestCheckEnvVarsBeyla(t *testing.T) {
	tests := []utils.EnvVarTestCase{
		{
			Name: "beyla component with required env vars",
			EnvVars: map[string]string{
				"BEYLA_SERVICE_NAME":        "test-service",
				"BEYLA_OPEN_PORT":           "8080",
				"GRAFANA_CLOUD_SUBMIT":      "metrics,traces",
				"GRAFANA_CLOUD_INSTANCE_ID": "test-instance",
				"GRAFANA_CLOUD_API_KEY":     "test-key",
			},
			Language:       "python",
			Components:     []string{"beyla"},
			ExpectedChecks: []string{},
		},
		{
			Name:           "nothing set",
			EnvVars:        map[string]string{},
			Language:       "python",
			Components:     []string{"beyla"},
			ExpectedErrors: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			utils.RunEnvVarComponentTest(t, tt, "Grafana Cloud",
				func(reporter utils.Reporter, c *utils.ComponentReporter, language string, components []string) {
					CheckBeylaSetup(c, language)
				})
		})
	}
}
