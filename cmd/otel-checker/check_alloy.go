package main

import (
	"strings"

	"github.com/grafana/otel-checker/checks/utils"

	"github.com/spf13/cobra"
)

func newCheckAlloyCmd(c *utils.Commands) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "alloy",
		Short:        "Check Grafana Alloy configuration",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cc *cobra.Command, _ []string) error {
			c.Components = []string{"alloy"}
			return runChecks(cc.Context(), *c)
		},
	}
	cmd.Flags().StringVar(&c.Language, "language", "",
		"Language used for instrumentation. Possible values: "+strings.Join(utils.SupportedLanguages, ", "))
	_ = cmd.MarkFlagRequired("language")
	_ = cmd.RegisterFlagCompletionFunc("language", staticCompletion(utils.SupportedLanguages))
	return cmd
}
