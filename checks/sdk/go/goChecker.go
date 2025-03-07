package _go

import "otel-checker/checks/utils"

func CheckGoSetup(reporter *utils.ComponentReporter, commands utils.Commands) {
	checkGoVersion(reporter)
	if commands.ManualInstrumentation {
		checkGoCodeBasedInstrumentation(reporter)
	} else {
		checkGoAutoInstrumentation(reporter)
	}
	CheckSupportedLibraries(reporter, commands)
}

func checkGoVersion(reporter *utils.ComponentReporter) {}

func checkGoAutoInstrumentation(reporter *utils.ComponentReporter) {}

func checkGoCodeBasedInstrumentation(reporter *utils.ComponentReporter) {}
