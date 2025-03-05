package sdk

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReadSupportedPythonLibraries(t *testing.T) {
	libs, err := supportedPythonLibraries()
	require.NoError(t, err)
	assert.Equal(t,
		[]string{"https://github.com/open-telemetry/opentelemetry-python-contrib/tree/main/instrumentation/opentelemetry-instrumentation-botocore"},
		findSupportedPythonLibraries(PythonLibrary{
			Name:    "botocore",
			Version: "1.5.16",
		}, libs))
}

func TestIncreaseLastPart(t *testing.T) {
	part, err := upperBoundForTilde("1.4.5")
	require.NoError(t, err)
	require.Equal(t, "1.5", part)
}
