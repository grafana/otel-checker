package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:          "otel-checker",
		Short:        "Validate OpenTelemetry instrumentation in a repository",
		Long:         "otel-checker scans a code repository, checks environment variables, validates Grafana Cloud tokens, and verifies correct SDK/collector/Beyla/Alloy configuration.",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
	}

	root.AddCommand(newCheckCmd())
	root.AddCommand(newServeCmd())
	root.AddCommand(newExplainCmd())
	root.AddCommand(newVersionCmd())

	return root
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := newRootCmd().ExecuteContext(ctx); err != nil {
		// cobra already printed the error; exit non-zero so callers see it.
		os.Exit(1)
	}
}
