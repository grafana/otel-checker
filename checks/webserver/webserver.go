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
	"time"

	"github.com/grafana/otel-checker/checks/utils"
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

// ComponentGroup is a per-component bundle of messages, used by the
// template to render section bodies without repeating the component name on
// every line.
type ComponentGroup struct {
	Component string
	Messages  []string
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
			groups = append(groups, ComponentGroup{Component: it.Component, Messages: []string{it.Message}})
			continue
		}
		groups[idx].Messages = append(groups[idx].Messages, it.Message)
	}
	return groups
}

// Loader returns the current Snapshot to display. Called on every HTTP
// request, so it should be cheap (e.g. a file read).
type Loader func() Snapshot

// Static wraps a fixed Snapshot in a Loader, for callers that don't need
// per-request updates.
func Static(s Snapshot) Loader { return func() Snapshot { return s } }

// Run starts an HTTP server on addr that renders the snapshot returned by
// loader on each request. Blocks until ctx is cancelled or the server exits.
func Run(ctx context.Context, addr string, loader Loader) error {
	t, err := template.ParseFS(tmpls, "tmpl/*.tmpl")
	if err != nil {
		return fmt.Errorf("parse templates: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.FileServer(http.FS(static)))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
