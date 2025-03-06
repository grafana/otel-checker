package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// EnvVarTestCase represents a test case for environment variable checks
type EnvVarTestCase struct {
	Name             string
	EnvVars          map[string]string
	Language         string
	Components       []string
	ExpectedErrors   []string
	ExpectedChecks   []string
	ExpectedWarnings []string
	IgnoreWarnings   bool
	IgnoreErrors     bool
	IgnoreChecks     bool
}

// RunEnvVarComponentTest runs a test case for environment variable checks with component support
func RunEnvVarComponentTest(
	t *testing.T,
	tt EnvVarTestCase,
	componentName string,
	checkFunc func(Reporter, *ComponentReporter, string, []string)) {
	// Set up environment variables for test
	for k, v := range tt.EnvVars {
		os.Setenv(k, v)
	}
	defer func() {
		// Clean up environment variables after test
		for k := range tt.EnvVars {
			os.Unsetenv(k)
		}
	}()

	reporter := Reporter{}
	component := reporter.Component(componentName)
	checkFunc(reporter, component, tt.Language, tt.Components)

	if !tt.IgnoreErrors {
		assert.Equal(t, tt.ExpectedErrors, component.Errors, "errors mismatch")
	}

	if !tt.IgnoreChecks {
		assert.Equal(t, tt.ExpectedChecks, component.Checks, "checks mismatch")
	}

	if !tt.IgnoreWarnings {
		assert.Equal(t, tt.ExpectedWarnings, component.Warnings, "warnings mismatch")
	}
}
