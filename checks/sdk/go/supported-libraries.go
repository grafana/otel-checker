package _go

import (
	_ "embed"
	"fmt"
	"os"
	"otel-checker/checks/sdk"
	"otel-checker/checks/sdk/supported"
	"otel-checker/checks/utils"
	"strings"

	"golang.org/x/mod/semver"
)

//go:embed supported-libraries.yaml
var file []byte

func readGoModFile(reporter *utils.ComponentReporter) []supported.Library {
	if utils.FileExists("go.mod") {
		return readGoMod(reporter, "go.mod")
	}
	return nil
}

func readGoMod(reporter *utils.ComponentReporter, path string) []supported.Library {
	dat, err := os.ReadFile(path)
	if err != nil {
		reporter.AddError(fmt.Sprintf("Could not read go.mod: %v", err))
		return nil
	}
	return readGoModFromContent(dat)
}

func readGoModFromContent(content []byte) []supported.Library {
	var deps []supported.Library
	lines := strings.Split(string(content), "\n")

	inRequire := false
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "require (" {
			inRequire = true
			continue
		} else if line == ")" && inRequire {
			inRequire = false
			continue
		}

		if inRequire && line != "" {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				name := parts[0]
				version := parts[1]
				deps = append(deps, supported.Library{
					Name:    name,
					Version: version,
				})
			}
		}
	}

	return deps
}

func supportedLibraries() (supported.SupportedModules, error) {
	return supported.LoadSupportedLibraries(file)
}

func findSupportedLibraries(library supported.Library, s supported.SupportedModules) []string {
	var links []string
	for _, module := range s {
		for _, instrumentation := range module.Instrumentations {
			for _, version := range instrumentation.TargetVersions[supported.TypeLibrary] {
				versionRange, err := sdk.ParseVersionRange(version)
				if err != nil {
					panic(fmt.Sprintf("error parsing version range: %v", err))
				}
				if library.Name == instrumentation.Name {
					v := sdk.FixVersion(library.Version)
					if semver.IsValid(v) {
						if versionRange.Matches(v) {
							links = append(links, instrumentation.Link)
						}
					}
				}
			}
		}
	}
	return links
}

func CheckSupportedLibraries(reporter *utils.ComponentReporter, commands utils.Commands) {
	supported, err := supportedLibraries()
	if err != nil {
		reporter.AddError(fmt.Sprintf("Error reading supported libraries: %v", err))
		return
	}

	deps := readGoModFile(reporter)
	if len(deps) == 0 {
		return
	}

	for _, dep := range deps {
		links := findSupportedLibraries(dep, supported)
		if len(links) > 0 {
			reporter.AddSuccessfulCheck(
				fmt.Sprintf("Found supported library: %s:%s at %s",
					dep.Name, dep.Version, strings.Join(links, ", ")))
		} else if commands.Debug {
			reporter.AddWarning(fmt.Sprintf("Found unsupported library: %s:%s", dep.Name, dep.Version))
		}
	}
}
