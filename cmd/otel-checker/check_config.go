package main

import (
	"github.com/grafana/otel-checker/checks/utils"

	"github.com/spf13/cobra"
)

func newCheckConfigCmd(c *utils.Commands) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "config",
		Short:        "Check the OpenTelemetry declarative configuration file",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cc *cobra.Command, _ []string) error {
			c.Components = []string{"config"}
			return runChecks(cc.Context(), *c)
		},
	}
	cmd.Flags().StringVar(&c.ConfigPath, "config-path", "",
		"Full path to the declarative config file. If unset, looks for otel-config.yaml then otel-config.yml in the current directory.")
	return cmd
}
