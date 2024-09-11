package sdk

func CheckPythonSetup(
	messages *map[string][]string,
	autoInstrumentation bool,
) {
	checkPythonVersion(messages)
	if autoInstrumentation {
		checkPythonAutoInstrumentation(messages)
	} else {
		checkPythonCodeBasedInstrumentation(messages)
	}
}

func checkPythonVersion(messages *map[string][]string) {}

func checkPythonAutoInstrumentation(messages *map[string][]string) {}

func checkPythonCodeBasedInstrumentation(messages *map[string][]string) {}
