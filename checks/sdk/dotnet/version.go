package dotnet

import (
	"fmt"
	"os/exec"
	"strings"
)

func readDotNetVersion() ([]string, error) {
	cmd := exec.Command("dotnet", "--version")
	stdout, err := cmd.Output()

	if err != nil {
		return nil, fmt.Errorf("Could not check .NET version: %s", err)
	}

	version := strings.TrimSpace(string(stdout))
	versionParts := strings.Split(version, ".")
	if len(versionParts) == 0 {
		return nil, fmt.Errorf("Could not parse .NET version: version string is empty")
	}

	return versionParts, nil
}
