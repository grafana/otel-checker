package java

import (
	"context"
	_ "embed"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/grafana/otel-checker/checks/sdk"
	"github.com/grafana/otel-checker/checks/sdk/supported"
	"github.com/grafana/otel-checker/checks/utils"
)

func CheckSetup(ctx context.Context, reporter *utils.ComponentReporter, commands utils.Commands) {
	javaVersion := checkJavaVersion(ctx, reporter)
	if commands.ManualInstrumentation {
		checkCodeBasedInstrumentation(ctx, reporter, commands.Debug, javaVersion)
	} else {
		checkAutoInstrumentation(ctx, reporter, commands.Debug, javaVersion)
	}
}

func checkJavaVersion(ctx context.Context, reporter *utils.ComponentReporter) int {
	out := sdk.RunCommand(reporter, exec.CommandContext(ctx, "java", "-version"))
	if out != "" {
		// openjdk version "21.0.2" 2024-01-16 LTS
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
		return major
	}
	return 0
}

func checkAutoInstrumentation(ctx context.Context, reporter *utils.ComponentReporter, debug bool, javaVersion int) {
	reportSupportedInstrumentations(ctx, reporter, debug, supported.TypeJavaagent, javaVersion)
}

func checkCodeBasedInstrumentation(ctx context.Context, reporter *utils.ComponentReporter, debug bool, javaVersion int) {
	reportSupportedInstrumentations(ctx, reporter, debug, supported.TypeLibrary, javaVersion)
}
