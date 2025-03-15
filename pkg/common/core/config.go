package core

import (
	"os"
)

const (
	LOG_LEVEL_ENV_VAR       = "LOG_LEVEL"
	LOG_LEVEL_DEFAULT_VALUE = "INFO"
)

func Getenv(name, defaultValue string) string {
	var result string
	if EnvPrefix != "" {
		result = os.Getenv(EnvPrefix + name)
		if result != "" {
			return result
		}
	}
	result = os.Getenv(name)
	if result != "" {
		return result
	}
	return defaultValue
}

func GetDefaultLogLevel() string {
	return Getenv(LOG_LEVEL_ENV_VAR, LOG_LEVEL_DEFAULT_VALUE)
}
