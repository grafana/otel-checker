// Package output renders a Reporter to a chosen format. Supported formats
// are "text" (colored stdout), "json", and "yaml". Library callers that
// want full control can implement Renderer themselves and call its Render
// method, or read Reporter.Results() and render manually.
package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/grafana/otel-checker/checks/utils"

	"github.com/fatih/color"
	"go.yaml.in/yaml/v3"
)

const (
	FormatText = "text"
	FormatJSON = "json"
	FormatYAML = "yaml"
)

// Renderer turns a Reporter into bytes written to w. Implementations should
// be cheap to construct; pass options as struct fields.
type Renderer interface {
	Render(w io.Writer, reporter *utils.Reporter) error
}

// TextRenderer prints a human-readable summary. ANSI color codes are emitted
// when fatih/color's runtime detection decides the output is a terminal
// (respects NO_COLOR, FORCE_COLOR, and isatty checks).
type TextRenderer struct{}

func (TextRenderer) Render(w io.Writer, reporter *utils.Reporter) error {
	res := reporter.Results()

	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	if len(res.Checks) > 0 {
		_, _ = green.Fprintf(w, "\n%d Successful Check(s)\n", len(res.Checks))
		for _, m := range res.Checks {
			_, _ = green.Fprintf(w, "✔ %s: %s \n", m.Component, m.Message)
		}
	}
	if len(res.Warnings) > 0 {
		_, _ = yellow.Fprintf(w, "\n%d Warning(s)\n", len(res.Warnings))
		for _, m := range res.Warnings {
			_, _ = yellow.Fprintf(w, "• %s: %s \n", m.Component, m.Message)
		}
	}
	if len(res.Errors) > 0 {
		_, _ = red.Fprintf(w, "\n%d Error(s)\n", len(res.Errors))
		for _, m := range res.Errors {
			_, _ = red.Fprintf(w, "✖ %s: %s \n", m.Component, m.Message)
		}
	}
	return nil
}

// JSONRenderer emits Reporter.Results() as JSON.
type JSONRenderer struct {
	// Indent is the per-level indent string. Empty means compact output.
	Indent string
}

func (j JSONRenderer) Render(w io.Writer, reporter *utils.Reporter) error {
	enc := json.NewEncoder(w)
	if j.Indent != "" {
		enc.SetIndent("", j.Indent)
	}
	return enc.Encode(reporter.Results())
}

// YAMLRenderer emits Reporter.Results() as YAML.
type YAMLRenderer struct{}

func (YAMLRenderer) Render(w io.Writer, reporter *utils.Reporter) error {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	defer func() { _ = enc.Close() }()
	return enc.Encode(reporter.Results())
}

// Render is a convenience that picks a Renderer based on the format string
// and writes to w. An empty format defaults to text.
func Render(w io.Writer, reporter *utils.Reporter, format string) error {
	if format == "" {
		format = FormatText
	}
	switch format {
	case FormatText:
		return TextRenderer{}.Render(w, reporter)
	case FormatJSON:
		return JSONRenderer{Indent: "  "}.Render(w, reporter)
	case FormatYAML:
		return YAMLRenderer{}.Render(w, reporter)
	default:
		return fmt.Errorf("format %q not supported. Possible values: text, json, yaml", format)
	}
}
