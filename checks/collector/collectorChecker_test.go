package collector

import (
	"os"
	"otel-checker/checks/utils"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckCollectorConfig(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "collector-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

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
      grpc: ""
      http: ""
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
				"collector: Value of exporter > otlphttp > endpoint on config.yaml set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
				"collector: Value of service > pipelines > traces > exporters on config.yaml contains otlphttp",
				"collector: Value of service > pipelines > traces > receivers on config.yaml contains otlp",
				"collector: Value of service > pipelines > logs > exporters on config.yaml contains otlphttp",
				"collector: Value of service > pipelines > logs > receivers on config.yaml contains otlp",
				"collector: Value of service > pipelines > metrics > exporters on config.yaml contains otlphttp",
				"collector: Value of service > pipelines > metrics > receivers on config.yaml contains otlp",
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
				"collector: Value of exporter > otlphttp > endpoint on config.yaml is set to localhost. Update to a Grafana endpoint similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp to be able to send telemetry to your Grafana Cloud instance",
			},
			expectedChecks: []string{
				"collector: Value of service > pipelines > traces > exporters on config.yaml contains otlphttp",
				"collector: Value of service > pipelines > traces > receivers on config.yaml contains otlp",
				"collector: Value of service > pipelines > logs > exporters on config.yaml contains otlphttp",
				"collector: Value of service > pipelines > logs > receivers on config.yaml contains otlp",
				"collector: Value of service > pipelines > metrics > exporters on config.yaml contains otlphttp",
				"collector: Value of service > pipelines > metrics > receivers on config.yaml contains otlp",
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
				"collector: Value of exporter > otlphttp > endpoint on config.yaml is not set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
			},
			expectedWarnings: []string{},
			expectedChecks: []string{
				"collector: Value of service > pipelines > traces > exporters on config.yaml contains otlphttp",
				"collector: Value of service > pipelines > traces > receivers on config.yaml contains otlp",
				"collector: Value of service > pipelines > logs > exporters on config.yaml contains otlphttp",
				"collector: Value of service > pipelines > logs > receivers on config.yaml contains otlp",
				"collector: Value of service > pipelines > metrics > exporters on config.yaml contains otlphttp",
				"collector: Value of service > pipelines > metrics > receivers on config.yaml contains otlp",
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
				"collector: The value of receivers > otlp > protocols > http is nil. Make sure the key exists on your config.yaml",
			},
			expectedChecks: []string{
				"collector: Value of exporter > otlphttp > endpoint on config.yaml set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
				"collector: Value of service > pipelines > traces > exporters on config.yaml contains otlphttp",
				"collector: Value of service > pipelines > traces > receivers on config.yaml contains otlp",
				"collector: Value of service > pipelines > logs > exporters on config.yaml contains otlphttp",
				"collector: Value of service > pipelines > logs > receivers on config.yaml contains otlp",
				"collector: Value of service > pipelines > metrics > exporters on config.yaml contains otlphttp",
				"collector: Value of service > pipelines > metrics > receivers on config.yaml contains otlp",
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
				"collector: Value of service > pipelines > traces > exporters on config.yaml does not contain otlphttp",
			},
			expectedChecks: []string{
				"collector: Value of exporter > otlphttp > endpoint on config.yaml set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
				"collector: Value of service > pipelines > traces > receivers on config.yaml contains otlp",
				"collector: Value of service > pipelines > logs > exporters on config.yaml contains otlphttp",
				"collector: Value of service > pipelines > logs > receivers on config.yaml contains otlp",
				"collector: Value of service > pipelines > metrics > exporters on config.yaml contains otlphttp",
				"collector: Value of service > pipelines > metrics > receivers on config.yaml contains otlp",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary config.yaml file
			configPath := tmpDir + "/"
			configFile := filepath.Join(tmpDir, "config.yaml")
			err := os.WriteFile(configFile, []byte(tt.configYAML), 0644)
			if err != nil {
				t.Fatalf("Failed to write config file: %v", err)
			}

			// Create a new reporter for testing
			reporter := utils.Reporter{}
			componentReporter := reporter.Component("collector")

			// Call the function under test
			checkCollectorConfig(componentReporter, configPath)

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
	defer os.RemoveAll(tmpDir)

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
	configPath := tmpDir + "/"
	configFile := filepath.Join(tmpDir, "config.yaml")
	err = os.WriteFile(configFile, []byte(validConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Create a new reporter for testing
	reporter := utils.Reporter{}
	componentReporter := reporter.Component("collector")

	// Call the function under test
	CheckCollectorSetup(componentReporter, "go", configPath)

	// Expected results
	expectedChecks := []string{
		"collector: Value of exporter > otlphttp > endpoint on config.yaml set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp",
		"collector: Value of service > pipelines > traces > exporters on config.yaml contains otlphttp",
		"collector: Value of service > pipelines > traces > receivers on config.yaml contains otlp",
		"collector: Value of service > pipelines > logs > exporters on config.yaml contains otlphttp",
		"collector: Value of service > pipelines > logs > receivers on config.yaml contains otlp",
		"collector: Value of service > pipelines > metrics > exporters on config.yaml contains otlphttp",
		"collector: Value of service > pipelines > metrics > receivers on config.yaml contains otlp",
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
	defer os.RemoveAll(tmpDir)

	// Create a new reporter for testing
	reporter := utils.Reporter{}
	componentReporter := reporter.Component("collector")

	// Call the function under test with a path that doesn't have a config.yaml
	checkCollectorConfig(componentReporter, tmpDir+"/")

	// Expect an error about not being able to find the config file
	assert.Len(t, componentReporter.Errors, 1, "expected one error")
	assert.Contains(t, componentReporter.Errors[0], "Could not check file", "error should mention the missing file")
}
