package sdk

import "otel-checker/checks/utils"

func CheckGoSetup(reporter *utils.ComponentReporter, manualInstrumentation bool) {
	checkGoVersion(reporter)
	if !manualInstrumentation {
		checkGoAutoInstrumentation(reporter)
	} else {
		checkGoCodeBasedInstrumentation(reporter)
	}
}

func checkGoVersion(reporter *utils.ComponentReporter) {}

func checkGoAutoInstrumentation(reporter *utils.ComponentReporter) {}

func checkGoCodeBasedInstrumentation(reporter *utils.ComponentReporter) {}
