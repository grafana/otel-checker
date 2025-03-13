package python

import (
	_ "embed"
	"fmt"
	"github.com/grafana/otel-checker/checks/sdk"
	"github.com/grafana/otel-checker/checks/utils"
	"os"
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

// readDependencies reads the Python dependencies from the requirements.txt file
func readDependencies(reporter *utils.ComponentReporter) []Library {
	path := "requirements.txt"
	if utils.FileExists(path) {
		return readRequirementsTxt(reporter, path)
	}
	// TODO: Add support for other package formats like poetry.lock
	return nil
}

// readRequirementsTxt reads and parses a requirements.txt file to extract dependencies
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

// parseRequirementsTxt parses the contents of a requirements.txt file and extracts libraries
func parseRequirementsTxt(reporter *utils.ComponentReporter, lines string) []Library {
	var deps []Library
	for _, line := range strings.Split(lines, "\n") {
		lib, ok := parseRequirementLine(line)
		if !ok {
			if line != "" {
				reporter.AddWarning(fmt.Sprintf("Could not parse line: %s", line))
			}
			continue
		}
		deps = append(deps, lib)
	}
	return deps
}

// parseRequirementLine parses a single line from a requirements.txt file
// Returns the library info and a boolean indicating success
func parseRequirementLine(line string) (Library, bool) {
	line = strings.TrimSpace(line)
	if line == "" {
		return Library{}, false
	}

	split := strings.Split(line, "==")
	if len(split) != 2 {
		return Library{}, false
	}

	return Library{
		Name:    strings.ToLower(strings.TrimSpace(split[0])),
		Version: strings.TrimSpace(split[1]),
	}, true
}

// outputSupportedLibraries reports which dependencies are supported by OpenTelemetry
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

// findSupportedLibraries finds OpenTelemetry instrumentation links for a given library
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

// supportedLibraries loads and parses the Python instrumentation libraries list from GitHub
func supportedLibraries() ([]SupportedLibrary, error) {
	// Load the README from GitHub that contains the supported libraries list
	readme, err := loadSupportedLibrariesReadme()
	if err != nil {
		return nil, err
	}

	// Parse the README to extract supported libraries
	return parseSupportedLibraries(readme)
}

// loadSupportedLibrariesReadme loads the README file from GitHub that contains the list of supported Python libraries
func loadSupportedLibrariesReadme() (string, error) {
	url := "https://raw.githubusercontent.com/open-telemetry/opentelemetry-python-contrib/refs/heads/main/instrumentation/README.md"
	bytes, err := sdk.LoadUrl(url)
	if err != nil {
		return "", fmt.Errorf("failed to load supported libraries: %w", err)
	}
	return string(bytes), nil
}

// parseSupportedLibraries parses the README content to extract supported libraries information
func parseSupportedLibraries(readme string) ([]SupportedLibrary, error) {
	var result []SupportedLibrary

	// Skip the header rows and start processing from row 4
	lines := strings.Split(readme, "\n")
	if len(lines) < 4 {
		return nil, fmt.Errorf("invalid README format: too few lines")
	}

	for _, line := range lines[3:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		library, err := parseSupportedLibraryLine(line)
		if err != nil {
			return nil, err
		}

		result = append(result, library)
	}

	return result, nil
}

// parseSupportedLibraryLine parses a single line from the README table of supported libraries
func parseSupportedLibraryLine(line string) (SupportedLibrary, error) {
	// Split by the table column separator
	columns := strings.Split(line, "|")
	if len(columns) < 3 {
		return SupportedLibrary{}, fmt.Errorf("invalid library line format: %s", line)
	}

	// Extract the library name from the Markdown link
	mdLink := strings.TrimSpace(columns[1])
	matches := linkRegex.FindStringSubmatch(mdLink)
	if len(matches) < 2 {
		return SupportedLibrary{}, fmt.Errorf("could not extract library name from: %s", mdLink)
	}

	name := matches[1]
	url := fmt.Sprintf("https://github.com/open-telemetry/opentelemetry-python-contrib/tree/main/instrumentation/opentelemetry-instrumentation-%s", name)

	// Parse version ranges
	versionRangesText := strings.TrimSpace(columns[2])
	ranges, err := versionRanges(versionRangesText)
	if err != nil {
		return SupportedLibrary{}, fmt.Errorf("error parsing version ranges for %s: %w", name, err)
	}

	return SupportedLibrary{
		Name:         name,
		Link:         url,
		VersionRange: ranges,
	}, nil
}

// parseVersionConstraint parses a single version constraint like "<1.0" or "< 1.0"
// and extracts the operator and version
func parseVersionConstraint(constraint string) (operator string, version string, err error) {
	// Check for combined operator and version (like "<1" instead of "< 1")
	for _, op := range []string{"<=", ">=", "~=", "<", ">"} {
		if strings.HasPrefix(constraint, op) {
			version := strings.TrimSpace(constraint[len(op):])
			if len(version) > 0 {
				return op, version, nil
			}
		}
	}

	// Handle space-separated format
	parts := strings.Split(constraint, " ")
	if len(parts) == 2 {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("invalid version constraint format: %s", constraint)
}

// handleVersionStatement processes a statement that contains version constraints
func handleVersionStatement(res map[string]sdk.VersionRange, name string, statement string) error {
	operator, version, err := parseVersionConstraint(statement)
	if err != nil {
		return err
	}

	return addVersionRange(res, name, operator, version)
}

func versionRanges(list string) (map[string]sdk.VersionRange, error) {
	res := map[string]sdk.VersionRange{}
	name := ""

	for _, s := range strings.Split(list, ",") {
		s = strings.TrimSpace(s)

		// Split by space to handle different formats
		parts := strings.Split(s, " ")

		// Determine the structure based on the number of parts
		switch {
		case len(parts) >= 3:
			// Format: "name operator version [extra]"
			name = parts[0]
			// Join the remaining parts back together
			statement := strings.Join(parts[1:], " ")
			operator, version, err := parseVersionConstraint(statement)
			if err != nil {
				return nil, fmt.Errorf("error parsing version %s: %v", list, err)
			}

			err = addVersionRange(res, name, operator, version)
			if err != nil {
				return nil, fmt.Errorf("error parsing version %s: %v", list, err)
			}

		case len(parts) == 2:
			// Format could be "name operator<version>" or "operator version"
			if name == "" {
				// This is the first part with the name
				name = parts[0]

				// Check if the second part has a combined operator and version
				if hasOperatorPrefix(parts[1]) {
					err := handleVersionStatement(res, name, parts[1])
					if err != nil {
						return nil, fmt.Errorf("error parsing version %s: %v", list, err)
					}
				} else {
					// Standard format: name operator version (split across statements)
					err := addVersionRange(res, name, parts[1], "")
					if err != nil {
						return nil, fmt.Errorf("error parsing version %s: %v", list, err)
					}
				}
			} else {
				// This is a continuation with just operator and version
				err := addVersionRange(res, name, parts[0], parts[1])
				if err != nil {
					return nil, fmt.Errorf("error parsing version %s: %v", list, err)
				}
			}

		case len(parts) == 1:
			// Could be a single term (just the library name) or a combined constraint without spaces
			if name == "" {
				// Just the library name
				name = parts[0]
				res[name] = sdk.VersionRange{} // no version constraint
			} else {
				// Check for constraints without spaces
				err := handleVersionStatement(res, name, parts[0])
				if err != nil {
					return nil, fmt.Errorf("error parsing version %s: %v", list, err)
				}
			}

		default:
			return nil, fmt.Errorf("error parsing version %s: invalid version range statement: %s", list, s)
		}
	}
	return res, nil
}

// hasOperatorPrefix checks if a string starts with a version operator
func hasOperatorPrefix(s string) bool {
	for _, op := range []string{"<=", ">=", "~=", "<", ">"} {
		if strings.HasPrefix(s, op) {
			return true
		}
	}
	return false
}

// addVersionRange adds a version range constraint to the result map
// It validates the version format, creates a new range, and merges it with existing constraints
func addVersionRange(res map[string]sdk.VersionRange, name string, op string, version string) error {
	// Handle empty version for cases where it might be in a separate statement
	if version == "" {
		return fmt.Errorf("empty version not supported")
	}

	// Validate version format
	if err := validateVersionFormat(version); err != nil {
		return fmt.Errorf("invalid version: %s", version)
	}

	// Create new version range based on operator
	r, err := createVersionRange(op, version)
	if err != nil {
		return err
	}

	// Merge with existing range if any
	old, exists := res[name]
	if exists {
		r = mergeVersionRanges(r, old)
	}

	res[name] = r
	return nil
}

// validateVersionFormat checks if the version string starts with a digit
func validateVersionFormat(version string) error {
	if len(version) == 0 {
		return fmt.Errorf("empty version")
	}
	_, err := strconv.Atoi(fmt.Sprintf("%c", version[0]))
	return err
}

// mergeVersionRanges combines two version ranges, preserving the most restrictive bounds
func mergeVersionRanges(r1 sdk.VersionRange, r2 sdk.VersionRange) sdk.VersionRange {
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

// createVersionRange creates a VersionRange based on the operator and version
func createVersionRange(op string, version string) (sdk.VersionRange, error) {
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
		upperBound, err := calculateTildeUpperBound(version)
		if err != nil {
			return sdk.VersionRange{}, err
		}
		return sdk.VersionRange{
			Lower:          version,
			Upper:          upperBound,
			LowerInclusive: true,
		}, nil
	}
	return sdk.VersionRange{}, fmt.Errorf("invalid version range operation: '%s'", op)
}

// calculateTildeUpperBound computes the upper bound for the ~= operator
// For ~=X.Y.Z it creates a range from X.Y.Z (inclusive) to X.(Y+1) (exclusive)
func calculateTildeUpperBound(version string) (string, error) {
	split := strings.Split(version, ".")
	// Need at least two parts for a meaningful upper bound
	if len(split) < 2 {
		return "", fmt.Errorf("version needs at least major.minor format for tilde operator: %s", version)
	}

	// Remove the last part (patch level) and increment the next-to-last part
	split = split[:len(split)-1]
	last, err := strconv.Atoi(split[len(split)-1])
	if err != nil {
		return "", err
	}
	split[len(split)-1] = strconv.Itoa(last + 1)
	return strings.Join(split, "."), nil
}
