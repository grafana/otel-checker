package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/grafana/otel-checker/checks/utils"
	"github.com/grafana/otel-checker/checks/webserver"

	"github.com/spf13/cobra"
)

func newServeCmd() *cobra.Command {
	var (
		data   string
		listen string
	)
	cmd := &cobra.Command{
		Use:          "serve",
		Short:        "Serve a previously-captured set of check results via the web UI",
		Long:         `Serve a previously-captured set of check results via the web UI. Reads the results map (JSON) from --data, or from stdin when --data=-. Pair it with "otel-checker check ... --format=json" to capture results, then serve them later.`,
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cc *cobra.Command, _ []string) error {
			results, err := loadResults(data)
			if err != nil {
				return err
			}
			return webserver.Run(cc.Context(), listen, results)
		},
	}
	cmd.Flags().StringVar(&data, "data", "",
		`Path to a JSON file with check results (as produced by --format=json). Use "-" to read from stdin.`)
	cmd.Flags().StringVar(&listen, "listen", utils.DefaultListen,
		"host:port the web server binds to")
	_ = cmd.MarkFlagRequired("data")
	return cmd
}

func loadResults(path string) (map[string][]string, error) {
	var r io.Reader
	if path == "-" {
		r = os.Stdin
	} else {
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("open %s: %w", path, err)
		}
		defer func() { _ = f.Close() }()
		r = f
	}
	var msgs map[string][]string
	if err := json.NewDecoder(r).Decode(&msgs); err != nil {
		return nil, fmt.Errorf("decode results: %w", err)
	}
	return msgs, nil
}
