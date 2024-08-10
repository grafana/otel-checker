package grafana

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"

	utils "otel-checker/utils"
)

func CheckGrafanaSetup(
	messages *map[string][]string,
	language string,
	components []string,
) {
	checkEnvVarsGrafana(messages, language, components)
	checkAuth(messages)
}

func checkEnvVarsGrafana(
	messages *map[string][]string,
	language string,
	components []string,
) {
	if os.Getenv("OTEL_SERVICE_NAME") == "" {
		utils.AddWarning(messages, "Grafana Cloud", "It's recommended the environment variable OTEL_SERVICE_NAME to be set to your service name, for easier identification")
	} else {
		utils.AddSuccessfulCheck(messages, "Grafana Cloud", "OTEL_SERVICE_NAME is set")
	}

	if os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL") != "http/protobuf" {
		utils.AddError(messages, "Grafana Cloud", "OTEL_EXPORTER_OTLP_PROTOCOL is not set to 'http/protobuf'")
	} else {
		utils.AddSuccessfulCheck(messages, "Grafana Cloud", "OTEL_EXPORTER_OTLP_PROTOCOL set to 'http/protobuf'")
	}

	match, _ := regexp.MatchString("https:\\/\\/.+\\.grafana\\.net\\/otlp", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if match {
		utils.AddSuccessfulCheck(messages, "Grafana Cloud", "OTEL_EXPORTER_OTLP_ENDPOINT set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp")
	} else {
		if strings.Contains(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), "localhost") {
			utils.AddWarning(messages, "Grafana Cloud", "OTEL_EXPORTER_OTLP_ENDPOINT is set to localhost. Update to a Grafana endpoint similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp to be able to send telemetry to your Grafana Cloud instance")
		} else {
			utils.AddError(messages, "Grafana Cloud", "OTEL_EXPORTER_OTLP_ENDPOINT is not set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp")
		}
	}

	tokenStart := "Authorization=Basic "
	if language == "python" {
		tokenStart = "Authorization=Basic%20"
	}
	if strings.Contains(os.Getenv("OTEL_EXPORTER_OTLP_HEADERS"), tokenStart) {
		utils.AddSuccessfulCheck(messages, "Grafana Cloud", "OTEL_EXPORTER_OTLP_HEADERS is set correctly")
	} else {
		utils.AddError(messages, "Grafana Cloud", fmt.Sprintf("OTEL_EXPORTER_OTLP_HEADERS is not set. Value should have '%s'", tokenStart))
	}

	if slices.Contains(components, "beyla") {
		if os.Getenv("BEYLA_SERVICE_NAME") == "" {
			utils.AddWarning(messages, "Beyla", "It's recommended the environment variable BEYLA_SERVICE_NAME to be set to your service name")
		} else {
			utils.AddSuccessfulCheck(messages, "Beyla", "BEYLA_SERVICE_NAME is set")
		}

		if os.Getenv("BEYLA_OPEN_PORT") == "" {
			utils.AddError(messages, "Beyla", "BEYLA_OPEN_PORT must be set")
		} else {
			utils.AddSuccessfulCheck(messages, "Beyla", "BEYLA_SERVICE_NAME is set")
		}

		if os.Getenv("GRAFANA_CLOUD_SUBMIT") == "" {
			utils.AddError(messages, "Beyla", "GRAFANA_CLOUD_SUBMIT must be set to 'metrics' and/or 'traces'")
		} else {
			utils.AddSuccessfulCheck(messages, "Beyla", "GRAFANA_CLOUD_SUBMIT is set correctly")
		}

		if os.Getenv("GRAFANA_CLOUD_INSTANCE_ID") == "" {
			utils.AddError(messages, "Beyla", "GRAFANA_CLOUD_INSTANCE_ID must be set")
		} else {
			utils.AddSuccessfulCheck(messages, "Beyla", "GRAFANA_CLOUD_INSTANCE_ID is set")
		}

		if os.Getenv("GRAFANA_CLOUD_API_KEY") == "" {
			utils.AddError(messages, "Beyla", "GRAFANA_CLOUD_API_KEY must be set")
		} else {
			utils.AddSuccessfulCheck(messages, "Beyla", "GRAFANA_CLOUD_API_KEY is set")
		}
	}

}

func checkAuth(messages *map[string][]string) {
	if strings.Contains(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), "localhost") {
		utils.AddWarning(messages, "Grafana Cloud", "Credentials not checked, since OTEL_EXPORTER_OTLP_ENDPOINT is using localhost")
		return
	}
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") == "" || os.Getenv("OTEL_EXPORTER_OTLP_HEADERS") == "" {
		utils.AddWarning(messages, "Grafana Cloud", "Credentials not checked: Both environment variables OTEL_EXPORTER_OTLP_ENDPOINT and OTEL_EXPORTER_OTLP_HEADERS need to be set")
	} else {
		endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") + "/v1/metrics"
		req, err := http.NewRequest("POST", endpoint, nil)
		if err != nil {
			utils.AddError(messages, "Grafana Cloud", fmt.Sprintf("Error while testing credentials of OTEL_EXPORTER_OTLP_ENDPOINT: %s", err))
		}
		authValue := ""
		for _, h := range strings.SplitN(os.Getenv("OTEL_EXPORTER_OTLP_HEADERS"), ",", -1) {
			key, value, _ := strings.Cut(h, "=")
			if key == "Authorization" {
				authValue = value
			}
		}
		req.Header.Set("Authorization", authValue)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			utils.AddError(messages, "Grafana Cloud", fmt.Sprintf("Error while testing credentials of OTEL_EXPORTER_OTLP_ENDPOINT: %s", err))
		}

		if resp.StatusCode == 401 {
			utils.AddError(messages, "Grafana Cloud", fmt.Sprintf("Error while testing credentials of OTEL_EXPORTER_OTLP_ENDPOINT: %s", resp.Status))
		} else {
			utils.AddSuccessfulCheck(messages, "Grafana Cloud", "Credentials for OTEL_EXPORTER_OTLP_ENDPOINT are correct")
		}
		defer resp.Body.Close()
	}
}
