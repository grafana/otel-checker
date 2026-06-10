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
// reporter component that produced it. ExplainID, when non-empty, references
// a document in the explain package that describes the finding and how to
// address it.
type ComponentResult struct {
	Component string `json:"component" yaml:"component"`
	Message   string `json:"message" yaml:"message"`
	ExplainID string `json:"explain_id,omitempty" yaml:"explain_id,omitempty"`
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
	name string
	// Checks, Warnings, Errors hold the human-readable message for each
	// finding. Public for back-compat with existing tests; new code should
	// consume Reporter.Results() to get the typed ComponentResult with ExplainID.
	Checks   []string
	Warnings []string
	Errors   []string
	// checkIDs, warningIDs, errorIDs are parallel to the slices above and
	// hold the explain-doc ID for each finding ("" when no explain applies).
	// Kept unexported so the lockstep invariant (one ID per message) can only
	// be maintained through the Add* methods.
	checkIDs   []string
	warningIDs []string
	errorIDs   []string
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
	for _, c := range r.components {
		res.Checks = append(res.Checks, c.toResults(c.Checks, c.checkIDs)...)
		res.Warnings = append(res.Warnings, c.toResults(c.Warnings, c.warningIDs)...)
		res.Errors = append(res.Errors, c.toResults(c.Errors, c.errorIDs)...)
	}
	return res
}

// toResults zips a message slice with its parallel ID slice into a
// []ComponentResult tagged with the component's name. The two slices are
// always in lockstep because all mutation goes through AddXxxWithExplain.
func (c *ComponentReporter) toResults(messages, ids []string) []ComponentResult {
	out := make([]ComponentResult, len(messages))
	for i, m := range messages {
		out[i] = ComponentResult{Component: c.name, Message: m, ExplainID: ids[i]}
	}
	return out
}

func (r *ComponentReporter) AddSuccessfulCheck(message string) {
	r.AddSuccessfulCheckWithExplain("", message)
}

func (r *ComponentReporter) AddWarning(message string) {
	r.AddWarningWithExplain("", message)
}

func (r *ComponentReporter) AddInternalError(message string) {
	r.AddInternalErrorWithExplain("", message)
}

func (r *ComponentReporter) AddError(message string) {
	r.AddErrorWithExplain("", message)
}

// AddSuccessfulCheckWithExplain records a successful check with an optional
// explain-doc ID. Successful checks rarely need explain docs, but the
// variant is provided for symmetry.
func (r *ComponentReporter) AddSuccessfulCheckWithExplain(explainID, message string) {
	r.Checks = append(r.Checks, message)
	r.checkIDs = append(r.checkIDs, explainID)
}

// AddWarningWithExplain records a warning with the given explain-doc ID.
// Pass "" if no explain doc applies.
func (r *ComponentReporter) AddWarningWithExplain(explainID, message string) {
	r.Warnings = append(r.Warnings, message)
	r.warningIDs = append(r.warningIDs, explainID)
}

// AddInternalErrorWithExplain records an internal-error message, prefixed
// with "Internal Error: " and reported as a warning. The explain-doc ID, if
// any, typically points at a "report this bug upstream" doc.
func (r *ComponentReporter) AddInternalErrorWithExplain(explainID, message string) {
	r.AddWarningWithExplain(explainID, "Internal Error: "+message)
}

// AddErrorWithExplain records an error with the given explain-doc ID. Pass
// "" if no explain doc applies.
func (r *ComponentReporter) AddErrorWithExplain(explainID, message string) {
	r.Errors = append(r.Errors, message)
	r.errorIDs = append(r.errorIDs, explainID)
}

func FileExists(path string) bool {
	_, err := os.ReadFile(path)
	return err == nil
}
