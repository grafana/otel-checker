package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grafana/otel-checker/checks/explain"
)

func TestExplainShowKnownID(t *testing.T) {
	// Pick the first registered ID so this test stays valid as the registry grows.
	ids := explain.All()
	if len(ids) == 0 {
		t.Fatal("no explain IDs registered; expected at least one")
	}
	id := ids[0]
	doc, _ := explain.Lookup(id)

	var buf bytes.Buffer
	if err := runExplainShow(&buf, id); err != nil {
		t.Fatalf("runExplainShow(%q) returned err: %v", id, err)
	}
	out := buf.String()
	if !strings.Contains(out, doc.Title) {
		t.Errorf("explain show output missing title %q. Got:\n%s", doc.Title, out)
	}
	if !strings.Contains(out, doc.Body) {
		t.Errorf("explain show output missing body. Got:\n%s", out)
	}
}

func TestExplainShowUnknownID(t *testing.T) {
	var buf bytes.Buffer
	err := runExplainShow(&buf, "definitely.not.a.real.id")
	if err == nil {
		t.Fatal("runExplainShow(unknown) returned nil, want error")
	}
	if !strings.Contains(err.Error(), "definitely.not.a.real.id") {
		t.Errorf("error %q should mention the offending id", err.Error())
	}
	if !strings.Contains(err.Error(), "explain list") {
		t.Errorf("error %q should point users at `explain list`", err.Error())
	}
}

func TestExplainList(t *testing.T) {
	var buf bytes.Buffer
	if err := runExplainList(&buf); err != nil {
		t.Fatalf("runExplainList: %v", err)
	}
	out := buf.String()
	for _, id := range explain.All() {
		if !strings.Contains(out, id) {
			t.Errorf("explain list output missing %q", id)
		}
	}
}

func TestExplainAll(t *testing.T) {
	// Write a results.json with two errors (different IDs), one warning
	// (a third ID, but duplicated), and one finding without a ExplainID.
	dir := t.TempDir()
	path := filepath.Join(dir, "results.json")
	body := `{
  "checks": [],
  "warnings": [
    {"component": "X", "message": "duplicate", "explain_id": "env.otel-service-name.unset"},
    {"component": "X", "message": "duplicate again", "explain_id": "env.otel-service-name.unset"}
  ],
  "errors": [
    {"component": "X", "message": "untagged"},
    {"component": "X", "message": "first", "explain_id": "collector.config.unreadable"},
    {"component": "X", "message": "second", "explain_id": "js.node-version.too-old"}
  ]
}`
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	var buf bytes.Buffer
	if err := runExplainAll(&buf, path); err != nil {
		t.Fatalf("runExplainAll: %v", err)
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
			t.Fatalf("output missing explain id %q.\nGot:\n%s", id, out)
		}
		if i < prev {
			t.Errorf("explain id %q appears before the previous expected id", id)
		}
		prev = i
	}

	// Only one occurrence of the duplicate id.
	if got := strings.Count(out, "env.otel-service-name.unset"); got != 1 {
		t.Errorf("env.otel-service-name.unset appears %d times, want 1 (deduped)", got)
	}
}

func TestExplainAllMissingFile(t *testing.T) {
	dir := t.TempDir()
	missing := filepath.Join(dir, "nope.json")
	var buf bytes.Buffer
	err := runExplainAll(&buf, missing)
	if err == nil {
		t.Fatal("runExplainAll on missing file returned nil, want error")
	}
	if !strings.Contains(err.Error(), "no results file found") {
		t.Errorf("error %q should explain the missing-file case", err.Error())
	}
}
