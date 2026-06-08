package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/grafana/otel-checker/checks/utils"

	"github.com/fatih/color"
	"go.yaml.in/yaml/v3"
)

func newReporter(t *testing.T) *utils.Reporter {
	t.Helper()
	r := &utils.Reporter{}
	sdk := r.Component("SDK")
	sdk.AddSuccessfulCheck("foo")
	sdk.AddWarning("bar")
	collector := r.Component("Collector")
	collector.AddError("baz")
	return r
}

func TestRenderJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := Render(&buf, newReporter(t), FormatJSON); err != nil {
		t.Fatalf("Render(json): %v", err)
	}
	var got utils.Results
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, buf.String())
	}
	if len(got.Checks) != 1 || got.Checks[0].Component != "SDK" || got.Checks[0].Message != "foo" {
		t.Errorf("checks = %+v, want one SDK/foo entry", got.Checks)
	}
	if len(got.Warnings) != 1 || got.Warnings[0].Component != "SDK" || got.Warnings[0].Message != "bar" {
		t.Errorf("warnings = %+v, want one SDK/bar entry", got.Warnings)
	}
	if len(got.Errors) != 1 || got.Errors[0].Component != "Collector" || got.Errors[0].Message != "baz" {
		t.Errorf("errors = %+v, want one Collector/baz entry", got.Errors)
	}
}

func TestRenderYAML(t *testing.T) {
	var buf bytes.Buffer
	if err := Render(&buf, newReporter(t), FormatYAML); err != nil {
		t.Fatalf("Render(yaml): %v", err)
	}
	var got utils.Results
	if err := yaml.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid YAML: %v\n%s", err, buf.String())
	}
	if len(got.Checks) != 1 || got.Checks[0].Message != "foo" {
		t.Errorf("checks = %+v, want one entry with foo", got.Checks)
	}
}

func TestRenderText(t *testing.T) {
	// Disable color codes so substring assertions don't trip on ANSI escapes.
	saved := color.NoColor
	color.NoColor = true
	t.Cleanup(func() { color.NoColor = saved })

	var buf bytes.Buffer
	if err := Render(&buf, newReporter(t), FormatText); err != nil {
		t.Fatalf("Render(text): %v", err)
	}
	out := buf.String()
	for _, want := range []string{"1 Successful Check", "SDK: foo", "1 Warning", "SDK: bar", "1 Error", "Collector: baz"} {
		if !strings.Contains(out, want) {
			t.Errorf("text output missing %q.\nGot:\n%s", want, out)
		}
	}
}

func TestRenderDefaultsToText(t *testing.T) {
	saved := color.NoColor
	color.NoColor = true
	t.Cleanup(func() { color.NoColor = saved })

	var buf bytes.Buffer
	if err := Render(&buf, newReporter(t), ""); err != nil {
		t.Fatalf("Render(empty format): %v", err)
	}
	if !strings.Contains(buf.String(), "1 Successful Check") {
		t.Errorf("empty format did not produce text output. Got:\n%s", buf.String())
	}
}

func TestRenderUnknownFormat(t *testing.T) {
	err := Render(&bytes.Buffer{}, newReporter(t), "xml")
	if err == nil {
		t.Fatal("Render(xml) returned nil, want error")
	}
	if !strings.Contains(err.Error(), `"xml" not supported`) {
		t.Errorf("error = %q, want it to mention xml", err.Error())
	}
}
