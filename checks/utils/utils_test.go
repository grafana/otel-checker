package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReporterResults(t *testing.T) {
	r := &Reporter{}
	sdk := r.Component("SDK")
	sdk.AddSuccessfulCheck("foo")
	sdk.AddWarning("bar")
	collector := r.Component("Collector")
	collector.AddError("baz")

	got := r.Results()

	assert.ElementsMatch(t, []string{"SDK: foo"}, got[CHECKS])
	assert.ElementsMatch(t, []string{"SDK: bar"}, got[WARNINGS])
	assert.ElementsMatch(t, []string{"Collector: baz"}, got[ERRORS])
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		in      Commands
		wantErr string
	}{
		{
			name: "valid go sdk",
			in:   Commands{Language: "go", Components: []string{"sdk"}},
		},
		{
			name: "valid js collector with trimmed component",
			in:   Commands{Language: "js", Components: []string{"sdk", " collector"}},
		},
		{
			name:    "unsupported language",
			in:      Commands{Language: "rust", Components: []string{"sdk"}},
			wantErr: `language "rust" not supported`,
		},
		{
			name:    "empty language with sdk",
			in:      Commands{Components: []string{"sdk"}},
			wantErr: "language required for components",
		},
		{
			name: "collector-only without language",
			in:   Commands{Components: []string{"collector"}},
		},
		{
			name:    "no components",
			in:      Commands{Language: "go"},
			wantErr: "at least one component required",
		},
		{
			name:    "unsupported component",
			in:      Commands{Language: "go", Components: []string{"sdk", "bogus"}},
			wantErr: `component "bogus" not supported`,
		},
		{
			name: "js manual without instrumentation file",
			in: Commands{
				Language:              "js",
				Components:            []string{"sdk"},
				ManualInstrumentation: true,
			},
			wantErr: "manual-instrumentation is set",
		},
		{
			name: "js manual with instrumentation file",
			in: Commands{
				Language:              "js",
				Components:            []string{"sdk"},
				ManualInstrumentation: true,
				InstrumentationFile:   "src/inst.js",
			},
		},
		{
			name: "web-server with malformed listen address",
			in: Commands{
				Language:   "go",
				Components: []string{"sdk"},
				WebServer:  true,
				Listen:     "no-port-here",
			},
			wantErr: "not a valid host:port",
		},
		{
			name: "web-server with valid listen address",
			in: Commands{
				Language:   "go",
				Components: []string{"sdk"},
				WebServer:  true,
				Listen:     "127.0.0.1:8080",
			},
		},
		{
			name: "listen address ignored when web-server is off",
			in: Commands{
				Language:   "go",
				Components: []string{"sdk"},
				Listen:     "bogus",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := Validate(tc.in)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("Validate() returned unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("Validate() returned nil, want error containing %q", tc.wantErr)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("Validate() error = %q, want it to contain %q", err.Error(), tc.wantErr)
			}
		})
	}
}
