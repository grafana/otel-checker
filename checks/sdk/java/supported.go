package java

import (
	"context"
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

func reportSupportedInstrumentations(ctx context.Context, reporter *utils.ComponentReporter, debug bool, instrumentationType supported.InstrumentationType, javaVersion int) {
	s, err := supportedLibraries(ctx)
	if err != nil {
		reporter.AddErrorWithFix("java.supported-libs.fetch-failed",
			fmt.Sprintf("Error reading supported libraries: %v", err))
	}

	deps := readDependencies(ctx, reporter)
	outputSupportedLibraries(deps, s, reporter, debug, instrumentationType, javaVersion)
}

func readDependencies(ctx context.Context, reporter *utils.ComponentReporter) []Library {
	if utils.FileExists("pom.xml") {
		return checkMaven(ctx, reporter)
	}
	for _, file := range gradleFiles {
		if utils.FileExists(file) {
			return checkGradle(ctx, file, reporter)
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

func outputSupportedLibraries(deps []Library, supportedModules supported.SupportedModules, reporter *utils.ComponentReporter, debug bool, instrumentationType supported.InstrumentationType, javaVersion int) {
	for _, dep := range deps {
		links := findSupportedLibraries(dep, supportedModules, instrumentationType, javaVersion, reporter)
		if len(links) > 0 {
			reporter.AddSuccessfulCheck(
				fmt.Sprintf("Found supported library: %s:%s:%s at %s",
					dep.Group, dep.Artifact, dep.Version, strings.Join(links, ", ")))
		} else if debug {
			reporter.AddWarningWithFix("java.library.unsupported",
				fmt.Sprintf("Found unsupported library: %s:%s:%s", dep.Group, dep.Artifact, dep.Version))
		}
		outputSupportedLibraries(dep.Children, supportedModules, reporter, false, instrumentationType, 0)
	}
}

func findSupportedLibraries(library Library, supportedModules supported.SupportedModules, instrumentationType supported.InstrumentationType, javaVersion int, reporter *utils.ComponentReporter) []string {
	var links []string
	for moduleName, instrumentations := range supportedModules {
		for _, instrumentation := range instrumentations {
			var versions []string
			if instrumentationType == supported.TypeJavaagent {
				versions = instrumentation.Versions
			} else if instrumentationType == supported.TypeLibrary && instrumentation.SupportsManualInstrumentation {
				// Manual instrumentation supports the same versions
				versions = instrumentation.Versions
			}

			for _, version := range versions {
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
			reporter.AddErrorWithFix("java.version.parse-error",
				fmt.Sprintf("Error parsing Java version %s: %v", version, err))
			return false
		}
		return wantJavaVersion <= javaVersion
	}

	// e.g. com.amazonaws:aws-lambda-java-core:[1.0.0,)
	split := strings.Split(version, ":")
	if len(split) != 3 {
		reporter.AddInternalErrorWithFix("internal.java.semver",
			fmt.Sprintf("Invalid java version for module %s: %s", moduleName, version))
		return false
	}
	versionRange, err := sdk.ParseVersionRange(split[2])
	if err != nil {
		reporter.AddInternalErrorWithFix("internal.java.version-range",
			fmt.Sprintf("Error parsing version range for module %s: %s", moduleName, version))
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

func supportedLibraries(ctx context.Context) (supported.SupportedModules, error) {
	bytes, err := sdk.LoadUrl(ctx, "https://raw.githubusercontent.com/open-telemetry/opentelemetry-java-instrumentation/refs/heads/main/docs/instrumentation-list.yaml")
	if err != nil {
		return nil, err
	}
	return LoadSupportedJavaLibraries(bytes)
}

// JavaInstrumentation matches the upstream Java YAML format
type JavaInstrumentation struct {
	Name                     string   `yaml:"name"`
	DisplayName              string   `yaml:"display_name"`
	Description              string   `yaml:"description"`
	SrcPath                  string   `yaml:"source_path"`
	Link                     string   `yaml:"link,omitempty"`
	LibraryLink              string   `yaml:"library_link,omitempty"`
	JavavagentTargetVersions []string `yaml:"javaagent_target_versions"`
	HasStandaloneLibrary     bool     `yaml:"has_standalone_library"`
	HasJavaagent             bool     `yaml:"has_javaagent"`
}

type supportedJavaModulesMap struct {
	Libraries map[string][]JavaInstrumentation `yaml:"libraries"`
}

type supportedJavaModulesList struct {
	Libraries []JavaInstrumentation `yaml:"libraries"`
}

// LoadSupportedJavaLibraries loads supported libraries from a YAML file and maps to generic format
func LoadSupportedJavaLibraries(data []byte) (supported.SupportedModules, error) {
	javaModulesMap := supportedJavaModulesMap{}
	if err := yaml.Unmarshal(data, &javaModulesMap); err == nil && len(javaModulesMap.Libraries) > 0 {
		return mapSupportedJavaModules(javaModulesMap.Libraries), nil
	}

	javaModulesList := supportedJavaModulesList{}
	if err := yaml.Unmarshal(data, &javaModulesList); err != nil {
		return nil, err
	}

	grouped := make(map[string][]JavaInstrumentation)
	for _, library := range javaModulesList.Libraries {
		grouped[library.Name] = append(grouped[library.Name], library)
	}
	return mapSupportedJavaModules(grouped), nil
}

func mapSupportedJavaModules(modules map[string][]JavaInstrumentation) supported.SupportedModules {
	result := make(supported.SupportedModules)
	for moduleName, javaInstrumentations := range modules {
		instrumentations := make([]supported.Instrumentation, 0, len(javaInstrumentations))
		for _, javaInst := range javaInstrumentations {
			versions := javaInst.JavavagentTargetVersions
			if !javaInst.HasJavaagent && len(versions) == 0 {
				continue
			}

			link := javaInst.Link
			if link == "" {
				link = javaInst.LibraryLink
			}

			instrumentations = append(instrumentations, supported.Instrumentation{
				Name:                          javaInst.Name,
				Description:                   javaInst.Description,
				SrcPath:                       javaInst.SrcPath,
				Link:                          link,
				Versions:                      versions,
				SupportsManualInstrumentation: javaInst.HasStandaloneLibrary,
			})
		}
		if len(instrumentations) > 0 {
			result[moduleName] = instrumentations
		}
	}
	return result
}
