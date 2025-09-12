package java

import (
	"encoding/base64"
	"github.com/grafana/otel-checker/checks/sdk/supported"
	"github.com/grafana/otel-checker/checks/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestFindSupportedLibrary(t *testing.T) {
	modules, err := supportedLibraries()
	require.NoError(t, err)
	assert.Equal(t,
		[]string{
			"https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/instrumentation/executors/javaagent",
			"https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/instrumentation/http-url-connection/javaagent",
			"https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/instrumentation/java-http-client/javaagent",
			"https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/instrumentation/java-http-server/javaagent",
			"https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/instrumentation/jdbc/javaagent",
			"https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/instrumentation/logback/logback-appender-1.0/javaagent",
			"https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/instrumentation/rmi/javaagent",
		},
		findSupportedLibraries(Library{
			Group:    "ch.qos.logback",
			Artifact: "logback-classic",
			Version:  "1.5.16",
		}, modules, supported.TypeJavaagent, 8, nil))
}

func TestExtractInstrumentationNames(t *testing.T) {
	links := []string{
		"https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/instrumentation/executors/javaagent",
		"https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/instrumentation/http-url-connection/javaagent",
		"https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/instrumentation/java-http-client/javaagent",
		"https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/instrumentation/java-http-server/javaagent",
		"https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/instrumentation/jdbc/javaagent",
		"https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/instrumentation/rmi/javaagent",
		"https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/instrumentation/kafka/kafka-clients/kafka-clients-0.11/javaagent",
		"https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/instrumentation/spring/spring-core-2.0/javaagent",
		"https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/instrumentation/apache-httpclient/apache-httpclient-5.0/javaagent",
	}
	
	expected := []string{
		"apache-httpclient-5.0",
		"executors",
		"http-url-connection",
		"java-http-client", 
		"java-http-server",
		"jdbc",
		"kafka-clients-0.11",
		"rmi",
		"spring-core-2.0",
	}
	
	result := extractInstrumentationNames(links)
	assert.Equal(t, expected, result)
}

func TestExtractInstrumentationNamesWithDuplicates(t *testing.T) {
	links := []string{
		"https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/instrumentation/executors/javaagent",
		"https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/instrumentation/executors/library",
		"https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/instrumentation/kafka/kafka-clients/kafka-clients-0.11/javaagent",
		"https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/instrumentation/kafka/kafka-clients/kafka-clients-0.11/library",
		"https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/instrumentation/jdbc/javaagent",
	}
	
	expected := []string{
		"executors",
		"jdbc",
		"kafka-clients-0.11",
	}
	
	result := extractInstrumentationNames(links)
	assert.Equal(t, expected, result)
}

func TestOutputSupportedLibrariesReturnsInstrumentations(t *testing.T) {
	modules, err := supportedLibraries()
	require.NoError(t, err)
	
	// Create mock reporter
	reporter := &utils.ComponentReporter{}
	
	// Test with a library that has known instrumentations
	deps := []Library{
		{
			Group:    "ch.qos.logback",
			Artifact: "logback-classic",
			Version:  "1.5.16",
		},
	}
	
	result := outputSupportedLibraries(deps, modules, reporter, false, false, supported.TypeJavaagent, 8)
	
	// Should return a list of instrumentations
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "executors")
	assert.Contains(t, result, "jdbc")
	assert.Contains(t, result, "rmi")
	
	// Should be sorted
	for i := 1; i < len(result); i++ {
		assert.True(t, result[i-1] <= result[i], "Result should be sorted")
	}
}

func TestGenerateInstrumentationExplorerLink(t *testing.T) {
	tests := []struct {
		name            string
		instrumentations []string
		expectedContains []string
	}{
		{
			name:            "empty list",
			instrumentations: []string{},
			expectedContains: nil, // should return empty string
		},
		{
			name:            "single instrumentation",
			instrumentations: []string{"jdbc"},
			expectedContains: []string{"https://jaydeluca.github.io/instrumentation-explorer/analyze?instrumentations=", "&version=2.19"},
		},
		{
			name:            "multiple instrumentations",
			instrumentations: []string{"executors", "jdbc", "spring"},
			expectedContains: []string{"https://jaydeluca.github.io/instrumentation-explorer/analyze?instrumentations=", "&version=2.19"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateInstrumentationExplorerLink(tt.instrumentations)
			
			if len(tt.expectedContains) == 0 {
				assert.Empty(t, result)
				return
			}
			
			// Check that the link contains expected parts
			for _, expected := range tt.expectedContains {
				assert.Contains(t, result, expected)
			}
			
			// Verify base64 encoding works correctly
			if len(tt.instrumentations) > 0 {
				// Extract the base64 part from the URL
				parts := strings.Split(result, "instrumentations=")
				assert.Len(t, parts, 2)
				base64Part := strings.Split(parts[1], "&")[0]
				
				// Decode and verify
				decoded, err := base64.StdEncoding.DecodeString(base64Part)
				require.NoError(t, err)
				
				expected := strings.Join(tt.instrumentations, ",")
				assert.Equal(t, expected, string(decoded))
			}
		})
	}
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
