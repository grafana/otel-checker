package collector

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/grafana/otel-checker/checks/utils"

	"go.yaml.in/yaml/v3"
)

func CheckCollectorSetup(reporter *utils.ComponentReporter, language string, configPath string) {
	checkCollectorConfig(reporter, configPath)
}

type configFile struct {
	Receivers map[string]receiverConfig `yaml:"receivers"`
	Exporters map[string]exporterConfig `yaml:"exporters"`
	Service   struct {
		Pipelines struct {
			Traces struct {
				Receivers  []string `yaml:"receivers"`
				Processors []string `yaml:"processors"`
				Exporters  []string `yaml:"exporters"`
			} `yaml:"traces"`
			Logs struct {
				Receivers  []string `yaml:"receivers"`
				Processors []string `yaml:"processors"`
				Exporters  []string `yaml:"exporters"`
			} `yaml:"logs"`
			Metrics struct {
				Receivers  []string `yaml:"receivers"`
				Processors []string `yaml:"processors"`
				Exporters  []string `yaml:"exporters"`
			} `yaml:"metrics"`
		} `yaml:"pipelines"`
	} `yaml:"service"`
}

type receiverConfig struct {
	Protocols map[string]any `yaml:"protocols"`
}

type exporterConfig struct {
	Endpoint string                 `yaml:"endpoint"`
	Auth     map[string]interface{} `yaml:"auth"`
}

var otlpHTTPGrafanaEndpointPattern = regexp.MustCompile(`https://.+\.grafana\.net/otlp`)

// checkCollectorConfig reads the Collector configuration file and runs the
// downstream checks against it. configPath, when non-empty, must be the full
// path to the file the user wants checked (any filename). When configPath is
// empty, the checker falls back to config.yaml then config.yml in the current
// working directory.
func checkCollectorConfig(reporter *utils.ComponentReporter, configPath string) {
	var candidates []string
	if configPath != "" {
		candidates = []string{configPath}
	} else {
		candidates = []string{"config.yaml", "config.yml"}
	}
	var (
		filePath string
		yamlFile []byte
	)
	for _, p := range candidates {
		b, err := os.ReadFile(p)
		if err == nil {
			filePath = p
			yamlFile = b
			break
		}
	}
	if yamlFile == nil {
		reporter.AddErrorWithExplain("collector.config.unreadable",
			fmt.Sprintf("Could not find a Collector config file. Tried: %s", strings.Join(candidates, ", ")))
		return
	}
	var c configFile
	if err := yaml.Unmarshal(yamlFile, &c); err != nil {
		reporter.AddErrorWithExplain("collector.config.parse-error",
			fmt.Sprintf("Could not parse file %s: %s", filePath, err))
		return
	}

	checkOTLPReceiverHTTPProtocol(reporter, c.Receivers)

	checkOTLPHTTPExporterEndpoint(reporter, c.Exporters)

	// Traces
	if containsOTLPHTTPExporter(c.Service.Pipelines.Traces.Exporters) {
		reporter.AddSuccessfulCheck("Value of service > pipelines > traces > exporters on config.yaml contains otlphttp")
	} else {
		reporter.AddWarningWithExplain("collector.pipelines.traces-otlphttp-missing",
			"Value of service > pipelines > traces > exporters on config.yaml does not contain otlphttp")
	}
	if containsOTLPReceiver(c.Service.Pipelines.Traces.Receivers) {
		reporter.AddSuccessfulCheck("Value of service > pipelines > traces > receivers on config.yaml contains otlp")
	} else {
		reporter.AddWarningWithExplain("collector.pipelines.traces-otlp-missing",
			"Value of service > pipelines > traces > receivers on config.yaml does not contain otlp")
	}

	// Logs
	if containsOTLPHTTPExporter(c.Service.Pipelines.Logs.Exporters) {
		reporter.AddSuccessfulCheck("Value of service > pipelines > logs > exporters on config.yaml contains otlphttp")
	} else {
		reporter.AddWarningWithExplain("collector.pipelines.logs-otlphttp-missing",
			"Value of service > pipelines > logs > exporters on config.yaml does not contain otlphttp")
	}
	if containsOTLPReceiver(c.Service.Pipelines.Logs.Receivers) {
		reporter.AddSuccessfulCheck("Value of service > pipelines > logs > receivers on config.yaml contains otlp")
	} else {
		reporter.AddWarningWithExplain("collector.pipelines.logs-otlp-missing",
			"Value of service > pipelines > logs > receivers on config.yaml does not contain otlp")
	}

	// Metrics
	if containsOTLPHTTPExporter(c.Service.Pipelines.Metrics.Exporters) {
		reporter.AddSuccessfulCheck("Value of service > pipelines > metrics > exporters on config.yaml contains otlphttp")
	} else {
		reporter.AddWarningWithExplain("collector.pipelines.metrics-otlphttp-missing",
			"Value of service > pipelines > metrics > exporters on config.yaml does not contain otlphttp")
	}
	if containsOTLPReceiver(c.Service.Pipelines.Metrics.Receivers) {
		reporter.AddSuccessfulCheck("Value of service > pipelines > metrics > receivers on config.yaml contains otlp")
	} else {
		reporter.AddWarningWithExplain("collector.pipelines.metrics-otlp-missing",
			"Value of service > pipelines > metrics > receivers on config.yaml does not contain otlp")
	}
}

