package common

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

const LOG_LEVEL_NAME = "LOG_LEVEL"

// NOTE: use this LogLevelVar if you are going to create a new Logger for replacing the default logger
var LogLevelVar = new(slog.LevelVar)

func defaultConsoleLogger() *slog.Logger {
	return slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		NoColor: !isatty.IsTerminal(os.Stderr.Fd()),
		Level:   LogLevelVar,
	}))
}

func defaultLogLevel() slog.Level {
	var level slog.Level
	if err := level.UnmarshalText([]byte(GetEnv(LOG_LEVEL_NAME, "INFO"))); err != nil {
		return slog.LevelInfo
	}
	return level
}

func init() {
	LogLevelVar.Set(defaultLogLevel())
	slog.SetDefault(defaultConsoleLogger())
}
