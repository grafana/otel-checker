package supported

import (
	"fmt"
	"github.com/grafana/otel-checker/checks/sdk"
	"github.com/grafana/otel-checker/checks/utils"
	"strings"

	"golang.org/x/mod/semver"
)

// CheckLibraries is a common function to check supported libraries
func CheckLibraries(reporter *utils.ComponentReporter,
	commands utils.Commands,
	supportedLibs SupportedModules,
	dependencies []Library,
	instrumentationType InstrumentationType) {
	if len(dependencies) == 0 {
		return
	}

	for _, dep := range dependencies {
		links := FindSupportedLibraries(dep, supportedLibs, instrumentationType)
		if len(links) > 0 {
			reporter.AddSuccessfulCheck(
				fmt.Sprintf("Found supported library: %s:%s at %s",
					dep.Name, dep.Version, strings.Join(links, ", ")))
		} else if commands.Debug {
			reporter.AddWarning(fmt.Sprintf("Found unsupported library: %s:%s", dep.Name, dep.Version))
		}
	}
}

// FindSupportedLibraries checks if a library is supported by any instrumentation
func FindSupportedLibraries(library Library, supportedModules SupportedModules, instrumentationType InstrumentationType) []string {
	var links []string
	for _, module := range supportedModules {
		for _, instrumentation := range module.Instrumentations {
			for _, version := range instrumentation.TargetVersions[instrumentationType] {
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
