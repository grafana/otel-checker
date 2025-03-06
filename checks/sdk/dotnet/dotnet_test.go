package dotnet

import (
	"testing"

	"otel-checker/checks/utils"
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
			Language:       "csharp",
			ExpectedChecks: []string{"dotnet: All required environment variables for .NET auto-instrumentation are set with correct values."},
		},
		{
			Name: "missing required env vars",
			EnvVars: map[string]string{
				"CORECLR_ENABLE_PROFILING": "1",
			},
			Language: "csharp",
			ExpectedErrors: []string{
				"dotnet: Missing required environment variables: CORECLR_PROFILER, CORECLR_PROFILER_PATH, OTEL_DOTNET_AUTO_HOME",
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
			Language: "csharp",
			ExpectedErrors: []string{
				"dotnet: Incorrect values for environment variables: CORECLR_ENABLE_PROFILING: 0, CORECLR_PROFILER: wrong-guid",
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
