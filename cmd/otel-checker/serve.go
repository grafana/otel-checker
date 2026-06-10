package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/grafana/otel-checker/checks/utils"
	"github.com/grafana/otel-checker/checks/webserver"

	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v3"
)

// candidateNames lists the default result-file names that `serve` searches
// for when --data is not supplied, in priority order.
var candidateNames = []string{"results.json", "results.yaml", "results.yml"}

func newServeCmd() *cobra.Command {
	var (
		data   string
		listen string
	)
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve a results file (JSON or YAML) via the web UI, polling for updates",
		Long: `Serve a results file (JSON or YAML) via the web UI.

If --data is not provided, looks for ./results.json, ./results.yaml, or ./results.yml
in the current directory. The page auto-reloads every few seconds and picks up new
content the moment the file is written, so capturing results into the watched path
(e.g. "otel-checker check ... --format=json > results.json") refreshes the UI live.

If the file does not exist yet, the server still starts and shows a placeholder
message pointing at the expected path.`,
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cc *cobra.Command, _ []string) error {
			path := resolveDataPath(data)
			abs, err := filepath.Abs(path)
			if err != nil {
				abs = path
			}
			loader := func() webserver.Snapshot {
				return loadFileSnapshot(abs)
			}
			return webserver.Run(cc.Context(), listen, loader)
		},
	}
	cmd.Flags().StringVar(&data, "data", "",
		"Path to a JSON or YAML results file. When omitted, looks for ./results.json, ./results.yaml, or ./results.yml.")
	cmd.Flags().StringVar(&listen, "listen", utils.DefaultListen,
		"host:port the web server binds to")
	return cmd
}

// resolveDataPath returns the path the server should watch. If explicit is
// set, it wins. Otherwise we look for any of the default candidate names in
// the current directory; if none exists yet, we still return the canonical
// results.json so the placeholder page mentions a sensible path.
func resolveDataPath(explicit string) string {
	if explicit != "" {
		return explicit
	}
	for _, name := range candidateNames {
		if _, err := os.Stat(name); err == nil {
			return name
		}
	}

	return candidateNames[0]
}

// loadFileSnapshot reads path and decodes it as JSON or YAML based on its
// extension. A missing file is not an error — the returned Snapshot has
// Available=false and Source set to the path so the template can display
// the placeholder.
func loadFileSnapshot(path string) webserver.Snapshot {
	snap := webserver.Snapshot{
		Source: path,
		Reload: true,
	}
	results, err := readResultsFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			snap.Source = fmt.Sprintf("%s (parse error: %v)", path, err)
		}
		return snap
	}
	snap.Results = results
	snap.Available = true
	return snap
}

// readResultsFile opens path and decodes it as JSON or YAML based on its
// extension. Returns os.IsNotExist-compatible errors when the file is
// missing so callers can distinguish "no file" from "bad file".
func readResultsFile(path string) (utils.Results, error) {
	f, err := os.Open(path)
	if err != nil {
		return utils.Results{}, err
	}
	defer func() { _ = f.Close() }()

	var results utils.Results
	switch strings.ToLower(filepath.Ext(path)) {
	case ".yaml", ".yml":
		if err := yaml.NewDecoder(f).Decode(&results); err != nil {
			return utils.Results{}, fmt.Errorf("parse %s: %w", path, err)
		}
	default:
		if err := json.NewDecoder(f).Decode(&results); err != nil {
			return utils.Results{}, fmt.Errorf("parse %s: %w", path, err)
		}
	}
	return results, nil
}
