package beyla

import (
	"github.com/grafana/otel-checker/checks/env"
	"github.com/grafana/otel-checker/checks/utils"
)

var (
	ServiceName = env.EnvVar{
		Name:        "BEYLA_SERVICE_NAME",
		Required:    false,
		Description: "Service name for Beyla",
	}

	OpenPort = env.EnvVar{
		Name:        "BEYLA_OPEN_PORT",
		Required:    true,
		Description: "Port for Beyla to listen on",
	}

	GrafanaCloudSubmit = env.EnvVar{
		Name:        "GRAFANA_CLOUD_SUBMIT",
		Required:    true,
		Description: "Types of telemetry to submit to Grafana Cloud",
	}

	GrafanaCloudInstanceID = env.EnvVar{
		Name:        "GRAFANA_CLOUD_INSTANCE_ID",
		Required:    true,
		Description: "Grafana Cloud instance ID",
	}

	GrafanaCloudAPIKey = env.EnvVar{
		Name:        "GRAFANA_CLOUD_API_KEY",
		Required:    true,
		Description: "Grafana Cloud API key",
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
