package beyla

import (
	"testing"

	"github.com/grafana/otel-checker/checks/utils"
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
			Language:   "python",
			Components: []string{"beyla"},
			ExpectedChecks: []string{
				"BEYLA_SERVICE_NAME is set to 'test-service'",
				"BEYLA_OPEN_PORT is set to '8080'",
				"GRAFANA_CLOUD_SUBMIT is set to 'metrics,traces'",
				"GRAFANA_CLOUD_INSTANCE_ID is set to 'test-instance'",
				"GRAFANA_CLOUD_API_KEY is set to 'test-key'",
			},
		},
		{
			Name:       "nothing set",
			EnvVars:    map[string]string{},
			Language:   "python",
			Components: []string{"beyla"},
			ExpectedErrors: []string{
				"BEYLA_OPEN_PORT is not set",
				"GRAFANA_CLOUD_SUBMIT is not set",
				"GRAFANA_CLOUD_INSTANCE_ID is not set",
				"GRAFANA_CLOUD_API_KEY is not set",
			},
			ExpectedWarnings: []string{
				"BEYLA_SERVICE_NAME is not set",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			utils.RunEnvVarComponentTest(t, tt, "Beyla",
				func(reporter utils.Reporter, c *utils.ComponentReporter, language string, components []string) {
					CheckBeylaSetup(c, language)
				})
		})
	}
}
