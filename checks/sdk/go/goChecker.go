package _go

import (
	"otel-checker/checks/sdk/supported"
	"otel-checker/checks/utils"
)

func CheckGoSetup(reporter *utils.ComponentReporter, commands utils.Commands) {
	checkGoVersion(reporter)
	if commands.ManualInstrumentation {
		checkSupportedLibraries(reporter, commands, supported.TypeLibrary)
		checkGoCodeBasedInstrumentation(reporter)
	} else {
		checkGoAutoInstrumentation(reporter)
	}
}

func checkGoVersion(reporter *utils.ComponentReporter) {}

func checkGoAutoInstrumentation(reporter *utils.ComponentReporter) {}

func checkGoCodeBasedInstrumentation(reporter *utils.ComponentReporter) {}
