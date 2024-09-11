package sdk

func CheckGoSetup(
	messages *map[string][]string,
	autoInstrumentation bool,
) {
	checkGoVersion(messages)
	if autoInstrumentation {
		checkGoAutoInstrumentation(messages)
	} else {
		checkGoCodeBasedInstrumentation(messages)
	}
}

func checkGoVersion(messages *map[string][]string) {}

func checkGoAutoInstrumentation(messages *map[string][]string) {}

func checkGoCodeBasedInstrumentation(messages *map[string][]string) {}
