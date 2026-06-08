package main

import (
	"github.com/grafana/otel-checker/checks/utils"

	"github.com/spf13/cobra"
)

func newCheckCollectorCmd(c *utils.Commands) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "collector",
		Short:        "Check the OpenTelemetry Collector config.yaml",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cc *cobra.Command, _ []string) error {
			c.Components = []string{"collector"}
			return runChecks(cc.Context(), *c)
		},
	}
	cmd.Flags().StringVar(&c.CollectorConfigPath, "collector-config-path", "",
		"Path to the directory containing the collector's config.yaml")
	return cmd
}
