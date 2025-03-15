package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetenv(t *testing.T) {
	const givenKey = "TEST_KEY"
	testCases := []struct {
		name                string
		envPrefixedValue    string
		envNotPrefixedValue string
		defaultValue        string
		expectedValue       string
	}{
		{"default", "", "", "baz", "baz"},
		{"env-without", "foo", "", "baz", "foo"},
		{"env-with", "", "bar", "baz", "bar"},
		{"env-with-and-without", "foo", "bar", "baz", "foo"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Setenv(givenKey, testCase.envNotPrefixedValue)
			if EnvPrefix != "" {
				t.Setenv(EnvPrefix+givenKey, testCase.envPrefixedValue)
			} else if testCase.envPrefixedValue != "" {
				t.Setenv(givenKey, testCase.envPrefixedValue)
			}
			resultValue := Getenv(givenKey, testCase.defaultValue)
			assert.Equal(t, testCase.expectedValue, resultValue)
		})
	}
}

func TestGetDefaultLogLevel(t *testing.T) {
	testCases := []struct {
		name                string
		envPrefixedValue    string
		envNotPrefixedValue string
		expectedValue       string
	}{
		{"default", "", "", LOG_LEVEL_DEFAULT_VALUE},
		{"env-without", "DEBUG", "", "DEBUG"},
		{"env-with", "", "DEBUG", "DEBUG"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Setenv(LOG_LEVEL_ENV_VAR, testCase.envNotPrefixedValue)
			if EnvPrefix != "" {
				t.Setenv(EnvPrefix+LOG_LEVEL_ENV_VAR, testCase.envPrefixedValue)
			} else if testCase.envPrefixedValue != "" {
				t.Setenv(LOG_LEVEL_ENV_VAR, testCase.envPrefixedValue)
			}
			resultValue := GetDefaultLogLevel()
			assert.Equal(t, testCase.expectedValue, resultValue)
		})
	}
}
