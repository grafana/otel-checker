package js

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/grafana/otel-checker/checks/sdk/supported"
	"github.com/grafana/otel-checker/checks/utils"
	"os"
	"strings"
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
	return readPackageLockFromContent(dat)
}

func readPackageLockFromContent(content []byte) []supported.Library {
	var lock struct {
		Packages map[string]struct {
			Version string `json:"version"`
		} `json:"packages"`
	}

	if err := json.Unmarshal(content, &lock); err != nil {
		return nil
	}

	var deps []supported.Library
	for path, pkg := range lock.Packages {
		// Skip the root package
		if path == "" {
			continue
		}
		// Extract package name from path (e.g. "node_modules/@fastify/ajv-compiler" -> "@fastify/ajv-compiler")
		name := strings.TrimPrefix(path, "node_modules/")
		deps = append(deps, supported.Library{
			Name:    name,
			Version: pkg.Version,
		})
	}

	return deps
}

func readPackageJson(reporter *utils.ComponentReporter) []supported.Library {
	dat, err := os.ReadFile("package.json")
	if err != nil {
		reporter.AddError(fmt.Sprintf("Could not read package.json: %v", err))
		return nil
	}
	return readPackageJsonFromContent(dat)
}

func readPackageJsonFromContent(content []byte) []supported.Library {
	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}

	if err := json.Unmarshal(content, &pkg); err != nil {
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

	return deps
}

func supportedLibraries() (supported.SupportedModules, error) {
	return supported.LoadSupportedLibraries(file)
}

// CheckSupportedLibraries checks if JS dependencies are supported by OpenTelemetry
func CheckSupportedLibraries(reporter *utils.ComponentReporter, commands utils.Commands) {
	supportedLibs, err := supportedLibraries()
	if err != nil {
		reporter.AddError(fmt.Sprintf("Error reading supported libraries: %v", err))
		return
	}

	deps := readDependencies(reporter)
	supported.CheckLibraries(reporter, commands, supportedLibs, deps, supported.TypeLibrary)
}
