package sdk

import (
	_ "embed"
	"fmt"
	"otel-checker/checks/utils"
	"regexp"
	"strconv"
	"strings"
)

func CheckPythonSetup(reporter *utils.ComponentReporter, manualInstrumentation bool) {
	checkPythonVersion(reporter)
	if !manualInstrumentation {
		checkPythonAutoInstrumentation(reporter)
	} else {
		checkPythonCodeBasedInstrumentation(reporter)
	}
}

func checkPythonVersion(reporter *utils.ComponentReporter) {}

func checkPythonAutoInstrumentation(reporter *utils.ComponentReporter) {}

func checkPythonCodeBasedInstrumentation(reporter *utils.ComponentReporter) {}

type PythonLibrary struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type SupportedPythonLibrary struct {
	Name         string
	Link         string
	VersionRange map[string]VersionRange
}

//go:embed python-instrumentations.md
var supportedLibraries string

var linkRegex = regexp.MustCompile("\\[opentelemetry-instrumentation-(.*)]")

func supportedPythonLibraries() ([]SupportedPythonLibrary, error) {
	var res []SupportedPythonLibrary
	for _, library := range strings.Split(supportedLibraries, "\n")[2:] {
		library = strings.TrimSpace(library)
		if library == "" {
			continue
		}
		l := strings.Split(library, "|")
		mdLink := strings.TrimSpace(l[1])
		name := linkRegex.FindStringSubmatch(mdLink)[1]
		url := fmt.Sprintf("https://github.com/open-telemetry/opentelemetry-python-contrib/tree/main/instrumentation/opentelemetry-instrumentation-%s", name)
		versionRange := strings.TrimSpace(l[2])
		ranges, err := pythonVersionRanges(versionRange)
		if err != nil {
			return nil, err
		}
		res = append(res, SupportedPythonLibrary{
			Name:         name,
			Link:         url,
			VersionRange: ranges,
		})
	}
	return res, nil
}

func pythonVersionRanges(list string) (map[string]VersionRange, error) {
	res := map[string]VersionRange{}
	name := ""
	for _, s := range strings.Split(list, ",") {
		statement := strings.Split(s, " ")
		if len(statement) == 3 {
			name = statement[0]
			err := addVersionRange(res, name, statement[1], statement[2])
			if err != nil {
				return nil, err
			}
		} else if len(statement) == 2 {
			err := addVersionRange(res, name, statement[0], statement[1])
			if err != nil {
				return nil, err
			}
		} else if len(statement) == 1 {
			// no version range => all versions
			res[name] = VersionRange{}
		} else {
			return nil, fmt.Errorf("invalid version range statement: %s", s)
		}
	}
	return res, nil
}

func addVersionRange(res map[string]VersionRange, name string, op string, version string) error {
	r, err := newRange(op, version)
	if err != nil {
		return fmt.Errorf("error parsing version range for %s: %v", name, err)
	}
	old, ok := res[name]
	if ok {
		r = mergeRanges(r, old)
	}
	res[name] = r
	return nil
}

func mergeRanges(r1 VersionRange, r2 VersionRange) VersionRange {
	if r1.lower == "" {
		r1.lower = r2.lower
		r1.lowerInclusive = r2.lowerInclusive
	}
	if r1.upper == "" {
		r1.upper = r2.upper
		r1.upperInclusive = r2.upperInclusive
	}
	return r1
}

func newRange(op string, version string) (VersionRange, error) {
	switch strings.TrimSpace(op) {
	case "<":
		return VersionRange{
			upper: version,
		}, nil
	case "<=":
		return VersionRange{
			upper:          version,
			upperInclusive: true,
		}, nil
	case ">=":
		return VersionRange{
			lower:          version,
			lowerInclusive: true,
		}, nil
	case "~=":
		part, err := upperBoundForTilde(version)
		if err != nil {
			return VersionRange{}, err
		}
		return VersionRange{
			lower:          version,
			upper:          part,
			lowerInclusive: true,
		}, nil
	}
	return VersionRange{}, fmt.Errorf("invalid version range operation: '%s'", op)
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

func findSupportedPythonLibraries(want PythonLibrary, supported []SupportedPythonLibrary) []string {
	var links []string
	for _, lib := range supported {
		for dep, versionRange := range lib.VersionRange {
			if dep == want.Name && versionRange.matches(want.Version) {
				links = append(links, lib.Link)
			}
		}
	}
	return links
}
