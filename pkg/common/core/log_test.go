package core

import (
	"log/slog"
	"testing"

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

func TestInitServiceLogging(t *testing.T) {
	InitServiceLogging()
}
