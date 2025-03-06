package dotnet

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"

	"otel-checker/checks/utils"
)

const minDotNetVersion = 8

func CheckDotNetSetup(reporter *utils.ComponentReporter, commands utils.Commands) {
	checkDotNetVersion(reporter)

	project, err := findAndLoadProject()

	if err != nil {
		reporter.AddError(fmt.Sprintf("Failed to find and load project: %s", err))
		return
	}

	reporter.AddSuccessfulCheck(fmt.Sprintf("Found project: %s", project.path))

	reportDotNetSupportedInstrumentations(reporter, project.SDK)

	if commands.ManualInstrumentation {
		checkDotNetCodeBasedInstrumentation(reporter)
	} else {
		checkDotNetAutoInstrumentation(reporter)
	}
}

func checkDotNetVersion(reporter *utils.ComponentReporter) {
	versionParts, err := readDotNetVersion()

	if err != nil {
		reporter.AddError(fmt.Sprintf("Could not check .NET version: %s", err))
		return
	}

	if len(versionParts) == 0 {
		reporter.AddError("Could not parse .NET version: version string is empty")
		return
	}
	majorVersion := versionParts[0]
	v, err := strconv.Atoi(majorVersion)

	if err != nil {
		reporter.AddError(fmt.Sprintf("Could not parse .NET version: %s", err))
		return
	}

	if v >= minDotNetVersion {
		reporter.AddSuccessfulCheck(fmt.Sprintf("Using .NET version equal or greater than minimum recommended (%d.0)", minDotNetVersion))
	} else {
		reporter.AddError(fmt.Sprintf("Not using recommended .NET version. Update your .NET SDK to at least version %d.0", minDotNetVersion))
	}
}

type EnvVarValidator func(string, string) error

func checkEnvironmentVariables(
	reporter *utils.ComponentReporter,
	requiredVars []string,
	expectedValues map[string]string,
	customValidators map[string]EnvVarValidator,
) bool {
	// Check for missing required variables
	missingVars := []string{}
	for _, envVar := range requiredVars {
		if _, exists := os.LookupEnv(envVar); !exists {
			missingVars = append(missingVars, envVar)
		}
	}

	if len(missingVars) > 0 {
		reporter.AddError(fmt.Sprintf("Missing required environment variables: %s", strings.Join(missingVars, ", ")))
		return false
	}

	// Check for incorrect values
	wrongValues := make(map[string]string)
	for envVar, expectedValue := range expectedValues {
		envVarValue := os.Getenv(envVar)
		if envVarValue != expectedValue {
			wrongValues[envVar] = envVarValue
		}
	}

	if len(wrongValues) > 0 {
		s := make([]string, 0, len(wrongValues))
		v := slices.Collect(maps.Keys(wrongValues))
		// Sort keys to make the output deterministic
		sort.Strings(v)
		for _, k := range v {
			s = append(s, fmt.Sprintf("%s: %s", k, wrongValues[k]))
		}
		reporter.AddError(fmt.Sprintf("Incorrect values for environment variables: %s", strings.Join(s, ", ")))
		return false
	}

	// Run custom validators
	for envVar, validator := range customValidators {
		if envVarValue, exists := os.LookupEnv(envVar); exists {
			if err := validator(envVar, envVarValue); err != nil {
				reporter.AddError(fmt.Sprintf("Validation failed for %s: %s", envVar, err))
				return false
			}
		}
	}

	return true
}

func checkDotNetAutoInstrumentation(reporter *utils.ComponentReporter) {
	requiredEnvVars := []string{
		"CORECLR_ENABLE_PROFILING",
		"CORECLR_PROFILER",
		"CORECLR_PROFILER_PATH",
		"OTEL_DOTNET_AUTO_HOME",
	}

	expectedValues := map[string]string{
		"CORECLR_ENABLE_PROFILING": "1",
		"CORECLR_PROFILER":         "{918728DD-259F-4A6A-AC2B-B85E1B658318}",
	}

	if success := checkEnvironmentVariables(reporter, requiredEnvVars, expectedValues, nil); success {
		reporter.AddSuccessfulCheck("All required environment variables for .NET auto-instrumentation are set with correct values.")
	}
}

func checkDotNetCodeBasedInstrumentation(reporter *utils.ComponentReporter) {}

func findAndLoadProject() (*CSharpProject, error) {
	projectPath, err := FindCSharpProject(".")
	if err != nil {
		return nil, err
	}

	project, err := LoadCSharpProject(projectPath)

	if err != nil {
		return nil, err
	}

	return project, nil
}

func reportDotNetSupportedInstrumentations(reporter *utils.ComponentReporter, sdk string) {
	deps, err := ReadDependenciesFromCli()

	if err != nil {
		reporter.AddError(fmt.Sprintf("Failed to read dependencies: %s", err))
		return
	}

	instr := ReadAvailableInstrumentations()

	implicit, err := ImplicitPackagesForSdk(sdk)

	if err != nil {
		reporter.AddError(fmt.Sprintf("Unrecognized SDK: %s", sdk))
		return
	}

	if len(implicit) == 0 {
		reporter.AddWarning(fmt.Sprintf("No implicit packages found for SDK: %s", sdk))
	} else {
		for _, pkg := range implicit {
			lib, ok := instr[pkg]

			if !ok {
				continue
			}

			reporter.AddSuccessfulCheck(fmt.Sprintf("Found supported instrumentation for %s: %s", pkg, lib))
		}
	}

	for _, project := range deps.Projects {
		for _, framework := range project.Frameworks {
			packages := append(framework.TopLevelPackages, framework.TransitivePackages...)
			for _, pkg := range packages {
				lib, ok := instr[pkg.ID]

				if !ok {
					continue
				}

				reporter.AddSuccessfulCheck(fmt.Sprintf("Found supported instrumentation for %s: %s", pkg.ID, lib))
			}
		}
	}
	if len(deps.Projects) == 0 {
		reporter.AddError("No dependencies found in project")
		return
	}
}
