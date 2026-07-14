package dotnet

import (
	"context"
	"fmt"
	"strconv"

	"github.com/grafana/otel-checker/checks/env"
	"github.com/grafana/otel-checker/checks/utils"
)

const minDotNetVersion = 8

func CheckDotNetSetup(ctx context.Context, reporter *utils.ComponentReporter, commands utils.Commands) {
	checkDotNetVersion(ctx, reporter)

	project, err := findAndLoadProject()

	if err != nil {
		reporter.AddErrorWithExplain("dotnet.project.not-found",
			fmt.Sprintf("Failed to find and load project: %s", err))
		return
	}

	reporter.AddSuccessfulCheck(fmt.Sprintf("Found project: %s", project.path))

	reportDotNetSupportedInstrumentations(ctx, reporter, project.SDK)

	if commands.ManualInstrumentation {
		checkDotNetCodeBasedInstrumentation(reporter)
	} else {
		checkDotNetAutoInstrumentation(reporter)
	}
}

func checkDotNetVersion(ctx context.Context, reporter *utils.ComponentReporter) {
	versionParts, err := readDotNetVersion(ctx)

	if err != nil {
		reporter.AddErrorWithExplain("dotnet.version.unknown",
			fmt.Sprintf("Could not check .NET version: %s", err))
		return
	}

	if len(versionParts) == 0 {
		reporter.AddErrorWithExplain("dotnet.version.empty",
			"Could not parse .NET version: version string is empty")
		return
	}
	majorVersion := versionParts[0]
	v, err := strconv.Atoi(majorVersion)

	if err != nil {
		reporter.AddErrorWithExplain("dotnet.version.unknown",
			fmt.Sprintf("Could not parse .NET version: %s", err))
		return
	}

	if v >= minDotNetVersion {
		reporter.AddSuccessfulCheck(fmt.Sprintf("Using .NET version equal or greater than minimum recommended (%d.0)", minDotNetVersion))
	} else {
		reporter.AddErrorWithExplain("dotnet.version.too-old",
			fmt.Sprintf("Not using recommended .NET version. Update your .NET SDK to at least version %d.0", minDotNetVersion))
	}
}

func checkDotNetAutoInstrumentation(reporter *utils.ComponentReporter) {
	env.CheckEnvVars(reporter, "dotnet",
		env.EnvVar{
			Name:          "CORECLR_ENABLE_PROFILING",
			RequiredValue: "1",
			ExplainID:     "dotnet.coreclr-enable-profiling.value-mismatch",
		},
		env.EnvVar{
			Name:          "CORECLR_PROFILER",
			RequiredValue: "{918728DD-259F-4A6A-AC2B-B85E1B658318}",
			ExplainID:     "dotnet.coreclr-profiler.value-mismatch",
		},
		env.EnvVar{
			Name:      "CORECLR_PROFILER_PATH",
			Required:  true,
			ExplainID: "dotnet.coreclr-profiler-path.unset",
		},
		env.EnvVar{
			Name:      "OTEL_DOTNET_AUTO_HOME",
			Required:  true,
			ExplainID: "dotnet.otel-dotnet-auto-home.unset",
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

func reportDotNetSupportedInstrumentations(ctx context.Context, reporter *utils.ComponentReporter, sdk string) {
	deps, err := ReadDependenciesFromCli(ctx)

	if err != nil {
		reporter.AddErrorWithExplain("dotnet.dependencies.unreadable",
			fmt.Sprintf("Failed to read dependencies: %s", err))
		return
	}

	instr := ReadAvailableInstrumentations()

	implicit, err := ImplicitPackagesForSdk(sdk)

	if err != nil {
		reporter.AddErrorWithExplain("dotnet.sdk.unrecognized",
			fmt.Sprintf("Unrecognized SDK: %s", sdk))
		return
	}

	if len(implicit) == 0 {
		reporter.AddWarningWithExplain("dotnet.sdk.no-implicit-packages",
			fmt.Sprintf("No implicit packages found for SDK: %s", sdk))
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
		reporter.AddErrorWithExplain("dotnet.project.no-dependencies",
			"No dependencies found in project")
		return
	}
}
