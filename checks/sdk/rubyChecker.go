package sdk

import (
	"os"
	"os/exec"
	utils "otel-checker/checks/utils"
	"strings"

	"golang.org/x/mod/semver"
)

func CheckRubySetup(
	messages *map[string][]string,
	autoInstrumentation bool,
) {
	checkRubyVersion(messages)
	checkBundlerInstalled(messages)
	if autoInstrumentation {
		checkRubyAutoInstrumentation(messages)
	} else {
		checkRubyCodeBasedInstrumentation(messages)
	}
}

// While tested, support for jruby and truffleruby are on a best-effort basis at this time.
func checkRubyVersion(messages *map[string][]string) {
	hasCRuby := checkCRubyVersion(messages)
	hasJRuby := checkJRubyVersion(messages)
	hasTruffleRuby := checkJRubyVersion(messages)

	if hasCRuby || hasJRuby || hasTruffleRuby {
		utils.AddSuccessfulCheck(messages, "SDK", "Ruby setup successful")
	} else {
		utils.AddError(messages, "SDK", "No Ruby found, install CRuby >= 3.0, JRuby >= 9.3.2.0, or TruffleRuby >= 22.1")
	}
}

func checkCRubyVersion(messages *map[string][]string) bool {
	cmd := exec.Command("ruby", "-v")
	stdout, err := cmd.Output()

	if err != nil {
		return false
	}

	if strings.Contains(string(stdout), "ruby 3") {
		utils.AddSuccessfulCheck(messages, "SDK", "Using CRuby >= 3.0")
		return true
	} else {
		utils.AddError(messages, "SDK", "Not using recommended CRuby version, update to CRuby >= 3.0")
		return false
	}
}

func checkJRubyVersion(messages *map[string][]string) bool {
	cmd := exec.Command("jruby", "--version")
	stdout, err := cmd.Output()

	if err != nil {
		return false
	}

	version := strings.Fields(string(stdout))[2]

	if semver.Compare(version, "9.3.2.0") >= 0 {
		utils.AddSuccessfulCheck(messages, "SDK", "Using JRuby >= 9.3.2.0")
		return true
	} else {
		utils.AddError(messages, "SDK", "Not using recommended JRuby version, update to JRuby >= 9.3.2.0")
		return false
	}
}

func checkBundlerInstalled(messages *map[string][]string) {
	cmd := exec.Command("bundle", "-v")
	_, err := cmd.Output()

	if err != nil {
		utils.AddError(messages, "SDK", "Bundler not found. Run 'gem install bundler' to install it.")
	} else {
		utils.AddSuccessfulCheck(messages, "SDK", "Bundler found. Run 'bundle install' to install dependencies.")
	}
}

// TruffleRuby check is not supported yet
//func checkTruffleRubyVersion(messages *map[string][]string) bool {
//	return false
//}

func checkRubyAutoInstrumentation(messages *map[string][]string) {
	content, err := os.ReadFile("Gemfile.lock")
	if err != nil {
		utils.AddError(messages, "SDK", "Could not find Gemfile.lock. Run 'bundle install' to generate it.")
		return
	}

	contentStr := string(content)

	if strings.Contains(contentStr, "opentelemetry-sdk") {
		utils.AddSuccessfulCheck(messages, "SDK", "Found required dependency: opentelemetry-sdk")
	} else {
		utils.AddError(messages, "SDK", "Missing required OpenTelemetry Ruby dependency: opentelemetry-sdk. Add it to your Gemfile and run 'bundle install'.")
	}

	if strings.Contains(contentStr, "opentelemetry-instrumentation-all") {
		utils.AddSuccessfulCheck(messages, "SDK", "Found required dependency: opentelemetry-instrumentation-all")
	} else {
		utils.AddError(messages, "SDK", "Missing required OpenTelemetry Ruby dependency: opentelemetry-instrumentation-all. Add it to your Gemfile and run 'bundle install'.")
	}
}

func checkRubyCodeBasedInstrumentation(messages *map[string][]string) {
	content, err := os.ReadFile("Gemfile.lock")
	if err != nil {
		utils.AddError(messages, "SDK", "Could not find Gemfile.lock. Run 'bundle install' to generate it.")
		return
	}

	contentStr := string(content)

	if strings.Contains(contentStr, "opentelemetry-sdk") {
		utils.AddSuccessfulCheck(messages, "SDK", "Found required dependency: opentelemetry-sdk")
	} else {
		utils.AddError(messages, "SDK", "Missing required OpenTelemetry Ruby dependency: opentelemetry-sdk. Add it to your Gemfile and run 'bundle install'.")
	}

	if strings.Contains(contentStr, "opentelemetry-api") {
		utils.AddSuccessfulCheck(messages, "SDK", "Found required dependency: opentelemetry-api")
	} else {
		utils.AddError(messages, "SDK", "Missing required OpenTelemetry Ruby dependency: opentelemetry-api. Add it to your Gemfile and run 'bundle install'.")
	}

	if strings.Contains(contentStr, "opentelemetry-common") {
		utils.AddSuccessfulCheck(messages, "SDK", "Found required dependency: opentelemetry-common")
	} else {
		utils.AddError(messages, "SDK", "Missing required OpenTelemetry Ruby dependency: opentelemetry-common. Add it to your Gemfile and run 'bundle install'.")
	}
}

// possibly check instrumentation file for presence of following strings:
// OpenTelemetry.tracer_provider.tracer
// OpenTelemetry.propagation.extract
// OpenTelemetry::Context.with_current
// tracer.in_span
// span.set_attribute
