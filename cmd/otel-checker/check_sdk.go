package main

import (
	"strings"

	"github.com/grafana/otel-checker/checks/utils"

	"github.com/spf13/cobra"
)

func newCheckSDKCmd(c *utils.Commands) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "sdk",
		Short:        "Check the language SDK setup",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cc *cobra.Command, _ []string) error {
			c.Components = []string{"sdk"}
			return runChecks(cc.Context(), *c)
		},
	}
	f := cmd.Flags()
	f.StringVar(&c.Language, "language", "",
		"Language used for instrumentation. Possible values: "+strings.Join(utils.SupportedLanguages, ", "))
	f.BoolVar(&c.ManualInstrumentation, "manual-instrumentation", false,
		"Use manual instrumentation (auto-instrumentation is the default)")
	f.StringVar(&c.InstrumentationFile, "instrumentation-file", "",
		`Path to the instrumentation file. Required with --manual-instrumentation for JS, e.g. "src/inst/instrumentation.js"`)
	f.StringVar(&c.PackageJsonPath, "package-json-path", "",
		"Path to the directory containing package.json (JS only)")
	_ = cmd.MarkFlagRequired("language")
	_ = cmd.RegisterFlagCompletionFunc("language", staticCompletion(utils.SupportedLanguages))
	return cmd
}
