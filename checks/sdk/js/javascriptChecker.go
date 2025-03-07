package js

import (
	"fmt"
	"os"
	"os/exec"
	"otel-checker/checks/env"
	"otel-checker/checks/utils"
	"strconv"
	"strings"
)

func CheckJSSetup(reporter *utils.ComponentReporter, commands utils.Commands) {
	checkResourceDetectors(reporter)
	checkNodeVersion(reporter)
	if commands.ManualInstrumentation {
		checkJSCodeBasedInstrumentation(reporter, commands.PackageJsonPath, commands.InstrumentationFile)
	} else {
		checkJSAutoInstrumentation(reporter, commands.PackageJsonPath)
	}
	CheckSupportedLibraries(reporter, commands)
}

func checkResourceDetectors(reporter *utils.ComponentReporter) {
	env.CheckEnvVar("", env.EnvVar{
		Name: "OTEL_NODE_RESOURCE_DETECTORS",
		Validator: func(value string, language string, reporter *utils.ComponentReporter) {
			if value == "" ||
				!strings.Contains(value, "env") ||
				!strings.Contains(value, "host") ||
				!strings.Contains(value, "os") ||
				!strings.Contains(value, "serviceinstance") {
				reporter.AddWarning("It's recommended the environment variable OTEL_NODE_RESOURCE_DETECTORS to be set to at least `env,host,os,serviceinstance`")
			} else {
				reporter.AddSuccessfulCheck("OTEL_NODE_RESOURCE_DETECTORS has recommended values")
			}
		},
		Description: "at least `env,host,os,serviceinstance`",
	}, reporter)
}

func checkNodeVersion(reporter *utils.ComponentReporter) {
	cmd := exec.Command("node", "-v")
	stdout, err := cmd.Output()

	if err != nil {
		reporter.AddError(fmt.Sprintf("Could not check minimum node version: %s", err))
		return
	}
	versionInfo := strings.Split(string(stdout), ".")
	v, err := strconv.Atoi(versionInfo[0][1:])
	if err != nil {
		reporter.AddError(fmt.Sprintf("Could not check minimum node version: %s", err))
		return
	}
	if v >= 16 {
		reporter.AddSuccessfulCheck("Using node version equal or greater than minimum recommended")
	} else {
		reporter.AddError("Not using recommended node version. Update your node to at least version 16")
	}
}

func checkJSAutoInstrumentation(
	reporter *utils.ComponentReporter,
	packageJsonPath string,
) {
	checkAutoInstrumentationNodeOptions(reporter)

	// Dependencies for auto instrumentation on package.json
	filePath := packageJsonPath + "package.json"
	dat, err := os.ReadFile(filePath)
	if err != nil {
		reporter.AddError(fmt.Sprintf("Could not check file %s: %s", filePath, err))
	} else {
		if strings.Contains(string(dat), `"@opentelemetry/auto-instrumentations-node"`) {
			reporter.AddSuccessfulCheck("Dependency @opentelemetry/auto-instrumentations-node added on package.json")
		} else {
			reporter.AddError("Dependency @opentelemetry/auto-instrumentations-node missing on package.json. Install the dependency with `npm install @opentelemetry/auto-instrumentations-node`")
		}

		if strings.Contains(string(dat), `"@opentelemetry/api"`) {
			reporter.AddSuccessfulCheck("Dependency @opentelemetry/api added on package.json")
		} else {
			reporter.AddError("Dependency @opentelemetry/api missing on package.json. Install the dependency with `npm install @opentelemetry/auto-instrumentations-node`")
		}
	}
}

func checkAutoInstrumentationNodeOptions(reporter *utils.ComponentReporter) {
	env.CheckEnvVar("", env.EnvVar{
		Name:          "NODE_OPTIONS",
		Recommended:   true,
		RequiredValue: "--require @opentelemetry/auto-instrumentations-node/register",
		Message:       `NODE_OPTIONS not set. You can set it by running 'export NODE_OPTIONS="--require @opentelemetry/auto-instrumentations-node/register"' or add the same '--require ...' when starting your application`,
	}, reporter)
}

func checkJSCodeBasedInstrumentation(
	reporter *utils.ComponentReporter,
	packageJsonPath string,
	instrumentationFile string,
) {
	if os.Getenv("NODE_OPTIONS") == "--require @opentelemetry/auto-instrumentations-node/register" {
		reporter.AddError(`The flag "-manual-instrumentation" was set, but the value of NODE_OPTIONS is set to require auto-instrumentation. Run "unset NODE_OPTIONS" to remove the requirement that can cause a conflict with manual instrumentations`)
	}

	// Dependencies for auto instrumentation on package.json
	filePath := packageJsonPath + "package.json"
	packageJsonContent, err := os.ReadFile(filePath)
	if err != nil {
		reporter.AddError(fmt.Sprintf("Could not check file %s: %s", filePath, err))
	} else {
		if strings.Contains(string(packageJsonContent), `"@opentelemetry/api"`) {
			reporter.AddSuccessfulCheck("Dependency @opentelemetry/api added on package.json")
		} else {
			reporter.AddError("Dependency @opentelemetry/api missing on package.json")
		}

		if strings.Contains(string(packageJsonContent), `"@opentelemetry/exporter-trace-otlp-proto"`) {
			reporter.AddError(`Dependency @opentelemetry/exporter-trace-otlp-proto added on package.json, which is not supported by Grafana. Switch the dependency to "@opentelemetry/exporter-trace-otlp-http" instead`)
		}
	}

	// Check Exporter
	instrumentationFileContent, err := os.ReadFile(instrumentationFile)
	if err != nil {
		reporter.AddError(fmt.Sprintf("Could not check file %s: %s", instrumentationFile, err))
	} else {
		if strings.Contains(string(instrumentationFileContent), "ConsoleSpanExporter") {
			reporter.AddWarning("Instrumentation file is using ConsoleSpanExporter. This exporter is useful during debugging, but replace with OTLPTraceExporter to send to Grafana Cloud")
		}
		if strings.Contains(string(instrumentationFileContent), "ConsoleMetricExporter") {
			reporter.AddWarning("Instrumentation file is using ConsoleMetricExporter. This exporter is useful during debugging, but replace with OTLPMetricExporter to send to Grafana Cloud")
		}
	}
}
