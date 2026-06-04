package utils

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/fatih/color"
)

const ERRORS = "errors"
const WARNINGS = "warnings"
const CHECKS = "checks"

type Commands struct {
	Language              string
	Components            []string
	ManualInstrumentation bool
	WebServer             bool
	InstrumentationFile   string
	PackageJsonPath       string
	CollectorConfigPath   string
	Debug                 bool
}

var (
	SupportedLanguages  = []string{"dotnet", "go", "java", "js", "python", "ruby", "php"}
	SupportedComponents = []string{"sdk", "beyla", "alloy", "collector", "grafana-cloud"}
)

func Validate(c Commands) error {
	if !slices.Contains(SupportedLanguages, c.Language) {
		return fmt.Errorf("language %q not supported. Possible values: %s", c.Language, strings.Join(SupportedLanguages, ", "))
	}
	if len(c.Components) == 0 {
		return fmt.Errorf("at least one component required. Possible values: %s", strings.Join(SupportedComponents, ", "))
	}
	for _, comp := range c.Components {
		if !slices.Contains(SupportedComponents, strings.TrimSpace(comp)) {
			return fmt.Errorf("component %q not supported. Possible values: %s", comp, strings.Join(SupportedComponents, ", "))
		}
	}
	if c.Language == "js" && c.ManualInstrumentation && c.InstrumentationFile == "" {
		return fmt.Errorf(`when manual-instrumentation is set, an instrumentation file is required (InstrumentationFile or -instrumentation-file=path/to/file.js)`)
	}
	return nil
}

func GetArguments() Commands {
	if len(os.Args[1:]) < 1 {
		fmt.Println(color.RedString("You must pass a language used for your instrumentation, such as -language=js"))
		os.Exit(1)
	}

	languageValue := flag.String("language", "", "Language used for instrumentation (required). Possible values: dotnet, go, java, js, python")
	componentsString := flag.String("components", "", "Instrumentation components to test, separated by ',' (required). Possible values: sdk, collector, beyla, alloy")
	manualInstrumentation := flag.Bool("manual-instrumentation", false, "Provide if your application is using manual instrumentation")
	debug := flag.Bool("debug", false, "Output debug information")
	webServer := flag.Bool("web-server", false, "Set if you would like the results served in a web server in addition to console output")

	// javascript
	instrumentationFile := flag.String("instrumentation-file", "", `Name (including path) to instrumentation file. Required if using manual-instrumentation. E.g."-instrumentation-file=src/inst/instrumentation.js"`)
	packageJsonPath := flag.String("package-json-path", "", `Path to package.json file. Required if instrumentation is in JavaScript and the file is not in the same location as the otel-checker is being executed from. E.g. "-package-json-path=src/inst/"`)

	// collector
	collectorConfigPath := flag.String("collector-config-path", "", `Path to collector's config.yaml file. Required if using Collector and the config file is not in the same location as the otel-checker is being executed from. E.g. "-collector-config-path=src/inst/"`)
	flag.Parse()

	var components []string
	if *componentsString != "" {
		components = strings.Split(*componentsString, ",")
	}

	command := Commands{
		Language:              *languageValue,
		Components:            components,
		WebServer:             *webServer,
		ManualInstrumentation: *manualInstrumentation,
		InstrumentationFile:   *instrumentationFile,
		PackageJsonPath:       *packageJsonPath,
		CollectorConfigPath:   *collectorConfigPath,
		Debug:                 *debug,
	}

	if err := Validate(command); err != nil {
		fmt.Println(color.RedString(err.Error()))
		os.Exit(1)
	}

	if command.PackageJsonPath != "" && !strings.HasSuffix(command.PackageJsonPath, "/") {
		command.PackageJsonPath = command.PackageJsonPath + "/"
	}
	if command.CollectorConfigPath != "" && !strings.HasSuffix(command.CollectorConfigPath, "/") {
		command.CollectorConfigPath = command.CollectorConfigPath + "/"
	}
	return command
}

type Reporter struct {
	components []*ComponentReporter
}

type ComponentReporter struct {
	name     string
	Checks   []string
	Warnings []string
	Errors   []string
}

func (r *Reporter) Component(name string) *ComponentReporter {
	for _, component := range r.components {
		if component.name == name {
			return component
		}
	}
	c := &ComponentReporter{name: name}
	r.components = append(r.components, c)
	return c
}

// Results aggregates the checks, warnings, and errors across all components
// without producing any output. Callers that want to render results their own
// way (e.g. gcx) should use this; CLI callers should use PrintResults.
func (r *Reporter) Results() map[string][]string {
	res := make(map[string][]string)
	var checks []string
	for _, component := range r.components {
		checks = append(checks, component.Checks...)
	}
	res[CHECKS] = checks
	var warnings []string
	for _, component := range r.components {
		warnings = append(warnings, component.Warnings...)
	}
	res[WARNINGS] = warnings
	var errors []string
	for _, component := range r.components {
		errors = append(errors, component.Errors...)
	}
	res[ERRORS] = errors
	return res
}

func (r *Reporter) PrintResults() map[string][]string {
	res := r.Results()
	checks := res[CHECKS]
	warnings := res[WARNINGS]
	errors := res[ERRORS]

	if len(checks) > 0 {
		green := color.New(color.FgGreen)
		_, _ = green.Printf("\n%d Successful Check(s)\n", len(checks))
		for _, m := range checks {
			_, _ = green.Printf("✔ %s \n", m)
		}
	}
	if len(warnings) > 0 {
		yellow := color.New(color.FgYellow)
		_, _ = yellow.Printf("\n%d Warning(s)\n", len(warnings))
		for _, m := range warnings {
			_, _ = yellow.Printf("• %s \n", m)
		}
	}
	if len(errors) > 0 {
		red := color.New(color.FgRed)
		_, _ = red.Printf("\n%d Error(s)\n", len(errors))
		for _, m := range errors {
			_, _ = red.Printf("✖ %s \n", m)
		}
	}
	return res
}

func (r *ComponentReporter) AddSuccessfulCheck(message string) {
	r.Checks = append(r.Checks, fmt.Sprintf(`%s: %s`, r.name, message))
}

func (r *ComponentReporter) AddWarning(message string) {
	r.Warnings = append(r.Warnings, fmt.Sprintf(`%s: %s`, r.name, message))
}

func (r *ComponentReporter) AddInternalError(message string) {
	r.Warnings = append(r.Warnings, fmt.Sprintf(`%s: Internal Error: %s`, r.name, message))
}

func (r *ComponentReporter) AddError(message string) {
	r.Errors = append(r.Errors, fmt.Sprintf(`%s: %s`, r.name, message))
}

func FileExists(path string) bool {
	_, err := os.ReadFile(path)
	return err == nil
}
