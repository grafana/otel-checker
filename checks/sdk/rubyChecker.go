package sdk

import (
	"os"
	"os/exec"
	"otel-checker/checks/utils"
	"strings"

	"golang.org/x/mod/semver"
)

func CheckRubySetup(reporter *utils.ComponentReporter, autoInstrumentation bool) {
	checkRubyVersion(reporter)
	checkBundlerInstalled(reporter)
	if autoInstrumentation {
		checkRubyAutoInstrumentation(reporter)
	} else {
		checkRubyCodeBasedInstrumentation(reporter)
	}
}

// While tested, support for jruby and truffleruby are on a best-effort basis at this time.
func checkRubyVersion(reporter *utils.ComponentReporter) {
	hasCRuby := checkCRubyVersion(reporter)
	hasJRuby := checkJRubyVersion(reporter)
	hasTruffleRuby := checkJRubyVersion(reporter)

	if hasCRuby || hasJRuby || hasTruffleRuby {
		reporter.AddSuccessfulCheck("Ruby setup successful")
	} else {
		reporter.AddError("No Ruby found, install CRuby >= 3.0, JRuby >= 9.3.2.0, or TruffleRuby >= 22.1")
	}
}

func checkCRubyVersion(reporter *utils.ComponentReporter) bool {
	cmd := exec.Command("ruby", "-v")
	stdout, err := cmd.Output()

	if err != nil {
		return false
	}

	if strings.Contains(string(stdout), "ruby 3") {
		reporter.AddSuccessfulCheck("Using CRuby >= 3.0")
		return true
	} else {
		reporter.AddError("Not using recommended CRuby version, update to CRuby >= 3.0")
		return false
	}
}

func checkJRubyVersion(reporter *utils.ComponentReporter) bool {
	cmd := exec.Command("jruby", "--version")
	stdout, err := cmd.Output()

	if err != nil {
		return false
	}

	version := strings.Fields(string(stdout))[2]

	if semver.Compare(version, "9.3.2.0") >= 0 {
		reporter.AddSuccessfulCheck("Using JRuby >= 9.3.2.0")
		return true
	} else {
		reporter.AddError("Not using recommended JRuby version, update to JRuby >= 9.3.2.0")
		return false
	}
}

func checkBundlerInstalled(reporter *utils.ComponentReporter) {
	cmd := exec.Command("bundle", "-v")
	_, err := cmd.Output()

	if err != nil {
		reporter.AddError("Bundler not found. Run 'gem install bundler' to install it.")
	} else {
		reporter.AddSuccessfulCheck("Bundler found. Run 'bundle install' to install dependencies.")
	}
}

// TruffleRuby check is not supported yet
//func checkTruffleRubyVersion(reporter *utils.ComponentReporter) bool {
//	return false
//}

func checkRubyAutoInstrumentation(reporter *utils.ComponentReporter) {
	content, err := os.ReadFile("Gemfile.lock")
	if err != nil {
		reporter.AddError("Could not find Gemfile.lock. Run 'bundle install' to generate it.")
		return
	}

	contentStr := string(content)

	if strings.Contains(contentStr, "opentelemetry-sdk") {
		reporter.AddSuccessfulCheck("Found required dependency: opentelemetry-sdk")
	} else {
		reporter.AddError("Missing required OpenTelemetry Ruby dependency: opentelemetry-sdk. Add it to your Gemfile and run 'bundle install'.")
	}

	if strings.Contains(contentStr, "opentelemetry-instrumentation-all") {
		reporter.AddSuccessfulCheck("Found required dependency: opentelemetry-instrumentation-all")
	} else {
		reporter.AddError("Missing required OpenTelemetry Ruby dependency: opentelemetry-instrumentation-all. Add it to your Gemfile and run 'bundle install'.")
	}
}

func checkRubyCodeBasedInstrumentation(reporter *utils.ComponentReporter) {
	content, err := os.ReadFile("Gemfile.lock")
	if err != nil {
		reporter.AddError("Could not find Gemfile.lock. Run 'bundle install' to generate it.")
		return
	}

	contentStr := string(content)

	if strings.Contains(contentStr, "opentelemetry-sdk") {
		reporter.AddSuccessfulCheck("Found required dependency: opentelemetry-sdk")
	} else {
		reporter.AddError("Missing required OpenTelemetry Ruby dependency: opentelemetry-sdk. Add it to your Gemfile and run 'bundle install'.")
	}

	if strings.Contains(contentStr, "opentelemetry-api") {
		reporter.AddSuccessfulCheck("Found required dependency: opentelemetry-api")
	} else {
		reporter.AddError("Missing required OpenTelemetry Ruby dependency: opentelemetry-api. Add it to your Gemfile and run 'bundle install'.")
	}

	if strings.Contains(contentStr, "opentelemetry-common") {
		reporter.AddSuccessfulCheck("Found required dependency: opentelemetry-common")
	} else {
		reporter.AddError("Missing required OpenTelemetry Ruby dependency: opentelemetry-common. Add it to your Gemfile and run 'bundle install'.")
	}
}

// possibly check instrumentation file for presence of following strings:
// OpenTelemetry.tracer_provider.tracer
// OpenTelemetry.propagation.extract
// OpenTelemetry::Context.with_current
// tracer.in_span
// span.set_attribute
