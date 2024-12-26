package cmd

import (
	"log/slog"

	"github.com/gms1/go-project-template/pkg/common"
	"github.com/spf13/cobra"
)

var (
	verbose bool = false
	quiet   bool = false
)

var rootCmd = &cobra.Command{
	Use:     common.Program,
	Short:   common.DescriptionShort,
	Long:    common.DescriptionLong,
	Version: common.Version,
	Args:    cobra.MinimumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cmd.SilenceUsage = true
		if verbose {
			common.LogLevelVar.Set(slog.LevelDebug)
		}
		if quiet {
			common.LogLevelVar.Set(slog.LevelWarn)
		}
		slog.Debug("start")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.EnableCommandSorting = false
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose mode")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "quiet mode")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(docsCmd)
}
