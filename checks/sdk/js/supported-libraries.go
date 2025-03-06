package js

import (
	_ "embed"
	"encoding/json"
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

func readDependencies(reporter *utils.ComponentReporter) []supported.Library {
	// Try package-lock.json first
	if utils.FileExists("package-lock.json") {
		return readPackageLock(reporter)
	}
	// Fall back to package.json
	if utils.FileExists("package.json") {
		return readPackageJson(reporter)
	}
	return nil
}

func readPackageLock(reporter *utils.ComponentReporter) []supported.Library {
	dat, err := os.ReadFile("package-lock.json")
	if err != nil {
		reporter.AddError(fmt.Sprintf("Could not read package-lock.json: %v", err))
		return nil
	}

	var lock struct {
		Dependencies map[string]struct {
			Version string `json:"version"`
		} `json:"dependencies"`
	}

	if err := json.Unmarshal(dat, &lock); err != nil {
		reporter.AddError(fmt.Sprintf("Could not parse package-lock.json: %v", err))
		return nil
	}

	var deps []supported.Library
	for name, dep := range lock.Dependencies {
		deps = append(deps, supported.Library{
			Name:    name,
			Version: dep.Version,
		})
	}

	if len(deps) == 0 {
		reporter.AddWarning("No dependencies found in package-lock.json")
	}
	return deps
}

func readPackageJson(reporter *utils.ComponentReporter) []supported.Library {
	dat, err := os.ReadFile("package.json")
	if err != nil {
		reporter.AddError(fmt.Sprintf("Could not read package.json: %v", err))
		return nil
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}

	if err := json.Unmarshal(dat, &pkg); err != nil {
		reporter.AddError(fmt.Sprintf("Could not parse package.json: %v", err))
		return nil
	}

	var deps []supported.Library
	for name, version := range pkg.Dependencies {
		// Remove ^ or ~ from version
		version = strings.TrimPrefix(version, "^")
		version = strings.TrimPrefix(version, "~")
		deps = append(deps, supported.Library{
			Name:    name,
			Version: version,
		})
	}

	for name, version := range pkg.DevDependencies {
		version = strings.TrimPrefix(version, "^")
		version = strings.TrimPrefix(version, "~")
		deps = append(deps, supported.Library{
			Name:    name,
			Version: version,
		})
	}

	if len(deps) == 0 {
		reporter.AddWarning("No dependencies found in package.json")
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
