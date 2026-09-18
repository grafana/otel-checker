// Package output renders a Reporter to a chosen format. Supported formats
// are "text" (colored stdout), "json", and "yaml". Library callers that
// want full control can implement Renderer themselves and call its Render
// method, or read Reporter.Results() and render manually.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/grafana/otel-checker/checks/utils"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/fatih/color"
	"go.yaml.in/yaml/v3"
	"golang.org/x/term"
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

// TextRenderer prints a table with columns: STATUS | COMPONENT | MESSAGE | EXPLAIN_ID,
// bordered and sized to the terminal width when styling is available,
// tab-aligned fallback otherwise (piped output, NO_COLOR, non-TTY).
type TextRenderer struct{}

var textHeaders = []string{"STATUS", "COMPONENT", "MESSAGE", "EXPLAIN_ID"}

func (TextRenderer) Render(w io.Writer, reporter *utils.Reporter) error {
	res := reporter.Results()
	total := len(res.Errors) + len(res.Warnings) + len(res.Checks)
	if total == 0 {
		return nil
	}

	rows := make([][]string, 0, total)
	anyExplainID := false
	appendRow := func(status string, m utils.ComponentResult) {
		if m.ExplainID != "" {
			anyExplainID = true
		}
		rows = append(rows, []string{status, m.Component, m.Message, m.ExplainID})
	}

	for _, m := range res.Errors {
		appendRow("FAIL", m)
	}
	for _, m := range res.Warnings {
		appendRow("WARN", m)
	}
	for _, m := range res.Checks {
		appendRow("OK", m)
	}

	if stylingEnabled() {
		if err := renderStyledTable(w, rows); err != nil {
			return err
		}
	} else {
		if err := renderPlainTable(w, rows); err != nil {
			return err
		}
	}

	_, _ = fmt.Fprintf(w, "%d %s, %d %s, %d successful %s.\n",
		len(res.Errors), pluralize("error", len(res.Errors)),
		len(res.Warnings), pluralize("warning", len(res.Warnings)),
		len(res.Checks), pluralize("check", len(res.Checks)),
	)
	if anyExplainID {
		_, _ = fmt.Fprintln(w, `Run "otel-checker explain <id>" for guidance on any finding above.`)
	}
	return nil
}

func renderStyledTable(w io.Writer, rows [][]string) error {
	width := terminalWidth()

	borderColor := lipgloss.Color("#44474E")
	primaryColor := lipgloss.Color("#6E9FFF")
	errorColor := lipgloss.Color("#E24D42")
	warnColor := lipgloss.Color("#EAB839")
	okColor := lipgloss.Color("#508642")

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(primaryColor).
		Padding(0, 1)
	cellStyle := lipgloss.NewStyle().Padding(0, 1)
	evenRow := cellStyle.Foreground(lipgloss.Color("#CCCCCC"))
	oddRow := cellStyle.Foreground(lipgloss.Color("#999999"))

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(borderColor)).
		Headers(textHeaders...).
		Rows(rows...).
		Width(width).
		StyleFunc(func(row, col int) lipgloss.Style {
			var s lipgloss.Style
			switch {
			case row == table.HeaderRow:
				s = headerStyle
			case row%2 == 0:
				s = evenRow
			default:
				s = oddRow
			}

			if col == 0 && row >= 0 && row < len(rows) {
				switch rows[row][0] {
				case "FAIL":
					s = s.Foreground(errorColor).Bold(true)
				case "WARN":
					s = s.Foreground(warnColor)
				case "OK":
					s = s.Foreground(okColor)
				}
			}
			return s
		})

	_, err := fmt.Fprintln(w, t)
	return err
}

// renderPlainTable is the plain fallback used when stdout is piped,
// NO_COLOR is set, or the terminal doesn't support styling.
func renderPlainTable(w io.Writer, rows [][]string) error {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', tabwriter.TabIndent|tabwriter.DiscardEmptyColumns)
	for i, h := range textHeaders {
		if i > 0 {
			_, _ = fmt.Fprint(tw, "\t")
		}
		_, _ = fmt.Fprint(tw, h)
	}
	_, _ = fmt.Fprintln(tw)
	for _, row := range rows {
		for i, v := range row {
			if i > 0 {
				_, _ = fmt.Fprint(tw, "\t")
			}
			_, _ = fmt.Fprint(tw, v)
		}
		_, _ = fmt.Fprintln(tw)
	}
	return tw.Flush()
}

func stylingEnabled() bool {
	if color.NoColor {
		return false
	}
	return term.IsTerminal(int(os.Stdout.Fd()))
}

func terminalWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return 100
	}
	return w
}

func pluralize(singular string, count int) string {
	if count == 1 {
		return singular
	}
	return singular + "s"
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
