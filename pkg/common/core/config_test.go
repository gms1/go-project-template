package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetenv(t *testing.T) {
	testCases := []struct {
		name          string
		withDefault   bool
		expectedValue string
	}{
		{"default", true, "foo"},
		{"explicit", false, "foo"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			const givenKey = "TEST_KEY"
			var resultValue string
			if testCase.withDefault {
				t.Setenv(givenKey, "")
				resultValue = Getenv(givenKey, testCase.expectedValue)
			} else {
				t.Setenv(givenKey, testCase.expectedValue)
				resultValue = Getenv(givenKey, "bar")
			}
			assert.Equal(t, testCase.expectedValue, resultValue)
		})
	}
}
