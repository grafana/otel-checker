package checks

import (
	"context"
	"slices"

	"github.com/grafana/otel-checker/checks/alloy"
	"github.com/grafana/otel-checker/checks/beyla"
	"github.com/grafana/otel-checker/checks/collector"
	"github.com/grafana/otel-checker/checks/config"
	"github.com/grafana/otel-checker/checks/env"
	"github.com/grafana/otel-checker/checks/grafana"
	"github.com/grafana/otel-checker/checks/sdk"
	"github.com/grafana/otel-checker/checks/sdk/dotnet"
	_go "github.com/grafana/otel-checker/checks/sdk/go"
	"github.com/grafana/otel-checker/checks/sdk/java"
	"github.com/grafana/otel-checker/checks/sdk/js"
	"github.com/grafana/otel-checker/checks/sdk/python"
	"github.com/grafana/otel-checker/checks/utils"
)

// Run executes all configured checks and returns the populated reporter.
// It does not produce any output; pair it with output.Render or
// reporter.Results() to display or consume the results. ctx is plumbed to
// HTTP and exec calls so callers can cancel a slow check.
func Run(ctx context.Context, commands utils.Commands) *utils.Reporter {
	reporter := utils.Reporter{}

	// Declarative config parsing is done once, up front, so both the
	// `config` component (which reports on the parse itself) and the
	// `grafana-cloud` component (which reads endpoints from it when
	// present) can share the result.
	resolvedConfigPath, configCandidates := config.Resolve(commands.ConfigPath)
	var (
		parsedConfig  *config.File
		configLoadErr error
	)
	if resolvedConfigPath != "" {
		parsedConfig, configLoadErr = config.Load(resolvedConfigPath)
	} else {
		configLoadErr = errNoConfigFile
	}

	env.CheckCommon(reporter.Component("Common Environment Variables"), commands.Language)

	for _, c := range commands.Components {
		switch c {
		case "sdk":
			SDKSetup(ctx, reporter.Component("SDK"), commands)
		case "beyla":
			beyla.CheckBeylaSetup(reporter.Component("Beyla"), commands.Language)
		case "alloy":
			alloy.CheckAlloySetup(reporter.Component("Alloy"), commands.Language)
		case "collector":
			collector.CheckCollectorSetup(
				reporter.Component("Collector"),
				commands.Language,
				commands.CollectorConfigPath,
			)
		case "grafana-cloud":
			grafana.CheckGrafanaSetup(ctx, reporter, reporter.Component("Grafana Cloud"), commands, parsedConfig)
		case "config":
			config.CheckConfigSetup(
				reporter.Component("Declarative Config"),
				resolvedConfigPath,
				configCandidates,
				parsedConfig,
				configLoadErr,
			)
			// Endpoint format validation is the grafana check's job, but
			// it lives under the same file the config component just
			// parsed — so drive it from here when grafana-cloud isn't
			// already in the components list (which would run it too and
			// duplicate every finding).
			if parsedConfig != nil && !slices.Contains(commands.Components, "grafana-cloud") {
				grafana.CheckEndpointsFromConfig(reporter.Component("Grafana Cloud"), parsedConfig)
			}
		}
	}

	return &reporter
}

// errNoConfigFile is the sentinel returned to CheckConfigSetup when neither
// an explicit --config-path nor a default file could be located. It's
// distinguished from a real load error by the empty resolvedConfigPath.
var errNoConfigFile = errNoFile{}

type errNoFile struct{}

func (errNoFile) Error() string { return "no declarative config file found" }

func SDKSetup(ctx context.Context, reporter *utils.ComponentReporter, commands utils.Commands) {
	switch commands.Language {
	case "dotnet":
		dotnet.CheckDotNetSetup(ctx, reporter, commands)
	case "go":
		_go.CheckGoSetup(ctx, reporter, commands)
	case "java":
		java.CheckSetup(ctx, reporter, commands)
	case "js":
		js.CheckJSSetup(ctx, reporter, commands)
	case "python":
		python.CheckSetup(ctx, reporter, commands)
	case "ruby":
		sdk.CheckRubySetup(ctx, reporter, commands)
	case "php":
		sdk.CheckPHPSetup(ctx, reporter, commands)
	}
}
