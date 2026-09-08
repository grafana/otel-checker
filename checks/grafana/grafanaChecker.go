package grafana

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/grafana/otel-checker/checks/config"
	"github.com/grafana/otel-checker/checks/env"
	"github.com/grafana/otel-checker/checks/utils"
)

var credentialCheckClient = &http.Client{Timeout: 10 * time.Second}

var (
	baseEndpointRegex   = regexp.MustCompile(`^https://.+\.grafana\.net/otlp/?$`)
	signalEndpointRegex = map[string]*regexp.Regexp{
		"traces":  regexp.MustCompile(`^https://.+\.grafana\.net/otlp/v1/traces/?$`),
		"metrics": regexp.MustCompile(`^https://.+\.grafana\.net/otlp/v1/metrics/?$`),
		"logs":    regexp.MustCompile(`^https://.+\.grafana\.net/otlp/v1/logs/?$`),
	}
)

var (
	OtelExporterOTLPProtocol = env.EnvVar{
		Name:          "OTEL_EXPORTER_OTLP_PROTOCOL",
		RequiredValue: "http/protobuf",
		Description:   "Protocol for OTLP exporter",
		Message:       "OTEL_EXPORTER_OTLP_PROTOCOL must be set to 'http/protobuf'",
		ExplainID:     "grafana-cloud.protocol.invalid",
	}

	// OtelExporterOTLPEndpoint is only kept as a var so checkAuth and other
	// callers can look up its value by name. Endpoint validation happens in
	// checkEndpoints — it's context-dependent on the signal-specific vars.
	OtelExporterOTLPEndpoint = env.EnvVar{
		Name:        "OTEL_EXPORTER_OTLP_ENDPOINT",
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
				reporter.AddErrorWithExplain("grafana-cloud.headers.missing-auth",
					fmt.Sprintf("OTEL_EXPORTER_OTLP_HEADERS is not set. Value should have '%s...'", tokenStart))
			}
		},
		Description: "OTLP exporter headers",
	}
)

// CheckGrafanaSetup runs the Grafana Cloud connectivity checks. When
// parsedConfig is non-nil, endpoints come from the declarative config
// file and the environment-variable-based endpoint / header / credential
// checks are skipped — the two configuration paths are mutually
// exclusive per OpenTelemetry declarative-configuration semantics.
func CheckGrafanaSetup(ctx context.Context, reporter utils.Reporter, grafanaReporter *utils.ComponentReporter, commands utils.Commands, parsedConfig *config.File) {
	if parsedConfig != nil {
		CheckEndpointsFromConfig(grafanaReporter, parsedConfig)
		grafanaReporter.AddWarningWithExplain("grafana-cloud.credentials.skipped",
			"Credentials and headers not checked — endpoints were read from the declarative config file")
		return
	}
	checkEnvVarsGrafana(reporter, grafanaReporter, commands.Language, commands.Components)
	checkAuth(ctx, grafanaReporter)
}

func checkEnvVarsGrafana(reporter utils.Reporter, grafana *utils.ComponentReporter, language string, components []string) {
	env.CheckEnvVars(grafana, language,
		OtelExporterOTLPProtocol,
		OtelExporterOTLPHeaders)
	checkEndpoints(grafana)
}

// CheckEndpointsFromConfig validates each OTLP HTTP endpoint declared in
// the declarative config, applying the same Grafana Cloud regex as the
// env-var flow. Endpoints undergo env-var substitution first so
// ${VAR:-default} placeholders resolve to concrete URLs.
//
// Exported so both check paths can call it: `check grafana-cloud` (via
// CheckGrafanaSetup) and `check config` (dispatched directly from
// checks.Run so a config-only invocation still validates endpoints,
// which is grafana's responsibility).
func CheckEndpointsFromConfig(reporter *utils.ComponentReporter, f *config.File) {
	endpoints := f.SignalEndpoints()
	for _, signal := range []string{"traces", "metrics", "logs"} {
		list := endpoints[signal]
		if len(list) == 0 {
			reporter.AddWarningWithExplain("config.endpoint.missing",
				fmt.Sprintf("No otlp_http endpoint declared for %s in the config file", signal))
			continue
		}
		re := signalEndpointRegex[signal]
		for _, raw := range list {
			expanded, unresolved := config.ExpandEnv(raw)
			if len(unresolved) > 0 {
				reporter.AddErrorWithExplain("config.env-var.unresolved",
					fmt.Sprintf("%s endpoint references unset environment variable(s) with no default: %s (raw: %q)",
						signal, strings.Join(unresolved, ", "), raw))
				continue
			}
			describe := fmt.Sprintf("%s endpoint (%q)", signal, raw)
			switch {
			case re.MatchString(expanded):
				reporter.AddSuccessfulCheck(fmt.Sprintf("%s set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/%s", describe, signal))
			case strings.Contains(expanded, "localhost"):
				reporter.AddWarningWithExplain("grafana-cloud.endpoint.localhost",
					fmt.Sprintf("%s resolves to a localhost URL (%q). Update to a Grafana endpoint similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/%s", describe, expanded, signal))
			case isValidURL(expanded):
				reporter.AddWarningWithExplain("grafana-cloud.endpoint.not-grafana",
					fmt.Sprintf("%s resolves to a valid URL (%q) but is not a Grafana Cloud endpoint. Data will not reach Grafana Cloud", describe, expanded))
			default:
				reporter.AddErrorWithExplain("grafana-cloud.signal-endpoint.invalid-format",
					fmt.Sprintf("%s resolves to %q, not a valid URL matching https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/%s", describe, expanded, signal))
			}
		}
	}
}

