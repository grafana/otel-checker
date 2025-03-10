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
	Language                string
	Components              []string
	ManualInstrumentation   bool
	WebServer               bool
	InstrumentationFile     string
	PackageJsonPath         string
	CollectorConfigPath     string
	Debug                   bool
	EnableGrafanaCloudCheck bool
}

func GetArguments() Commands {
	command := Commands{}
	args := os.Args[1:]
	if len(args) < 1 {
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

	possibleLanguages := []string{"dotnet", "go", "java", "js", "python", "ruby", "php"}
	if !slices.Contains(possibleLanguages, *languageValue) {
		fmt.Println(color.RedString(fmt.Sprintf("Language %s not supported. Possible values: dotnet, go, java, js, python, ruby", *languageValue)))
		os.Exit(1)
	}

	if *componentsString == "" {
		fmt.Println(color.RedString(`Component flag required. Possible values: sdk, beyla, alloy, collector. E.g. -components="sdk,collector"`))
		os.Exit(1)
	}

	possibleComponents := []string{"sdk", "beyla", "alloy", "collector", "grafana-cloud"}
	components := strings.Split(*componentsString, ",")
	for _, c := range components {
		if !slices.Contains(possibleComponents, strings.Trim(c, " ")) {
			fmt.Println(color.RedString(fmt.Sprintf(`Component %s not supported. Possible values: sdk, collector, beyla, alloy. E.g. -components="sdk,collector"`, c)))
			os.Exit(1)
		}
	}

	// javascript
	if *languageValue == "js" && *instrumentationFile == "" && *manualInstrumentation {
		fmt.Println(color.RedString(`When manual-instrumentation is being used, a instrumentation file is required. Remove "-manual-instrumentation" or "-instrumentation-file=path/to/file/file.js"`))
		os.Exit(1)
	}
	if *packageJsonPath != "" && !strings.HasSuffix(*packageJsonPath, "/") {
		*packageJsonPath = *packageJsonPath + "/"
	}

	// collector
	if *collectorConfigPath != "" && !strings.HasSuffix(*collectorConfigPath, "/") {
		*collectorConfigPath = *collectorConfigPath + "/"
	}

	command.Language = *languageValue
	command.Components = components
	command.WebServer = *webServer
	command.ManualInstrumentation = *manualInstrumentation
	command.InstrumentationFile = *instrumentationFile
	command.PackageJsonPath = *packageJsonPath
	command.CollectorConfigPath = *collectorConfigPath
	command.Debug = *debug
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

func (r *Reporter) PrintResults() map[string][]string {
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

	if len(checks) > 0 {
		green := color.New(color.FgGreen)
		green.Printf("\n%d Successful Check(s)\n", len(checks))
		for _, m := range checks {
			green.Printf("✔ %s \n", m)
		}
	}
	if len(warnings) > 0 {
		yellow := color.New(color.FgYellow)
		yellow.Printf("\n%d Warning(s)\n", len(warnings))
		for _, m := range warnings {
			yellow.Printf("• %s \n", m)
		}
	}
	if len(errors) > 0 {
		red := color.New(color.FgRed)
		red.Printf("\n%d Error(s)\n", len(errors))
		for _, m := range errors {
			red.Printf("✖ %s \n", m)
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

func (r *ComponentReporter) AddError(message string) {
	r.Errors = append(r.Errors, fmt.Sprintf(`%s: %s`, r.name, message))
}

func FileExists(path string) bool {
	_, err := os.ReadFile(path)
	return err == nil
}
