package common

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	slogotel "github.com/remychantenay/slog-otel"
)

const LOG_LEVEL_NAME = "LOG_LEVEL"

// NOTE: use this LogLevelVar if you are going to create a new Logger for replacing the default logger
var LogLevelVar = new(slog.LevelVar)

func defaultLogLevel() slog.Level {
	var level slog.Level
	if err := level.UnmarshalText([]byte(GetEnv(LOG_LEVEL_NAME, "INFO"))); err != nil {
		return slog.LevelInfo
	}
	return level
}

func NewConsoleHandler(target *os.File) slog.Handler {
	return tint.NewHandler(target, &tint.Options{
		NoColor: !isatty.IsTerminal(target.Fd()),
		Level:   LogLevelVar,
	})
}

func defaultConsoleLogger() *slog.Logger {
	return slog.New(NewConsoleHandler(os.Stderr))
}

func defaultServiceLogger() *slog.Logger {
	// NOTE: using github.com/remychantenay/slog-otel instead of go.opentelemetry.io/contrib/bridges/otelslog
	// because the later does not seem to support chaining to a console handler; see below
	return slog.New(slogotel.OtelHandler{
		Next: NewConsoleHandler(os.Stdout),
	})
}

func InitConsoleLogging() {
	LogLevelVar.Set(defaultLogLevel())
	slog.SetDefault(defaultConsoleLogger())
}

func InitServiceLogging() {
	slog.SetDefault(defaultServiceLogger())
}

func init() {
	InitConsoleLogging()
}