// isValidURL is true when s parses as a URL with both a scheme and a
// host — matches the collector checker's helper of the same name.
func isValidURL(s string) bool {
	if s == "" {
		return false
	}
	u, err := url.Parse(s)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}

// checkEndpoints validates the OTLP endpoint variables. Each signal (traces,
// metrics, logs) can be pointed at Grafana Cloud in one of two ways: via the
// signal-specific `OTEL_EXPORTER_OTLP_<SIGNAL>_ENDPOINT` (must include the
// `/v1/<signal>` path) or via the base `OTEL_EXPORTER_OTLP_ENDPOINT` (must
// NOT include the signal path — the SDK appends it). The base is required
// only for signals that don't have a signal-specific override.
func checkEndpoints(reporter *utils.ComponentReporter) {
	base := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")

	signals := []string{"traces", "metrics", "logs"}
	var missing []string
	for _, signal := range signals {
		varName := fmt.Sprintf("OTEL_EXPORTER_OTLP_%s_ENDPOINT", strings.ToUpper(signal))
		value := os.Getenv(varName)
		if value == "" {
			missing = append(missing, signal)
			continue
		}
		if signalEndpointRegex[signal].MatchString(value) {
			reporter.AddSuccessfulCheck(fmt.Sprintf("%s set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/%s", varName, signal))
		} else {
			reporter.AddErrorWithExplain("grafana-cloud.signal-endpoint.invalid-format",
				fmt.Sprintf("%s is not set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/%s", varName, signal))
		}
	}

	baseRequired := len(missing) > 0

	if base == "" {
		if baseRequired {
			reporter.AddErrorWithExplain("grafana-cloud.endpoint.unset",
				fmt.Sprintf("OTEL_EXPORTER_OTLP_ENDPOINT is not set — required because signal-specific endpoint(s) missing for: %s", strings.Join(missing, ", ")))
		} else {
			reporter.AddSuccessfulCheck("OTEL_EXPORTER_OTLP_ENDPOINT is unset — all signals are covered by signal-specific endpoints")
		}
		return
	}

	if baseEndpointRegex.MatchString(base) {
		reporter.AddSuccessfulCheck("OTEL_EXPORTER_OTLP_ENDPOINT set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp")
	} else if strings.Contains(base, "localhost") {
		reporter.AddWarningWithExplain("grafana-cloud.endpoint.localhost",
			"OTEL_EXPORTER_OTLP_ENDPOINT is set to localhost. Update to a Grafana endpoint similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp to be able to send telemetry to your Grafana Cloud instance")
	} else {
		reporter.AddErrorWithExplain("grafana-cloud.endpoint.invalid-format",
			"OTEL_EXPORTER_OTLP_ENDPOINT is not set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp (no signal suffix like /v1/traces)")
	}
}

func checkAuth(ctx context.Context, reporter *utils.ComponentReporter) {
	endpoint := env.GetValue(OtelExporterOTLPEndpoint)
	if strings.Contains(endpoint, "localhost") {
		reporter.AddWarningWithExplain("grafana-cloud.credentials.skipped",
			"Credentials not checked, since OTEL_EXPORTER_OTLP_ENDPOINT is using localhost")
		return
	}

	headers := env.GetValue(OtelExporterOTLPHeaders)
	if endpoint == "" || headers == "" {
		reporter.AddWarningWithExplain("grafana-cloud.credentials.skipped",
			"Credentials not checked, since both environment variables OTEL_EXPORTER_OTLP_ENDPOINT and OTEL_EXPORTER_OTLP_HEADERS need to be set for this check")
		return
	}

	// Test credentials
	testEndpoint := endpoint + "/v1/metrics"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, testEndpoint, nil)
	if err != nil {
		reporter.AddErrorWithExplain("grafana-cloud.credentials.network-error",
			fmt.Sprintf("Error while testing credentials of OTEL_EXPORTER_OTLP_ENDPOINT: %s", err))
		return
	}

	// Extract auth value from headers
	authValue := ""
	for _, h := range strings.Split(headers, ",") {
		key, value, _ := strings.Cut(h, "=")
		if key == "Authorization" {
			authValue = value
		}
	}
	req.Header.Set("Authorization", authValue)

	resp, err := credentialCheckClient.Do(req)
	if err != nil {
		reporter.AddErrorWithExplain("grafana-cloud.credentials.network-error",
			fmt.Sprintf("Error while testing credentials of OTEL_EXPORTER_OTLP_ENDPOINT: %s", err))
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode == 401 {
		reporter.AddErrorWithExplain("grafana-cloud.credentials.unauthorized",
			fmt.Sprintf("Error while testing credentials of OTEL_EXPORTER_OTLP_ENDPOINT: %s", resp.Status))
	} else {
		reporter.AddSuccessfulCheck("Credentials for OTEL_EXPORTER_OTLP_ENDPOINT are correct")
	}
}
