package checks

import (
	"otel-checker/checks/alloy"
	"otel-checker/checks/beyla"
	"otel-checker/checks/collector"
	"otel-checker/checks/grafana"
	"otel-checker/checks/sdk"
	"otel-checker/checks/utils"
)

func RunAllChecks() map[string][]string {
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
	return messages
}
