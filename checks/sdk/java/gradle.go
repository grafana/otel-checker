package java

import (
	"fmt"
	"github.com/grafana/otel-checker/checks/sdk"
	"github.com/grafana/otel-checker/checks/utils"
	"os/exec"
	"slices"
	"strings"
)

var gradleFiles = []string{
	"build.gradle",
	"build.gradle.kts",
}

func checkGradle(file string, reporter *utils.ComponentReporter) []Library {
	println("Reading Gradle dependencies")

	out := sdk.RunCommand(reporter, exec.Command(searchWrapper("gradle", "gradlew"),
		fmt.Sprintf("--build-file=%s", file), "dependencies", "--configuration=runtimeClasspath"))
	if out == "" {
		return []Library{}
	}
	deps := parseGradleDeps(out)
	if len(deps) == 0 {
		reporter.AddWarning("No Gradle dependencies found")
	}
	return deps
}

func parseGradleDeps(out string) []Library {
	lines := strings.Split(out, "\n")
	var deps []Library
	for _, l := range lines {
		if strings.Contains(l, "---") {
			index := strings.Index(l, "---")
			dep := strings.TrimSpace(l[index+4:])
			split := strings.Split(dep, ":")
			if len(split) == 3 {
				d := Library{
					Group:    split[0],
					Artifact: split[1],
					Version:  split[2],
				}
				s := d.String()
				if !slices.ContainsFunc(
					deps,
					func(l Library) bool {
						return l.String() == s
					}) {
					deps = append(deps, d)
				}
			}
		}
	}

	return deps
}
