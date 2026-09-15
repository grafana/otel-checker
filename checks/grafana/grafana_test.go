package grafana

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grafana/otel-checker/checks/config"
	"github.com/grafana/otel-checker/checks/utils"
)

func TestCheckEndpointsFromConfig(t *testing.T) {
	// strPtr returns a *string; the generated schema types wrap
	// endpoints as `*string` so this is what tests use to set them.
	strPtr := func(s string) *string { return &s }

	setTraceEndpoint := func(f *config.File, endpoint string) {
		f.TracerProvider.Processors[0].Batch.Exporter.OTLPHTTP.Endpoint = strPtr(endpoint)
	}

	// Baseline file — three Grafana Cloud endpoints, one per signal, no
	// env-var substitution required.
	grafanaCloudFile := func() *config.File {
		return &config.File{
			TracerProvider: &config.TracerProvider{
				Processors: []config.SpanProcessor{{
					Batch: &config.BatchSpanProcessor{Exporter: config.SpanExporter{
						OTLPHTTP: &config.OTLPHTTPExporter{
							Endpoint: strPtr("https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/traces"),
						},
					}},
				}},
			},
			MeterProvider: &config.MeterProvider{
				Readers: []config.MetricReader{{
					Periodic: &config.PeriodicMetricReader{Exporter: config.PushMetricExporter{
						OTLPHTTP: &config.OTLPHTTPMetricExporter{
							Endpoint: strPtr("https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/metrics"),
						},
					}},
				}},
			},
			LoggerProvider: &config.LoggerProvider{
				Processors: []config.LogRecordProcessor{{
					Batch: &config.BatchLogRecordProcessor{Exporter: config.LogRecordExporter{
						OTLPHTTP: &config.OTLPHTTPExporter{
							Endpoint: strPtr("https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/logs"),
						},
					}},
				}},
			},
		}
	}

	t.Run("all three signals valid", func(t *testing.T) {
		r := &utils.ComponentReporter{}
		CheckEndpointsFromConfig(r, grafanaCloudFile())
		assert.Empty(t, r.Errors)
		assert.Empty(t, r.Warnings)
		assert.Len(t, r.Checks, 3)
	})

	t.Run("localhost endpoint warns", func(t *testing.T) {
		f := grafanaCloudFile()
		setTraceEndpoint(f, "http://localhost:4318/v1/traces")
		r := &utils.ComponentReporter{}
		CheckEndpointsFromConfig(r, f)
		assert.Empty(t, r.Errors)
		assert.Contains(t, r.Warnings[0], "localhost")
	})

	t.Run("valid non-Grafana URL warns", func(t *testing.T) {
		f := grafanaCloudFile()
		setTraceEndpoint(f, "https://otel.example.com/v1/traces")
		r := &utils.ComponentReporter{}
		CheckEndpointsFromConfig(r, f)
		assert.Empty(t, r.Errors)
		assert.Contains(t, r.Warnings[0], "not a Grafana Cloud endpoint")
	})

	t.Run("invalid URL errors", func(t *testing.T) {
		f := grafanaCloudFile()
		setTraceEndpoint(f, "not-a-url")
		r := &utils.ComponentReporter{}
		CheckEndpointsFromConfig(r, f)
		assert.NotEmpty(t, r.Errors)
	})

	t.Run("missing signal provider warns", func(t *testing.T) {
		f := grafanaCloudFile()
		f.LoggerProvider = nil
		r := &utils.ComponentReporter{}
		CheckEndpointsFromConfig(r, f)
		assert.Contains(t, r.Warnings[0], "No otlp_http endpoint declared for logs")
	})

	t.Run("env-var substitution resolves via process env", func(t *testing.T) {
		t.Setenv("MY_ENDPOINT", "https://otlp-gateway-prod-us-east-0.grafana.net/otlp")
		f := grafanaCloudFile()
		setTraceEndpoint(f, "${MY_ENDPOINT}/v1/traces")
		r := &utils.ComponentReporter{}
		CheckEndpointsFromConfig(r, f)
		assert.Empty(t, r.Errors)
		assert.Empty(t, r.Warnings)
		assert.Len(t, r.Checks, 3)
	})

	t.Run("env-var default used when var unset", func(t *testing.T) {
		f := grafanaCloudFile()
		setTraceEndpoint(f, "${MISSING_VAR:-http://localhost:4318}/v1/traces")
		r := &utils.ComponentReporter{}
		CheckEndpointsFromConfig(r, f)
		assert.Empty(t, r.Errors)
		assert.Contains(t, r.Warnings[0], "localhost")
	})

	t.Run("unresolved env var errors", func(t *testing.T) {
		f := grafanaCloudFile()
		setTraceEndpoint(f, "${OTHER_MISSING_VAR}/v1/traces")
		r := &utils.ComponentReporter{}
		CheckEndpointsFromConfig(r, f)
		assert.NotEmpty(t, r.Errors)
		assert.Contains(t, r.Errors[0], "OTHER_MISSING_VAR")
	})
}

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
