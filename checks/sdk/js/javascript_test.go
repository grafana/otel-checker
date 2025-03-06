package js

import (
	"testing"

	"otel-checker/checks/utils"
)

func TestCheckEnvVars(t *testing.T) {
	tests := []utils.EnvVarTestCase{
		{
			Name: "all recommended env vars set correctly",
			EnvVars: map[string]string{
				"OTEL_NODE_RESOURCE_DETECTORS": "env,host,os,serviceinstance",
			},
			Language:       "javascript",
			ExpectedChecks: []string{"js: OTEL_NODE_RESOURCE_DETECTORS has recommended values"},
		},
		{
			Name:             "missing recommended env vars",
			EnvVars:          map[string]string{},
			Language:         "javascript",
			ExpectedWarnings: []string{"js: It's recommended the environment variable OTEL_NODE_RESOURCE_DETECTORS to be set to at least `env,host,os,serviceinstance`"},
		},
		{
			Name: "incomplete resource detectors",
			EnvVars: map[string]string{
				"OTEL_NODE_RESOURCE_DETECTORS": "env,host",
			},
			Language:         "javascript",
			ExpectedWarnings: []string{"js: It's recommended the environment variable OTEL_NODE_RESOURCE_DETECTORS to be set to at least `env,host,os,serviceinstance`"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			utils.RunEnvVarComponentTest(t, tt, "js",
				func(reporter utils.Reporter, c *utils.ComponentReporter, language string, components []string) {
					checkEnvVars(c)
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
			Language:       "javascript",
			ExpectedChecks: []string{"js: NODE_OPTIONS set correctly"},
			ExpectedErrors: []string{"js: Could not check file ./package.json: open ./package.json: no such file or directory"},
		},
		{
			Name:     "NODE_OPTIONS not set",
			EnvVars:  map[string]string{},
			Language: "javascript",
			ExpectedWarnings: []string{
				"js: NODE_OPTIONS not set. You can set it by running 'export NODE_OPTIONS=\"--require @opentelemetry/auto-instrumentations-node/register\"' or add the same '--require ...' when starting your application",
			},
			ExpectedErrors: []string{"js: Could not check file ./package.json: open ./package.json: no such file or directory"},
		},
		{
			Name: "NODE_OPTIONS set incorrectly",
			EnvVars: map[string]string{
				"NODE_OPTIONS": "--require something-else",
			},
			Language: "javascript",
			ExpectedWarnings: []string{
				"js: NODE_OPTIONS not set. You can set it by running 'export NODE_OPTIONS=\"--require @opentelemetry/auto-instrumentations-node/register\"' or add the same '--require ...' when starting your application",
			},
			ExpectedErrors: []string{"js: Could not check file ./package.json: open ./package.json: no such file or directory"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			utils.RunEnvVarComponentTest(t, tt, "js",
				func(reporter utils.Reporter, c *utils.ComponentReporter, language string, components []string) {
					checkJSAutoInstrumentation(c, "./")
				})
		})
	}
}
