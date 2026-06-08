package utils

import (
	"fmt"
	"net"
	"os"
	"slices"
	"strings"
)

const ERRORS = "errors"
const WARNINGS = "warnings"
const CHECKS = "checks"

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
		return fmt.Errorf("at least one component required. Possible values: %s", strings.Join(SupportedComponents, ", "))
	}
	needsLanguage := false
	for _, comp := range c.Components {
		comp = strings.TrimSpace(comp)
		if !slices.Contains(SupportedComponents, comp) {
			return fmt.Errorf("component %q not supported. Possible values: %s", comp, strings.Join(SupportedComponents, ", "))
		}
		if slices.Contains(LanguageRequiredFor, comp) {
			needsLanguage = true
		}
	}
	if needsLanguage && c.Language == "" {
		return fmt.Errorf("language required for components: %s", strings.Join(LanguageRequiredFor, ", "))
	}
	if c.Language != "" && !slices.Contains(SupportedLanguages, c.Language) {
		return fmt.Errorf("language %q not supported. Possible values: %s", c.Language, strings.Join(SupportedLanguages, ", "))
	}
	if c.Language == "js" && c.ManualInstrumentation && c.InstrumentationFile == "" {
		return fmt.Errorf(`when manual-instrumentation is set, an instrumentation file is required (InstrumentationFile or -instrumentation-file=path/to/file.js)`)
	}
	if c.Format != "" && !slices.Contains(SupportedFormats, c.Format) {
		return fmt.Errorf("format %q not supported. Possible values: %s", c.Format, strings.Join(SupportedFormats, ", "))
	}
	if c.WebServer && c.Listen != "" {
		if _, _, err := net.SplitHostPort(c.Listen); err != nil {
			return fmt.Errorf("listen address %q is not a valid host:port: %w", c.Listen, err)
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
// way should use this; CLI callers should use PrintResults.
func (r *Reporter) Results() map[string][]string {
	res := make(map[string][]string)
	var checks []string
	for _, component := range r.components {
		checks = append(checks, component.Checks...)
	}
	res[CHECKS] = checks
	var warnings []string
	for _, component := range r.components {
		warnings = append(warnings, component.Warnings...)
	}
	res[WARNINGS] = warnings
	var errors []string
	for _, component := range r.components {
		errors = append(errors, component.Errors...)
	}
	res[ERRORS] = errors
	return res
}

func (r *ComponentReporter) AddSuccessfulCheck(message string) {
	r.Checks = append(r.Checks, fmt.Sprintf(`%s: %s`, r.name, message))
}

func (r *ComponentReporter) AddWarning(message string) {
	r.Warnings = append(r.Warnings, fmt.Sprintf(`%s: %s`, r.name, message))
}

func (r *ComponentReporter) AddInternalError(message string) {
	r.Warnings = append(r.Warnings, fmt.Sprintf(`%s: Internal Error: %s`, r.name, message))
}

func (r *ComponentReporter) AddError(message string) {
	r.Errors = append(r.Errors, fmt.Sprintf(`%s: %s`, r.name, message))
}

func FileExists(path string) bool {
	_, err := os.ReadFile(path)
	return err == nil
}
