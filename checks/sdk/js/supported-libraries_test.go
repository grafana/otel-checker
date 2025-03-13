package js

import (
	"github.com/grafana/otel-checker/checks/sdk/supported"
	"testing"
)

func TestReadPackageJson(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []supported.Library
	}{
		{
			name: "parse package.json with dependencies and devDependencies",
			content: `{
				"dependencies": {
					"express": "^4.18.2",
					"@opentelemetry/instrumentation-express": "~0.35.0"
				},
				"devDependencies": {
					"typescript": "~5.3.3"
				}
			}`,
			expected: []supported.Library{
				{Name: "express", Version: "4.18.2"},
				{Name: "@opentelemetry/instrumentation-express", Version: "0.35.0"},
				{Name: "typescript", Version: "5.3.3"},
			},
		},
		{
			name: "parse package.json with only dependencies",
			content: `{
				"dependencies": {
					"express": "4.18.2",
					"@opentelemetry/instrumentation-express": "0.35.0"
				}
			}`,
			expected: []supported.Library{
				{Name: "express", Version: "4.18.2"},
				{Name: "@opentelemetry/instrumentation-express", Version: "0.35.0"},
			},
		},
		{
			name: "parse package.json with only devDependencies",
			content: `{
				"devDependencies": {
					"typescript": "5.3.3"
				}
			}`,
			expected: []supported.Library{
				{Name: "typescript", Version: "5.3.3"},
			},
		},
		{
			name:     "parse empty package.json",
			content:  `{}`,
			expected: []supported.Library{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := readPackageJsonFromContent([]byte(tt.content))
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

func TestReadPackageLock(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []supported.Library
	}{
		{
			name: "parse package-lock.json with multiple dependencies",
			content: `{
				"name": "example-fastify",
				"version": "0.26.0",
				"lockfileVersion": 2,
				"requires": true,
				"packages": {
					"node_modules/@fastify/ajv-compiler": {
						"version": "3.5.0",
						"resolved": "https://registry.npmjs.org/@fastify/ajv-compiler/-/ajv-compiler-3.5.0.tgz",
						"integrity": "sha512-ebbEtlI7dxXF5ziNdr05mOY8NnDiPB1XvAlLHctRt/Rc+C3LCOVW5imUVX+mhvUhnNzmPBHewUkOFgGlCxgdAA==",
						"dependencies": {
							"ajv": "^8.11.0",
							"ajv-formats": "^2.1.1",
							"fast-uri": "^2.0.0"
						}
					},
					"node_modules/express": {
						"version": "4.18.2",
						"resolved": "https://registry.npmjs.org/express/-/express-4.18.2.tgz",
						"integrity": "sha512-5/PsL9iGiiH9nMG2WLbQJCszTs+AwHom0DPv8O3dsdrbXuuP9PWUJ5RhdleT3f1wOTn7d2x4mO1Qw8OHtZHwINg==",
						"dependencies": {
							"accepts": "~1.3.8",
							"body-parser": "1.20.2"
						}
					},
					"node_modules/@opentelemetry/instrumentation-express": {
						"version": "0.35.0",
						"resolved": "https://registry.npmjs.org/@opentelemetry/instrumentation-express/-/instrumentation-express-0.35.0.tgz",
						"integrity": "sha512-xyz123",
						"dependencies": {
							"@opentelemetry/api": "^1.7.0",
							"@opentelemetry/semantic-conventions": "^1.21.0"
						}
					}
				}
			}`,
			expected: []supported.Library{
				{Name: "@fastify/ajv-compiler", Version: "3.5.0"},
				{Name: "express", Version: "4.18.2"},
				{Name: "@opentelemetry/instrumentation-express", Version: "0.35.0"},
			},
		},
		{
			name: "parse package-lock.json with root package",
			content: `{
				"name": "example",
				"version": "1.0.0",
				"lockfileVersion": 2,
				"requires": true,
				"packages": {
					"": {
						"version": "1.0.0"
					},
					"node_modules/express": {
						"version": "4.18.2"
					}
				}
			}`,
			expected: []supported.Library{
				{Name: "express", Version: "4.18.2"},
			},
		},
		{
			name:     "parse empty package-lock.json",
			content:  `{}`,
			expected: []supported.Library{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := readPackageLockFromContent([]byte(tt.content))
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
			got := supported.FindSupportedLibraries(tt.library, s, supported.TypeLibrary)
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
