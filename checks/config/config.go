// Package config parses the OpenTelemetry declarative configuration YAML
// (https://opentelemetry.io/docs/specs/otel/configuration/) and exposes helpers
// for extracting the fields the checker cares about.
package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"go.yaml.in/yaml/v3"
)

// File is the top-level structure of an OTel declarative configuration
// file.
type File struct {
	FileFormat     string          `yaml:"file_format"`
	Resource       *Resource       `yaml:"resource,omitempty"`
	TracerProvider *TracerProvider `yaml:"tracer_provider,omitempty"`
	MeterProvider  *MeterProvider  `yaml:"meter_provider,omitempty"`
	LoggerProvider *LoggerProvider `yaml:"logger_provider,omitempty"`
}

type Resource struct {
	Attributes     []ResourceAttribute `yaml:"attributes,omitempty"`
	AttributesList string              `yaml:"attributes_list,omitempty"`
}

type ResourceAttribute struct {
	Name  string `yaml:"name"`
	Value any    `yaml:"value"`
	Type  string `yaml:"type,omitempty"`
}

type TracerProvider struct {
	Processors []SpanProcessor `yaml:"processors,omitempty"`
}

type SpanProcessor struct {
	Batch  *WithExporter `yaml:"batch,omitempty"`
	Simple *WithExporter `yaml:"simple,omitempty"`
}

type MeterProvider struct {
	Readers []MetricReader `yaml:"readers,omitempty"`
}

type MetricReader struct {
	Periodic *WithExporter `yaml:"periodic,omitempty"`
	Pull     *WithExporter `yaml:"pull,omitempty"`
}

type LoggerProvider struct {
	Processors []LogProcessor `yaml:"processors,omitempty"`
}

type LogProcessor struct {
	Batch  *WithExporter `yaml:"batch,omitempty"`
	Simple *WithExporter `yaml:"simple,omitempty"`
}

type WithExporter struct {
	Exporter Exporter `yaml:"exporter"`
}

type Exporter struct {
	OTLPHTTP *OTLPHTTPExporter `yaml:"otlp_http,omitempty"`
}

type OTLPHTTPExporter struct {
	Endpoint string `yaml:"endpoint,omitempty"`
}

// Load reads and parses the YAML at path. Env-var substitutions inside
// string values are left as-is, callers use ExpandEnv when they need
// the resolved value.
func Load(path string) (*File, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	var f File
	if err := yaml.Unmarshal(raw, &f); err != nil {
		return nil, fmt.Errorf("could not parse config file: %w", err)
	}
	return &f, nil
}

func (f *File) SignalEndpoints() map[string][]string {
	out := map[string][]string{"traces": nil, "metrics": nil, "logs": nil}
	if f == nil {
		return out
	}

	if f.TracerProvider != nil {
		for _, p := range f.TracerProvider.Processors {
			if e := endpointOf(p.Batch); e != "" {
				out["traces"] = append(out["traces"], e)
			}
			if e := endpointOf(p.Simple); e != "" {
				out["traces"] = append(out["traces"], e)
			}
		}
	}
	if f.MeterProvider != nil {
		for _, r := range f.MeterProvider.Readers {
			if e := endpointOf(r.Periodic); e != "" {
				out["metrics"] = append(out["metrics"], e)
			}
			if e := endpointOf(r.Pull); e != "" {
				out["metrics"] = append(out["metrics"], e)
			}
		}
	}
	if f.LoggerProvider != nil {
		for _, p := range f.LoggerProvider.Processors {
			if e := endpointOf(p.Batch); e != "" {
				out["logs"] = append(out["logs"], e)
			}
			if e := endpointOf(p.Simple); e != "" {
				out["logs"] = append(out["logs"], e)
			}
		}
	}
	return out
}

func endpointOf(w *WithExporter) string {
	if w == nil || w.Exporter.OTLPHTTP == nil {
		return ""
	}
	return w.Exporter.OTLPHTTP.Endpoint
}

func (f *File) ResourceAttributes() map[string]string {
	out := map[string]string{}
	if f == nil || f.Resource == nil {
		return out
	}

	// attributes_list (lower priority) parses first.
	if raw := f.Resource.AttributesList; raw != "" {
		expanded, _ := ExpandEnv(raw)
		for pair := range strings.SplitSeq(expanded, ",") {
			key, value, ok := strings.Cut(pair, "=")
			if !ok {
				continue
			}
			key = strings.TrimSpace(key)
			value = strings.TrimSpace(value)
			if key == "" || value == "" {
				continue
			}
			out[key] = value
		}
	}

	// attributes (higher priority) overrides.
	for _, attr := range f.Resource.Attributes {
		if attr.Value == nil {
			continue
		}
		raw := fmt.Sprintf("%v", attr.Value)
		expanded, _ := ExpandEnv(raw)
		expanded = strings.TrimSpace(expanded)
		if expanded == "" {
			continue
		}
		out[attr.Name] = expanded
	}
	return out
}

// envVarPattern matches either the `$$` escape sequence or a full
// substitution reference. The alternation is ordered so `$$` wins at
// positions where the two could overlap. Submatch indices:
//
//	1 - optional "env:" prefix
//	2 - variable name
//	3 - optional ":-default" content
//
// When submatch 2 is empty, the overall match is the `$$` escape.
var envVarPattern = regexp.MustCompile(`\$\$|\$\{(env:)?([A-Za-z_][A-Za-z0-9_]*)(?::-([^}]*))?\}`)

// ExpandEnv resolves every ${VAR}, ${env:VAR}, and ${VAR:-default}
// occurrence in s using os.LookupEnv. When a referenced variable is
// unset (or set to the empty string) and no default is provided, the
// placeholder is left intact and the variable name is appended to
// unresolved so the caller can surface the failure.
func ExpandEnv(s string) (result string, unresolved []string) {
	result = envVarPattern.ReplaceAllStringFunc(s, func(match string) string {
		if match == "$$" {
			return "$"
		}
		sub := envVarPattern.FindStringSubmatch(match)
		name := sub[2]
		hasDefault := strings.Contains(match, ":-")
		defaultVal := sub[3]

		value, present := os.LookupEnv(name)
		if present && value != "" {
			return value
		}
		if hasDefault {
			return defaultVal
		}
		unresolved = append(unresolved, name)
		return match
	})
	return result, unresolved
}
