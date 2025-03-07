package sdk

import (
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"otel-checker/checks/utils"
	"strings"

	"golang.org/x/mod/semver"
)

type VersionRange struct {
	Lower          string
	LowerInclusive bool
	Upper          string
	UpperInclusive bool
}

var ignoreVersionSuffixes = []string{
	".Final",
	".RELEASE",
	".GA",
	"-M1",
}

// ParseVersionRange parses versions of the form "[5.0,)", "[2.0,4.0)".
func ParseVersionRange(version string) (VersionRange, error) {
	split := strings.Split(version, ",")
	if len(split) == 1 {
		return VersionRange{
			Lower:          version,
			LowerInclusive: true,
			Upper:          version,
			UpperInclusive: true,
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
	for _, suffix := range ignoreVersionSuffixes {
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

func (r *VersionRange) Matches(version string) bool {
	return checkBound(r.Lower, r.LowerInclusive, version, -1) && checkBound(r.Upper, r.UpperInclusive, version, 1)
}

func checkBound(bound string, inclusive bool, version string, sgn int) bool {
	if bound == "" {
		return true
	}
	cmp := semver.Compare(FixVersion(bound), FixVersion(version))
	if cmp == 0 {
		return inclusive
	}
	return cmp == sgn
}

func RunCommand(reporter *utils.ComponentReporter, cmd *exec.Cmd) string {
	println("Running command:", cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		reporter.AddError(fmt.Sprintf("Error running %s:\n%v\n%s", cmd.String(), err, output))
		return ""
	}
	return string(output)
}

func LoadUrl(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching instrumentation list: %v", err)
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	return bytes, nil
}
