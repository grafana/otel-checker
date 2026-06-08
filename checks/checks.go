package checks

import (
	"context"

	"github.com/grafana/otel-checker/checks/alloy"
	"github.com/grafana/otel-checker/checks/beyla"
	"github.com/grafana/otel-checker/checks/collector"
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
			grafana.CheckGrafanaSetup(ctx, reporter, reporter.Component("Grafana Cloud"), commands)
		}
	}

	return &reporter
}

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
