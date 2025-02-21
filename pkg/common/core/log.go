package core

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	slogotel "github.com/remychantenay/slog-otel"
)

const LOG_LEVEL_NAME = "LOG_LEVEL"

// NOTE: use this LogLevelVar if you are going to create a new Logger for replacing the default logger
var LogLevelVar = new(slog.LevelVar) //nolint:gochecknoglobals

func defaultLogLevel() slog.Level {
	var level slog.Level
	if err := level.UnmarshalText([]byte(Getenv(LOG_LEVEL_NAME, "INFO"))); err != nil {
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

func NewServiceHandler(target *os.File) slog.Handler {
	// NOTE: using github.com/remychantenay/slog-otel instead of go.opentelemetry.io/contrib/bridges/otelslog
	// because the later does not seem to support chaining to a console handler; see below
	return slogotel.OtelHandler{
		Next: NewConsoleHandler(target),
	}
}

func InitConsoleLogging() {
	slog.SetDefault(slog.New(NewConsoleHandler(os.Stderr)))
}

func InitServiceLogging() {
	slog.SetDefault(slog.New(NewServiceHandler(os.Stdout)))
}

func init() { //nolint:gochecknoinits
	LogLevelVar.Set(defaultLogLevel())
	InitConsoleLogging()
}
