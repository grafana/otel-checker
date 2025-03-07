package python

import (
	_ "embed"
	"fmt"
	"os"
	"otel-checker/checks/sdk"
	"otel-checker/checks/utils"
	"regexp"
	"strconv"
	"strings"
)

func CheckSetup(reporter *utils.ComponentReporter, commands utils.Commands) {
	checkPythonVersion(reporter)
	if commands.ManualInstrumentation {
		checkCodeBasedInstrumentation(reporter, commands.Debug)
	} else {
		checkAutoInstrumentation(reporter, commands.Debug)
	}
}

func checkPythonVersion(reporter *utils.ComponentReporter) {}

func checkAutoInstrumentation(reporter *utils.ComponentReporter, debug bool) {
	reportSupportedLibraries(reporter, debug)
}

func checkCodeBasedInstrumentation(reporter *utils.ComponentReporter, debug bool) {
	reportSupportedLibraries(reporter, debug)
}

func reportSupportedLibraries(reporter *utils.ComponentReporter, debug bool) {
	supported, err := supportedLibraries()
	if err != nil {
		reporter.AddError(fmt.Sprintf("Error reading supported libraries: %v", err))
	}

	deps := readDependencies(reporter)
	outputSupportedLibraries(deps, supported, reporter, debug)
}

func readDependencies(reporter *utils.ComponentReporter) []Library {
	path := "requirements.txt"
	if utils.FileExists(path) {
		return readRequirementsTxt(reporter, path)
	}
	// also https://github.com/fpgmaas/cookiecutter-poetry-example/blob/main/poetry.lock
	return nil
}

func readRequirementsTxt(reporter *utils.ComponentReporter, path string) []Library {
	readFile, err := os.ReadFile(path)
	if err != nil {
		reporter.AddError(fmt.Sprintf("Could not read file %s: %v", path, err))
		return nil
	}

	deps := parseRequirementsTxt(reporter, string(readFile))
	if len(deps) == 0 {
		reporter.AddWarning(fmt.Sprintf("No dependencies found in %s", path))
	}
	return deps
}

func outputSupportedLibraries(deps []Library, supported []SupportedLibrary, reporter *utils.ComponentReporter, debug bool) {
	for _, dep := range deps {
		links := findSupportedLibraries(dep, supported)
		if len(links) > 0 {
			reporter.AddSuccessfulCheck(
				fmt.Sprintf("Found supported library: %s:%s at %s",
					dep.Name, dep.Version, strings.Join(links, ", ")))
		} else if debug {
			reporter.AddWarning(fmt.Sprintf("Found unsupported library: %s:%s", dep.Name, dep.Version))
		}
	}
}

func parseRequirementsTxt(reporter *utils.ComponentReporter, lines string) []Library {
	var deps []Library
	for _, line := range strings.Split(lines, "\n") {
		if line == "" {
			continue
		}
		// e.g. blinker==1.9.0
		split := strings.Split(line, "==")
		if len(split) != 2 {
			reporter.AddWarning(fmt.Sprintf("Could not parse line: %s", line))
			continue
		}
		deps = append(deps, Library{
			Name:    strings.ToLower(strings.TrimSpace(split[0])),
			Version: strings.TrimSpace(split[1]),
		})
	}
	return deps
}

type Library struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type SupportedLibrary struct {
	Name         string
	Link         string
	VersionRange map[string]sdk.VersionRange
}

var linkRegex = regexp.MustCompile(`\[opentelemetry-instrumentation-(.*)]`)

