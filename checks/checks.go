package checks

import (
	"otel-checker/checks/alloy"
	"otel-checker/checks/beyla"
	"otel-checker/checks/collector"
	"otel-checker/checks/grafana"
	"otel-checker/checks/sdk"
	"otel-checker/checks/sdk/java"
	"otel-checker/checks/sdk/python"
	"otel-checker/checks/utils"
)

func RunAllChecks(commands utils.Commands) map[string][]string {
	reporter := utils.Reporter{}

	grafana.CheckGrafanaSetup(reporter, reporter.Component("Grafana Cloud"), commands.Language, commands.Components)

	for _, c := range commands.Components {
		if c == "alloy" {
			alloy.CheckAlloySetup(reporter.Component("Alloy"), commands.Language)
		}

		if c == "beyla" {
			beyla.CheckBeylaSetup(reporter.Component("Beyla"), commands.Language)
		}

		if c == "collector" {
			collector.CheckCollectorSetup(
				reporter.Component("Collector"),
				commands.Language,
				commands.CollectorConfigPath,
			)
		}

		if c == "sdk" {
			SDKSetup(reporter.Component("SDK"), commands)
		}
	}

	return reporter.PrintResults()
}

func SDKSetup(reporter *utils.ComponentReporter, commands utils.Commands) {
	switch commands.Language {
	case "dotnet":
		sdk.CheckDotNetSetup(reporter, commands)
	case "go":
		sdk.CheckGoSetup(reporter, commands)
	case "java":
		java.CheckSetup(reporter, commands)
	case "js":
		sdk.CheckJSSetup(reporter, commands)
	case "python":
		python.CheckSetup(reporter, commands)
	case "ruby":
		sdk.CheckRubySetup(reporter, commands)
	}
}
