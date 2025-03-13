package dotnet

import (
	"testing"

	"github.com/grafana/otel-checker/checks/utils"
)

func TestCheckDotNetAutoInstrumentation(t *testing.T) {
	tests := []utils.EnvVarTestCase{
		{
			Name: "all required env vars set correctly",
			EnvVars: map[string]string{
				"CORECLR_ENABLE_PROFILING": "1",
				"CORECLR_PROFILER":         "{918728DD-259F-4A6A-AC2B-B85E1B658318}",
				"CORECLR_PROFILER_PATH":    "/path/to/profiler",
				"OTEL_DOTNET_AUTO_HOME":    "/path/to/auto",
			},
			Language: "csharp",
			ExpectedChecks: []string{
				"dotnet: CORECLR_ENABLE_PROFILING is set to '1'",
				"dotnet: CORECLR_PROFILER is set to '{918728DD-259F-4A6A-AC2B-B85E1B658318}'",
				"dotnet: CORECLR_PROFILER_PATH is set to '/path/to/profiler'",
				"dotnet: OTEL_DOTNET_AUTO_HOME is set to '/path/to/auto'",
			},
		},
		{
			Name: "missing required env vars",
			EnvVars: map[string]string{
				"CORECLR_ENABLE_PROFILING": "1",
			},
			Language:     "csharp",
			IgnoreChecks: true,
			ExpectedErrors: []string{
				"dotnet: CORECLR_PROFILER must be set to '{918728DD-259F-4A6A-AC2B-B85E1B658318}'",
				"dotnet: CORECLR_PROFILER_PATH is not set",
				"dotnet: OTEL_DOTNET_AUTO_HOME is not set",
			},
		},
		{
			Name: "incorrect values for env vars",
			EnvVars: map[string]string{
				"CORECLR_ENABLE_PROFILING": "0",
				"CORECLR_PROFILER":         "wrong-guid",
				"CORECLR_PROFILER_PATH":    "/path/to/profiler",
				"OTEL_DOTNET_AUTO_HOME":    "/path/to/auto",
			},
			Language:     "csharp",
			IgnoreChecks: true,
			ExpectedErrors: []string{
				"dotnet: CORECLR_ENABLE_PROFILING must be set to '1'",
				"dotnet: CORECLR_PROFILER must be set to '{918728DD-259F-4A6A-AC2B-B85E1B658318}'",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			utils.RunEnvVarComponentTest(t, tt, "dotnet",
				func(reporter utils.Reporter, c *utils.ComponentReporter, language string, components []string) {
					checkDotNetAutoInstrumentation(c)
				})
		})
	}
}
