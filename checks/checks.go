package checks

import (
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
// reporter.Results() to display or consume the results.
func Run(commands utils.Commands) *utils.Reporter {
	reporter := utils.Reporter{}

	env.CheckCommon(reporter.Component("Common Environment Variables"), commands.Language)

	for _, c := range commands.Components {
		switch c {
		case "sdk":
			SDKSetup(reporter.Component("SDK"), commands)
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
			grafana.CheckGrafanaSetup(reporter, reporter.Component("Grafana Cloud"), commands)
		}
	}

	return &reporter
}

func SDKSetup(reporter *utils.ComponentReporter, commands utils.Commands) {
	switch commands.Language {
	case "dotnet":
		dotnet.CheckDotNetSetup(reporter, commands)
	case "go":
		_go.CheckGoSetup(reporter, commands)
	case "java":
		java.CheckSetup(reporter, commands)
	case "js":
		js.CheckJSSetup(reporter, commands)
	case "python":
		python.CheckSetup(reporter, commands)
	case "ruby":
		sdk.CheckRubySetup(reporter, commands)
	case "php":
		sdk.CheckPHPSetup(reporter, commands)
	}
}
