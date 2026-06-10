// Package explain is the registry of fix documentation for findings emitted
// by the checks/* packages. Each markdown file under docs/ describes a
// single fix and is keyed by a stable kebab-namespaced ID
// (e.g. "env.otel-service-name.unset"). The CLI's `explain` subcommand and
// the web UI's `/explain/{id}` route both look up documents through this
// package; downstream consumers can import Lookup and All to build their
// own surfaces.
package explain

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"go.yaml.in/yaml/v3"
)

//go:embed docs/*.md
var docsFS embed.FS

// Doc is the structured form of a single fix entry. Body holds the raw
// markdown beneath the front-matter block.
type Doc struct {
	ID       string
	Title    string
	Severity string
	Body     string
}

// frontMatter is the YAML block at the top of each doc file.
type frontMatter struct {
	ID       string `yaml:"id"`
	Title    string `yaml:"title"`
	Severity string `yaml:"severity"`
}

var registry = map[string]Doc{}

func init() {
	if err := load(docsFS); err != nil {
		panic(fmt.Sprintf("explain: failed to load doc registry: %v", err))
	}
}

// load walks the embedded docs directory, parses every markdown file, and
// populates the registry. Returns an error on duplicate IDs or malformed
// front-matter so init() can fail loudly at startup.
func load(fsys fs.FS) error {
	entries, err := fs.ReadDir(fsys, "docs")
	if err != nil {
		return fmt.Errorf("read docs dir: %w", err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		path := "docs/" + e.Name()
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		doc, err := parse(data)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		if _, dup := registry[doc.ID]; dup {
			return fmt.Errorf("duplicate fix id %q (file %s)", doc.ID, path)
		}
		registry[doc.ID] = doc
	}
	return nil
}

// parse extracts the YAML front-matter and body from a markdown file.
// Files MUST start with a "---\n" line, contain a closing "---\n" within the
// first ~20 lines, and have non-empty id/title/severity fields.
func parse(data []byte) (Doc, error) {
	s := string(data)
	if !strings.HasPrefix(s, "---\n") {
		return Doc{}, fmt.Errorf("missing front-matter (file must start with ---)")
	}
	rest := s[len("---\n"):]
	end := strings.Index(rest, "\n---\n")
	if end < 0 {
		return Doc{}, fmt.Errorf("missing closing --- on front-matter")
	}
	var fm frontMatter
	if err := yaml.Unmarshal([]byte(rest[:end]), &fm); err != nil {
		return Doc{}, fmt.Errorf("yaml unmarshal: %w", err)
	}
	if fm.ID == "" {
		return Doc{}, fmt.Errorf("front-matter: missing id")
	}
	if fm.Title == "" {
		return Doc{}, fmt.Errorf("front-matter: missing title")
	}
	if fm.Severity == "" {
		return Doc{}, fmt.Errorf("front-matter: missing severity")
	}
	body := strings.TrimLeft(rest[end+len("\n---\n"):], "\n")
	return Doc{ID: fm.ID, Title: fm.Title, Severity: fm.Severity, Body: body}, nil
}

// Lookup returns the registered Doc for id, or (Doc{}, false) if no such
// fix is known.
func Lookup(id string) (Doc, bool) {
	d, ok := registry[id]
	return d, ok
}

// All returns every registered fix ID, sorted alphabetically. Useful for
// shell completion and the `explain list` subcommand.
func All() []string {
	out := make([]string, 0, len(registry))
	for id := range registry {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}
