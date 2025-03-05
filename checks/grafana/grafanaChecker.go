package grafana

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"

	utils "otel-checker/checks/utils"
)

func CheckGrafanaSetup(reporter utils.Reporter, grafanaReporter *utils.ComponentReporter, language string, components []string) {
	checkEnvVarsGrafana(reporter, grafanaReporter, language, components)
	checkAuth(grafanaReporter)
}

func checkEnvVarsGrafana(reporter utils.Reporter, grafana *utils.ComponentReporter, language string, components []string) {
	if os.Getenv("OTEL_SERVICE_NAME") == "" {
		grafana.AddWarning("It's recommended the environment variable OTEL_SERVICE_NAME to be set to your service name, for easier identification")
	} else {
		grafana.AddSuccessfulCheck("OTEL_SERVICE_NAME is set")
	}

	if os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL") != "http/protobuf" {
		grafana.AddError("OTEL_EXPORTER_OTLP_PROTOCOL is not set to 'http/protobuf'")
	} else {
		grafana.AddSuccessfulCheck("OTEL_EXPORTER_OTLP_PROTOCOL set to 'http/protobuf'")
	}

	if os.Getenv("OTEL_METRICS_EXPORTER") == "none" {
		grafana.AddError("The value of OTEL_METRICS_EXPORTER cannot be 'none'. Change the value to 'otlp' or leave it unset")
	} else {
		if os.Getenv("OTEL_METRICS_EXPORTER") == "" {
			grafana.AddSuccessfulCheck("OTEL_METRICS_EXPORTER is unset, with a default value of 'otlp'")
		} else {
			grafana.AddSuccessfulCheck(fmt.Sprintf("The value of OTEL_METRICS_EXPORTER is set to '%s'", os.Getenv("OTEL_METRICS_EXPORTER")))
		}
	}
	if os.Getenv("OTEL_TRACES_EXPORTER") == "none" {
		grafana.AddError("The value of OTEL_TRACES_EXPORTER cannot be 'none'. Change the value to 'otlp' or leave it unset")
	} else {
		if os.Getenv("OTEL_TRACES_EXPORTER") == "" {
			grafana.AddSuccessfulCheck("OTEL_TRACES_EXPORTER is unset, with a default value of 'otlp'")
		} else {
			grafana.AddSuccessfulCheck(fmt.Sprintf("The value of OTEL_TRACES_EXPORTER is set to '%s'", os.Getenv("OTEL_TRACES_EXPORTER")))
		}
	}
	if os.Getenv("OTEL_LOGS_EXPORTER") == "none" {
		grafana.AddError("The value of OTEL_LOGS_EXPORTER cannot be 'none'. Change the value to 'otlp' or leave it unset")
	} else {
		if os.Getenv("OTEL_LOGS_EXPORTER") == "" {
			grafana.AddSuccessfulCheck("OTEL_LOGS_EXPORTER is unset, with a default value of 'otlp'")
		} else {
			grafana.AddSuccessfulCheck(fmt.Sprintf("The value of OTEL_LOGS_EXPORTER is set to '%s'", os.Getenv("OTEL_LOGS_EXPORTER")))
		}
	}

	match, _ := regexp.MatchString("https:\\/\\/.+\\.grafana\\.net\\/otlp", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if match {
		grafana.AddSuccessfulCheck("OTEL_EXPORTER_OTLP_ENDPOINT set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp")
	} else {
		if strings.Contains(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), "localhost") {
			grafana.AddWarning("OTEL_EXPORTER_OTLP_ENDPOINT is set to localhost. Update to a Grafana endpoint similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp to be able to send telemetry to your Grafana Cloud instance")
		} else {
			grafana.AddError("OTEL_EXPORTER_OTLP_ENDPOINT is not set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp")
		}
	}

	tokenStart := "Authorization=Basic "
	if language == "python" {
		tokenStart = "Authorization=Basic%20"
	}
	if strings.Contains(os.Getenv("OTEL_EXPORTER_OTLP_HEADERS"), tokenStart) {
		grafana.AddSuccessfulCheck("OTEL_EXPORTER_OTLP_HEADERS is set correctly")
	} else {
		grafana.AddError(fmt.Sprintf("OTEL_EXPORTER_OTLP_HEADERS is not set. Value should have '%s...'", tokenStart))
	}

	if slices.Contains(components, "beyla") {
		beyla := reporter.Component("Beyla")
		if os.Getenv("BEYLA_SERVICE_NAME") == "" {
			beyla.AddWarning("It's recommended the environment variable BEYLA_SERVICE_NAME to be set to your service name")
		} else {
			beyla.AddSuccessfulCheck("BEYLA_SERVICE_NAME is set")
		}

		if os.Getenv("BEYLA_OPEN_PORT") == "" {
			beyla.AddError("BEYLA_OPEN_PORT must be set")
		} else {
			beyla.AddSuccessfulCheck("BEYLA_SERVICE_NAME is set")
		}

		if os.Getenv("GRAFANA_CLOUD_SUBMIT") == "" {
			beyla.AddError("GRAFANA_CLOUD_SUBMIT must be set to 'metrics' and/or 'traces'")
		} else {
			beyla.AddSuccessfulCheck("GRAFANA_CLOUD_SUBMIT is set correctly")
		}

		if os.Getenv("GRAFANA_CLOUD_INSTANCE_ID") == "" {
			beyla.AddError("GRAFANA_CLOUD_INSTANCE_ID must be set")
		} else {
			beyla.AddSuccessfulCheck("GRAFANA_CLOUD_INSTANCE_ID is set")
		}

		if os.Getenv("GRAFANA_CLOUD_API_KEY") == "" {
			beyla.AddError("GRAFANA_CLOUD_API_KEY must be set")
		} else {
			beyla.AddSuccessfulCheck("GRAFANA_CLOUD_API_KEY is set")
		}
	}

}

func checkAuth(reporter *utils.ComponentReporter) {
	if strings.Contains(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"), "localhost") {
		reporter.AddWarning("Credentials not checked, since OTEL_EXPORTER_OTLP_ENDPOINT is using localhost")
		return
	}
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") == "" || os.Getenv("OTEL_EXPORTER_OTLP_HEADERS") == "" {
		reporter.AddWarning("Credentials not checked, since both environment variables OTEL_EXPORTER_OTLP_ENDPOINT and OTEL_EXPORTER_OTLP_HEADERS need to be set for this check")
	} else {
		endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") + "/v1/metrics"
		req, err := http.NewRequest("POST", endpoint, nil)
		if err != nil {
			reporter.AddError(fmt.Sprintf("Error while testing credentials of OTEL_EXPORTER_OTLP_ENDPOINT: %s", err))
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
			reporter.AddError(fmt.Sprintf("Error while testing credentials of OTEL_EXPORTER_OTLP_ENDPOINT: %s", err))
		}

		if resp.StatusCode == 401 {
			reporter.AddError(fmt.Sprintf("Error while testing credentials of OTEL_EXPORTER_OTLP_ENDPOINT: %s", resp.Status))
		} else {
			reporter.AddSuccessfulCheck("Credentials for OTEL_EXPORTER_OTLP_ENDPOINT are correct")
		}
		defer resp.Body.Close()
	}
}
