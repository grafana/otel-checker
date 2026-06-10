package explain_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/grafana/otel-checker/checks/explain"
)

// explainIDLiteral matches the string-literal first argument of any
// AddXxxWithExplain call. It deliberately matches conservatively
// (kebab-namespaced tokens only) to avoid grabbing format strings or
// unrelated literals.
var explainIDLiteral = regexp.MustCompile(`Add(?:Successful|Internal)?(?:Check|Warning|Error)WithExplain\("([a-z][a-z0-9-]*(?:\.[a-z0-9-]+)+)"`)

// TestEveryExplainIDUsedInCodeIsRegistered walks the checks/ source tree,
// finds every AddXxxWithExplain call-site, and asserts the literal explain
// ID argument resolves through explain.Lookup. Catches typos in call sites
// at test time instead of at runtime.
func TestEveryExplainIDUsedInCodeIsRegistered(t *testing.T) {
	root, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	seen := map[string]struct{}{}
	walkErr := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		// Skip the explain package itself — it has no AddXxxWithExplain
		// callers, and the regex would otherwise match anything resembling
		// one.
		if strings.Contains(path, "/checks/explain/") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, m := range explainIDLiteral.FindAllStringSubmatch(string(data), -1) {
			seen[m[1]] = struct{}{}
		}
		return nil
	})
	if walkErr != nil {
		t.Fatalf("walk: %v", walkErr)
	}
	if len(seen) == 0 {
		t.Fatal("regex matched no explain IDs anywhere — either no call sites were migrated or the regex is broken")
	}
	for id := range seen {
		if _, ok := explain.Lookup(id); !ok {
			t.Errorf("call-site references explain ID %q but no docs/%s.md is registered", id, id)
		}
	}
}
