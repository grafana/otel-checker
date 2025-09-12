package java

import (
	"encoding/base64"
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/grafana/otel-checker/checks/sdk"
	"github.com/grafana/otel-checker/checks/sdk/supported"
	"github.com/grafana/otel-checker/checks/utils"
	"go.yaml.in/yaml/v3"
	"golang.org/x/mod/semver"
)

var javaVersionRegex = regexp.MustCompile(`Java (\d)+\+`)

type Library struct {
	Group    string    `json:"groupId"`
	Artifact string    `json:"artifactId"`
	Version  string    `json:"version"`
	Children []Library `json:"children"`
}

func (l *Library) String() string {
	return fmt.Sprintf("%s:%s:%s", l.Group, l.Artifact, l.Version)
}

func reportSupportedInstrumentations(reporter *utils.ComponentReporter, debug bool, explorer bool, instrumentationType supported.InstrumentationType, javaVersion int) {
	s, err := supportedLibraries()
	if err != nil {
		reporter.AddError(fmt.Sprintf("Error reading supported libraries: %v", err))
		return
	}

	deps := readDependencies(reporter)
	// Initialize global collection for this run
	globalInstrumentations = make(map[string]bool)
	isFirstCall = true // Ensure it's reset for each top-level call to reportSupportedInstrumentations

	outputSupportedLibraries(deps, s, reporter, debug, explorer, instrumentationType, javaVersion) // This populates globalInstrumentations

	// Output final summary using the global collection
	var globalList []string
	for name := range globalInstrumentations {
		globalList = append(globalList, name)
	}
	slices.Sort(globalList)

	if len(globalList) > 0 {
		reporter.AddSuccessfulCheck(fmt.Sprintf("🎯 Summary: Found %d unique OpenTelemetry instrumentations: %s",
			len(globalList), strings.Join(globalList, ", ")))

		if explorer {
			explorerLink := generateInstrumentationExplorerLink(globalList)
			if explorerLink != "" {
				reporter.AddSuccessfulCheck(fmt.Sprintf("🔗 Explore these instrumentations: %s", explorerLink))
			}
		}
	} else {
		reporter.AddSuccessfulCheck("❌ No instrumentations found")
	}
}

func readDependencies(reporter *utils.ComponentReporter) []Library {
	if utils.FileExists("pom.xml") {
		return checkMaven(reporter)
	}
	for _, file := range gradleFiles {
		if utils.FileExists(file) {
			return checkGradle(file, reporter)
		}
	}
	return nil
}

func searchWrapper(base string, wrapper string) string {
	tool := getWrapper(wrapper, []string{"."})
	if tool == "" {
		return base
	}
	return tool
}

func getWrapper(wrapper string, level []string) string {
	if len(level) > 10 {
		return ""
	}
	p := filepath.Join(filepath.Join(level...), wrapper)
	if utils.FileExists(p) {
		// the . is needed to run the wrapper in the current directory
		return fmt.Sprintf(".%c%s", filepath.Separator, p)
	}
	return getWrapper(wrapper, append(level, ".."))
}

// extractInstrumentationNames extracts instrumentation names from GitHub URLs
// e.g. "https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/instrumentation/kafka/kafka-clients/kafka-clients-0.11/javaagent" -> "kafka-clients-0.11"
func extractInstrumentationNames(links []string) []string {
	var names []string
	seen := make(map[string]bool)
	
	for _, link := range links {
		// Extract the instrumentation name from the URL path
		// The instrumentation name is the second-to-last directory before /javaagent
		parts := strings.Split(link, "/")
		if len(parts) >= 9 && parts[7] == "instrumentation" {
			// Find the last part that is "javaagent" and get the part before it
			for i := len(parts) - 1; i >= 8; i-- {
				if parts[i] == "javaagent" && i > 8 {
					name := parts[i-1]
					if !seen[name] {
						names = append(names, name)
						seen[name] = true
					}
					break
				}
			}
		}
	}
	
	slices.Sort(names)
	return names
}

// generateInstrumentationExplorerLink creates a link to the instrumentation explorer with base64 encoded instrumentations
func generateInstrumentationExplorerLink(instrumentations []string) string {
	if len(instrumentations) == 0 {
		return ""
	}
	
	// Join instrumentations without spaces
	instrumentationsList := strings.Join(instrumentations, ",")
	
	// Base64 encode the list
	encoded := base64.StdEncoding.EncodeToString([]byte(instrumentationsList))
	
	// Create the link with a reasonable default version (latest stable)
	return fmt.Sprintf("https://jaydeluca.github.io/instrumentation-explorer/analyze?instrumentations=%s&version=2.19", encoded)
}

// Global variable to collect all instrumentations across calls
var globalInstrumentations = make(map[string]bool)
var isFirstCall = true

