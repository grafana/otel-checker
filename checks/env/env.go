package env

import (
	"fmt"
	"os"

	"otel-checker/checks/utils"
)

// EnvVar represents an environment variable configuration
type EnvVar struct {
	Name          string
	Required      bool
	Recommended   bool
	DefaultValue  string
	RequiredValue string
	Validator     func(value string, language string, reporter *utils.ComponentReporter)
	Description   string
	Message       string
}

// CheckEnvVar validates an environment variable against its configuration and reports the result
func CheckEnvVar(language string, envVar EnvVar, reporter *utils.ComponentReporter) {
	value := GetValue(envVar)
	if envVar.Validator != nil {
		envVar.Validator(value, language, reporter)
	} else {
		if envVar.RequiredValue != "" && !envVar.Recommended {
			envVar.Required = true
		}

		if envVar.Required && checkValue(envVar, value, reporter.AddError) {
			return
		}
		if envVar.Recommended && checkValue(envVar, value, reporter.AddWarning) {
			return
		}
		reporter.AddSuccessfulCheck(fmt.Sprintf("%s is set to '%s'", envVar.Name, value))
	}
}

func checkValue(e EnvVar, value string, report func(string)) bool {
	if e.RequiredValue != "" {
		if value != e.RequiredValue {
			if e.Message == "" {
				report(fmt.Sprintf("%s must be set to '%s'", e.Name, e.RequiredValue))
			} else {
				report(e.Message)
			}
			return true
		}
	} else {
		if value == "" {
			description := e.Message
			if description == "" {
				description = fmt.Sprintf("%s is not set", e.Name)
			}
			report(description)
			return true
		}
	}
	return false
}

// CheckEnvVars validates multiple environment variables and reports the results
func CheckEnvVars(reporter *utils.ComponentReporter, language string, envVars ...EnvVar) {
	for _, envVar := range envVars {
		CheckEnvVar(language, envVar, reporter)
	}
}

// GetValue returns the value of an environment variable with its default value if not set
func GetValue(envVar EnvVar) string {
	value := os.Getenv(envVar.Name)
	if value == "" && envVar.DefaultValue != "" {
		return envVar.DefaultValue
	}
	return value
}

// IsEnvVarSet checks if an environment variable is set
func IsEnvVarSet(envVar EnvVar) bool {
	return os.Getenv(envVar.Name) != ""
}
