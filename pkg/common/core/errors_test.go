package core

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorRootCause(t *testing.T) {
	errNotWrapped := errors.New("not wrapped")
	errWrappedOneLevel := fmt.Errorf("wrapped one level: %w", errNotWrapped)
	errWrappedSecondLevel := fmt.Errorf("wrapped second level: %w", errWrappedOneLevel)

	testCases := []struct {
		name          string
		givenError    error
		expectedError error
	}{
		{"no error", nil, nil},
		{"error not wrapped", errNotWrapped, errNotWrapped},
		{"error wrapped one level", errWrappedOneLevel, errNotWrapped},
		{"error wrapped second level", errWrappedSecondLevel, errNotWrapped},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expectedError, ErrorRootCause(testCase.givenError))
		})
	}
}

func TestNewStackTraceError(t *testing.T) {
	givenStackTraceText := "my error stack"
	givenErrorMessage := "my error message"
	err := NewStackTraceError([]byte(givenStackTraceText), givenErrorMessage)
	assert.Equal(t, givenStackTraceText, Stack(err).Error())
	assert.Equal(t, givenErrorMessage, err.Error())
	assert.Nil(t, errors.Unwrap(err))
}

func TestStackTraceErrorf(t *testing.T) {
	givenStackTraceText := "my error stack"
	givenErrorMessage := "my error message"
	err := StackTraceErrorf([]byte(givenStackTraceText), "%s", givenErrorMessage)
	assert.Equal(t, givenStackTraceText, Stack(err).Error())
	assert.Equal(t, givenErrorMessage, err.Error())
	assert.Nil(t, errors.Unwrap(err))
}

func TestToStackTraceError(t *testing.T) {
	givenStackTraceText := "my error stack"
	givenErrorMessage := "my error message"
	givenWrappedError := errors.New("my wrapped error")
	givenError := fmt.Errorf("%s: %w", givenErrorMessage, givenWrappedError)
	err := ToStackTraceError([]byte(givenStackTraceText), givenError)
	assert.Equal(t, givenStackTraceText, Stack(err).Error())
	assert.Contains(t, err.Error(), givenErrorMessage)
	assert.Equal(t, givenWrappedError, errors.Unwrap(err))
}
