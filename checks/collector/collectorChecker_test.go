package collector

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/grafana/otel-checker/checks/utils"

	"github.com/stretchr/testify/assert"
)

func TestCheckCollectorConfig(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "collector-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(tmpDir)

	tests := []struct {
		name             string
		configYAML       string
		expectedErrors   []string
		expectedWarnings []string
		expectedChecks   []string
	}{
		{
			name: "Valid Grafana Cloud configuration",
			configYAML: `
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318 
exporters:
  otlphttp:
    endpoint: https://otlp-gateway-prod-us-east-0.grafana.net/otlp
    auth:
      headers:
        Authorization: "Basic base64-encoded-token"
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: []
      exporters: [otlphttp]
    logs:
      receivers: [otlp]
      processors: []
      exporters: [otlphttp]
    metrics:
      receivers: [otlp]
      processors: []
      exporters: [otlphttp]
`,
			expectedErrors:   []string{},
			expectedWarnings: []string{},
			expectedChecks: []string{
				"Value of exporter > otlphttp > endpoint on config.yaml set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
				"Value of service > pipelines > traces > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > traces > receivers on config.yaml contains otlp",
				"Value of service > pipelines > logs > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > logs > receivers on config.yaml contains otlp",
				"Value of service > pipelines > metrics > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > metrics > receivers on config.yaml contains otlp",
			},
		},
		{
			name: "Valid Grafana Cloud configuration with named receiver and exporter",
			configYAML: `
receivers:
  otlp/app:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318
exporters:
  otlphttp/grafana_cloud:
    endpoint: https://otlp-gateway-prod-us-east-0.grafana.net/otlp
  otlphttp/local:
    endpoint: http://localhost:4318
service:
  pipelines:
    traces:
      receivers: [otlp/app]
      processors: []
      exporters: [otlphttp/grafana_cloud, otlphttp/local]
    logs:
      receivers: [otlp/app]
      processors: []
      exporters: [otlphttp/grafana_cloud]
    metrics:
      receivers: [otlp/app]
      processors: []
      exporters: [otlphttp/grafana_cloud]
`,
			expectedErrors:   []string{},
			expectedWarnings: []string{},
			expectedChecks: []string{
				"Value of exporter > otlphttp/grafana_cloud > endpoint on config.yaml set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
				"Value of service > pipelines > traces > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > traces > receivers on config.yaml contains otlp",
				"Value of service > pipelines > logs > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > logs > receivers on config.yaml contains otlp",
				"Value of service > pipelines > metrics > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > metrics > receivers on config.yaml contains otlp",
			},
		},
		{
			name: "Valid Grafana Cloud configuration with legacy underscore exporter type",
			configYAML: `
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318
exporters:
  otlp_http/grafana_cloud:
    endpoint: https://otlp-gateway-prod-us-east-0.grafana.net/otlp
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: []
      exporters: [otlp_http/grafana_cloud]
    logs:
      receivers: [otlp]
      processors: []
      exporters: [otlp_http/grafana_cloud]
    metrics:
      receivers: [otlp]
      processors: []
      exporters: [otlp_http/grafana_cloud]
`,
			expectedErrors:   []string{},
			expectedWarnings: []string{},
			expectedChecks: []string{
				"Value of exporter > otlp_http/grafana_cloud > endpoint on config.yaml set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
				"Value of service > pipelines > traces > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > traces > receivers on config.yaml contains otlp",
				"Value of service > pipelines > logs > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > logs > receivers on config.yaml contains otlp",
				"Value of service > pipelines > metrics > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > metrics > receivers on config.yaml contains otlp",
			},
		},
		{
			name: "Invalid Grafana Cloud endpoint with named otlp_http exporter",
			configYAML: `
receivers:
  otlp:
    protocols:
      grpc: ""
      http: ""
exporters:
  otlp_http/invalid:
    endpoint: invalid_endpoint
service:
  pipelines:
    metrics:
      receivers: [otlp]
      exporters: [otlp_http/invalid]
    logs:
      receivers: [otlp]
      exporters: [otlp_http/invalid]
    traces:
      receivers: [otlp]
      exporters: [otlp_http/invalid]
`,
			expectedErrors: []string{
				"Value of exporter > otlphttp > endpoint on config.yaml is not set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
			},
			expectedWarnings: []string{},
			expectedChecks: []string{
				"Value of service > pipelines > traces > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > traces > receivers on config.yaml contains otlp",
				"Value of service > pipelines > logs > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > logs > receivers on config.yaml contains otlp",
				"Value of service > pipelines > metrics > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > metrics > receivers on config.yaml contains otlp",
			},
		},
		{
			name: "Localhost configuration",
			configYAML: `
receivers:
  otlp:
    protocols:
      grpc: ""
      http: ""
exporters:
  otlphttp:
    endpoint: http://localhost:4318
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: []
      exporters: [otlphttp]
    logs:
      receivers: [otlp]
      processors: []
      exporters: [otlphttp]
    metrics:
      receivers: [otlp]
      processors: []
      exporters: [otlphttp]
`,
			expectedErrors: []string{},
			expectedWarnings: []string{
				"Value of exporter > otlphttp > endpoint on config.yaml is set to localhost. Update to a Grafana endpoint similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp to be able to send telemetry to your Grafana Cloud instance",
			},
			expectedChecks: []string{
				"Value of service > pipelines > traces > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > traces > receivers on config.yaml contains otlp",
				"Value of service > pipelines > logs > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > logs > receivers on config.yaml contains otlp",
				"Value of service > pipelines > metrics > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > metrics > receivers on config.yaml contains otlp",
			},
		},
		{
			name: "Invalid endpoint format",
			configYAML: `
receivers:
  otlp:
    protocols:
      grpc: ""
      http: ""
exporters:
  otlphttp:
    endpoint: http://invalid-endpoint.com
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: []
      exporters: [otlphttp]
    logs:
      receivers: [otlp]
      processors: []
      exporters: [otlphttp]
    metrics:
      receivers: [otlp]
      processors: []
      exporters: [otlphttp]
`,
			expectedErrors: []string{
				"Value of exporter > otlphttp > endpoint on config.yaml is not set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
			},
			expectedWarnings: []string{},
			expectedChecks: []string{
				"Value of service > pipelines > traces > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > traces > receivers on config.yaml contains otlp",
				"Value of service > pipelines > logs > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > logs > receivers on config.yaml contains otlp",
				"Value of service > pipelines > metrics > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > metrics > receivers on config.yaml contains otlp",
			},
		},
		{
			name: "Missing http protocol",
			configYAML: `
receivers:
  otlp:
    protocols:
      grpc: ""
exporters:
  otlphttp:
    endpoint: https://otlp-gateway-prod-us-east-0.grafana.net/otlp
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: []
      exporters: [otlphttp]
    logs:
      receivers: [otlp]
      processors: []
      exporters: [otlphttp]
    metrics:
      receivers: [otlp]
      processors: []
      exporters: [otlphttp]
`,
			expectedErrors: []string{},
			expectedWarnings: []string{
				"The value of receivers > otlp > protocols > http is nil. Make sure the key exists on your config.yaml",
			},
			expectedChecks: []string{
				"Value of exporter > otlphttp > endpoint on config.yaml set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
				"Value of service > pipelines > traces > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > traces > receivers on config.yaml contains otlp",
				"Value of service > pipelines > logs > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > logs > receivers on config.yaml contains otlp",
				"Value of service > pipelines > metrics > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > metrics > receivers on config.yaml contains otlp",
			},
		},
		{
			name: "Missing otlphttp exporter in traces pipeline",
			configYAML: `
receivers:
  otlp:
    protocols:
      grpc: ""
      http: ""
exporters:
  otlphttp:
    endpoint: https://otlp-gateway-prod-us-east-0.grafana.net/otlp
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: []
      exporters: [otlp]
    logs:
      receivers: [otlp]
      processors: []
      exporters: [otlphttp]
    metrics:
      receivers: [otlp]
      processors: []
      exporters: [otlphttp]
`,
			expectedErrors: []string{},
			expectedWarnings: []string{
				"Value of service > pipelines > traces > exporters on config.yaml does not contain otlphttp",
			},
			expectedChecks: []string{
				"Value of exporter > otlphttp > endpoint on config.yaml set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
				"Value of service > pipelines > traces > receivers on config.yaml contains otlp",
				"Value of service > pipelines > logs > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > logs > receivers on config.yaml contains otlp",
				"Value of service > pipelines > metrics > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > metrics > receivers on config.yaml contains otlp",
			},
		},
		{
			// Every pipeline uses a non-OTLP receiver — this should now warn
			// per pipeline rather than silently "succeed" as before.
			name: "Pipelines fed only by non-OTLP receivers",
			configYAML: `
receivers:
  otlp:
    protocols:
      grpc: ""
      http: ""
  prometheus:
    config:
      scrape_configs:
        - job_name: 'app'
  filelog:
    include: [/var/log/app/*.log]
exporters:
  otlphttp:
    endpoint: https://otlp-gateway-prod-us-east-0.grafana.net/otlp
service:
  pipelines:
    traces:
      receivers: [zipkin]
      processors: []
      exporters: [otlphttp]
    logs:
      receivers: [filelog]
      processors: []
      exporters: [otlphttp]
    metrics:
      receivers: [prometheus]
      processors: []
      exporters: [otlphttp]
`,
			expectedErrors: []string{},
			expectedWarnings: []string{
				"Value of service > pipelines > traces > receivers on config.yaml does not contain otlp",
				"Value of service > pipelines > logs > receivers on config.yaml does not contain otlp",
				"Value of service > pipelines > metrics > receivers on config.yaml does not contain otlp",
			},
			expectedChecks: []string{
				"Value of exporter > otlphttp > endpoint on config.yaml set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
				"Value of service > pipelines > traces > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > logs > exporters on config.yaml contains otlphttp",
				"Value of service > pipelines > metrics > exporters on config.yaml contains otlphttp",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Write the config to the temp dir and pass its full path.
			configFile := filepath.Join(tmpDir, "config.yaml")
			err := os.WriteFile(configFile, []byte(tt.configYAML), 0644)
			if err != nil {
				t.Fatalf("Failed to write config file: %v", err)
			}

			// Create a new reporter for testing
			reporter := utils.Reporter{}
			componentReporter := reporter.Component("collector")

			// Call the function under test with the full file path.
			checkCollectorConfig(componentReporter, configFile)

			// Compare the results
			assert.ElementsMatch(t, tt.expectedErrors, componentReporter.Errors, "errors mismatch")
			assert.ElementsMatch(t, tt.expectedWarnings, componentReporter.Warnings, "warnings mismatch")
			assert.ElementsMatch(t, tt.expectedChecks, componentReporter.Checks, "checks mismatch")
		})
	}
}

