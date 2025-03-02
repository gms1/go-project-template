package core

import (
	"context"
	"errors"
	"log/slog"
	"runtime/debug"
	"testing"

	"github.com/gms1/go-project-template/test"
	"github.com/lmittmann/tint"
	slogotel "github.com/remychantenay/slog-otel"
	"github.com/stretchr/testify/assert"
)

func TestDefaultLogLevel(t *testing.T) {
	testCases := []struct {
		name             string
		givenLogLevel    string
		expectedLogLevel slog.Level
	}{
		{"default", "", slog.LevelInfo},
		{"with DEBUG", "DEBUG", slog.LevelDebug},
		{"with WARN+4", "WARN+4", slog.LevelWarn + 4},
		{"with unknown", "unknown", slog.LevelInfo},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Setenv(LOG_LEVEL_NAME, testCase.givenLogLevel)
			level := defaultLogLevel()
			assert.Equal(t, testCase.expectedLogLevel, level)
		})
	}
}

func TestInitLogging(t *testing.T) {
	InitServiceLogging()
	assert.IsType(t, slogotel.OtelHandler{}, slog.Default().Handler())
	InitConsoleLogging()
	assert.IsType(t, tint.NewHandler(nil, &tint.Options{}), slog.Default().Handler())
}

func TestLogErrorAndStackTrace(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name       string
		givenError error
		hasStack   bool
	}{
		{"error without stack", errors.New("error{no stack}"), false},
		{"error with stack", NewStackTraceError(debug.Stack(), "{failed}"), true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, stderr, _ := test.CaptureOutput(func() error {
				InitConsoleLogging()
				LogErrorAndStackTrace(ctx, testCase.name, testCase.givenError)
				return nil
			})
			assert.Contains(t, stderr, testCase.name)
			assert.Contains(t, stderr, testCase.givenError.Error())
			if testCase.hasStack {
				assert.Contains(t, stderr, STACK_TRACE_MARKER)
			} else {
				assert.NotContains(t, stderr, STACK_TRACE_MARKER)
			}
		})
	}
}
