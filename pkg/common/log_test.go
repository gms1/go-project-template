package common

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultLogLevel(t *testing.T) {
	t.Setenv(LOG_LEVEL_NAME, "")
	level := defaultLogLevel()
	assert.Equal(t, slog.LevelInfo, level)
}

func TestKnownLogLevelString(t *testing.T) {
	t.Setenv(LOG_LEVEL_NAME, "DEBUG")
	level := defaultLogLevel()
	assert.Equal(t, slog.LevelDebug, level)
}

func TestKnownLogLevelNumber(t *testing.T) {
	t.Setenv(LOG_LEVEL_NAME, "WARN+4")
	level := defaultLogLevel()
	assert.Equal(t, slog.LevelWarn+4, level)
}

func TestUnknownLogLevel(t *testing.T) {
	t.Setenv(LOG_LEVEL_NAME, "unknown")
	level := defaultLogLevel()
	assert.Equal(t, slog.LevelInfo, level)
}

func TestInitServiceLogging(t *testing.T) {
	InitServiceLogging()
}
