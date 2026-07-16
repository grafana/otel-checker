package explain

import (
	"strings"
	"testing"
)

func TestAllIDsAreLookupable(t *testing.T) {
	ids := All()
	if len(ids) == 0 {
		t.Fatal("All() returned no IDs; the docs directory may not be embedded")
	}
	for _, id := range ids {
		doc, ok := Lookup(id)
		if !ok {
			t.Errorf("Lookup(%q) returned !ok, but All() included it", id)
			continue
		}
		if doc.ID != id {
			t.Errorf("Lookup(%q).ID = %q, want %q", id, doc.ID, id)
		}
		if doc.Title == "" {
			t.Errorf("Lookup(%q).Title is empty", id)
		}
		if doc.Severity == "" {
			t.Errorf("Lookup(%q).Severity is empty", id)
		}
	}
}

func TestSeverityIsValid(t *testing.T) {
	valid := map[string]bool{"warning": true, "error": true, "internal": true}
	for _, id := range All() {
		doc, _ := Lookup(id)
		if !valid[doc.Severity] {
			t.Errorf("%s: severity %q not in {warning, error, internal}", id, doc.Severity)
		}
	}
}

func TestLookupUnknown(t *testing.T) {
	_, ok := Lookup("definitely.not.a.real.id")
	if ok {
		t.Error("Lookup of an unknown ID returned ok=true")
	}
}

func TestAllSorted(t *testing.T) {
	ids := All()
	for i := 1; i < len(ids); i++ {
		if ids[i-1] > ids[i] {
			t.Fatalf("All() not sorted: %s > %s", ids[i-1], ids[i])
		}
	}
}

func TestParseRejectsMissingFrontMatter(t *testing.T) {
	_, err := parse([]byte("body only, no front matter"))
	if err == nil || !strings.Contains(err.Error(), "front-matter") {
		t.Errorf("parse() without front-matter, got err = %v", err)
	}
}

func TestParseRejectsMissingClosingDelimiter(t *testing.T) {
	_, err := parse([]byte("---\nid: x\ntitle: y\nseverity: warning\nbody\n"))
	if err == nil || !strings.Contains(err.Error(), "closing") {
		t.Errorf("parse() without closing ---, got err = %v", err)
	}
}

func TestParseAcceptsCRLF(t *testing.T) {
	// Docs edited via the GitHub web UI on Windows come back with CRLF
	// line endings. The parser must normalize before checking the ---
	// prefix; regression from #416.
	src := "---\r\nid: x\r\ntitle: y\r\nseverity: warning\r\n---\r\nbody\r\n"
	doc, err := parse([]byte(src))
	if err != nil {
		t.Fatalf("parse(CRLF) returned err = %v, want nil", err)
	}
	if doc.ID != "x" || doc.Title != "y" || doc.Severity != "warning" {
		t.Errorf("parse(CRLF) returned wrong front-matter: %+v", doc)
	}
}

func TestParseRejectsMissingFields(t *testing.T) {
	cases := []string{
		"---\ntitle: y\nseverity: warning\n---\nbody", // no id
		"---\nid: x\nseverity: warning\n---\nbody",    // no title
		"---\nid: x\ntitle: y\n---\nbody",             // no severity
	}
	for i, src := range cases {
		if _, err := parse([]byte(src)); err == nil {
			t.Errorf("case %d: parse() accepted missing required field", i)
		}
	}
}
