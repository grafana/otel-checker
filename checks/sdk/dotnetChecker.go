package sdk

import "otel-checker/checks/utils"

func CheckDotNetSetup(reporter *utils.ComponentReporter, manualInstrumentation bool) {
	checkDotNetVersion(reporter)
	if !manualInstrumentation {
		checkDotNetAutoInstrumentation(reporter)
	} else {
		checkDotNetCodeBasedInstrumentation(reporter)
	}
}

func checkDotNetVersion(reporter *utils.ComponentReporter) {}

func checkDotNetAutoInstrumentation(reporter *utils.ComponentReporter) {}

func checkDotNetCodeBasedInstrumentation(reporter *utils.ComponentReporter) {}
