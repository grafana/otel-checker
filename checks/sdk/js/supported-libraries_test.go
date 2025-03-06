package js

import (
	"os"
	"otel-checker/checks/sdk/supported"
	"testing"
)

func TestReadPackageJson(t *testing.T) {
	// Create a temporary package.json
	content := `{
		"dependencies": {
			"express": "^4.18.2",
			"@opentelemetry/instrumentation-express": "~0.35.0"
		},
		"devDependencies": {
			"typescript": "~5.3.3"
		}
	}`
	tmpfile, err := os.CreateTemp("", "package.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Change to the directory containing the temp file
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)
	if err := os.Chdir(os.TempDir()); err != nil {
		t.Fatal(err)
	}

	// Test reading dependencies
	deps := readPackageJson(nil)
	if len(deps) != 3 {
		t.Errorf("Expected 3 dependencies, got %d", len(deps))
	}

	// Check express dependency
	express := findDep(deps, "express")
	if express == nil {
		t.Error("express dependency not found")
	} else if express.Version != "4.18.2" {
		t.Errorf("Expected express version 4.18.2, got %s", express.Version)
	}

	// Check @opentelemetry/instrumentation-express dependency
	otel := findDep(deps, "@opentelemetry/instrumentation-express")
	if otel == nil {
		t.Error("@opentelemetry/instrumentation-express dependency not found")
	} else if otel.Version != "0.35.0" {
		t.Errorf("Expected @opentelemetry/instrumentation-express version 0.35.0, got %s", otel.Version)
	}

	// Check typescript dependency
	ts := findDep(deps, "typescript")
	if ts == nil {
		t.Error("typescript dependency not found")
	} else if ts.Version != "5.3.3" {
		t.Errorf("Expected typescript version 5.3.3, got %s", ts.Version)
	}
}

func TestReadPackageLock(t *testing.T) {
	// Create a temporary package-lock.json
	content := `{
		"dependencies": {
			"express": {
				"version": "4.18.2"
			},
			"@opentelemetry/instrumentation-express": {
				"version": "0.35.0"
			}
		}
	}`
	tmpfile, err := os.CreateTemp("", "package-lock.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Change to the directory containing the temp file
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldDir)
	if err := os.Chdir(os.TempDir()); err != nil {
		t.Fatal(err)
	}

	// Test reading dependencies
	deps := readPackageLock(nil)
	if len(deps) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(deps))
	}

	// Check express dependency
	express := findDep(deps, "express")
	if express == nil {
		t.Error("express dependency not found")
	} else if express.Version != "4.18.2" {
		t.Errorf("Expected express version 4.18.2, got %s", express.Version)
	}

	// Check @opentelemetry/instrumentation-express dependency
	otel := findDep(deps, "@opentelemetry/instrumentation-express")
	if otel == nil {
		t.Error("@opentelemetry/instrumentation-express dependency not found")
	} else if otel.Version != "0.35.0" {
		t.Errorf("Expected @opentelemetry/instrumentation-express version 0.35.0, got %s", otel.Version)
	}
}

func TestFindSupportedLibraries(t *testing.T) {
	// Create test data
	supported := supported.SupportedModules{
		"express": {
			Instrumentations: []supported.Instrumentation{
				{
					Name: "express",
					Link: "https://github.com/open-telemetry/opentelemetry-js-contrib/tree/main/plugins/node/opentelemetry-instrumentation-express",
					TargetVersions: map[supported.InstrumentationType][]string{
						supported.TypeLibrary: {"[4.0.0,)"},
					},
				},
			},
		},
		"@opentelemetry/instrumentation-express": {
			Instrumentations: []supported.Instrumentation{
				{
					Name: "@opentelemetry/instrumentation-express",
					Link: "https://github.com/open-telemetry/opentelemetry-js-contrib/tree/main/plugins/node/opentelemetry-instrumentation-express",
					TargetVersions: map[supported.InstrumentationType][]string{
						supported.TypeLibrary: {"[0.35.0,)"},
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
			name: "supported express version",
			library: supported.Library{
				Name:    "express",
				Version: "4.18.2",
			},
			expected: []string{"https://github.com/open-telemetry/opentelemetry-js-contrib/tree/main/plugins/node/opentelemetry-instrumentation-express"},
		},
		{
			name: "unsupported express version",
			library: supported.Library{
				Name:    "express",
				Version: "3.0.0",
			},
			expected: nil,
		},
		{
			name: "supported @opentelemetry/instrumentation-express version",
			library: supported.Library{
				Name:    "@opentelemetry/instrumentation-express",
				Version: "0.35.0",
			},
			expected: []string{"https://github.com/open-telemetry/opentelemetry-js-contrib/tree/main/plugins/node/opentelemetry-instrumentation-express"},
		},
		{
			name: "unsupported @opentelemetry/instrumentation-express version",
			library: supported.Library{
				Name:    "@opentelemetry/instrumentation-express",
				Version: "0.34.0",
			},
			expected: nil,
		},
		{
			name: "unknown library",
			library: supported.Library{
				Name:    "unknown",
				Version: "1.0.0",
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findSupportedLibraries(tt.library, supported)
			if len(got) != len(tt.expected) {
				t.Errorf("Expected %d links, got %d", len(tt.expected), len(got))
			}
			for i, link := range got {
				if link != tt.expected[i] {
					t.Errorf("Expected link %s, got %s", tt.expected[i], link)
				}
			}
		})
	}
}

func findDep(deps []supported.Library, name string) *supported.Library {
	for i := range deps {
		if deps[i].Name == name {
			return &deps[i]
		}
	}
	return nil
}
