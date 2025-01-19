package cmd

import (
	"log/slog"

	"github.com/gms1/go-project-template/pkg/common"
	"github.com/spf13/cobra"
)

const (
	FLAG_VERBOSE_NAME      = "verbose"
	FLAG_VERBOSE_SHORTHAND = "v"
	FLAG_QUIET_NAME        = "quiet"
	FLAG_QUIET_SHORTHAND   = "q"
)

var (
	verbose bool
	quiet   bool
)

var rootCmd = &cobra.Command{
	Use:     common.Program,
	Short:   common.DescriptionShort,
	Long:    common.DescriptionLong,
	Version: common.Version,
	Args:    cobra.MinimumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cmd.SilenceUsage = true
		v, _ := cmd.Flags().GetBool(FLAG_VERBOSE_NAME)
		q, _ := cmd.Flags().GetBool(FLAG_QUIET_NAME)
		if v && common.LogLevelVar.Level() > slog.LevelDebug {
			common.LogLevelVar.Set(slog.LevelDebug)
		}
		if q && common.LogLevelVar.Level() < slog.LevelWarn {
			common.LogLevelVar.Set(slog.LevelWarn)
		}
		slog.Debug("start")
	},
}

func Execute() error {
	initPersistenceFlagValues()
	return rootCmd.Execute()
}

func initPersistenceFlagValues() {
	verbose = false
	quiet = false
}

func init() {
	cobra.EnableCommandSorting = false
	rootCmd.PersistentFlags().BoolVarP(&verbose, FLAG_VERBOSE_NAME, FLAG_VERBOSE_SHORTHAND, false, "verbose mode")
	rootCmd.PersistentFlags().BoolVarP(&quiet, FLAG_QUIET_NAME, FLAG_QUIET_SHORTHAND, false, "quiet mode")
	rootCmd.AddCommand(serviceCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(docsCmd)
}
