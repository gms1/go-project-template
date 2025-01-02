package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultValue(t *testing.T) {
	const givenKey = "TEST_KEY"
	const expectedValue = "foo"
	t.Setenv(givenKey, "")
	resultValue := GetEnv(givenKey, expectedValue)
	assert.Equal(t, expectedValue, resultValue)
}

func TestProvidedValue(t *testing.T) {
	const givenKey = "TEST_KEY"
	const expectedValue = "foo"
	t.Setenv(givenKey, expectedValue)
	resultValue := GetEnv(givenKey, "bar")
	assert.Equal(t, expectedValue, resultValue)
}
