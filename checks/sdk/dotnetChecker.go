package sdk

import "otel-checker/checks/utils"

func CheckDotNetSetup(reporter *utils.ComponentReporter, commands utils.Commands) {
	checkDotNetVersion(reporter)
	if commands.ManualInstrumentation {
		checkDotNetCodeBasedInstrumentation(reporter)
	} else {
		checkDotNetAutoInstrumentation(reporter)
	}
}

func checkDotNetVersion(reporter *utils.ComponentReporter) {}

func checkDotNetAutoInstrumentation(reporter *utils.ComponentReporter) {}

func checkDotNetCodeBasedInstrumentation(reporter *utils.ComponentReporter) {}
