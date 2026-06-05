// Package webserver renders otel-checker results as an HTML page served
// from a local HTTP listener. The CLI uses this when --web-server is set;
// library callers can invoke Run directly with a custom address and the
// Reporter.Results() map.
package webserver

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

//go:embed static/*
var static embed.FS

//go:embed tmpl/*
var tmpls embed.FS

// Run starts an HTTP server on addr that renders the given messages.
// Blocks until the server exits.
func Run(addr string, messages map[string][]string) error {
	t, err := template.ParseFS(tmpls, "tmpl/*.tmpl")
	if err != nil {
		return fmt.Errorf("parse templates: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.FileServer(http.FS(static)))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := t.ExecuteTemplate(w, "index.html.tmpl", struct {
			Messages map[string][]string
		}{Messages: messages}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Printf("Application available on http://%s", addr)
	return http.ListenAndServe(addr, mux)
}
