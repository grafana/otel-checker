package dotnet

import (
	"fmt"
	"github.com/grafana/otel-checker/checks/env"
	"github.com/grafana/otel-checker/checks/utils"
	"strconv"
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

func checkDotNetAutoInstrumentation(reporter *utils.ComponentReporter) {
	env.CheckEnvVars(reporter, "dotnet",
		env.EnvVar{
			Name:          "CORECLR_ENABLE_PROFILING",
			RequiredValue: "1",
		},
		env.EnvVar{
			Name:          "CORECLR_PROFILER",
			RequiredValue: "{918728DD-259F-4A6A-AC2B-B85E1B658318}",
		},
		env.EnvVar{
			Name:     "CORECLR_PROFILER_PATH",
			Required: true,
		},
		env.EnvVar{
			Name:     "OTEL_DOTNET_AUTO_HOME",
			Required: true,
		})
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
