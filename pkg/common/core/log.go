package core

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
	slogotel "github.com/remychantenay/slog-otel"
)

const (
	LOG_LEVEL_NAME              = "LOG_LEVEL"
	STACK_TRACE_MARKER          = "Stacktrace:"
	DEFAULT_CONSOLE_TIME_FORMAT = time.TimeOnly
	DEFAULT_SERVICE_TIME_FORMAT = time.RFC3339
)

// NOTE: use this LogLevelVar if you are going to create a new Logger for replacing the default logger
var LogLevelVar = new(slog.LevelVar) //nolint:gochecknoglobals

func defaultLogLevel() slog.Level {
	var level slog.Level
	if err := level.UnmarshalText([]byte(Getenv(LOG_LEVEL_NAME, "INFO"))); err != nil {
		return slog.LevelInfo
	}
	return level
}

func NewConsoleHandler(target *os.File, timeFormat string) slog.Handler {
	return tint.NewHandler(target, &tint.Options{
		NoColor:    !isatty.IsTerminal(target.Fd()),
		Level:      LogLevelVar,
		TimeFormat: timeFormat,
	})
}

func NewServiceHandler(target *os.File) slog.Handler {
	// NOTE: using github.com/remychantenay/slog-otel instead of go.opentelemetry.io/contrib/bridges/otelslog
	// because the later does not seem to support chaining to a console handler; see below
	return slogotel.OtelHandler{
		Next: NewConsoleHandler(target, DEFAULT_SERVICE_TIME_FORMAT),
	}
}

func InitConsoleLogging() {
	slog.SetDefault(slog.New(NewConsoleHandler(os.Stderr, DEFAULT_CONSOLE_TIME_FORMAT)))
}

func InitServiceLogging() {
	slog.SetDefault(slog.New(NewServiceHandler(os.Stdout)))
}

func init() { //nolint:gochecknoinits
	LogLevelVar.Set(defaultLogLevel())
	InitConsoleLogging()
}

func LogErrorAndStackTrace(ctx context.Context, text string, err error) {
	if stack := Stack(err); stack != nil {
		//nolint:errorlint
		slog.ErrorContext(ctx, fmt.Errorf("%s: %w: %s %v", text, err, STACK_TRACE_MARKER, stack).Error())
	} else {
		slog.ErrorContext(ctx, fmt.Errorf("%s: %w", text, err).Error())
	}
}
