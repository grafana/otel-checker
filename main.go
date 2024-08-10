package main

import (
	"otel-checker/alloy"
	"otel-checker/beyla"
	"otel-checker/collector"
	"otel-checker/grafana"
	"otel-checker/sdk"
	"otel-checker/utils"
)

func main() {
	messages := utils.CreateMessagesMap()
	commands := utils.GetArguments()

	grafana.CheckGrafanaSetup(&messages, commands.Language, commands.Components)

	for _, c := range commands.Components {
		if c == "alloy" {
			alloy.CheckAlloySetup(&messages, commands.Language)
		}

		if c == "beyla" {
			beyla.CheckBeylaSetup(&messages, commands.Language)
		}

		if c == "collector" {
			collector.CheckCollectorSetup(
				&messages,
				commands.Language,
				commands.CollectorConfigPath,
			)
		}

		if c == "sdk" {
			sdk.CheckSDKSetup(
				&messages,
				commands.Language,
				commands.AutoInstrumentation,
				commands.PackageJsonPath,
				commands.InstrumentationFile,
			)
		}
	}

	utils.PrintResults(messages)
}
