package sdk

func CheckSDKSetup(
	messages *map[string][]string,
	language string,
	autoInstrumentation bool,
	packageJsonPath string,
) {
	switch language {
	case "dotnet":
		CheckDotNetSetup(messages, autoInstrumentation)
	case "go":
		CheckGoSetup(messages, autoInstrumentation)
	case "java":
		CheckJavaSetup(messages, autoInstrumentation)
	case "js":
		CheckJSSetup(messages, autoInstrumentation, packageJsonPath)
	case "python":
		CheckPythonSetup(messages, autoInstrumentation)
	}
}