func TestCheckCollectorSetup(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "collector-setup-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(tmpDir)

	// Create a valid config.yaml file
	validConfig := `
receivers:
  otlp:
    protocols:
      grpc: ""
      http: ""
exporters:
  otlphttp:
    endpoint: https://otlp-gateway-prod-us-east-0.grafana.net/otlp
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: []
      exporters: [otlphttp]
    logs:
      receivers: [otlp]
      processors: []
      exporters: [otlphttp]
    metrics:
      receivers: [otlp]
      processors: []
      exporters: [otlphttp]
`
	configFile := filepath.Join(tmpDir, "config.yaml")
	err = os.WriteFile(configFile, []byte(validConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Create a new reporter for testing
	reporter := utils.Reporter{}
	componentReporter := reporter.Component("collector")

	// Call the function under test with the full file path.
	CheckCollectorSetup(componentReporter, "go", configFile)

	// Expected results
	expectedChecks := []string{
		"Value of exporter > otlphttp > endpoint on config.yaml set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
		"Value of service > pipelines > traces > exporters on config.yaml contains otlphttp",
		"Value of service > pipelines > traces > receivers on config.yaml contains otlp",
		"Value of service > pipelines > logs > exporters on config.yaml contains otlphttp",
		"Value of service > pipelines > logs > receivers on config.yaml contains otlp",
		"Value of service > pipelines > metrics > exporters on config.yaml contains otlphttp",
		"Value of service > pipelines > metrics > receivers on config.yaml contains otlp",
	}

	// Verify the results
	assert.Empty(t, componentReporter.Errors, "no errors expected")
	assert.Empty(t, componentReporter.Warnings, "no warnings expected")
	assert.ElementsMatch(t, expectedChecks, componentReporter.Checks, "checks mismatch")
}

func TestMissingConfigFile(t *testing.T) {
	// Create a temporary directory for test files (without writing a config file)
	tmpDir, err := os.MkdirTemp("", "collector-missing-config-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(tmpDir)

	// Point the checker at a specific non-existent file inside tmpDir.
	missing := filepath.Join(tmpDir, "does-not-exist.yaml")

	reporter := utils.Reporter{}
	componentReporter := reporter.Component("collector")
	checkCollectorConfig(componentReporter, missing)

	// The tried-paths list is just the explicit file the user asked for.
	assert.Len(t, componentReporter.Errors, 1, "expected one error")
	msg := componentReporter.Errors[0]
	assert.Contains(t, msg, "Could not find a Collector config file", "error should say no config was found")
	assert.Contains(t, msg, missing, "error should list the path the user provided")
}

func TestMissingConfigFileFallback(t *testing.T) {
	// With no path passed, the checker falls back to config.yaml then
	// config.yml in the current working directory. Chdir to an empty tmp
	// dir so neither exists.
	tmpDir, err := os.MkdirTemp("", "collector-fallback-missing-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(tmpDir)
	origWD, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origWD) }()

	reporter := utils.Reporter{}
	componentReporter := reporter.Component("collector")
	checkCollectorConfig(componentReporter, "")

	assert.Len(t, componentReporter.Errors, 1, "expected one error")
	msg := componentReporter.Errors[0]
	assert.Contains(t, msg, "Could not find a Collector config file")
	assert.Contains(t, msg, "config.yaml", "fallback list should include config.yaml")
	assert.Contains(t, msg, "config.yml", "fallback list should include config.yml")
}

func TestCheckCollectorConfigYmlExtension(t *testing.T) {
	// Exercise the empty-flag fallback: config.yml (short extension) in the
	// current working directory must be picked up when the user doesn't pass
	// --collector-config-path.
	tmpDir, err := os.MkdirTemp("", "collector-yml-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(tmpDir)

	config := `
receivers:
  otlp:
    protocols:
      grpc: ""
      http: ""
exporters:
  otlphttp:
    endpoint: https://otlp-gateway-prod-us-east-0.grafana.net/otlp
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: []
      exporters: [otlphttp]
    logs:
      receivers: [otlp]
      processors: []
      exporters: [otlphttp]
    metrics:
      receivers: [otlp]
      processors: []
      exporters: [otlphttp]
`
	if err := os.WriteFile(filepath.Join(tmpDir, "config.yml"), []byte(config), 0644); err != nil {
		t.Fatalf("Failed to write config.yml: %v", err)
	}

	origWD, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origWD) }()

	reporter := utils.Reporter{}
	componentReporter := reporter.Component("collector")

	checkCollectorConfig(componentReporter, "")

	assert.Empty(t, componentReporter.Errors, "no errors expected for a valid config.yml")
	assert.NotEmpty(t, componentReporter.Checks, "expected successful checks against config.yml")
}

func TestCheckCollectorConfigCustomFilename(t *testing.T) {
	// The Collector accepts any filename via --config; the checker should too
	// when the user passes a full path to --collector-config-path.
	tmpDir, err := os.MkdirTemp("", "collector-custom-name-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(tmpDir)

	config := `
receivers:
  otlp:
    protocols:
      grpc: ""
      http: ""
exporters:
  otlphttp:
    endpoint: https://otlp-gateway-prod-us-east-0.grafana.net/otlp
service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: []
      exporters: [otlphttp]
    logs:
      receivers: [otlp]
      processors: []
      exporters: [otlphttp]
    metrics:
      receivers: [otlp]
      processors: []
      exporters: [otlphttp]
`
	customPath := filepath.Join(tmpDir, "my-collector.yaml")
	if err := os.WriteFile(customPath, []byte(config), 0644); err != nil {
		t.Fatalf("Failed to write %s: %v", customPath, err)
	}

	reporter := utils.Reporter{}
	componentReporter := reporter.Component("collector")

	checkCollectorConfig(componentReporter, customPath)

	assert.Empty(t, componentReporter.Errors, "no errors expected when a valid config with a non-standard filename is passed")
	assert.NotEmpty(t, componentReporter.Checks, "expected successful checks")
}
