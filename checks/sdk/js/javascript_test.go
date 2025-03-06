package js

import (
	"os"
	"testing"

	"otel-checker/checks/utils"
)

func TestCheckEnvVars(t *testing.T) {
	tests := []struct {
		name             string
		envVars          map[string]string
		expectedErrors   []string
		expectedChecks   []string
		expectedWarnings []string
	}{
		{
			name: "all recommended env vars set correctly",
			envVars: map[string]string{
				"OTEL_NODE_RESOURCE_DETECTORS": "env,host,os,serviceinstance",
			},
			expectedErrors: []string{},
			expectedChecks: []string{
				"OTEL_NODE_RESOURCE_DETECTORS has recommended values",
			},
			expectedWarnings: []string{},
		},
		{
			name:           "missing recommended env vars",
			envVars:        map[string]string{},
			expectedErrors: []string{},
			expectedChecks: []string{},
			expectedWarnings: []string{
				"It's recommended the environment variable OTEL_NODE_RESOURCE_DETECTORS to be set to at least `env,host,os,serviceinstance`",
			},
		},
		{
			name: "incomplete resource detectors",
			envVars: map[string]string{
				"OTEL_NODE_RESOURCE_DETECTORS": "env,host",
			},
			expectedErrors: []string{},
			expectedChecks: []string{},
			expectedWarnings: []string{
				"It's recommended the environment variable OTEL_NODE_RESOURCE_DETECTORS to be set to at least `env,host,os,serviceinstance`",
			},
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
			checkEnvVars(reporter)

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

			// Check warnings
			if len(reporter.Warnings) != len(tt.expectedWarnings) {
				t.Errorf("expected %d warnings, got %d", len(tt.expectedWarnings), len(reporter.Warnings))
			}
			for i, warning := range reporter.Warnings {
				if warning != tt.expectedWarnings[i] {
					t.Errorf("warning %d: expected %q, got %q", i, tt.expectedWarnings[i], warning)
				}
			}
		})
	}
}

func TestCheckJSAutoInstrumentation(t *testing.T) {
	tests := []struct {
		name             string
		envVars          map[string]string
		expectedErrors   []string
		expectedChecks   []string
		expectedWarnings []string
	}{
		{
			name: "NODE_OPTIONS set correctly",
			envVars: map[string]string{
				"NODE_OPTIONS": "--require @opentelemetry/auto-instrumentations-node/register",
			},
			expectedErrors: []string{},
			expectedChecks: []string{
				"NODE_OPTIONS set correctly",
			},
			expectedWarnings: []string{},
		},
		{
			name:           "NODE_OPTIONS not set",
			envVars:        map[string]string{},
			expectedErrors: []string{},
			expectedChecks: []string{},
			expectedWarnings: []string{
				"NODE_OPTIONS not set. You can set it by running 'export NODE_OPTIONS=\"--require @opentelemetry/auto-instrumentations-node/register\"' or add the same '--require ...' when starting your application",
			},
		},
		{
			name: "NODE_OPTIONS set incorrectly",
			envVars: map[string]string{
				"NODE_OPTIONS": "--require something-else",
			},
			expectedErrors: []string{},
			expectedChecks: []string{},
			expectedWarnings: []string{
				"NODE_OPTIONS not set. You can set it by running 'export NODE_OPTIONS=\"--require @opentelemetry/auto-instrumentations-node/register\"' or add the same '--require ...' when starting your application",
			},
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
			checkJSAutoInstrumentation(reporter, "")

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

			// Check warnings
			if len(reporter.Warnings) != len(tt.expectedWarnings) {
				t.Errorf("expected %d warnings, got %d", len(tt.expectedWarnings), len(reporter.Warnings))
			}
			for i, warning := range reporter.Warnings {
				if warning != tt.expectedWarnings[i] {
					t.Errorf("warning %d: expected %q, got %q", i, tt.expectedWarnings[i], warning)
				}
			}
		})
	}
}
