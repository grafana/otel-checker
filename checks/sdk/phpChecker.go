package sdk

import (
	"github.com/grafana/otel-checker/checks/utils"
	"os"
	"os/exec"
	"strings"
)

func CheckPHPSetup(reporter *utils.ComponentReporter, commands utils.Commands) {
	checkPHPVersion(reporter)
	checkComposerInstalled(reporter)

	composerFile, err := checkComposerFileExists(reporter)
	if err != nil {
		return
	}

	// shared instrumentation checks
	checkPHPRequiredInstrumentation(reporter, &composerFile)

	if commands.ManualInstrumentation {
		checkPHPManualInstrumentation(reporter, &composerFile)
	} else {
		checkPHPAutoInstrumentation(reporter, &composerFile)
	}
}

func checkPHPVersion(reporter *utils.ComponentReporter) {
	cmd := exec.Command("php", "-v")
	stdout, err := cmd.Output()

	if err != nil {
		reporter.AddError("PHP not found, install PHP >= 8.0")
		return
	}

	if strings.Contains(string(stdout), "PHP 8") {
		reporter.AddSuccessfulCheck("Using PHP >= 8.0")
	} else {
		reporter.AddError("Not using recommended PHP version, update to PHP >= 8.0")
	}
}

func checkComposerInstalled(reporter *utils.ComponentReporter) {
	cmd := exec.Command("composer", "-v")
	_, err := cmd.Output()

	if err != nil {
		reporter.AddError("Composer not found. Run 'curl -sS https://getcomposer.org/installer | php' to install it.")
	} else {
		reporter.AddSuccessfulCheck("Composer found. Run 'composer install' to install dependencies.")
	}
}

func checkComposerFileExists(reporter *utils.ComponentReporter) (string, error) {
	_, err := os.ReadFile("composer.json")
	if err != nil {
		reporter.AddError("Could not find composer.json, create one, add dependencies, and run 'composer install'")
		return "", err
	}

	content, err := os.ReadFile("composer.lock")
	if err != nil {
		reporter.AddError("Could not find composer.lock, run 'composer install' to generate it")
		return "", err
	}

	composerFile := string(content)
	reporter.AddSuccessfulCheck("Found composer.lock")

	return composerFile, nil
}

func checkPHPRequiredInstrumentation(reporter *utils.ComponentReporter, composerFile *string) {
	requiredPackages := []string{
		"open-telemetry/api",
		"open-telemetry/sem-conv",
		"open-telemetry/sdk",
		"open-telemetry/exporter-otlp",
	}

	// loop through requiredPackages and check if they are in composer.lock
	for _, pkg := range requiredPackages {
		if strings.Contains(*composerFile, pkg) {
			reporter.AddSuccessfulCheck("Found required dependency: " + pkg)
		} else {
			reporter.AddError("Missing required dependency: " + pkg + ", add it to your composer.json and run 'composer install'")
		}
	}
}

func checkPHPAutoInstrumentation(reporter *utils.ComponentReporter, composerFile *string) {
	autoPackages := []string{
		"open-telemetry/opentelemetry-auto-symfony",
		"open-telemetry/opentelemetry-auto-pdo",
		"open-telemetry/opentelemetry-auto-laravel",
		"open-telemetry/opentelemetry-auto-wordpress",
		"open-telemetry/opentelemetry-auto-slim",
		"open-telemetry/opentelemetry-auto-psr3",
		"open-telemetry/opentelemetry-auto-psr18",
		"open-telemetry/opentelemetry-auto-psr15",
		"open-telemetry/opentelemetry-auto-guzzle",
		"open-telemetry/opentelemetry-auto-yii",
		"open-telemetry/opentelemetry-auto-psr6",
		"open-telemetry/opentelemetry-auto-psr16",
		"open-telemetry/opentelemetry-auto-psr14",
		"open-telemetry/opentelemetry-auto-mongodb",
		"open-telemetry/opentelemetry-auto-io",
		"open-telemetry/opentelemetry-auto-http-async",
		"open-telemetry/opentelemetry-auto-ext-rdkafka",
		"open-telemetry/opentelemetry-auto-ext-amqp",
		"open-telemetry/opentelemetry-auto-curl",
		"open-telemetry/opentelemetry-auto-openai-php",
		"open-telemetry/opentelemetry-auto-mysqli",
		"open-telemetry/opentelemetry-auto-codeigniter",
		"open-telemetry/opentelemetry-auto-cakephp",
	}

	found := false

	// loop through optionalPackages and check if they are in composer.lock
	for _, pkg := range autoPackages {
		if strings.Contains(*composerFile, pkg) {
			found = true
			reporter.AddSuccessfulCheck("Found optional instrumentation dependency: " + pkg)
		}
	}

	// if not optionalFound then add error
	if !found {
		reporter.AddError("Missing instrumentation dependencies, add them to your composer.json and run 'composer install'")
	}
}

func checkPHPManualInstrumentation(reporter *utils.ComponentReporter, composerFile *string) {
	// Empty function for future implementation
}
