package sdk

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/grafana/otel-checker/checks/utils"

	"golang.org/x/mod/semver"
)

func CheckRubySetup(ctx context.Context, reporter *utils.ComponentReporter, commands utils.Commands) {
	checkRubyVersion(ctx, reporter)
	checkBundlerInstalled(ctx, reporter)

	gemfile, err := checkGemfileExists(reporter)
	if err != nil {
		return
	}

	// shared instrumentation checks
	checkRubyRequiredInstrumentation(reporter, &gemfile)

	if !commands.ManualInstrumentation {
		checkRubyAutoInstrumentation(reporter, &gemfile)
	}

	// code based instrumentation uses the required instrumentation gems
}

// While tested, support for jruby and truffleruby are on a best-effort basis at this time.
func checkRubyVersion(ctx context.Context, reporter *utils.ComponentReporter) {
	hasCRuby := checkCRubyVersion(ctx, reporter)
	hasJRuby := checkJRubyVersion(ctx, reporter)
	hasTruffleRuby := checkTruffleRubyVersion(reporter)

	if hasCRuby || hasJRuby || hasTruffleRuby {
		reporter.AddSuccessfulCheck("Ruby setup successful")
	} else {
		reporter.AddErrorWithExplain("ruby.runtime.not-found",
			"No Ruby found, install CRuby >= 3.0, JRuby >= 9.3.2.0, or TruffleRuby >= 22.1")
	}
}

func checkBundlerInstalled(ctx context.Context, reporter *utils.ComponentReporter) {
	cmd := exec.CommandContext(ctx, "bundle", "-v")
	_, err := cmd.Output()

	if err != nil {
		reporter.AddErrorWithExplain("ruby.bundler.not-found",
			"Bundler not found. Run 'gem install bundler' to install it.")
	} else {
		reporter.AddSuccessfulCheck("Bundler found. Run 'bundle install' to install dependencies.")
	}
}

func checkGemfileExists(reporter *utils.ComponentReporter) (string, error) {
	_, err := os.ReadFile("Gemfile")
	if err != nil {
		reporter.AddErrorWithExplain("ruby.gemfile.missing",
			"Could not find Gemfile, create one, add dependencies, and run 'bundle install'")
		return "", err
	}

	content, err := os.ReadFile("Gemfile.lock")
	if err != nil {
		reporter.AddErrorWithExplain("ruby.gemfile-lock.missing",
			"Could not find Gemfile.lock run 'bundle install' to generate it")
		return "", err
	}

	gemfile := string(content)
	reporter.AddSuccessfulCheck("Found Gemfile.lock")

	return gemfile, nil
}

func checkCRubyVersion(ctx context.Context, reporter *utils.ComponentReporter) bool {
	cmd := exec.CommandContext(ctx, "ruby", "-v")
	stdout, err := cmd.Output()

	if err != nil {
		return false
	}

	if strings.Contains(string(stdout), "ruby 3") {
		reporter.AddSuccessfulCheck("Using CRuby >= 3.0")
		return true
	} else {
		reporter.AddErrorWithExplain("ruby.cruby.too-old",
			"Not using recommended CRuby version, update to CRuby >= 3.0")
		return false
	}
}

func checkJRubyVersion(ctx context.Context, reporter *utils.ComponentReporter) bool {
	cmd := exec.CommandContext(ctx, "jruby", "--version")
	stdout, err := cmd.Output()

	if err != nil {
		return false
	}

	version := strings.Fields(string(stdout))[2]

	if semver.Compare(version, "9.3.2.0") >= 0 {
		reporter.AddSuccessfulCheck("Using JRuby >= 9.3.2.0")
		return true
	} else {
		reporter.AddErrorWithExplain("ruby.jruby.too-old",
			"Not using recommended JRuby version, update to JRuby >= 9.3.2.0")
		return false
	}
}

// not implemented yet
func checkTruffleRubyVersion(reporter *utils.ComponentReporter) bool {
	return false
}

func checkRubyRequiredInstrumentation(reporter *utils.ComponentReporter, gemfile *string) {
	requiredGems := []string{
		"opentelemetry-api",
		"opentelemetry-sdk",
		"opentelemetry-exporter-otlp",
	}

	// loop through requiredGs and check if they are in Gemfile.lock
	for _, gem := range requiredGems {
		if strings.Contains(*gemfile, gem) {
			reporter.AddSuccessfulCheck("Found required dependency: " + gem)
		} else {
			reporter.AddErrorWithExplain("ruby.gem.missing-required",
				"Missing required dependency: "+gem+", add it to your Gemfile and run 'bundle install'")
		}
	}
}

func checkRubyAutoInstrumentation(reporter *utils.ComponentReporter, gemfile *string) {
	allFound := false
	optionalFound := false

	allGem := "opentelemetry-instrumentation-all"

	optionalGems := []string{
		"opentelemetry-instrumentation-active_model_serializers",
		"opentelemetry-instrumentation-aws_lambda",
		"opentelemetry-instrumentation-aws_sdk",
		"opentelemetry-instrumentation-bunny",
		"opentelemetry-instrumentation-concurrent_ruby",
		"opentelemetry-instrumentation-dalli",
		"opentelemetry-instrumentation-delayed_job",
		"opentelemetry-instrumentation-ethon",
		"opentelemetry-instrumentation-excon",
		"opentelemetry-instrumentation-faraday",
		"opentelemetry-instrumentation-grape",
		"opentelemetry-instrumentation-graphql",
		"opentelemetry-instrumentation-gruf",
		"opentelemetry-instrumentation-http",
		"opentelemetry-instrumentation-http_client",
		"opentelemetry-instrumentation-koala",
		"opentelemetry-instrumentation-lmdb",
		"opentelemetry-instrumentation-mongo",
		"opentelemetry-instrumentation-mysql2",
		"opentelemetry-instrumentation-net_http",
		"opentelemetry-instrumentation-pg",
		"opentelemetry-instrumentation-que",
		"opentelemetry-instrumentation-racecar",
		"opentelemetry-instrumentation-rack",
		"opentelemetry-instrumentation-rails",
		"opentelemetry-instrumentation-rake",
		"opentelemetry-instrumentation-rdkafka",
		"opentelemetry-instrumentation-redis",
		"opentelemetry-instrumentation-resque",
		"opentelemetry-instrumentation-restclient",
		"opentelemetry-instrumentation-ruby_kafka",
		"opentelemetry-instrumentation-sidekiq",
		"opentelemetry-instrumentation-sinatra",
		"opentelemetry-instrumentation-trilogy",
	}

	if strings.Contains(*gemfile, allGem) {
		allFound = true
		reporter.AddSuccessfulCheck("Found optional instrumentation dependency: " + allGem)
	}

	// loop through requiredGs and check if they are in Gemfile.lock
	for _, gem := range optionalGems {
		if strings.Contains(*gemfile, gem) {
			optionalFound = true
			reporter.AddSuccessfulCheck("Found optional instrumentation dependency: " + gem)
		}
	}

	// if not allFound or not optionalFound then add error
	if !allFound || !optionalFound {
		reporter.AddErrorWithExplain("ruby.gem.missing-instrumentation",
			"Missing instrumentation dependencies, add them to your Gemfile and run 'bundle install'")
	}
}
