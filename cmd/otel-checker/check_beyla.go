package main

import (
	"strings"

	"github.com/grafana/otel-checker/checks/utils"

	"github.com/spf13/cobra"
)

func newCheckBeylaCmd(c *utils.Commands) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "beyla",
		Short:        "Check Beyla configuration",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cc *cobra.Command, _ []string) error {
			c.Components = []string{"beyla"}
			return runChecks(cc.Context(), *c)
		},
	}
	cmd.Flags().StringVar(&c.Language, "language", "",
		"Language used for instrumentation. Possible values: "+strings.Join(utils.SupportedLanguages, ", "))
	_ = cmd.MarkFlagRequired("language")
	_ = cmd.RegisterFlagCompletionFunc("language", staticCompletion(utils.SupportedLanguages))
	return cmd
}
