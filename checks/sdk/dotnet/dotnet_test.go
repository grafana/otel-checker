package dotnet

import (
	"os"
	"testing"

	"otel-checker/checks/utils"
)

func TestCheckDotNetAutoInstrumentation(t *testing.T) {
	tests := []struct {
		name           string
		envVars        map[string]string
		expectedErrors []string
		expectedChecks []string
	}{
		{
			name: "all required env vars set correctly",
			envVars: map[string]string{
				"CORECLR_ENABLE_PROFILING": "1",
				"CORECLR_PROFILER":         "{918728DD-259F-4A6A-AC2B-B85E1B658318}",
				"CORECLR_PROFILER_PATH":    "/path/to/profiler",
				"OTEL_DOTNET_AUTO_HOME":    "/path/to/auto",
			},
			expectedErrors: []string{},
			expectedChecks: []string{
				"All required environment variables for .NET auto-instrumentation are set with correct values.",
			},
		},
		{
			name: "missing required env vars",
			envVars: map[string]string{
				"CORECLR_ENABLE_PROFILING": "1",
			},
			expectedErrors: []string{
				"Missing required environment variables: CORECLR_PROFILER, CORECLR_PROFILER_PATH, OTEL_DOTNET_AUTO_HOME",
			},
			expectedChecks: []string{},
		},
		{
			name: "incorrect values for env vars",
			envVars: map[string]string{
				"CORECLR_ENABLE_PROFILING": "0",
				"CORECLR_PROFILER":         "wrong-guid",
				"CORECLR_PROFILER_PATH":    "/path/to/profiler",
				"OTEL_DOTNET_AUTO_HOME":    "/path/to/auto",
			},
			expectedErrors: []string{
				"Incorrect values for environment variables: CORECLR_ENABLE_PROFILING: 0, CORECLR_PROFILER: wrong-guid",
			},
			expectedChecks: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables for test
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				// Clean up environment variables after test
				for k := range tt.envVars {
					os.Unsetenv(k)
				}
			}()

			reporter := &utils.ComponentReporter{}
			checkDotNetAutoInstrumentation(reporter)

			// Check errors
			if len(reporter.Errors) != len(tt.expectedErrors) {
				t.Errorf("expected %d errors, got %d", len(tt.expectedErrors), len(reporter.Errors))
			}
			for i, err := range reporter.Errors {
				if err != tt.expectedErrors[i] {
					t.Errorf("error %d: expected %q, got %q", i, tt.expectedErrors[i], err)
				}
			}

			// Check successful checks
			if len(reporter.Checks) != len(tt.expectedChecks) {
				t.Errorf("expected %d checks, got %d", len(tt.expectedChecks), len(reporter.Checks))
			}
			for i, check := range reporter.Checks {
				if check != tt.expectedChecks[i] {
					t.Errorf("check %d: expected %q, got %q", i, tt.expectedChecks[i], check)
				}
			}
		})
	}
}
