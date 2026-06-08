package utils

import (
	"net"
	"os"
	"slices"
	"strings"
)

// Results is the typed snapshot of a Reporter's aggregated state.
// Library callers should consume this directly; the CLI marshals it as
// JSON/YAML and renders it as text via the output package.
type Results struct {
	Checks   []ComponentResult `json:"checks" yaml:"checks"`
	Warnings []ComponentResult `json:"warnings" yaml:"warnings"`
	Errors   []ComponentResult `json:"errors" yaml:"errors"`
}

// ComponentResult is a single check/warning/error message, tagged with the
// reporter component that produced it.
type ComponentResult struct {
	Component string `json:"component" yaml:"component"`
	Message   string `json:"message" yaml:"message"`
}

type Commands struct {
	Language              string
	Components            []string
	ManualInstrumentation bool
	WebServer             bool
	Listen                string
	InstrumentationFile   string
	PackageJsonPath       string
	CollectorConfigPath   string
	Debug                 bool
	Format                string
}

const DefaultListen = "127.0.0.1:8080"

var (
	SupportedLanguages  = []string{"dotnet", "go", "java", "js", "python", "ruby", "php"}
	SupportedComponents = []string{"sdk", "beyla", "alloy", "collector", "grafana-cloud"}
	SupportedFormats    = []string{"text", "json", "yaml"}
	// LanguageRequiredFor lists the components whose checks need a language hint.
	// "collector" is intentionally omitted — its YAML schema is language-agnostic.
	LanguageRequiredFor = []string{"sdk", "beyla", "alloy", "grafana-cloud"}
)

func Validate(c Commands) error {
	if len(c.Components) == 0 {
		return ErrNoComponents
	}
	needsLanguage := false
	for _, comp := range c.Components {
		comp = strings.TrimSpace(comp)
		if !slices.Contains(SupportedComponents, comp) {
			return &UnsupportedComponentError{Component: comp}
		}
		if slices.Contains(LanguageRequiredFor, comp) {
			needsLanguage = true
		}
	}
	if needsLanguage && c.Language == "" {
		return ErrLanguageRequired
	}
	if c.Language != "" && !slices.Contains(SupportedLanguages, c.Language) {
		return &UnsupportedLanguageError{Language: c.Language}
	}
	if c.Language == "js" && c.ManualInstrumentation && c.InstrumentationFile == "" {
		return ErrManualInstrumentationFile
	}
	if c.Format != "" && !slices.Contains(SupportedFormats, c.Format) {
		return &UnsupportedFormatError{Format: c.Format}
	}
	if c.WebServer && c.Listen != "" {
		if _, _, err := net.SplitHostPort(c.Listen); err != nil {
			return &InvalidListenError{Listen: c.Listen, Err: err}
		}
	}
	return nil
}

type Reporter struct {
	components []*ComponentReporter
}

type ComponentReporter struct {
	name     string
	Checks   []string
	Warnings []string
	Errors   []string
}

func (r *Reporter) Component(name string) *ComponentReporter {
	for _, component := range r.components {
		if component.name == name {
			return component
		}
	}
	c := &ComponentReporter{name: name}
	r.components = append(r.components, c)
	return c
}

// Results aggregates the checks, warnings, and errors across all components
// without producing any output. Callers that want to render results their own
// way should use this; the CLI marshals it via the output package.
func (r *Reporter) Results() Results {
	var res Results
	for _, component := range r.components {
		for _, m := range component.Checks {
			res.Checks = append(res.Checks, ComponentResult{Component: component.name, Message: m})
		}
		for _, m := range component.Warnings {
			res.Warnings = append(res.Warnings, ComponentResult{Component: component.name, Message: m})
		}
		for _, m := range component.Errors {
			res.Errors = append(res.Errors, ComponentResult{Component: component.name, Message: m})
		}
	}
	return res
}

func (r *ComponentReporter) AddSuccessfulCheck(message string) {
	r.Checks = append(r.Checks, message)
}

func (r *ComponentReporter) AddWarning(message string) {
	r.Warnings = append(r.Warnings, message)
}

func (r *ComponentReporter) AddInternalError(message string) {
	r.Warnings = append(r.Warnings, "Internal Error: "+message)
}

func (r *ComponentReporter) AddError(message string) {
	r.Errors = append(r.Errors, message)
}

func FileExists(path string) bool {
	_, err := os.ReadFile(path)
	return err == nil
}
