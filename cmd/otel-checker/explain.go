package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/grafana/otel-checker/checks/explain"
	"github.com/grafana/otel-checker/checks/utils"

	"charm.land/glamour/v2"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// newExplainCmd returns the `explain` parent command:
//   - With an ID positional arg: prints the explain doc for that ID.
//   - Without args: reads a results file (./results.json by default) and
//     prints explain docs for every distinct ID flagged on an error or warning.
//   - `explain list` (child subcommand): prints every registered ID.
func newExplainCmd() *cobra.Command {
	var data string
	cmd := &cobra.Command{
		Use:   "explain [id]",
		Short: "Show explanation for findings",
		Long: `Show explanation for findings emitted by ` + "`otel-checker check`" + `.

With an ID argument, prints the explain doc for that ID. With no arguments,
reads a results file and prints the explain doc for every distinct ID flagged
on an error or warning, deduplicating along the way.

When invoked with no ID and --data is omitted, looks for ./results.json,
./results.yaml, or ./results.yml in the current directory — the same lookup
as ` + "`otel-checker serve`" + `. Use ` + "`explain list`" + ` to see every available ID.`,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		ValidArgsFunction: func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
			return explain.All(), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cc *cobra.Command, args []string) error {
			if len(args) == 1 {
				return runExplainShow(cc.OutOrStdout(), args[0])
			}
			return runExplainAll(cc.OutOrStdout(), data)
		},
	}
	cmd.Flags().StringVar(&data, "data", "",
		"Path to a JSON or YAML results file. Only used when no ID is given. When omitted, looks for ./results.json, ./results.yaml, or ./results.yml.")
	cmd.AddCommand(newExplainListCmd())
	return cmd
}

func newExplainListCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "list",
		Short:        "List every available explain ID with its title",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cc *cobra.Command, _ []string) error {
			return runExplainList(cc.OutOrStdout())
		},
	}
}

func runExplainShow(w io.Writer, id string) error {
	doc, ok := explain.Lookup(id)
	if !ok {
		return fmt.Errorf("unknown explain ID %q. Run \"otel-checker explain list\" to see every available ID", id)
	}
	return renderMarkdown(w, fmt.Sprintf("# %s\n\n%s", doc.Title, doc.Body))
}

// renderMarkdown writes source to w. When w is a terminal, the markdown is
// styled through glamour (headers, links, code blocks, lists). When w is a
// pipe, file, or in-memory buffer, the raw markdown is written instead so
// piping (`otel-checker explain id | grep`) and tests stay clean.
func renderMarkdown(w io.Writer, source string) error {
	if f, ok := w.(*os.File); ok && term.IsTerminal(int(f.Fd())) {
		if out, err := styleMarkdown(source); err == nil {
			_, err := fmt.Fprint(w, out)
			return err
		}
		// Fall through to raw on render error.
	}
	_, err := fmt.Fprint(w, source)
	return err
}

// styleMarkdown is split out so tests can assert the glamour call still
// works — renderMarkdown swallows errors and falls back to raw output, which
// would otherwise hide a glamour API break (e.g. a removed style name).
func styleMarkdown(source string) (string, error) {
	return glamour.RenderWithEnvironmentConfig(source)
}

func runExplainList(w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, id := range explain.All() {
		doc, _ := explain.Lookup(id)
		_, _ = fmt.Fprintf(tw, "%s\t%s\n", id, doc.Title)
	}
	return tw.Flush()
}

func runExplainAll(w io.Writer, dataPath string) error {
	path := resolveDataPath(dataPath)
	results, err := readResultsFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no results file found at %s. Capture results first, e.g. `otel-checker check ... --format=json > %s`", path, path)
		}
		return err
	}

	ids := flaggedExplainIDs(results)
	if len(ids) == 0 {
		_, _ = fmt.Fprintf(w, "No explain IDs found in %s — either nothing failed or the findings predate the explain-ID system.\n", path)
		return nil
	}

	_, _ = fmt.Fprintf(w, "Showing %d explanation(s) for findings in %s:\n", len(ids), path)
	separator := strings.Repeat("─", 72)
	for _, id := range ids {
		_, _ = fmt.Fprintln(w)
		_, _ = fmt.Fprintln(w, separator)
		_, _ = fmt.Fprintf(w, "[%s]\n\n", id)
		if err := runExplainShow(w, id); err != nil {
			// Shouldn't happen — the coverage test guarantees every emitted
			// ID resolves — but degrade gracefully in case the registry and
			// results file disagree (e.g. results from a newer binary).
			_, _ = fmt.Fprintf(w, "(could not load explain doc: %v)\n", err)
		}
	}
	return nil
}

// flaggedExplainIDs returns the distinct explain IDs that appear on any error or
// warning result, preserving first-occurrence ordering (errors first, then
// warnings) so output matches the order of the original check.
func flaggedExplainIDs(results utils.Results) []string {
	seen := map[string]struct{}{}
	var out []string
	collect := func(items []utils.ComponentResult) {
		for _, r := range items {
			if r.ExplainID == "" {
				continue
			}
			if _, ok := seen[r.ExplainID]; ok {
				continue
			}
			seen[r.ExplainID] = struct{}{}
			out = append(out, r.ExplainID)
		}
	}
	collect(results.Errors)
	collect(results.Warnings)
	return out
}
