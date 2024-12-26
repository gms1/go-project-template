package common

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/mattn/go-isatty"
)

// NOTE: use this LogLevelVar if you are going to create a new Logger for replacing the default logger
var LogLevelVar = new(slog.LevelVar)

func init() {
	logger := slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		NoColor: !isatty.IsTerminal(os.Stderr.Fd()),
		Level:   LogLevelVar,
	}))
	slog.SetDefault(logger)
}
