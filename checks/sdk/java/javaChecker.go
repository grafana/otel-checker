package java

import (
	_ "embed"
	"fmt"
	"github.com/grafana/otel-checker/checks/sdk"
	"github.com/grafana/otel-checker/checks/sdk/supported"
	"github.com/grafana/otel-checker/checks/utils"
	"os/exec"
	"strconv"
	"strings"
)

func CheckSetup(reporter *utils.ComponentReporter, commands utils.Commands) {
	checkJavaVersion(reporter)
	if commands.ManualInstrumentation {
		checkCodeBasedInstrumentation(reporter, commands.Debug)
	} else {
		checkAutoInstrumentation(reporter, commands.Debug)
	}
}

func checkJavaVersion(reporter *utils.ComponentReporter) {
	out := sdk.RunCommand(reporter, exec.Command("java", "-version"))
	if out != "" {
		//openjdk version "21.0.2" 2024-01-16 LTS
		line := strings.Split(out, "\n")[0]
		field := strings.Split(line, " ")[2]
		version := strings.Trim(field, "\"")
		major, err := strconv.Atoi(strings.Split(version, ".")[0])
		if err != nil {
			reporter.AddError(fmt.Sprintf("Error parsing Java version %s: %v", out, err))
		}
		if strings.HasPrefix(version, "1.8") {
			major = 8
		}
		if major < 8 {
			reporter.AddError(fmt.Sprintf("Java version %s is not supported. Please use Java 8 or higher", version))
		} else {
			reporter.AddSuccessfulCheck(fmt.Sprintf("Java version %s is supported", version))
		}
	}
}

func checkAutoInstrumentation(reporter *utils.ComponentReporter, debug bool) {
	reportSupportedInstrumentations(reporter, debug, supported.TypeJavaagent)
}

func checkCodeBasedInstrumentation(reporter *utils.ComponentReporter, debug bool) {
	reportSupportedInstrumentations(reporter, debug, supported.TypeLibrary)
}
