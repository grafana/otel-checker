package dotnet

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func readDotNetVersion(ctx context.Context) ([]string, error) {
	cmd := exec.CommandContext(ctx, "dotnet", "--version")
	stdout, err := cmd.Output()

	if err != nil {
		return nil, fmt.Errorf("could not check .NET version: %s", err)
	}

	version := strings.TrimSpace(string(stdout))
	versionParts := strings.Split(version, ".")
	if len(versionParts) == 0 {
		return nil, fmt.Errorf("could not parse .NET version: version string is empty")
	}

	return versionParts, nil
}
