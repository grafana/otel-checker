package checks

import (
	"otel-checker/checks/alloy"
	"otel-checker/checks/beyla"
	"otel-checker/checks/collector"
	"otel-checker/checks/grafana"
	"otel-checker/checks/sdk"
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
			sdk.CheckSDKSetup(
				reporter.Component("SDK"),
				commands.Language,
				commands.AutoInstrumentation,
				commands.PackageJsonPath,
				commands.InstrumentationFile,
				commands.Debug,
			)
		}
	}

	return reporter.PrintResults()
}
