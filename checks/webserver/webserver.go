// Package webserver renders otel-checker results as an HTML page served
// from a local HTTP listener. The CLI uses this when --web-server is set;
// library callers can invoke Run directly with a custom address and the
// Reporter.Results() value.
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

// Run starts an HTTP server on addr that renders the given results.
// Blocks until ctx is cancelled or the server exits.
func Run(ctx context.Context, addr string, results utils.Results) error {
	t, err := template.ParseFS(tmpls, "tmpl/*.tmpl")
	if err != nil {
		return fmt.Errorf("parse templates: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.FileServer(http.FS(static)))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := t.ExecuteTemplate(w, "index.html.tmpl", struct {
			Results utils.Results
		}{Results: results}); err != nil {
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
