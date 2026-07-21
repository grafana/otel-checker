package main

import (
	"context"
	"os"
	"strings"

	"github.com/grafana/otel-checker/checks"
	"github.com/grafana/otel-checker/checks/output"
	"github.com/grafana/otel-checker/checks/utils"
	"github.com/grafana/otel-checker/checks/webserver"

	"github.com/spf13/cobra"
)

// runChecks runs the validated check configuration, prints the result in the
// chosen format, and optionally launches the web server. Used by every
// check-* subcommand (and the legacy root alias).
func runChecks(ctx context.Context, c utils.Commands) error {
	if err := utils.Validate(c); err != nil {
		return err
	}
	if c.PackageJsonPath != "" && !strings.HasSuffix(c.PackageJsonPath, "/") {
		c.PackageJsonPath += "/"
	}
	reporter := checks.Run(ctx, c)
	if err := output.Render(os.Stdout, reporter, c.Format); err != nil {
		return err
	}
	if !c.WebServer {
		return nil
	}
	return webserver.Run(ctx, c.Listen, webserver.Static(webserver.Snapshot{
		Results:   reporter.Results(),
		Source:    "live check results",
		Available: true,
	}))
}

func staticCompletion(values []string) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return values, cobra.ShellCompDirectiveNoFileComp
	}
}

// bindOutputFlags attaches the four presentation flags (format, debug,
// web-server, listen) as persistent flags so every descendant of the
// `check` command tree inherits them.
func bindOutputFlags(cmd *cobra.Command, c *utils.Commands) {
	f := cmd.PersistentFlags()
	f.StringVar(&c.Format, "format", "text",
		"Output format. Possible values: "+strings.Join(utils.SupportedFormats, ", "))
	f.BoolVar(&c.Debug, "debug", false, "Output debug information")
	f.BoolVar(&c.WebServer, "web-server", false, "Serve the results via a local web server")
	f.StringVar(&c.Listen, "listen", utils.DefaultListen,
		"host:port the web server binds to when --web-server is set")
	_ = cmd.RegisterFlagCompletionFunc("format", staticCompletion(utils.SupportedFormats))
}
