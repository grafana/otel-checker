package webserver

import (
	"html/template"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/grafana/otel-checker/checks/fixes"
)

// helper: build a test server that uses the package's real templates +
// embedded fixes registry, with the given snapshot loader.
func newTestServer(t *testing.T, loader Loader) *httptest.Server {
	t.Helper()
	tmpl, err := template.ParseFS(tmpls, "tmpl/*.tmpl")
	if err != nil {
		t.Fatalf("parse templates: %v", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/static/", http.FileServer(http.FS(fs.FS(static))))
	mux.HandleFunc("/fix/", func(w http.ResponseWriter, r *http.Request) {
		serveFix(w, r, tmpl)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if err := tmpl.ExecuteTemplate(w, "index.html.tmpl", loader()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	return httptest.NewServer(mux)
}

func TestFixPageKnownID(t *testing.T) {
	srv := newTestServer(t, Static(Snapshot{Available: true}))
	defer srv.Close()

	// Pick the first registered ID so the test stays valid as the registry grows.
	ids := fixes.All()
	if len(ids) == 0 {
		t.Fatal("no fix IDs registered")
	}
	id := ids[0]
	doc, _ := fixes.Lookup(id)

	resp, err := http.Get(srv.URL + "/fix/" + id)
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	body := readAll(t, resp)
	for _, want := range []string{id, doc.Title} {
		if !strings.Contains(body, want) {
			t.Errorf("body missing %q.\nGot:\n%s", want, body)
		}
	}
}

func TestFixPageUnknownID(t *testing.T) {
	srv := newTestServer(t, Static(Snapshot{Available: true}))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/fix/definitely.not.a.real.id")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", resp.StatusCode)
	}
}

func TestFixPathMissingID(t *testing.T) {
	srv := newTestServer(t, Static(Snapshot{Available: true}))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/fix/")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", resp.StatusCode)
	}
}

func TestDashboardOnlyAtRoot(t *testing.T) {
	srv := newTestServer(t, Static(Snapshot{Available: true}))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/random/garbage")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404 for non-root path", resp.StatusCode)
	}
}

func readAll(t *testing.T, resp *http.Response) string {
	t.Helper()
	buf := make([]byte, 0, 4096)
	tmp := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(tmp)
		if n > 0 {
			buf = append(buf, tmp[:n]...)
		}
		if err != nil {
			break
		}
	}
	return string(buf)
}
