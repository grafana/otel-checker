package java

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"otel-checker/checks/sdk"
	"otel-checker/checks/utils"
	"strings"
)

func checkMaven(reporter *utils.ComponentReporter) []Library {
	println("Reading Maven dependencies")

	out := sdk.RunCommand(reporter, exec.Command(searchWrapper("mvn", "mvnw"),
		"dependency:tree", "-Dscope=runtime", "-DoutputType=json"))
	if out == "" {
		return []Library{}
	}
	deps := parseMavenDeps(out)
	if len(deps) == 0 {
		reporter.AddWarning("No Maven dependencies found")
	}
	return deps
}

func parseMavenDeps(out string) []Library {
	c := ""
	isJson := false
	for l := range strings.Lines(out) {
		if strings.Contains(l, "[INFO] {") {
			isJson = true
		}
		if isJson {
			c += strings.TrimPrefix(l, "[INFO] ")
		}
		if strings.Contains(l, "[INFO] }") {
			isJson = false
		}
	}

	var deps Library
	err := json.Unmarshal([]byte(c), &deps)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return nil
	}

	return []Library{deps}
}
