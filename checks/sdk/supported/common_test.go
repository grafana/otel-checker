package supported

import (
	"otel-checker/checks/utils"
	"testing"
)

func TestFindSupportedLibraries(t *testing.T) {
	s := SupportedModules{
		"test-library": {
			Instrumentations: []Instrumentation{
				{
					Name: "test-library",
					Link: "https://example.com/test-library",
					TargetVersions: map[InstrumentationType][]string{
						TypeLibrary: {"[1.0.0,)"},
					},
				},
			},
		},
		"another-library": {
			Instrumentations: []Instrumentation{
				{
					Name: "another-library",
					Link: "https://example.com/another-library",
					TargetVersions: map[InstrumentationType][]string{
						TypeLibrary: {"[2.0.0,3.0.0)"},
					},
				},
			},
		},
	}

	tests := []struct {
		name                string
		library             Library
		instrumentationType InstrumentationType
		expected            []string
	}{
		{
			name: "supported library version",
			library: Library{
				Name:    "test-library",
				Version: "1.5.0",
			},
			instrumentationType: TypeLibrary,
			expected:            []string{"https://example.com/test-library"},
		},
		{
			name: "unsupported library version",
			library: Library{
				Name:    "test-library",
				Version: "0.9.0",
			},
			instrumentationType: TypeLibrary,
			expected:            nil,
		},
		{
			name: "supported library version in range",
			library: Library{
				Name:    "another-library",
				Version: "2.5.0",
			},
			instrumentationType: TypeLibrary,
			expected:            []string{"https://example.com/another-library"},
		},
		{
			name: "unsupported library version outside range",
			library: Library{
				Name:    "another-library",
				Version: "3.0.0",
			},
			instrumentationType: TypeLibrary,
			expected:            nil,
		},
		{
			name: "unknown library",
			library: Library{
				Name:    "unknown-library",
				Version: "1.0.0",
			},
			instrumentationType: TypeLibrary,
			expected:            nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindSupportedLibraries(tt.library, s, tt.instrumentationType)
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

func TestCheckLibraries(t *testing.T) {
	s := SupportedModules{
		"test-library": {
			Instrumentations: []Instrumentation{
				{
					Name: "test-library",
					Link: "https://example.com/test-library",
					TargetVersions: map[InstrumentationType][]string{
						TypeLibrary: {"[1.0.0,)"},
					},
				},
			},
		},
	}

	deps := []Library{
		{
			Name:    "test-library",
			Version: "1.5.0",
		},
		{
			Name:    "unsupported-library",
			Version: "1.0.0",
		},
	}

	reporter := &utils.ComponentReporter{
		Checks:   []string{},
		Warnings: []string{},
		Errors:   []string{},
	}
	commands := utils.Commands{Debug: true}

	CheckLibraries(reporter, commands, s, deps, TypeLibrary)

	if len(reporter.Checks) != 1 {
		t.Errorf("Expected 1 successful check, got %d", len(reporter.Checks))
	}
	if len(reporter.Warnings) != 1 {
		t.Errorf("Expected 1 warning, got %d", len(reporter.Warnings))
	}
}
