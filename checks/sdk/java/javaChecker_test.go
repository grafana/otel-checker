package java

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"otel-checker/checks/sdk/supported"
	"testing"
)

func TestFindSupportedLibrary(t *testing.T) {
	modules, err := supportedLibraries()
	require.NoError(t, err)
	assert.Equal(t,
		[]string{"https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/instrumentation/logback/logback-appender-1.0/javaagent"},
		findSupportedLibraries(Library{
			Group:    "ch.qos.logback",
			Artifact: "logback-classic",
			Version:  "1.5.16",
		}, modules, supported.TypeJavaagent))
}

func TestParseGradleDependencies(t *testing.T) {
	out := `> Task :custom:dependencies
------------------------------------------------------------
Project ':custom'
------------------------------------------------------------

runtimeClasspath - Runtime classpath of source set 'main'.
\--- io.opentelemetry.instrumentation:opentelemetry-instrumentation-bom-alpha:2.13.3-alpha
     +--- io.opentelemetry:opentelemetry-bom:1.47.0
     +--- io.opentelemetry:opentelemetry-bom-alpha:1.47.0-alpha
     |    \--- io.opentelemetry:opentelemetry-bom:1.47.0
     \--- io.opentelemetry.instrumentation:opentelemetry-instrumentation-bom:2.13.3
          \--- io.opentelemetry:opentelemetry-bom:1.47.0

(*) - Indicates repeated occurrences of a transitive dependency subtree. Gradle expands transitive dependency subtrees only once per project; repeat occurrences only display the root of the subtree, followed by this annotation.

A web-based, searchable dependency report is available by adding the --scan option.

BUILD SUCCESSFUL in 1s
1 actionable task: 1 executed
`
	deps := parseGradleDeps(out)
	assert.ElementsMatch(t, []Library{
		{
			Group:    "io.opentelemetry",
			Artifact: "opentelemetry-bom",
			Version:  "1.47.0",
		},
		{
			Group:    "io.opentelemetry",
			Artifact: "opentelemetry-bom-alpha",
			Version:  "1.47.0-alpha",
		},
		{
			Group:    "io.opentelemetry.instrumentation",
			Artifact: "opentelemetry-instrumentation-bom-alpha",
			Version:  "2.13.3-alpha",
		},
		{
			Group:    "io.opentelemetry.instrumentation",
			Artifact: "opentelemetry-instrumentation-bom",
			Version:  "2.13.3",
		},
	}, deps)
}
