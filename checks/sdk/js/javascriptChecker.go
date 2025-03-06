package js

import (
	"fmt"
	"os"
	"os/exec"
	"otel-checker/checks/utils"
	"strconv"
	"strings"
)

func CheckJSSetup(reporter *utils.ComponentReporter, commands utils.Commands) {
	checkEnvVars(reporter)
	checkNodeVersion(reporter)
	if commands.ManualInstrumentation {
		checkJSCodeBasedInstrumentation(reporter, commands.PackageJsonPath, commands.InstrumentationFile)
	} else {
		checkJSAutoInstrumentation(reporter, commands.PackageJsonPath)
	}
	checkSupportedLibraries(reporter, commands)
}

func checkEnvVars(reporter *utils.ComponentReporter) {
	if os.Getenv("OTEL_NODE_RESOURCE_DETECTORS") == "" ||
		!strings.Contains(os.Getenv("OTEL_NODE_RESOURCE_DETECTORS"), "env") ||
		!strings.Contains(os.Getenv("OTEL_NODE_RESOURCE_DETECTORS"), "host") ||
		!strings.Contains(os.Getenv("OTEL_NODE_RESOURCE_DETECTORS"), "os") ||
		!strings.Contains(os.Getenv("OTEL_NODE_RESOURCE_DETECTORS"), "serviceinstance") {
		reporter.AddWarning("It's recommended the environment variable OTEL_NODE_RESOURCE_DETECTORS to be set to at least `env,host,os,serviceinstance`")
	} else {
		reporter.AddSuccessfulCheck("OTEL_NODE_RESOURCE_DETECTORS has recommended values")
	}
}

func checkNodeVersion(reporter *utils.ComponentReporter) {
	cmd := exec.Command("node", "-v")
	stdout, err := cmd.Output()

	if err != nil {
		reporter.AddError(fmt.Sprintf("Could not check minimum node version: %s", err))
	}
	versionInfo := strings.Split(string(stdout), ".")
	v, err := strconv.Atoi(versionInfo[0][1:])
	if err != nil {
		reporter.AddError(fmt.Sprintf("Could not check minimum node version: %s", err))
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
	// NODE_OPTIONS should be set or that requirement should be added when starting the app
	if os.Getenv("NODE_OPTIONS") == "--require @opentelemetry/auto-instrumentations-node/register" {
		reporter.AddSuccessfulCheck("NODE_OPTIONS set correctly")
	} else {
		reporter.AddWarning(`NODE_OPTIONS not set. You can set it by running 'export NODE_OPTIONS="--require @opentelemetry/auto-instrumentations-node/register"' or add the same '--require ...' when starting your application`)
	}

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

func checkSupportedLibraries(reporter *utils.ComponentReporter, commands utils.Commands) {
	supported, err := supportedLibraries()
	if err != nil {
		reporter.AddError(fmt.Sprintf("Error reading supported libraries: %v", err))
		return
	}

	deps := readDependencies(reporter)
	if len(deps) == 0 {
		return
	}

	for _, dep := range deps {
		links := findSupportedLibraries(dep, supported)
		if len(links) > 0 {
			reporter.AddSuccessfulCheck(
				fmt.Sprintf("Found supported library: %s:%s at %s",
					dep.Name, dep.Version, strings.Join(links, ", ")))
		} else if commands.Debug {
			reporter.AddWarning(fmt.Sprintf("Found unsupported library: %s:%s", dep.Name, dep.Version))
		}
	}
}
