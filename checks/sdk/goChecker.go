package sdk

import "otel-checker/checks/utils"

func CheckGoSetup(reporter *utils.ComponentReporter, autoInstrumentation bool) {
	checkGoVersion(reporter)
	if autoInstrumentation {
		checkGoAutoInstrumentation(reporter)
	} else {
		checkGoCodeBasedInstrumentation(reporter)
	}
}

func checkGoVersion(reporter *utils.ComponentReporter) {}

func checkGoAutoInstrumentation(reporter *utils.ComponentReporter) {}

func checkGoCodeBasedInstrumentation(reporter *utils.ComponentReporter) {}
