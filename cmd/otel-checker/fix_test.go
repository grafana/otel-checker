package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grafana/otel-checker/checks/fixes"
)

func TestFixShowKnownID(t *testing.T) {
	// Pick the first registered ID so this test stays valid as the registry grows.
	ids := fixes.All()
	if len(ids) == 0 {
		t.Fatal("no fix IDs registered; expected at least one")
	}
	id := ids[0]
	doc, _ := fixes.Lookup(id)

	var buf bytes.Buffer
	if err := runFixShow(&buf, id); err != nil {
		t.Fatalf("runFixShow(%q) returned err: %v", id, err)
	}
	out := buf.String()
	if !strings.Contains(out, doc.Title) {
		t.Errorf("fix show output missing title %q. Got:\n%s", doc.Title, out)
	}
	if !strings.Contains(out, doc.Body) {
		t.Errorf("fix show output missing body. Got:\n%s", out)
	}
}

func TestFixShowUnknownID(t *testing.T) {
	var buf bytes.Buffer
	err := runFixShow(&buf, "definitely.not.a.real.id")
	if err == nil {
		t.Fatal("runFixShow(unknown) returned nil, want error")
	}
	if !strings.Contains(err.Error(), "definitely.not.a.real.id") {
		t.Errorf("error %q should mention the offending id", err.Error())
	}
	if !strings.Contains(err.Error(), "fix list") {
		t.Errorf("error %q should point users at `fix list`", err.Error())
	}
}

func TestFixList(t *testing.T) {
	var buf bytes.Buffer
	if err := runFixList(&buf); err != nil {
		t.Fatalf("runFixList: %v", err)
	}
	out := buf.String()
	for _, id := range fixes.All() {
		if !strings.Contains(out, id) {
			t.Errorf("fix list output missing %q", id)
		}
	}
}

func TestFixAll(t *testing.T) {
	// Write a results.json with two errors (different IDs), one warning
	// (a third ID, but duplicated), and one finding without a FixID.
	dir := t.TempDir()
	path := filepath.Join(dir, "results.json")
	body := `{
  "checks": [],
  "warnings": [
    {"component": "X", "message": "duplicate", "fix_id": "env.otel-service-name.unset"},
    {"component": "X", "message": "duplicate again", "fix_id": "env.otel-service-name.unset"}
  ],
  "errors": [
    {"component": "X", "message": "untagged"},
    {"component": "X", "message": "first", "fix_id": "collector.config.unreadable"},
    {"component": "X", "message": "second", "fix_id": "js.node-version.too-old"}
  ]
}`
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	var buf bytes.Buffer
	if err := runFixAll(&buf, path); err != nil {
		t.Fatalf("runFixAll: %v", err)
	}
	out := buf.String()

	// Errors come first, then warnings; duplicates collapse; untagged is skipped.
	wantOrder := []string{
		"collector.config.unreadable",
		"js.node-version.too-old",
		"env.otel-service-name.unset",
	}
	prev := -1
	for _, id := range wantOrder {
		i := strings.Index(out, id)
		if i < 0 {
			t.Fatalf("output missing fix id %q.\nGot:\n%s", id, out)
		}
		if i < prev {
			t.Errorf("fix id %q appears before the previous expected id", id)
		}
		prev = i
	}

	// Only one occurrence of the duplicate id.
	if got := strings.Count(out, "env.otel-service-name.unset"); got != 1 {
		t.Errorf("env.otel-service-name.unset appears %d times, want 1 (deduped)", got)
	}
}

func TestFixAllMissingFile(t *testing.T) {
	dir := t.TempDir()
	missing := filepath.Join(dir, "nope.json")
	var buf bytes.Buffer
	err := runFixAll(&buf, missing)
	if err == nil {
		t.Fatal("runFixAll on missing file returned nil, want error")
	}
	if !strings.Contains(err.Error(), "no results file found") {
		t.Errorf("error %q should explain the missing-file case", err.Error())
	}
}