func checkOTLPReceiverHTTPProtocol(reporter *utils.ComponentReporter, receivers map[string]receiverConfig) {
	for _, id := range componentIDsWithType(receivers, "otlp") {
		if receivers[id].Protocols["http"] != nil {
			return
		}
	}

	reporter.AddWarningWithExplain("collector.receivers.http-protocol-missing",
		"The value of receivers > otlp > protocols > http is nil. Make sure the key exists on your config.yaml")
}

func checkOTLPHTTPExporterEndpoint(reporter *utils.ComponentReporter, exporters map[string]exporterConfig) {
	ids := otlpHTTPExporterIDs(exporters)
	for _, id := range ids {
		if otlpHTTPGrafanaEndpointPattern.MatchString(exporters[id].Endpoint) {
			reporter.AddSuccessfulCheck(fmt.Sprintf("Value of exporter > %s > endpoint on config.yaml set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp", id))
			return
		}
	}

	for _, id := range ids {
		if strings.Contains(exporters[id].Endpoint, "localhost") {
			reporter.AddWarningWithExplain("collector.endpoint.localhost",
				fmt.Sprintf("Value of exporter > %s > endpoint on config.yaml is set to localhost. Update to a Grafana endpoint similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp to be able to send telemetry to your Grafana Cloud instance", id))
			return
		}
	}

	reporter.AddErrorWithExplain("collector.endpoint.invalid-format",
		"Value of exporter > otlphttp > endpoint on config.yaml is not set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp")
}

func otlpHTTPExporterIDs(exporters map[string]exporterConfig) []string {
	return componentIDsWithType(exporters, "otlphttp", "otlp_http")
}

func componentIDsWithType[T any](components map[string]T, types ...string) []string {
	ids := make([]string, 0, len(components))
	for id := range components {
		if componentIDHasType(id, types...) {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	return ids
}

func containsOTLPHTTPExporter(exporters []string) bool {
	for _, exporter := range exporters {
		if componentIDHasType(exporter, "otlphttp", "otlp_http") {
			return true
		}
	}
	return false
}

func containsOTLPReceiver(receivers []string) bool {
	for _, receiver := range receivers {
		if componentIDHasType(receiver, "otlp") {
			return true
		}
	}
	return false
}

func componentIDHasType(id string, types ...string) bool {
	componentType, _, _ := strings.Cut(id, "/")
	componentType = strings.ToLower(componentType)
	for _, t := range types {
		if componentType == t {
			return true
		}
	}
	return false
}
