package utils

import (
	"strings"
	"testing"
)

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
			name:    "empty language",
			in:      Commands{Components: []string{"sdk"}},
			wantErr: `language "" not supported`,
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
