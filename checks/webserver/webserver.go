// Package webserver renders otel-checker results as an HTML page served
// from a local HTTP listener. The CLI uses it both for the live --web-server
// flow on a `check` invocation and for the `serve` subcommand that polls a
// results file.
package webserver

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/grafana/otel-checker/checks/explain"
	"github.com/grafana/otel-checker/checks/utils"

	"github.com/gomarkdown/markdown"
)

//go:embed static/*
var static embed.FS

//go:embed tmpl/*
var tmpls embed.FS

// Snapshot is the data passed to the template for a single render.
// Loaders return this so the server can re-read state on every request.
type Snapshot struct {
	// Results is the populated reporter snapshot to display.
	Results utils.Results
	// Source is a human-readable description of where the snapshot came from
	// (e.g. a file path). Shown in the "no results" page.
	Source string
	// Available signals whether Results was successfully loaded. When false,
	// the template renders a placeholder pointing at Source.
	Available bool
	// Reload tells the template to add an HTML meta-refresh so the browser
	// polls for updates. Set by `serve`; not by one-shot `check --web-server`.
	Reload bool
}

// ComponentGroup is a per-component bundle of findings, used by the
// template to render section bodies without repeating the component name on
// every line. Each Item carries its own Message and (optional) ExplainID so the
// template can render a per-row "Explain" link.
type ComponentGroup struct {
	Component string
	Items     []utils.ComponentResult
}

// GroupedChecks returns Results.Checks bucketed by component, preserving
// the original ordering.
func (s Snapshot) GroupedChecks() []ComponentGroup { return groupResults(s.Results.Checks) }

// GroupedWarnings returns Results.Warnings bucketed by component.
func (s Snapshot) GroupedWarnings() []ComponentGroup { return groupResults(s.Results.Warnings) }

// GroupedErrors returns Results.Errors bucketed by component.
func (s Snapshot) GroupedErrors() []ComponentGroup { return groupResults(s.Results.Errors) }

func groupResults(items []utils.ComponentResult) []ComponentGroup {
	if len(items) == 0 {
		return nil
	}
	var groups []ComponentGroup
	indexByComponent := map[string]int{}
	for _, it := range items {
		idx, ok := indexByComponent[it.Component]
		if !ok {
			indexByComponent[it.Component] = len(groups)
			groups = append(groups, ComponentGroup{Component: it.Component, Items: []utils.ComponentResult{it}})
			continue
		}
		groups[idx].Items = append(groups[idx].Items, it)
	}
	return groups
}

// Loader returns the current Snapshot to display. Called on every HTTP
// request, so it should be cheap (e.g. a file read).
type Loader func() Snapshot

// Static wraps a fixed Snapshot in a Loader, for callers that don't need
// per-request updates.
func Static(s Snapshot) Loader { return func() Snapshot { return s } }

// explainView is the data passed to the explain-detail template.
type explainView struct {
	ID       string
	Title    string
	Severity string
	Body     template.HTML // pre-rendered markdown → HTML
}

// serveExplain handles GET /explain/<id>. Returns 404 if the ID is missing or unknown.
func serveExplain(w http.ResponseWriter, r *http.Request, t *template.Template) {
	id := strings.TrimPrefix(r.URL.Path, "/explain/")
	if id == "" || strings.Contains(id, "/") {
		http.NotFound(w, r)
		return
	}
	doc, ok := explain.Lookup(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	view := explainView{
		ID:       doc.ID,
		Title:    doc.Title,
		Severity: doc.Severity,
		Body:     template.HTML(markdown.ToHTML([]byte(doc.Body), nil, nil)), //nolint:gosec // body is trusted; comes from embedded docs at build time
	}
	if err := t.ExecuteTemplate(w, "explain.html.tmpl", view); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Run starts an HTTP server on addr that renders the snapshot returned by
// loader on each request. Blocks until ctx is cancelled or the server exits.
func Run(ctx context.Context, addr string, loader Loader) error {
	t, err := template.ParseFS(tmpls, "tmpl/*.tmpl")
	if err != nil {
		return fmt.Errorf("parse templates: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.FileServer(http.FS(static)))
	mux.HandleFunc("/explain/", func(w http.ResponseWriter, r *http.Request) {
		serveExplain(w, r, t)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		snap := loader()
		if err := t.ExecuteTemplate(w, "index.html.tmpl", snap); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	serveErr := make(chan error, 1)
	go func() {
		log.Printf("Application available on http://%s", addr)
		err := srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		serveErr <- err
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown: %w", err)
		}
		return <-serveErr
	case err := <-serveErr:
		return err
	}
}
