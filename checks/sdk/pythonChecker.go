package sdk

import "otel-checker/checks/utils"

func CheckPythonSetup(reporter *utils.ComponentReporter, manualInstrumentation bool) {
	checkPythonVersion(reporter)
	if !manualInstrumentation {
		checkPythonAutoInstrumentation(reporter)
	} else {
		checkPythonCodeBasedInstrumentation(reporter)
	}
}

func checkPythonVersion(reporter *utils.ComponentReporter) {}

func checkPythonAutoInstrumentation(reporter *utils.ComponentReporter) {}

func checkPythonCodeBasedInstrumentation(reporter *utils.ComponentReporter) {}
