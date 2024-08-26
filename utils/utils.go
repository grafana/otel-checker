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
	Language            string
	Components          []string
	AutoInstrumentation bool
	InstrumentationFile string
	PackageJsonPath     string
	CollectorConfigPath string
}

func GetArguments() Commands {
	command := Commands{}
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println(color.RedString("You must pass a language used for your instrumentation, such as --language=js"))
		os.Exit(1)
	}

	languageValue := flag.String("language", "", "Language used for instrumentation (required). Possible values: dotnet, go, java, js, python")
	componentsString := flag.String("components", "", "Instrumentation components to test, separated by ',' (required). Possible values: sdk, collector, beyla, alloy")
	autoInstrumentation := flag.Bool("auto-instrumentation", false, "Provide if your application is using auto instrumentation")
	instrumentationFile := flag.String("instrumentation-file", "", `Name (including path) to instrumentation file. Required if not using auto-instrumentation. E.g."-instrumentation-file=src/inst/instrumentation.js"`)
	packageJsonPath := flag.String("package-json-path", "", `Path to package.json file. Required if instrumentation is in JavaScript and the file is not in the same location as the otel-checker is being executed from. E.g. "-package-json-path=src/inst/"`)
	collectorConfigPath := flag.String("collector-config-path", "", `Path to collector's config.yaml file. Required if using Collector and the config file is not in the same location as the otel-checker is being executed from. E.g. "-collector-config-path=src/inst/"`)
	flag.Parse()

	possibleLanguages := []string{"dotnet", "go", "java", "js", "python", "ruby"}
	if !slices.Contains(possibleLanguages, *languageValue) {
		fmt.Println(color.RedString(fmt.Sprintf("Language %s not supported. Possible values: dotnet, go, java, js, python, ruby", *languageValue)))
		os.Exit(1)
	}

	if *componentsString == "" {
		fmt.Println(color.RedString(`Component flag required. Possible values: sdk, beyla, alloy, collector. E.g. -components="sdk,collector"`))
		os.Exit(1)
	}

	possibleComponents := []string{"sdk", "beyla", "alloy", "collector"}
	components := strings.Split(*componentsString, ",")
	for _, c := range components {
		if !slices.Contains(possibleComponents, strings.Trim(c, " ")) {
			fmt.Println(color.RedString(fmt.Sprintf(`Component %s not supported. Possible values: sdk, collector, beyla, alloy. E.g. -components="sdk,collector"`, c)))
			os.Exit(1)
		}
	}

	if *instrumentationFile == "" && !*autoInstrumentation {
		fmt.Println(color.RedString(`When auto-instrumentation is not being used, a instrumentation file is required. Add "-auto-instrumentation" or "-instrumentation-file=path/to/file/file.js"`))
		os.Exit(1)
	}
	if *packageJsonPath != "" && !strings.HasSuffix(*packageJsonPath, "/") {
		*packageJsonPath = *packageJsonPath + "/"
	}
	if *collectorConfigPath != "" && !strings.HasSuffix(*collectorConfigPath, "/") {
		*collectorConfigPath = *collectorConfigPath + "/"
	}

	command.Language = *languageValue
	command.Components = components
	command.AutoInstrumentation = *autoInstrumentation
	command.InstrumentationFile = *instrumentationFile
	command.PackageJsonPath = *packageJsonPath
	command.CollectorConfigPath = *collectorConfigPath
	return command
}

func CreateMessagesMap() map[string][]string {
	messages := make(map[string][]string)
	messages[CHECKS] = make([]string, 0)
	messages[WARNINGS] = make([]string, 0)
	messages[ERRORS] = make([]string, 0)

	return messages
}

func PrintResults(messages map[string][]string) {
	if len(messages[CHECKS]) > 0 {
		green := color.New(color.FgGreen)
		green.Printf("\n%d Successful Check(s)\n", len(messages[CHECKS]))
		for _, m := range messages[CHECKS] {
			green.Printf("✔ %s \n", m)
		}
	}
	if len(messages[WARNINGS]) > 0 {
		yellow := color.New(color.FgYellow)
		yellow.Printf("\n%d Warning(s)\n", len(messages[WARNINGS]))
		for _, m := range messages[WARNINGS] {
			yellow.Printf("• %s \n", m)
		}
	}
	if len(messages[ERRORS]) > 0 {
		red := color.New(color.FgRed)
		red.Printf("\n%d Error(s)\n", len(messages[ERRORS]))
		for _, m := range messages[ERRORS] {
			red.Printf("✖ %s \n", m)
		}
	}
}

func AddSuccessfulCheck(messages *map[string][]string, component string, message string) {
	(*messages)[CHECKS] = append((*messages)[CHECKS], fmt.Sprintf(`%s: %s`, component, message))
}

func AddWarning(messages *map[string][]string, component string, message string) {
	(*messages)[WARNINGS] = append((*messages)[WARNINGS], fmt.Sprintf(`%s: %s`, component, message))
}

func AddError(messages *map[string][]string, component string, message string) {
	(*messages)[ERRORS] = append((*messages)[ERRORS], fmt.Sprintf(`%s: %s`, component, message))
}
