package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandEnv(t *testing.T) {
	t.Setenv("EXISTS", "hello")
	t.Setenv("EMPTY", "")

	tests := []struct {
		name           string
		input          string
		wantResult     string
		wantUnresolved []string
	}{
		{
			name:       "no placeholder",
			input:      "http://otlp.example.com/v1/traces",
			wantResult: "http://otlp.example.com/v1/traces",
		},
		{
			name:       "plain VAR resolves",
			input:      "${EXISTS}/tail",
			wantResult: "hello/tail",
		},
		{
			name:       "env: prefixed VAR resolves",
			input:      "${env:EXISTS}/tail",
			wantResult: "hello/tail",
		},
		{
			name:       "VAR with default, VAR set — uses VAR",
			input:      "${EXISTS:-fallback}/tail",
			wantResult: "hello/tail",
		},
		{
			name:       "VAR with default, VAR unset — uses default",
			input:      "${MISSING:-http://localhost:4318}/v1/traces",
			wantResult: "http://localhost:4318/v1/traces",
		},
		{
			name:       "empty VAR treated as unset — uses default",
			input:      "${EMPTY:-fallback}/tail",
			wantResult: "fallback/tail",
		},
		{
			name:           "VAR unset with no default — leaves placeholder and reports",
			input:          "${MISSING}/tail",
			wantResult:     "${MISSING}/tail",
			wantUnresolved: []string{"MISSING"},
		},
		{
			name:           "multiple placeholders, one unresolved",
			input:          "${EXISTS}/${MISSING}/tail",
			wantResult:     "hello/${MISSING}/tail",
			wantUnresolved: []string{"MISSING"},
		},
		{
			name:       "default containing forward slashes",
			input:      "${MISSING:-http://localhost:4318}/v1/traces",
			wantResult: "http://localhost:4318/v1/traces",
		},
		{
			name:       "empty string input",
			input:      "",
			wantResult: "",
		},
		{
			name:       "$$ escape produces literal $",
			input:      "$$",
			wantResult: "$",
		},
		{
			name:       "$${VAR} escapes the substitution reference",
			input:      "$${EXISTS}",
			wantResult: "${EXISTS}",
		},
		{
			name:       "$$$ then substitution — escape then expand",
			input:      "$$${EXISTS}",
			wantResult: "$hello",
		},
		{
			name:       "$$ inside a larger string",
			input:      "cost: 5$$",
			wantResult: "cost: 5$",
		},
		{
			name:       "escape does not disable a nearby substitution",
			input:      "${EXISTS}/$${LITERAL}/${EXISTS}",
			wantResult: "hello/${LITERAL}/hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, unresolved := ExpandEnv(tt.input)
			assert.Equal(t, tt.wantResult, got)
			assert.Equal(t, tt.wantUnresolved, unresolved)
		})
	}
}

func TestLoadAndSignalEndpoints(t *testing.T) {
	yaml := `file_format: "1.1"

tracer_provider:
  processors:
    - batch:
        exporter:
          otlp_http:
            endpoint: ${OTEL_EXPORTER_OTLP_ENDPOINT:-http://localhost:4318}/v1/traces
    - simple:
        exporter:
          otlp_http:
            endpoint: https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/traces

meter_provider:
  readers:
    - periodic:
        exporter:
          otlp_http:
            endpoint: ${OTEL_EXPORTER_OTLP_ENDPOINT:-http://localhost:4318}/v1/metrics

logger_provider:
  processors:
    - batch:
        exporter:
          otlp_http:
            endpoint: ${OTEL_EXPORTER_OTLP_ENDPOINT:-http://localhost:4318}/v1/logs
`
	path := filepath.Join(t.TempDir(), "otel-config.yaml")
	require.NoError(t, os.WriteFile(path, []byte(yaml), 0o600))

	f, err := Load(path)
	require.NoError(t, err)
	require.NotNil(t, f)
	assert.Equal(t, "1.1", f.FileFormat)

	endpoints := f.SignalEndpoints()
	assert.Equal(t, []string{
		"${OTEL_EXPORTER_OTLP_ENDPOINT:-http://localhost:4318}/v1/traces",
		"https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/traces",
	}, endpoints["traces"])
	assert.Equal(t, []string{
		"${OTEL_EXPORTER_OTLP_ENDPOINT:-http://localhost:4318}/v1/metrics",
	}, endpoints["metrics"])
	assert.Equal(t, []string{
		"${OTEL_EXPORTER_OTLP_ENDPOINT:-http://localhost:4318}/v1/logs",
	}, endpoints["logs"])
}

func TestSignalEndpoints_NilProviders(t *testing.T) {
	f := &File{FileFormat: "1.1"}
	got := f.SignalEndpoints()
	assert.Empty(t, got["traces"])
	assert.Empty(t, got["metrics"])
	assert.Empty(t, got["logs"])
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "does-not-exist.yaml"))
	require.Error(t, err)
}

func TestLoad_InvalidYAML(t *testing.T) {
	path := filepath.Join(t.TempDir(), "broken.yaml")
	require.NoError(t, os.WriteFile(path, []byte("tracer_provider:\n  processors: [not-a-map]"), 0o600))
	_, err := Load(path)
	require.Error(t, err)
}

func TestResolve(t *testing.T) {
	t.Run("explicit path returned as-is", func(t *testing.T) {
		got, cands := Resolve("/tmp/whatever.yaml")
		assert.Equal(t, "/tmp/whatever.yaml", got)
		assert.Equal(t, []string{"/tmp/whatever.yaml"}, cands)
	})

	t.Run("no default files present", func(t *testing.T) {
		dir := t.TempDir()
		wd, _ := os.Getwd()
		t.Cleanup(func() { _ = os.Chdir(wd) })
		require.NoError(t, os.Chdir(dir))

		got, cands := Resolve("")
		assert.Empty(t, got)
		assert.Equal(t, DefaultCandidatePaths, cands)
	})

	t.Run("picks first default that exists", func(t *testing.T) {
		dir := t.TempDir()
		wd, _ := os.Getwd()
		t.Cleanup(func() { _ = os.Chdir(wd) })
		require.NoError(t, os.Chdir(dir))
		require.NoError(t, os.WriteFile("otel-config.yml", []byte(""), 0o600))

		got, _ := Resolve("")
		assert.Equal(t, "otel-config.yml", got)
	})
}
