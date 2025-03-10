package grafana

import (
	"fmt"
	"net/http"
	"otel-checker/checks/env"
	"otel-checker/checks/utils"
	"regexp"
	"strings"
)

var (
	OtelExporterOTLPProtocol = env.EnvVar{
		Name:          "OTEL_EXPORTER_OTLP_PROTOCOL",
		RequiredValue: "http/protobuf",
		Description:   "Protocol for OTLP exporter",
		Message:       "OTEL_EXPORTER_OTLP_PROTOCOL must be set to 'http/protobuf'",
	}

	OtelExporterOTLPEndpoint = env.EnvVar{
		Name:     "OTEL_EXPORTER_OTLP_ENDPOINT",
		Required: true,
		Validator: func(value string, language string, reporter *utils.ComponentReporter) {
			match, _ := regexp.MatchString("https://.+\\.grafana\\.net/otlp", value)
			if match {
				reporter.AddSuccessfulCheck("OTEL_EXPORTER_OTLP_ENDPOINT set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp")
			} else {
				if strings.Contains(value, "localhost") {
					reporter.AddWarning("OTEL_EXPORTER_OTLP_ENDPOINT is set to localhost. Update to a Grafana endpoint similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp to be able to send telemetry to your Grafana Cloud instance")
				} else {
					reporter.AddError("OTEL_EXPORTER_OTLP_ENDPOINT is not set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp")
				}
			}
		},
		Description: "OTLP exporter endpoint",
	}

	OtelExporterOTLPHeaders = env.EnvVar{
		Name:     "OTEL_EXPORTER_OTLP_HEADERS",
		Required: true,
		Validator: func(value string, language string, reporter *utils.ComponentReporter) {
			tokenStart := "Authorization=Basic "
			if language == "python" {
				tokenStart = "Authorization=Basic%20"
			}
			if strings.Contains(value, tokenStart) {
				reporter.AddSuccessfulCheck("OTEL_EXPORTER_OTLP_HEADERS is set correctly")
			} else {
				reporter.AddError(fmt.Sprintf("OTEL_EXPORTER_OTLP_HEADERS is not set. Value should have '%s...'", tokenStart))
			}
		},
		Description: "OTLP exporter headers",
	}
)

func CheckGrafanaSetup(reporter utils.Reporter, grafanaReporter *utils.ComponentReporter, commands utils.Commands) {
	checkEnvVarsGrafana(reporter, grafanaReporter, commands.Language, commands.Components)

	// Only check authentication if explicitly enabled via flag
	if commands.EnableGrafanaCloudCheck {
		checkAuth(grafanaReporter)
	}
}

func checkEnvVarsGrafana(reporter utils.Reporter, grafana *utils.ComponentReporter, language string, components []string) {
	// Check common OpenTelemetry variables
	commonVars := []env.EnvVar{
		OtelExporterOTLPProtocol,
		OtelExporterOTLPEndpoint,
		OtelExporterOTLPHeaders,
	}

	env.CheckEnvVars(grafana, language, commonVars...)
}

func checkAuth(reporter *utils.ComponentReporter) {
	endpoint := env.GetValue(OtelExporterOTLPEndpoint)
	if strings.Contains(endpoint, "localhost") {
		reporter.AddWarning("Credentials not checked, since OTEL_EXPORTER_OTLP_ENDPOINT is using localhost")
		return
	}

	headers := env.GetValue(OtelExporterOTLPHeaders)
	if endpoint == "" || headers == "" {
		reporter.AddWarning("Credentials not checked, since both environment variables OTEL_EXPORTER_OTLP_ENDPOINT and OTEL_EXPORTER_OTLP_HEADERS need to be set for this check")
		return
	}

	// Test credentials
	testEndpoint := endpoint + "/v1/metrics"
	req, err := http.NewRequest("POST", testEndpoint, nil)
	if err != nil {
		reporter.AddError(fmt.Sprintf("Error while testing credentials of OTEL_EXPORTER_OTLP_ENDPOINT: %s", err))
		return
	}

	// Extract auth value from headers
	authValue := ""
	for _, h := range strings.SplitN(headers, ",", -1) {
		key, value, _ := strings.Cut(h, "=")
		if key == "Authorization" {
			authValue = value
		}
	}
	req.Header.Set("Authorization", authValue)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		reporter.AddError(fmt.Sprintf("Error while testing credentials of OTEL_EXPORTER_OTLP_ENDPOINT: %s", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		reporter.AddError(fmt.Sprintf("Error while testing credentials of OTEL_EXPORTER_OTLP_ENDPOINT: %s", resp.Status))
	} else {
		reporter.AddSuccessfulCheck("Credentials for OTEL_EXPORTER_OTLP_ENDPOINT are correct")
	}
}
