package _go

import (
	_ "embed"
	"fmt"
	"github.com/grafana/otel-checker/checks/sdk/supported"
	"github.com/grafana/otel-checker/checks/utils"
	"os"
	"strings"
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

// CheckSupportedLibraries checks if Go dependencies are supported by OpenTelemetry
func CheckSupportedLibraries(reporter *utils.ComponentReporter, commands utils.Commands) {
	supportedLibs, err := supportedLibraries()
	if err != nil {
		reporter.AddError(fmt.Sprintf("Error reading supported libraries: %v", err))
		return
	}

	deps := readGoModFile(reporter)
	supported.CheckLibraries(reporter, commands, supportedLibs, deps, supported.TypeLibrary)
}
