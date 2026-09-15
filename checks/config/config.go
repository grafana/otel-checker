// Package config parses the OpenTelemetry declarative configuration YAML
// (https://opentelemetry.io/docs/specs/otel/configuration/) and exposes helpers
// for extracting the fields the checker cares about.
//
// The type declarations for the config model live in config.gen.go and are
// generated from the upstream JSON schema via scripts/generate_config_schema.sh.
// Only helpers (Load, SignalEndpoints, ResourceAttributes, ExpandEnv) are
// hand-written here.
package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"go.yaml.in/yaml/v3"
)

// Load reads, env-var-expands, and parses the YAML at path.
//
// Substitution happens BEFORE unmarshal because the generated schema
// types many fields as int/bool (schedule_delay, disabled, …) and yaml
// can't unmarshal a `${...}` string into a typed field.
//
// Unresolved variables (no env value and no `:-default`) are stripped
// to the empty string so yaml sees a null value for the field.
func Load(path string) (*File, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}

	expanded, _ := ExpandEnv(string(raw))
	// Any ${...} still present is an unresolved reference with no
	// default; drop them so yaml can null-out the field rather than
	// failing type coercion.
	expanded = envVarPattern.ReplaceAllString(expanded, "")

	var f File
	if err := yaml.Unmarshal([]byte(expanded), &f); err != nil {
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
			if e := traceEndpoint(p.Batch); e != "" {
				out["traces"] = append(out["traces"], e)
			}
			if p.Simple != nil {
				if e := otlpHTTPEndpoint(p.Simple.Exporter.OTLPHTTP); e != "" {
					out["traces"] = append(out["traces"], e)
				}
			}
		}
	}
	if f.MeterProvider != nil {
		for _, r := range f.MeterProvider.Readers {
			if r.Periodic != nil {
				if e := metricEndpoint(r.Periodic.Exporter.OTLPHTTP); e != "" {
					out["metrics"] = append(out["metrics"], e)
				}
			}
		}
	}
	if f.LoggerProvider != nil {
		for _, p := range f.LoggerProvider.Processors {
			if e := logEndpoint(p.Batch); e != "" {
				out["logs"] = append(out["logs"], e)
			}
			if p.Simple != nil {
				if e := otlpHTTPEndpoint(p.Simple.Exporter.OTLPHTTP); e != "" {
					out["logs"] = append(out["logs"], e)
				}
			}
		}
	}
	return out
}

func traceEndpoint(b *BatchSpanProcessor) string {
	if b == nil {
		return ""
	}
	return otlpHTTPEndpoint(b.Exporter.OTLPHTTP)
}

func logEndpoint(b *BatchLogRecordProcessor) string {
	if b == nil {
		return ""
	}
	return otlpHTTPEndpoint(b.Exporter.OTLPHTTP)
}

func otlpHTTPEndpoint(exp *OTLPHTTPExporter) string {
	if exp == nil || exp.Endpoint == nil {
		return ""
	}
	return *exp.Endpoint
}

func metricEndpoint(exp *OTLPHTTPMetricExporter) string {
	if exp == nil || exp.Endpoint == nil {
		return ""
	}
	return *exp.Endpoint
}

func (f *File) ResourceAttributes() map[string]string {
	out := map[string]string{}
	if f == nil || f.Resource == nil {
		return out
	}

	// attributes_list (lower priority) parses first.
	if list := f.Resource.AttributesList; list != nil && *list != "" {
		expanded, _ := ExpandEnv(*list)
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