func supportedLibraries() ([]SupportedLibrary, error) {
	bytes, err := sdk.LoadUrl("https://raw.githubusercontent.com/open-telemetry/opentelemetry-python-contrib/refs/heads/main/instrumentation/README.md")
	if err != nil {
		return nil, err
	}
	var res []SupportedLibrary
	for _, library := range strings.Split(string(bytes), "\n")[3:] {
		library = strings.TrimSpace(library)
		if library == "" {
			continue
		}
		l := strings.Split(library, "|")
		mdLink := strings.TrimSpace(l[1])
		name := linkRegex.FindStringSubmatch(mdLink)[1]
		url := fmt.Sprintf("https://github.com/open-telemetry/opentelemetry-python-contrib/tree/main/instrumentation/opentelemetry-instrumentation-%s", name)
		ranges, err := versionRanges(strings.TrimSpace(l[2]))
		if err != nil {
			return nil, err
		}
		res = append(res, SupportedLibrary{
			Name:         name,
			Link:         url,
			VersionRange: ranges,
		})
	}
	return res, nil
}

func versionRanges(list string) (map[string]sdk.VersionRange, error) {
	res := map[string]sdk.VersionRange{}
	name := ""
	var err error
	for _, s := range strings.Split(list, ",") {
		if s == "<0.15" {
			s = "< 0.15"
		}

		statement := strings.Split(strings.TrimSpace(s), " ")
		if len(statement) == 3 {
			name = statement[0]
			err = addVersionRange(res, name, statement[1], statement[2])
			if err != nil {
				goto er
			}
		} else if len(statement) == 2 {
			err = addVersionRange(res, name, statement[0], statement[1])
			if err != nil {
				goto er
			}
		} else if len(statement) == 1 {
			// no version range => all versions
			res[name] = sdk.VersionRange{}
		} else {
			err = fmt.Errorf("invalid version range statement: %s", s)
			goto er
		}
	}
	return res, nil

er:
	return nil, fmt.Errorf("error parsing version %s: %v", list, err)
}

func addVersionRange(res map[string]sdk.VersionRange, name string, op string, version string) error {
	err := isVersion(version)
	if err != nil {
		return fmt.Errorf("invalid version: %s", version)
	}

	r, err := newRange(op, version)
	if err != nil {
		return err
	}
	old, ok := res[name]
	if ok {
		r = mergeRanges(r, old)
	}
	res[name] = r
	return nil
}

func isVersion(version string) error {
	_, err := strconv.Atoi(fmt.Sprintf("%c", version[0]))
	return err
}

func mergeRanges(r1 sdk.VersionRange, r2 sdk.VersionRange) sdk.VersionRange {
	if r1.Lower == "" {
		r1.Lower = r2.Lower
		r1.LowerInclusive = r2.LowerInclusive
	}
	if r1.Upper == "" {
		r1.Upper = r2.Upper
		r1.UpperInclusive = r2.UpperInclusive
	}
	return r1
}

func newRange(op string, version string) (sdk.VersionRange, error) {
	switch strings.TrimSpace(op) {
	case "<":
		return sdk.VersionRange{
			Upper: version,
		}, nil
	case "<=":
		return sdk.VersionRange{
			Upper:          version,
			UpperInclusive: true,
		}, nil
	case ">=":
		return sdk.VersionRange{
			Lower:          version,
			LowerInclusive: true,
		}, nil
	case "~=":
		part, err := upperBoundForTilde(version)
		if err != nil {
			return sdk.VersionRange{}, err
		}
		return sdk.VersionRange{
			Lower:          version,
			Upper:          part,
			LowerInclusive: true,
		}, nil
	}
	return sdk.VersionRange{}, fmt.Errorf("invalid version range operation: '%s'", op)
}

func upperBoundForTilde(version string) (string, error) {
	split := strings.Split(version, ".")
	split = split[:len(split)-1]
	last, err := strconv.Atoi(split[len(split)-1])
	if err != nil {
		return "", err
	}
	split[len(split)-1] = strconv.Itoa(last + 1)
	return strings.Join(split, "."), nil
}

func findSupportedLibraries(want Library, supported []SupportedLibrary) []string {
	var links []string
	for _, lib := range supported {
		for dep, versionRange := range lib.VersionRange {
			if dep == want.Name && versionRange.Matches(want.Version) {
				links = append(links, lib.Link)
			}
		}
	}
	return links
}