func outputSupportedLibraries(deps []Library, supported supported.SupportedModules, reporter *utils.ComponentReporter, debug bool, explorer bool, instrumentationType supported.InstrumentationType, javaVersion int) []string {
	var allInstrumentations []string
	seen := make(map[string]bool)
	
	// Reset global collection on first call
	if isFirstCall {
		globalInstrumentations = make(map[string]bool)
		isFirstCall = false
	}
	
	for _, dep := range deps {
		links := findSupportedLibraries(dep, supported, instrumentationType, javaVersion, reporter)
		if len(links) > 0 {
			instrumentationNames := extractInstrumentationNames(links)
			reporter.AddSuccessfulCheck(
				fmt.Sprintf("Found supported library: %s:%s:%s at %s (instrumentations: %s)",
					dep.Group, dep.Artifact, dep.Version, strings.Join(links, ", "), strings.Join(instrumentationNames, ", ")))
			// Collect unique instrumentations
			for _, name := range instrumentationNames {
				if !seen[name] {
					allInstrumentations = append(allInstrumentations, name)
					seen[name] = true
				}
				// Also add to global collection
				globalInstrumentations[name] = true
			}
		} else if debug {
			reporter.AddWarning(fmt.Sprintf("Found unsupported library: %s:%s:%s", dep.Group, dep.Artifact, dep.Version))
		}
		
		// Collect instrumentations from children
		childInstrumentations := outputSupportedLibraries(dep.Children, supported, reporter, false, explorer, instrumentationType, javaVersion)
		for _, name := range childInstrumentations {
			if !seen[name] {
				allInstrumentations = append(allInstrumentations, name)
				seen[name] = true
			}
		}
	}
	
	slices.Sort(allInstrumentations)
	return allInstrumentations
}

func findSupportedLibraries(library Library, supported supported.SupportedModules, instrumentationType supported.InstrumentationType, javaVersion int, reporter *utils.ComponentReporter) []string {
	var links []string
	for moduleName, instrumentations := range supported {
		for _, instrumentation := range instrumentations {
			for _, version := range instrumentation.TargetVersions[instrumentationType] {
				if matchVersion(moduleName, version, library, javaVersion, reporter) {
					l := fmt.Sprintf("https://github.com/open-telemetry/opentelemetry-java-instrumentation/tree/main/%s/%s",
						instrumentation.SrcPath, instrumentationType)
					if !slices.Contains(links, l) {
						links = append(links, l)
					}
				}
			}
		}
	}
	slices.Sort(links)
	return links
}

func matchVersion(moduleName string, version string, library Library, javaVersion int, reporter *utils.ComponentReporter) bool {
	javaVersionMatch := javaVersionRegex.FindStringSubmatch(version)

	if javaVersionMatch != nil {
		// e.g. Java 8+
		wantJavaVersion, err := strconv.Atoi(javaVersionMatch[1])
		if err != nil {
			reporter.AddError(fmt.Sprintf("Error parsing Java version %s: %v", version, err))
			return false
		}
		return wantJavaVersion <= javaVersion
	}

	// e.g. com.amazonaws:aws-lambda-java-core:[1.0.0,)
	split := strings.Split(version, ":")
	if len(split) != 3 {
		reporter.AddInternalError(fmt.Sprintf("Invalid java version for module %s: %s", moduleName, version))
		return false
	}
	versionRange, err := sdk.ParseVersionRange(split[2])
	if err != nil {
		reporter.AddInternalError(fmt.Sprintf("Error parsing version range for module %s: %s", moduleName, version))
		return false
	}

	if library.Group == split[0] && library.Artifact == split[1] {
		v := sdk.FixVersion(library.Version)
		if semver.IsValid(v) {
			// ignore invalid versions from applications
			return versionRange.Matches(v)
		}
	}
	return false
}

func supportedLibraries() (supported.SupportedModules, error) {
	bytes, err := sdk.LoadUrl("https://raw.githubusercontent.com/open-telemetry/opentelemetry-java-instrumentation/refs/heads/main/docs/instrumentation-list.yaml")
	if err != nil {
		return nil, err
	}
	return LoadSupportedJavaLibraries(bytes)
}

// SupportedJavaModules is a struct that holds the supported Java libraries
type SupportedJavaModules struct {
	Libraries supported.SupportedModules `json:"libraries"`
}

// LoadSupportedJavaLibraries loads supported libraries from a YAML file
func LoadSupportedJavaLibraries(data []byte) (supported.SupportedModules, error) {
	modules := SupportedJavaModules{}
	err := yaml.Unmarshal(data, &modules)
	if err != nil {
		return nil, err
	}
	return modules.Libraries, nil
}
