package utils

import (
	"errors"
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

	assert.Equal(t, []ComponentResult{{Component: "SDK", Message: "foo"}}, got.Checks)
	assert.Equal(t, []ComponentResult{{Component: "SDK", Message: "bar"}}, got.Warnings)
	assert.Equal(t, []ComponentResult{{Component: "Collector", Message: "baz"}}, got.Errors)
}

func TestReporterResultsWithExplainIDs(t *testing.T) {
	r := &Reporter{}
	sdk := r.Component("SDK")
	sdk.AddSuccessfulCheckWithExplain("ok.explain", "passed")
	sdk.AddWarningWithExplain("warn.explain", "watch out")
	sdk.AddError("no explain here") // legacy path → empty ExplainID
	col := r.Component("Collector")
	col.AddErrorWithExplain("err.explain", "broken")
	col.AddInternalErrorWithExplain("oops.explain", "internal blip")

	got := r.Results()

	assert.Equal(t, []ComponentResult{{Component: "SDK", Message: "passed", ExplainID: "ok.explain"}}, got.Checks)
	assert.Equal(t, []ComponentResult{
		{Component: "SDK", Message: "watch out", ExplainID: "warn.explain"},
		{Component: "Collector", Message: "Internal Error: internal blip", ExplainID: "oops.explain"},
	}, got.Warnings)
	assert.Equal(t, []ComponentResult{
		{Component: "SDK", Message: "no explain here"},
		{Component: "Collector", Message: "broken", ExplainID: "err.explain"},
	}, got.Errors)
}

func TestAddInternalErrorPrefixesMessage(t *testing.T) {
	r := &Reporter{}
	sdk := r.Component("SDK")
	sdk.AddInternalError("boom")
	got := r.Results()
	assert.Equal(t, []ComponentResult{{Component: "SDK", Message: "Internal Error: boom"}}, got.Warnings)
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
			wantErr: "language required",
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
			wantErr: "manual-instrumentation requires",
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

func TestValidateTypedErrors(t *testing.T) {
	t.Run("ErrNoComponents", func(t *testing.T) {
		err := Validate(Commands{Language: "go"})
		if !errors.Is(err, ErrNoComponents) {
			t.Errorf("errors.Is(err, ErrNoComponents) = false, want true; err = %v", err)
		}
	})

	t.Run("ErrLanguageRequired", func(t *testing.T) {
		err := Validate(Commands{Components: []string{"sdk"}})
		if !errors.Is(err, ErrLanguageRequired) {
			t.Errorf("errors.Is(err, ErrLanguageRequired) = false, want true; err = %v", err)
		}
	})

	t.Run("ErrManualInstrumentationFile", func(t *testing.T) {
		err := Validate(Commands{
			Language:              "js",
			Components:            []string{"sdk"},
			ManualInstrumentation: true,
		})
		if !errors.Is(err, ErrManualInstrumentationFile) {
			t.Errorf("errors.Is(err, ErrManualInstrumentationFile) = false, want true; err = %v", err)
		}
	})

	t.Run("UnsupportedLanguageError", func(t *testing.T) {
		err := Validate(Commands{Language: "rust", Components: []string{"sdk"}})
		var ule *UnsupportedLanguageError
		if !errors.As(err, &ule) {
			t.Fatalf("errors.As did not extract *UnsupportedLanguageError; err = %v", err)
		}
		if ule.Language != "rust" {
			t.Errorf("ule.Language = %q, want %q", ule.Language, "rust")
		}
	})

	t.Run("UnsupportedComponentError", func(t *testing.T) {
		err := Validate(Commands{Language: "go", Components: []string{"bogus"}})
		var uce *UnsupportedComponentError
		if !errors.As(err, &uce) {
			t.Fatalf("errors.As did not extract *UnsupportedComponentError; err = %v", err)
		}
		if uce.Component != "bogus" {
			t.Errorf("uce.Component = %q, want %q", uce.Component, "bogus")
		}
	})

	t.Run("UnsupportedFormatError", func(t *testing.T) {
		err := Validate(Commands{
			Language:   "go",
			Components: []string{"sdk"},
			Format:     "xml",
		})
		var ufe *UnsupportedFormatError
		if !errors.As(err, &ufe) {
			t.Fatalf("errors.As did not extract *UnsupportedFormatError; err = %v", err)
		}
		if ufe.Format != "xml" {
			t.Errorf("ufe.Format = %q, want %q", ufe.Format, "xml")
		}
	})

	t.Run("InvalidListenError", func(t *testing.T) {
		err := Validate(Commands{
			Language:   "go",
			Components: []string{"sdk"},
			WebServer:  true,
			Listen:     "no-port-here",
		})
		var ile *InvalidListenError
		if !errors.As(err, &ile) {
			t.Fatalf("errors.As did not extract *InvalidListenError; err = %v", err)
		}
		if ile.Listen != "no-port-here" {
			t.Errorf("ile.Listen = %q, want %q", ile.Listen, "no-port-here")
		}
		if ile.Unwrap() == nil {
			t.Error("ile.Unwrap() = nil, want underlying net.SplitHostPort error")
		}
	})
}
