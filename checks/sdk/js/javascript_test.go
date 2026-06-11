package js

import (
	"testing"

	"github.com/grafana/otel-checker/checks/utils"
)

func TestCheckEnvVars(t *testing.T) {
	missingWarning := "It's recommended the environment variable OTEL_NODE_RESOURCE_DETECTORS to be set to `all`, or to a list containing at least `env,host,os,serviceinstance`"
	tests := []utils.EnvVarTestCase{
		{
			Name: "all recommended env vars set correctly",
			EnvVars: map[string]string{
				"OTEL_NODE_RESOURCE_DETECTORS": "env,host,os,serviceinstance",
			},
			Language:       "js",
			ExpectedChecks: []string{"OTEL_NODE_RESOURCE_DETECTORS has recommended values"},
		},
		{
			Name: "value all accepted",
			EnvVars: map[string]string{
				"OTEL_NODE_RESOURCE_DETECTORS": "all",
			},
			Language:       "js",
			ExpectedChecks: []string{"OTEL_NODE_RESOURCE_DETECTORS has recommended values"},
		},
		{
			Name:             "missing recommended env vars",
			EnvVars:          map[string]string{},
			Language:         "js",
			ExpectedWarnings: []string{missingWarning},
		},
		{
			Name: "incomplete resource detectors",
			EnvVars: map[string]string{
				"OTEL_NODE_RESOURCE_DETECTORS": "env,host",
			},
			Language:         "js",
			ExpectedWarnings: []string{missingWarning},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			utils.RunEnvVarComponentTest(t, tt, "js",
				func(reporter utils.Reporter, c *utils.ComponentReporter, language string, components []string) {
					checkResourceDetectors(c)
				})
		})
	}
}

func TestCheckJSAutoInstrumentation(t *testing.T) {
	tests := []utils.EnvVarTestCase{
		{
			Name: "NODE_OPTIONS set correctly",
			EnvVars: map[string]string{
				"NODE_OPTIONS": "--require @opentelemetry/auto-instrumentations-node/register",
			},
			Language:       "js",
			ExpectedChecks: []string{"NODE_OPTIONS is set to '--require @opentelemetry/auto-instrumentations-node/register'"},
		},
		{
			Name:     "NODE_OPTIONS not set",
			EnvVars:  map[string]string{},
			Language: "js",
			ExpectedWarnings: []string{
				"NODE_OPTIONS not set. You can set it by running 'export NODE_OPTIONS=\"--require @opentelemetry/auto-instrumentations-node/register\"' or add the same '--require ...' when starting your application",
			},
		},
		{
			Name: "NODE_OPTIONS set incorrectly",
			EnvVars: map[string]string{
				"NODE_OPTIONS": "--require something-else",
			},
			Language: "js",
			ExpectedWarnings: []string{
				"NODE_OPTIONS not set. You can set it by running 'export NODE_OPTIONS=\"--require @opentelemetry/auto-instrumentations-node/register\"' or add the same '--require ...' when starting your application",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			utils.RunEnvVarComponentTest(t, tt, "js",
				func(reporter utils.Reporter, c *utils.ComponentReporter, language string, components []string) {
					checkAutoInstrumentationNodeOptions(c)
				})
		})
	}
}
