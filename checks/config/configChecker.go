package config

import (
	"fmt"
	"strings"

	"github.com/grafana/otel-checker/checks/utils"
)

// CheckConfigSetup reports on the outcome of loading and parsing the
// declarative config file. The heavy lifting (Load) is done once at
// the top of checks.Run so the parsed *File can be shared with other
// components (e.g. grafana-cloud). This function receives the
// results of that load and reports on them.
//
// candidatePaths is the list of paths the loader tried, in order. It's
// used to make the "file not found" message actionable.
func CheckConfigSetup(reporter *utils.ComponentReporter, resolvedPath string, candidatePaths []string, file *File, loadErr error) {
	if loadErr != nil {
		if resolvedPath == "" {
			reporter.AddErrorWithExplain("config.file.unreadable",
				fmt.Sprintf("Could not find a declarative config file. Tried: %s", strings.Join(candidatePaths, ", ")))
			return
		}
		reporter.AddErrorWithExplain("config.file.parse-error",
			fmt.Sprintf("Could not parse %s: %s", resolvedPath, loadErr))
		return
	}

	reporter.AddSuccessfulCheck(fmt.Sprintf("Parsed declarative config file: %s", resolvedPath))

	if file.FileFormat == "" {
		reporter.AddErrorWithExplain("config.file-format.missing",
			fmt.Sprintf("%s does not declare a file_format — set e.g. file_format: \"1.1\"", resolvedPath))
	} else {
		reporter.AddSuccessfulCheck(fmt.Sprintf("file_format is set to %q", file.FileFormat))
	}
}

var DefaultCandidatePaths = []string{"otel-config.yaml", "otel-config.yml"}

func Resolve(explicitPath string) (path string, candidates []string) {
	if explicitPath != "" {
		return explicitPath, []string{explicitPath}
	}
	for _, p := range DefaultCandidatePaths {
		if utils.FileExists(p) {
			return p, DefaultCandidatePaths
		}
	}
	return "", DefaultCandidatePaths
}
