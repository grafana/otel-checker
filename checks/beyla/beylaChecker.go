package beyla

import (
	"github.com/grafana/otel-checker/checks/env"
	"github.com/grafana/otel-checker/checks/utils"
)

var (
	ServiceName = env.EnvVar{
		Name:        "BEYLA_SERVICE_NAME",
		Recommended: true,
		Description: "Service name for Beyla",
		ExplainID:   "beyla.service-name.unset",
	}

	OpenPort = env.EnvVar{
		Name:        "BEYLA_OPEN_PORT",
		Required:    true,
		Description: "Port for Beyla to listen on",
		ExplainID:   "beyla.open-port.unset",
	}

	GrafanaCloudSubmit = env.EnvVar{
		Name:        "GRAFANA_CLOUD_SUBMIT",
		Required:    true,
		Description: "Types of telemetry to submit to Grafana Cloud",
		ExplainID:   "beyla.grafana-cloud-submit.unset",
	}

	GrafanaCloudInstanceID = env.EnvVar{
		Name:        "GRAFANA_CLOUD_INSTANCE_ID",
		Required:    true,
		Description: "Grafana Cloud instance ID",
		ExplainID:   "beyla.grafana-cloud-instance-id.unset",
	}

	GrafanaCloudAPIKey = env.EnvVar{
		Name:        "GRAFANA_CLOUD_API_KEY",
		Required:    true,
		Description: "Grafana Cloud API key",
		ExplainID:   "beyla.grafana-cloud-api-key.unset",
	}
)

func CheckBeylaSetup(reporter *utils.ComponentReporter, language string) {
	CheckEnvVars(reporter, language)
}

func CheckEnvVars(reporter *utils.ComponentReporter, language string) {
	env.CheckEnvVars(reporter, language,
		ServiceName,
		OpenPort,
		GrafanaCloudSubmit,
		GrafanaCloudInstanceID,
		GrafanaCloudAPIKey)
}
