package collector

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/grafana/otel-checker/checks/utils"

	"go.yaml.in/yaml/v3"
)

func CheckCollectorSetup(reporter *utils.ComponentReporter, language string, configPath string) {
	checkCollectorConfig(reporter, configPath)
}

type configFile struct {
	Receivers struct {
		Otlp struct {
			Protocols map[string]any `yaml:"protocols"`
		} `yaml:"otlp"`
	} `yaml:"receivers"`
	Exporters struct {
		Otlphttp struct {
			Endpoint string                 `yaml:"endpoint"`
			Auth     map[string]interface{} `yaml:"auth"`
		} `yaml:"otlphttp"`
	} `yaml:"exporters"`
	Service struct {
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

// checkCollectorConfig reads the Collector configuration file and runs the
// downstream checks against it. configPath, when non-empty, must be the full
// path to the file the user wants checked. When configPath is empty, the
// checker falls back to config.yaml then config.yml in the current
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

	if c.Receivers.Otlp.Protocols["http"] == nil {
		reporter.AddWarningWithExplain("collector.receivers.http-protocol-missing",
			"The value of receivers > otlp > protocols > http is nil. Make sure the key exists on your config.yaml")
	}

	match, _ := regexp.MatchString("https:\\/\\/.+\\.grafana\\.net\\/otlp", c.Exporters.Otlphttp.Endpoint)
	if match {
		reporter.AddSuccessfulCheck("Value of exporter > otlphttp > endpoint on config.yaml set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp")
	} else {
		if strings.Contains(c.Exporters.Otlphttp.Endpoint, "localhost") {
			reporter.AddWarningWithExplain("collector.endpoint.localhost",
				"Value of exporter > otlphttp > endpoint on config.yaml is set to localhost. Update to a Grafana endpoint similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp to be able to send telemetry to your Grafana Cloud instance")
		} else {
			reporter.AddErrorWithExplain("collector.endpoint.invalid-format",
				"Value of exporter > otlphttp > endpoint on config.yaml is not set in the format similar to https://otlp-gateway-prod-us-east-0.grafana.net/otlp")
		}
	}

	// Traces
	if slices.Contains(c.Service.Pipelines.Traces.Exporters, "otlphttp") {
		reporter.AddSuccessfulCheck("Value of service > pipelines > traces > exporters on config.yaml contains otlphttp")
	} else {
		reporter.AddWarningWithExplain("collector.pipelines.traces-otlphttp-missing",
			"Value of service > pipelines > traces > exporters on config.yaml does not contain otlphttp")
	}
	if slices.Contains(c.Service.Pipelines.Traces.Receivers, "otlp") {
		reporter.AddSuccessfulCheck("Value of service > pipelines > traces > receivers on config.yaml contains otlp")
	} else {
		reporter.AddSuccessfulCheck("Value of service > pipelines > traces > receivers on config.yaml does not contain otlp")
	}

	// Logs
	if slices.Contains(c.Service.Pipelines.Logs.Exporters, "otlphttp") {
		reporter.AddSuccessfulCheck("Value of service > pipelines > logs > exporters on config.yaml contains otlphttp")
	} else {
		reporter.AddWarningWithExplain("collector.pipelines.logs-otlphttp-missing",
			"Value of service > pipelines > logs > exporters on config.yaml does not contain otlphttp")
	}
	if slices.Contains(c.Service.Pipelines.Logs.Receivers, "otlp") {
		reporter.AddSuccessfulCheck("Value of service > pipelines > logs > receivers on config.yaml contains otlp")
	} else {
		reporter.AddSuccessfulCheck("Value of service > pipelines > logs > receivers on config.yaml does not contain otlp")
	}

	// Metrics
	if slices.Contains(c.Service.Pipelines.Metrics.Exporters, "otlphttp") {
		reporter.AddSuccessfulCheck("Value of service > pipelines > metrics > exporters on config.yaml contains otlphttp")
	} else {
		reporter.AddWarningWithExplain("collector.pipelines.metrics-otlphttp-missing",
			"Value of service > pipelines > metrics > exporters on config.yaml does not contain otlphttp")
	}
	if slices.Contains(c.Service.Pipelines.Metrics.Receivers, "otlp") {
		reporter.AddSuccessfulCheck("Value of service > pipelines > metrics > receivers on config.yaml contains otlp")
	} else {
		reporter.AddSuccessfulCheck("Value of service > pipelines > metrics > receivers on config.yaml does not contain otlp")
	}
}
