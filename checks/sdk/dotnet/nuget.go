package dotnet

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type Package struct {
	ID               string `json:"id"`
	RequestedVersion string `json:"requestedVersion,omitempty"`
	ResolvedVersion  string `json:"resolvedVersion"`
}

type Framework struct {
	Framework          string    `json:"framework"`
	TopLevelPackages   []Package `json:"topLevelPackages"`
	TransitivePackages []Package `json:"transitivePackages"`
}

type Project struct {
	Path       string      `json:"path"`
	Frameworks []Framework `json:"frameworks"`
}

type NuGetPackageList struct {
	Version    int       `json:"version"`
	Parameters string    `json:"parameters"`
	Projects   []Project `json:"projects"`
}

// ListPackageDependencies runs 'dotnet list package' to get all package dependencies (including transitive)
// for the current project in the working directory
func ReadDependenciesFromCli() (*NuGetPackageList, error) {
	cmd := exec.Command("dotnet", "list", "package", "--format", "json", "--include-transitive")
	stdout, err := cmd.Output()

	if err != nil {
		return nil, fmt.Errorf("failed to run dotnet list package: %w", err)
	}

	var deps NuGetPackageList
	if err := json.Unmarshal(stdout, &deps); err != nil {
		return nil, fmt.Errorf("failed to parse dependencies JSON: %w", err)
	}

	return &deps, nil
}
