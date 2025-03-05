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
			CheckSDKSetup(
				reporter.Component("SDK"),
				commands.Language,
				commands.ManualInstrumentation,
				commands.PackageJsonPath,
				commands.InstrumentationFile,
				commands.Debug,
			)
		}
	}

	return reporter.PrintResults()
}

func CheckSDKSetup(reporter *utils.ComponentReporter, language string, autoInstrumentation bool, packageJsonPath string, instrumentationFile string, debug bool) {
	switch language {
	case "dotnet":
		sdk.CheckDotNetSetup(reporter, autoInstrumentation)
	case "go":
		sdk.CheckGoSetup(reporter, autoInstrumentation)
	case "java":
		java.CheckSetup(reporter, autoInstrumentation, debug)
	case "js":
		sdk.CheckJSSetup(reporter, autoInstrumentation, packageJsonPath, instrumentationFile)
	case "python":
		python.CheckSetup(reporter, autoInstrumentation, debug)
	case "ruby":
		sdk.CheckRubySetup(reporter, autoInstrumentation)
	}
}
