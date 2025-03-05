package sdk

import (
	"fmt"
	"golang.org/x/mod/semver"
	"strings"
)

type VersionRange struct {
	lower          string
	lowerInclusive bool
	upper          string
	upperInclusive bool
}

var suffixes = []string{
	".Final",
	".RELEASE",
	".GA",
	"-M1",
}

// ParseVersionRange parses versions of the form "[5.0,)", "[2.0,4.0)".
func ParseVersionRange(version string) (VersionRange, error) {
	split := strings.Split(version, ",")
	if len(split) == 1 {
		v := FixVersion(version)
		return VersionRange{
			lower:          v,
			lowerInclusive: true,
			upper:          v,
			upperInclusive: true,
		}, nil
	}

	if len(split) != 2 {
		return VersionRange{}, fmt.Errorf("version has more than one comma: %s", version)
	}
	lowerInclusive := false
	if strings.HasPrefix(version, "[") {
		lowerInclusive = true
	} else if strings.HasPrefix(version, "(") {
		lowerInclusive = false
	} else {
		return VersionRange{}, fmt.Errorf("version does not start with '[' or '(': %s", version)
	}

	upperInclusive := false
	if strings.HasSuffix(version, "]") {
		upperInclusive = true
	} else if strings.HasSuffix(version, ")") {
		upperInclusive = false
	} else {
		return VersionRange{}, fmt.Errorf("version does not end with ']' or ')': %s", version)
	}

	l := FixVersion(strings.TrimLeft(split[0], "[("))
	if l != "" {
		if !semver.IsValid(l) {
			return VersionRange{}, fmt.Errorf("invalid semver: '%s'", l)
		}
	}
	u := FixVersion(strings.TrimRight(split[1], ")]"))
	if u != "" {
		if !semver.IsValid(u) {
			return VersionRange{}, fmt.Errorf("invalid semver: '%s'", u)
		}
	}
	return VersionRange{
		l,
		lowerInclusive,
		u,
		upperInclusive,
	}, nil
}

func FixVersion(version string) string {
	version = strings.TrimSpace(version)
	if version == "" {
		return version
	}
	for _, suffix := range suffixes {
		version = strings.TrimSuffix(version, suffix)
	}

	if !strings.HasPrefix(version, "v") {
		version = fmt.Sprintf("v%s", version)
	}

	split := strings.Split(version, ".")
	if len(split) == 1 {
		version = version + ".0"
	}
	if len(split) == 4 {
		// pretend build version
		version = fmt.Sprintf("%s.%s.%s+%s", split[0], split[1], split[2], split[3])
	}

	return version
}

func (r *VersionRange) matches(version string) bool {
	return checkBound(r.lower, r.lowerInclusive, version, -1) && checkBound(r.upper, r.upperInclusive, version, 1)
}

func checkBound(bound string, inclusive bool, version string, sgn int) bool {
	if bound == "" {
		return true
	}
	cmp := semver.Compare(bound, version)
	if cmp == 0 {
		return inclusive
	}
	return cmp == sgn
}

func CheckSDKSetup(messages *map[string][]string, language string, autoInstrumentation bool, packageJsonPath string, instrumentationFile string, debug bool) {
	switch language {
	case "dotnet":
		CheckDotNetSetup(messages, autoInstrumentation)
	case "go":
		CheckGoSetup(messages, autoInstrumentation)
	case "java":
		CheckJavaSetup(messages, autoInstrumentation, debug)
	case "js":
		CheckJSSetup(messages, autoInstrumentation, packageJsonPath, instrumentationFile)
	case "python":
		CheckPythonSetup(messages, autoInstrumentation)
	case "ruby":
		CheckRubySetup(messages, autoInstrumentation)
	}
}
