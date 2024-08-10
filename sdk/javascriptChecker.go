package sdk

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	utils "otel-checker/utils"
)

func CheckJSSetup(
	messages *map[string][]string,
	autoInstrumentation bool,
	packageJsonPath string,
	instrumentationFile string,
) {
	checkEnvVars(messages)
	checkNodeVersion(messages)
	if autoInstrumentation {
		checkJSAutoInstrumentation(messages, packageJsonPath)
	} else {
		checkJSCodeBasedInstrumentation(messages, packageJsonPath, instrumentationFile)
	}
}

func checkEnvVars(messages *map[string][]string) {
	if os.Getenv("OTEL_NODE_RESOURCE_DETECTORS") == "" ||
		!strings.Contains(os.Getenv("OTEL_NODE_RESOURCE_DETECTORS"), "env") ||
		!strings.Contains(os.Getenv("OTEL_NODE_RESOURCE_DETECTORS"), "host") ||
		!strings.Contains(os.Getenv("OTEL_NODE_RESOURCE_DETECTORS"), "os") ||
		!strings.Contains(os.Getenv("OTEL_NODE_RESOURCE_DETECTORS"), "serviceinstance") {
		utils.AddWarning(messages, "SDK", "It's recommended the environment variable OTEL_NODE_RESOURCE_DETECTORS to be set to at least `env,host,os,serviceinstance`")
	} else {
		utils.AddSuccessfulCheck(messages, "SDK", "OTEL_NODE_RESOURCE_DETECTORS has recommended values")
	}
}

func checkNodeVersion(messages *map[string][]string) {
	cmd := exec.Command("node", "-v")
	stdout, err := cmd.Output()

	if err != nil {
		utils.AddError(messages, "SDK", fmt.Sprintf("Could not check minimum node version: %s", err))
	}
	versionInfo := strings.Split(string(stdout), ".")
	v, err := strconv.Atoi(versionInfo[0][1:])
	if err != nil {
		utils.AddError(messages, "SDK", fmt.Sprintf("Could not check minimum node version: %s", err))
	}
	if v >= 16 {
		utils.AddSuccessfulCheck(messages, "SDK", "Using node version equal or greater than minimum recommended")
	} else {
		utils.AddError(messages, "SDK", "Not using recommended node version. Update your node to at least version 16")
	}
}

func checkJSAutoInstrumentation(
	messages *map[string][]string,
	packageJsonPath string,
) {
	// NODE_OPTIONS should be set or that requirement should be added when starting the app
	if os.Getenv("NODE_OPTIONS") == "--require @opentelemetry/auto-instrumentations-node/register" {
		utils.AddSuccessfulCheck(messages, "SDK", "NODE_OPTIONS set correctly")
	} else {
		utils.AddWarning(messages, "SDK", `NODE_OPTIONS not set. You can set it by running 'export NODE_OPTIONS="--require @opentelemetry/auto-instrumentations-node/register"' or add the same required when starting your application`)
	}

	// Dependencies for auto instrumentation on package.json
	filePath := packageJsonPath + "package.json"
	dat, err := os.ReadFile(filePath)
	if err != nil {
		utils.AddError(messages, "SDK", fmt.Sprintf("Could not check file %s: %s", filePath, err))
	} else {
		if strings.Contains(string(dat), `"@opentelemetry/auto-instrumentations-node"`) {
			utils.AddSuccessfulCheck(messages, "SDK", "Dependency @opentelemetry/auto-instrumentations-node added on package.json")
		} else {
			utils.AddError(messages, "SDK", "Dependency @opentelemetry/auto-instrumentations-node missing on package.json. Install the dependency with `npm install @opentelemetry/auto-instrumentations-node`")
		}

		if strings.Contains(string(dat), `"@opentelemetry/api"`) {
			utils.AddSuccessfulCheck(messages, "SDK", "Dependency @opentelemetry/api added on package.json")
		} else {
			utils.AddError(messages, "SDK", "Dependency @opentelemetry/api missing on package.json. Install the dependency with `npm install @opentelemetry/auto-instrumentations-node`")
		}
	}
}

func checkJSCodeBasedInstrumentation(
	messages *map[string][]string,
	packageJsonPath string,
	instrumentationFile string,
) {
	if os.Getenv("NODE_OPTIONS") == "--require @opentelemetry/auto-instrumentations-node/register" {
		utils.AddError(messages, "SDK", `The flag "-auto-instrumentation" was not passed to otel-checker, but the value of NODE_OPTIONS is set to require auto-instrumentation. Run "unset NODE_OPTIONS" to remove the requirement that can cause a conflict with manual instrumentations`)
	}

	// Dependencies for auto instrumentation on package.json
	filePath := packageJsonPath + "package.json"
	packageJsonContent, err := os.ReadFile(filePath)
	if err != nil {
		utils.AddError(messages, "SDK", fmt.Sprintf("Could not check file %s: %s", filePath, err))
	} else {
		if strings.Contains(string(packageJsonContent), `"@opentelemetry/api"`) {
			utils.AddSuccessfulCheck(messages, "SDK", "Dependency @opentelemetry/api added on package.json")
		} else {
			utils.AddError(messages, "SDK", "Dependency @opentelemetry/api missing on package.json")
		}

		if strings.Contains(string(packageJsonContent), `"@opentelemetry/exporter-trace-otlp-proto"`) {
			utils.AddError(messages, "SDK", `Dependency @opentelemetry/exporter-trace-otlp-proto added on package.json, which is not supported by Grafana. Switch the dependency to "@opentelemetry/exporter-trace-otlp-http" instead`)
		}
	}

	// Check Exporter
	instrumentationFileContent, err := os.ReadFile(instrumentationFile)
	if err != nil {
		utils.AddError(messages, "SDK", fmt.Sprintf("Could not check file %s: %s", instrumentationFile, err))
	} else {
		if strings.Contains(string(instrumentationFileContent), "ConsoleSpanExporter") {
			utils.AddWarning(messages, "SDK", "Instrumentation file is using ConsoleSpanExporter. This exporter is useful during debugging, but replace with OTLPTraceExporter to send to Grafana Cloud")
		}
		if strings.Contains(string(instrumentationFileContent), "ConsoleMetricExporter") {
			utils.AddWarning(messages, "SDK", "Instrumentation file is using ConsoleMetricExporter. This exporter is useful during debugging, but replace with OTLPMetricExporter to send to Grafana Cloud")
		}
	}
}
