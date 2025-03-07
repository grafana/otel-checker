package _go

import (
	"otel-checker/checks/sdk/supported"
	"testing"
)

func TestReadGoMod(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []supported.Library
	}{
		{
			name: "parse go.mod with dependencies",
			content: `module go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc

require (
	github.com/stretchr/testify v1.10.0
	go.opentelemetry.io/otel v1.35.0
	go.opentelemetry.io/otel/metric v1.35.0
	go.opentelemetry.io/otel/trace v1.35.0
	google.golang.org/grpc v1.71.0
	google.golang.org/protobuf v1.36.5
)`,
			expected: []supported.Library{
				{Name: "github.com/stretchr/testify", Version: "v1.10.0"},
				{Name: "go.opentelemetry.io/otel", Version: "v1.35.0"},
				{Name: "go.opentelemetry.io/otel/metric", Version: "v1.35.0"},
				{Name: "go.opentelemetry.io/otel/trace", Version: "v1.35.0"},
				{Name: "google.golang.org/grpc", Version: "v1.71.0"},
				{Name: "google.golang.org/protobuf", Version: "v1.36.5"},
			},
		},
		{
			name: "parse go.mod with empty dependencies",
			content: `module go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc

require (
)`,
			expected: []supported.Library{},
		},
		{
			name: "parse go.mod with inline require statement",
			content: `module go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc

require github.com/stretchr/testify v1.10.0
require google.golang.org/grpc v1.71.0
`,
			expected: []supported.Library{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := readGoModFromContent([]byte(tt.content))
			if len(deps) != len(tt.expected) {
				t.Errorf("Expected %d dependencies, got %d", len(tt.expected), len(deps))
			}
			for _, expected := range tt.expected {
				found := findDep(deps, expected.Name)
				if found == nil {
					t.Errorf("Expected dependency %s not found", expected.Name)
				} else if found.Version != expected.Version {
					t.Errorf("Expected version %s for %s, got %s", expected.Version, expected.Name, found.Version)
				}
			}
		})
	}
}

func TestFindSupportedLibraries(t *testing.T) {
	s := supported.SupportedModules{
		"google.golang.org/grpc": {
			Instrumentations: []supported.Instrumentation{
				{
					Name: "google.golang.org/grpc",
					Link: "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc",
					TargetVersions: map[supported.InstrumentationType][]string{
						supported.TypeLibrary: {"[1.50.0,)"},
					},
				},
			},
		},
		"github.com/aws/aws-lambda-go": {
			Instrumentations: []supported.Instrumentation{
				{
					Name: "github.com/aws/aws-lambda-go",
					Link: "go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda",
					TargetVersions: map[supported.InstrumentationType][]string{
						supported.TypeLibrary: {"[1.41.0,)"},
					},
				},
			},
		},
	}

	tests := []struct {
		name     string
		library  supported.Library
		expected []string
	}{
		{
			name: "supported gRPC version",
			library: supported.Library{
				Name:    "google.golang.org/grpc",
				Version: "v1.71.0",
			},
			expected: []string{"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"},
		},
		{
			name: "unsupported gRPC version",
			library: supported.Library{
				Name:    "google.golang.org/grpc",
				Version: "v1.49.0",
			},
			expected: nil,
		},
		{
			name: "supported aws-lambda-go version",
			library: supported.Library{
				Name:    "github.com/aws/aws-lambda-go",
				Version: "v1.41.0",
			},
			expected: []string{"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"},
		},
		{
			name: "unsupported aws-lambda-go version",
			library: supported.Library{
				Name:    "github.com/aws/aws-lambda-go",
				Version: "v1.40.0",
			},
			expected: nil,
		},
		{
			name: "unknown library",
			library: supported.Library{
				Name:    "unknown/package",
				Version: "v1.0.0",
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := supported.FindSupportedLibraries(tt.library, s, supported.TypeLibrary)
			if (got == nil && tt.expected != nil) || (got != nil && tt.expected == nil) {
				t.Errorf("Expected %v, got %v", tt.expected, got)
			}
			if len(got) != len(tt.expected) {
				t.Errorf("Expected %d links, got %d", len(tt.expected), len(got))
			}
			for i, link := range tt.expected {
				if i < len(got) && got[i] != link {
					t.Errorf("Expected link %s, got %s", link, got[i])
				}
			}
		})
	}
}

// Helper function to find a dependency by name
func findDep(deps []supported.Library, name string) *supported.Library {
	for _, dep := range deps {
		if dep.Name == name {
			return &dep
		}
	}
	return nil
}
