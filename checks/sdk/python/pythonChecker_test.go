package python

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"otel-checker/checks/utils"
	"testing"
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

func TestIncreaseLastPart(t *testing.T) {
	part, err := upperBoundForTilde("1.4.5")
	require.NoError(t, err)
	require.Equal(t, "1.5", part)
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
