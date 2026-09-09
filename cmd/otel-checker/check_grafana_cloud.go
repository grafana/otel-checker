package main

import (
	"strings"

	"github.com/grafana/otel-checker/checks/utils"

	"github.com/spf13/cobra"
)

func newCheckGrafanaCloudCmd(c *utils.Commands) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "grafana-cloud",
		Short:        "Check Grafana Cloud connectivity and credentials",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cc *cobra.Command, _ []string) error {
			c.Components = []string{"grafana-cloud"}
			return runChecks(cc.Context(), *c)
		},
	}
	cmd.Flags().StringVar(&c.Language, "language", "",
		"Language used for instrumentation. Optional — only affects the OTEL_EXPORTER_OTLP_HEADERS format (Python expects Basic%20, others expect a literal space). Possible values: "+strings.Join(utils.SupportedLanguages, ", "))
	_ = cmd.RegisterFlagCompletionFunc("language", staticCompletion(utils.SupportedLanguages))
	return cmd
}
