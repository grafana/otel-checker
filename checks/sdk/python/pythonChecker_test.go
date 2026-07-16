package python

import (
	"context"
	"testing"

	"github.com/grafana/otel-checker/checks/sdk"
	"github.com/grafana/otel-checker/checks/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadSupportedPythonLibraries(t *testing.T) {
	libs, err := supportedLibraries(context.Background())
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

func TestParseRequirementsTxtPipCompileHashes(t *testing.T) {
	// pip-compile --generate-hashes output.
	out := `requests==2.32.3 \
    --hash=sha256:1efdf4e867d4d8ba4a9f6cf9ce07cd182c4c41de77f23814feb27ca93ca9d877 \
    --hash=sha256:5a79cfd57c5c8ad0c00c8ee71c8ea92b1c22b18d0e3e5e3466d4b3d5c2b2c8a1
click==8.1.7 \
    --hash=sha256:ae74fb96c20a0277a1d615f1e4d73c8414f5a98db8b799a7931d1582f3390c28
`
	reporter := utils.Reporter{}
	c := reporter.Component("SDK")
	deps := parseRequirementsTxt(c, out)
	assert.ElementsMatch(t, []Library{
		{Name: "requests", Version: "2.32.3"},
		{Name: "click", Version: "8.1.7"},
	}, deps)
	assert.Empty(t, c.Warnings, "hash lines should be skipped silently, not warned on")
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
