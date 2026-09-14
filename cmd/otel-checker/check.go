package main

import (
	"strings"

	"github.com/grafana/otel-checker/checks/utils"

	"github.com/spf13/cobra"
)

// newCheckCmd builds the `check` parent command. It is itself runnable
// and hosts one subcommand per component.
//
// Components to check are supplied either via a positional comma-separated
// list (`check sdk,collector,beyla`) or via the per-component subcommands
// (`check sdk`, `check collector`, ...). When no components are supplied,
// every component is checked.
//
// Note: space-separated positional args (`check sdk collector beyla`) won't
// work — the first token would resolve to the `sdk` subcommand.
func newCheckCmd() *cobra.Command {
	c := &utils.Commands{}

	cmd := &cobra.Command{
		Use:          "check [components]",
		Short:        "Run instrumentation checks",
		Long:         "Run one or more component checks. Pass a comma-separated component list (`check sdk,collector,beyla`) or use the per-component subcommand. With no positional argument, every component is checked.",
		Args:         cobra.ArbitraryArgs,
		SilenceUsage: true,
		RunE: func(cc *cobra.Command, args []string) error {
			c.Components = parseComponentArgs(args)
			if len(c.Components) == 0 {
				c.Components = append([]string{}, utils.SupportedComponents...)
			}
			return runChecks(cc.Context(), *c)
		},
		ValidArgsFunction: func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
			return utils.SupportedComponents, cobra.ShellCompDirectiveNoFileComp
		},
	}

	bindOutputFlags(cmd, c)

	f := cmd.Flags()
	f.StringVar(&c.Language, "language", "",
		"Language used for instrumentation. Possible values: "+strings.Join(utils.SupportedLanguages, ", "))
	f.BoolVar(&c.ManualInstrumentation, "manual-instrumentation", false,
		"Use manual instrumentation (auto-instrumentation is the default)")
	f.StringVar(&c.InstrumentationFile, "instrumentation-file", "",
		`Path to the instrumentation file. Required with --manual-instrumentation for JS, e.g. "src/inst/instrumentation.js"`)
	f.StringVar(&c.PackageJsonPath, "package-json-path", "",
		"Path to the directory containing package.json (JS only)")
	f.StringVar(&c.CollectorConfigPath, "collector-config-path", "",
		"Full path to the Collector config file. If unset, looks for config.yaml then config.yml in the current directory.")
	f.StringVar(&c.ConfigPath, "config-path", "",
		"Full path to the declarative config file. If unset, looks for otel-config.yaml then otel-config.yml in the current directory.")

	_ = cmd.RegisterFlagCompletionFunc("language", staticCompletion(utils.SupportedLanguages))

	cmd.AddCommand(
		newCheckSDKCmd(c),
		newCheckCollectorCmd(c),
		newCheckBeylaCmd(c),
		newCheckAlloyCmd(c),
		newCheckGrafanaCloudCmd(c),
		newCheckConfigCmd(c),
	)

	return cmd
}

// parseComponentArgs flattens positional arguments into a clean component
// list. Each argument is split on commas and trimmed; empty entries are
// dropped. Accepts both `check sdk,collector,beyla` and `check sdk collector beyla`.
func parseComponentArgs(args []string) []string {
	var out []string
	for _, a := range args {
		for _, s := range strings.Split(a, ",") {
			s = strings.TrimSpace(s)
			if s != "" {
				out = append(out, s)
			}
		}
	}
	return out
}
