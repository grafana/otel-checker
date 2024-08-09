package sdk

func CheckDotNetSetup(
	messages *map[string][]string,
	autoInstrumentation bool,
) {
	checkDotNetVersion(messages)
	if autoInstrumentation {
		checkDotNetAutoInstrumentation(messages)
	} else {
		checkDotNetCodeBasedInstrumentation(messages)
	}
}

func checkDotNetVersion(messages *map[string][]string) {}

func checkDotNetAutoInstrumentation(messages *map[string][]string) {}

func checkDotNetCodeBasedInstrumentation(messages *map[string][]string) {}
