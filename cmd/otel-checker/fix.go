package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/grafana/otel-checker/checks/fixes"
	"github.com/grafana/otel-checker/checks/utils"

	"github.com/spf13/cobra"
)

// newFixCmd returns the `fix` parent command:
//   - With an ID positional arg: prints the fix doc for that ID.
//   - Without args: reads a results file (./results.json by default) and
//     prints fix docs for every distinct ID flagged on an error or warning.
//   - `fix list` (child subcommand): prints every registered ID.
func newFixCmd() *cobra.Command {
	var data string
	cmd := &cobra.Command{
		Use:   "fix [id]",
		Short: "Show fix guidance for findings",
		Long: `Show fix guidance for findings emitted by ` + "`otel-checker check`" + `.

With an ID argument, prints the fix doc for that ID. With no arguments,
reads a results file and prints the fix doc for every distinct ID flagged
on an error or warning, deduplicating along the way.

When invoked with no ID and --data is omitted, looks for ./results.json,
./results.yaml, or ./results.yml in the current directory — the same lookup
as ` + "`otel-checker serve`" + `. Use ` + "`fix list`" + ` to see every available ID.`,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
		ValidArgsFunction: func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
			return fixes.All(), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cc *cobra.Command, args []string) error {
			if len(args) == 1 {
				return runFixShow(cc.OutOrStdout(), args[0])
			}
			return runFixAll(cc.OutOrStdout(), data)
		},
	}
	cmd.Flags().StringVar(&data, "data", "",
		"Path to a JSON or YAML results file. Only used when no ID is given. When omitted, looks for ./results.json, ./results.yaml, or ./results.yml.")
	cmd.AddCommand(newFixListCmd())
	return cmd
}

func newFixListCmd() *cobra.Command {
	return &cobra.Command{
		Use:          "list",
		Short:        "List every available fix ID with its title",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cc *cobra.Command, _ []string) error {
			return runFixList(cc.OutOrStdout())
		},
	}
}

func runFixShow(w io.Writer, id string) error {
	doc, ok := fixes.Lookup(id)
	if !ok {
		return fmt.Errorf("unknown fix ID %q. Run \"otel-checker fix list\" to see every available ID", id)
	}
	_, _ = fmt.Fprintf(w, "# %s\n\n", doc.Title)
	_, _ = fmt.Fprint(w, doc.Body)
	return nil
}

func runFixList(w io.Writer) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for _, id := range fixes.All() {
		doc, _ := fixes.Lookup(id)
		_, _ = fmt.Fprintf(tw, "%s\t%s\n", id, doc.Title)
	}
	return tw.Flush()
}

func runFixAll(w io.Writer, dataPath string) error {
	path := resolveDataPath(dataPath)
	results, err := readResultsFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no results file found at %s. Capture results first, e.g. `otel-checker check ... --format=json > %s`", path, path)
		}
		return err
	}

	ids := flaggedFixIDs(results)
	if len(ids) == 0 {
		_, _ = fmt.Fprintf(w, "No fix IDs found in %s — either nothing failed or the findings predate the fix-ID system.\n", path)
		return nil
	}

	_, _ = fmt.Fprintf(w, "Showing %d fix(es) for findings in %s:\n", len(ids), path)
	separator := strings.Repeat("─", 72)
	for _, id := range ids {
		_, _ = fmt.Fprintln(w)
		_, _ = fmt.Fprintln(w, separator)
		_, _ = fmt.Fprintf(w, "[%s]\n\n", id)
		if err := runFixShow(w, id); err != nil {
			// Shouldn't happen — the coverage test guarantees every emitted
			// ID resolves — but degrade gracefully in case the registry and
			// results file disagree (e.g. results from a newer binary).
			_, _ = fmt.Fprintf(w, "(could not load fix doc: %v)\n", err)
		}
	}
	return nil
}

// flaggedFixIDs returns the distinct fix IDs that appear on any error or
// warning result, preserving first-occurrence ordering (errors first, then
// warnings) so output matches the order of the original check.
func flaggedFixIDs(results utils.Results) []string {
	seen := map[string]struct{}{}
	var out []string
	collect := func(items []utils.ComponentResult) {
		for _, r := range items {
			if r.FixID == "" {
				continue
			}
			if _, ok := seen[r.FixID]; ok {
				continue
			}
			seen[r.FixID] = struct{}{}
			out = append(out, r.FixID)
		}
	}
	collect(results.Errors)
	collect(results.Warnings)
	return out
}
