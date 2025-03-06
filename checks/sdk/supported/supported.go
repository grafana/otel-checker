package supported

import (
	"gopkg.in/yaml.v3"
)

// Library represents a library with its name and version
type Library struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// InstrumentationType represents the type of instrumentation
type InstrumentationType string

const (
	TypeJavaagent InstrumentationType = "JAVAAGENT"
	TypeLibrary   InstrumentationType = "LIBRARY"
)

// Instrumentation represents a single instrumentation with its metadata
type Instrumentation struct {
	Name           string                           `yaml:"name"`
	SrcPath        string                           `yaml:"srcPath"`
	Link           string                           `yaml:"link,omitempty"`
	TargetVersions map[InstrumentationType][]string `yaml:"target_versions"`
}

// SupportedModule represents a module containing instrumentations
type SupportedModule struct {
	Instrumentations []Instrumentation `yaml:"instrumentations"`
}

// SupportedModules represents a map of module names to their supported modules
type SupportedModules map[string]SupportedModule

// LoadSupportedLibraries loads supported libraries from a YAML file
func LoadSupportedLibraries(data []byte) (SupportedModules, error) {
	modules := SupportedModules{}
	err := yaml.Unmarshal(data, &modules)
	if err != nil {
		return nil, err
	}
	delete(modules, "internal")
	return modules, nil
}
