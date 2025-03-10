package python

import (
	"otel-checker/checks/sdk"
	"otel-checker/checks/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadSupportedPythonLibraries(t *testing.T) {
	libs, err := supportedLibraries()
	require.NoError(t, err)
	assert.Equal(t,
		[]string{"https://github.com/open-telemetry/opentelemetry-python-contrib/tree/main/instrumentation/opentelemetry-instrumentation-botocore"},
		findSupportedLibraries(Library{
			Name:    "botocore",
			Version: "1.5.16",
		}, libs))
}

func TestParseRequirementsTxt(t *testing.T) {
	out := `blinker==1.9.0
	click==8.1.8
	`
	reporter := utils.Reporter{}
	deps := parseRequirementsTxt(reporter.Component("SDK"), out)
	assert.ElementsMatch(t, []Library{
		{
			Name:    "blinker",
			Version: "1.9.0",
		},
		{
			Name:    "click",
			Version: "8.1.8",
		},
	}, deps)
}

func TestVersionRanges(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected map[string]sdk.VersionRange
	}{
		{
			name:  "With space between operator and version",
			input: "library < 1.0",
			expected: map[string]sdk.VersionRange{
				"library": {
					Upper: "1.0",
				},
			},
		},
		{
			name:  "Without space between operator and version",
			input: "library <1.0",
			expected: map[string]sdk.VersionRange{
				"library": {
					Upper: "1.0",
				},
			},
		},
		{
			name:  "Multiple constraints with spaces",
			input: "library >= 1.0, < 2.0",
			expected: map[string]sdk.VersionRange{
				"library": {
					Lower:          "1.0",
					Upper:          "2.0",
					LowerInclusive: true,
				},
			},
		},
		{
			name:  "Multiple constraints without spaces",
			input: "library >=1.0, <2.0",
			expected: map[string]sdk.VersionRange{
				"library": {
					Lower:          "1.0",
					Upper:          "2.0",
					LowerInclusive: true,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := versionRanges(tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}
