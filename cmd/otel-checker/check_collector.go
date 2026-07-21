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
		"Full path to the Collector config file. If unset, looks for config.yaml then config.yml in the current directory.")
	return cmd
}
