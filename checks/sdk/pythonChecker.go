package sdk

import "otel-checker/checks/utils"

func CheckPythonSetup(reporter *utils.ComponentReporter, autoInstrumentation bool) {
	checkPythonVersion(reporter)
	if autoInstrumentation {
		checkPythonAutoInstrumentation(reporter)
	} else {
		checkPythonCodeBasedInstrumentation(reporter)
	}
}

func checkPythonVersion(reporter *utils.ComponentReporter) {}

func checkPythonAutoInstrumentation(reporter *utils.ComponentReporter) {}

func checkPythonCodeBasedInstrumentation(reporter *utils.ComponentReporter) {}
