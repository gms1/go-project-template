package cmd

import (
	"log/slog"

	"github.com/gms1/go-project-template/pkg/common/core"
	"github.com/spf13/cobra"
)

const (
	FLAG_VERBOSE_NAME      = "verbose"
	FLAG_VERBOSE_SHORTHAND = "v"
	FLAG_QUIET_NAME        = "quiet"
	FLAG_QUIET_SHORTHAND   = "q"
)

var rootCmd = &cobra.Command{
	Use:     core.Program,
	Short:   core.DescriptionShort,
	Long:    core.DescriptionLong,
	Version: core.Version,
	Args:    cobra.MinimumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cmd.SilenceUsage = true
		v, _ := cmd.Flags().GetBool(FLAG_VERBOSE_NAME)
		q, _ := cmd.Flags().GetBool(FLAG_QUIET_NAME)
		if v && core.LogLevelVar.Level() > slog.LevelDebug {
			core.LogLevelVar.Set(slog.LevelDebug)
		}
		if q && core.LogLevelVar.Level() < slog.LevelWarn {
			core.LogLevelVar.Set(slog.LevelWarn)
		}
		slog.Debug("start")
	},
}
