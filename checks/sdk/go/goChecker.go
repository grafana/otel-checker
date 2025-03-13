package _go

import (
	"github.com/grafana/otel-checker/checks/utils"
)

func CheckGoSetup(reporter *utils.ComponentReporter, commands utils.Commands) {
	checkGoVersion(reporter)
	if commands.ManualInstrumentation {
		checkGoCodeBasedInstrumentation(reporter)
		CheckSupportedLibraries(reporter, commands)
	} else {
		checkGoAutoInstrumentation(reporter)
	}
}

func checkGoVersion(reporter *utils.ComponentReporter) {}

func checkGoAutoInstrumentation(reporter *utils.ComponentReporter) {}

func checkGoCodeBasedInstrumentation(reporter *utils.ComponentReporter) {}
