package sdk

func CheckJavaSetup(
	messages *map[string][]string,
	autoInstrumentation bool,
) {
	checkJavaVersion(messages)
	if autoInstrumentation {
		checkJavaAutoInstrumentation(messages)
	} else {
		checkJavaCodeBasedInstrumentation(messages)
	}
}

func checkJavaVersion(messages *map[string][]string) {}

func checkJavaAutoInstrumentation(messages *map[string][]string) {}

func checkJavaCodeBasedInstrumentation(messages *map[string][]string) {}
