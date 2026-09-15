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

	// Env-var substitution happens in Load (before yaml unmarshal) so
	// SignalEndpoints returns the fully-resolved URLs. Unset vars fall
	// through to the `:-` default.
	endpoints := f.SignalEndpoints()
	assert.Equal(t, []string{
		"http://localhost:4318/v1/traces",
		"https://otlp-gateway-prod-us-east-0.grafana.net/otlp/v1/traces",
	}, endpoints["traces"])
	assert.Equal(t, []string{
		"http://localhost:4318/v1/metrics",
	}, endpoints["metrics"])
	assert.Equal(t, []string{
		"http://localhost:4318/v1/logs",
	}, endpoints["logs"])
}

func TestSignalEndpoints_NilProviders(t *testing.T) {
	f := &File{FileFormat: "1.1"}
	got := f.SignalEndpoints()
	assert.Empty(t, got["traces"])
	assert.Empty(t, got["metrics"])
	assert.Empty(t, got["logs"])
}

// TestLoad_UpstreamSDKConfigExample loads the official
// otel-sdk-config.yaml example and asserts
// the parser accepts every field the file uses.
//
// If a future spec revision adds a field that our lenient YAML
// decoding can't tolerate, this test catches it before it reaches
// customer configs. Refresh the fixture with:
//
//	curl -sL https://raw.githubusercontent.com/open-telemetry/opentelemetry-configuration/main/examples/otel-sdk-config.yaml > testdata/otel-sdk-config.yaml
func TestLoad_UpstreamSDKConfigExample(t *testing.T) {
	f, err := Load("testdata/otel-sdk-config.yaml")
	require.NoError(t, err)
	require.NotNil(t, f)

	assert.Equal(t, "1.1", f.FileFormat)

	endpoints := f.SignalEndpoints()
	assert.Equal(t, []string{"http://localhost:4318/v1/traces"}, endpoints["traces"])
	assert.Equal(t, []string{"http://localhost:4318/v1/metrics"}, endpoints["metrics"])
	assert.Equal(t, []string{"http://localhost:4318/v1/logs"}, endpoints["logs"])

	assert.Equal(t, map[string]string{"service.name": "unknown_service"}, f.ResourceAttributes())
}

func TestResourceAttributes(t *testing.T) {
	// AttributesList is typed as ResourceAttributesList (*string in the
	// generated schema); tests use strPtr to satisfy the wrapper type.
	strPtr := func(s string) *string { return &s }

	t.Run("nil file returns empty map", func(t *testing.T) {
		var f *File
		assert.Empty(t, f.ResourceAttributes())
	})

	t.Run("nil resource returns empty map", func(t *testing.T) {
		f := &File{}
		assert.Empty(t, f.ResourceAttributes())
	})

	t.Run("attributes only", func(t *testing.T) {
		f := &File{Resource: &Resource{
			Attributes: []AttributeNameValue{
				{Name: "service.name", Value: "checkout"},
				{Name: "service.version", Value: "1.4.2"},
			},
		}}
		assert.Equal(t, map[string]string{
			"service.name":    "checkout",
			"service.version": "1.4.2",
		}, f.ResourceAttributes())
	})

	t.Run("attributes_list only", func(t *testing.T) {
		f := &File{Resource: &Resource{
			AttributesList: strPtr("service.name=shop,deployment.environment.name=prod"),
		}}
		assert.Equal(t, map[string]string{
			"service.name":                "shop",
			"deployment.environment.name": "prod",
		}, f.ResourceAttributes())
	})

	t.Run("attributes override attributes_list per spec", func(t *testing.T) {
		f := &File{Resource: &Resource{
			AttributesList: strPtr("service.name=from-list,service.version=1.0"),
			Attributes: []AttributeNameValue{
				{Name: "service.name", Value: "from-attributes"},
			},
		}}
		got := f.ResourceAttributes()
		assert.Equal(t, "from-attributes", got["service.name"])
		assert.Equal(t, "1.0", got["service.version"])
	})

	t.Run("env-var substitution resolved via process env", func(t *testing.T) {
		t.Setenv("MY_SERVICE", "billing")
		f := &File{Resource: &Resource{
			Attributes: []AttributeNameValue{
				{Name: "service.name", Value: "${MY_SERVICE}"},
			},
		}}
		assert.Equal(t, "billing", f.ResourceAttributes()["service.name"])
	})

	t.Run("env-var default used when var unset", func(t *testing.T) {
		f := &File{Resource: &Resource{
			Attributes: []AttributeNameValue{
				{Name: "service.name", Value: "${MISSING_VAR:-fallback}"},
			},
		}}
		assert.Equal(t, "fallback", f.ResourceAttributes()["service.name"])
	})

	t.Run("attributes_list env-var expansion", func(t *testing.T) {
		t.Setenv("EXTRA_ATTRS", "region=eu-west-2,tier=paid")
		f := &File{Resource: &Resource{
			AttributesList: strPtr("${EXTRA_ATTRS}"),
		}}
		got := f.ResourceAttributes()
		assert.Equal(t, "eu-west-2", got["region"])
		assert.Equal(t, "paid", got["tier"])
	})

	t.Run("unresolved env var leaves placeholder as literal — dropped as empty when only content", func(t *testing.T) {
		// attributes_list references an unset env var with no default —
		// ExpandEnv leaves the ${...} placeholder, which is then parsed as
		// key-value pairs and produces no valid entries.
		f := &File{Resource: &Resource{
			AttributesList: strPtr("${UNSET_VAR_ATTR_LIST}"),
		}}
		assert.Empty(t, f.ResourceAttributes())
	})

	t.Run("nil attribute value skipped per spec", func(t *testing.T) {
		f := &File{Resource: &Resource{
			Attributes: []AttributeNameValue{
				{Name: "service.name", Value: nil},
				{Name: "service.version", Value: "1.0"},
			},
		}}
		got := f.ResourceAttributes()
		assert.Equal(t, map[string]string{"service.version": "1.0"}, got)
	})

	t.Run("non-string attribute value stringified", func(t *testing.T) {
		// AttributeNameValue.value can be number/bool per schema —
		// we stringify for presence checks.
		f := &File{Resource: &Resource{
			Attributes: []AttributeNameValue{
				{Name: "service.version", Value: 2},
				{Name: "sampling.enabled", Value: true},
			},
		}}
		got := f.ResourceAttributes()
		assert.Equal(t, "2", got["service.version"])
		assert.Equal(t, "true", got["sampling.enabled"])
	})
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
