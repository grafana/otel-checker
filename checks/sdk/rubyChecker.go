package sdk

func CheckRubySetup(
	messages *map[string][]string,
	autoInstrumentation bool,
) {
	if autoInstrumentation {
		checkRubyAutoInstrumentation(messages)
	} else {
		checkRubyCodeBasedInstrumentation(messages)
	}
}

func checkRubyAutoInstrumentation(messages *map[string][]string) {}

func checkRubyCodeBasedInstrumentation(messages *map[string][]string) {}
