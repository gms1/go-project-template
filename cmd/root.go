package cmd

import (
	"log/slog"

	"github.com/gms1/go-project-template/pkg/common"
	"github.com/spf13/cobra"
)

var (
	Verbose bool = false
	Quiet   bool = false
)

var rootCmd = &cobra.Command{
	Use:     common.Program,
	Short:   common.DescriptionShort,
	Long:    common.DescriptionLong,
	Version: common.Version,
	Args:    cobra.MinimumNArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		cmd.SilenceUsage = true
		if Verbose {
			common.LogLevelVar.Set(slog.LevelDebug)
		}
		if Quiet {
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
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose mode")
	rootCmd.PersistentFlags().BoolVarP(&Quiet, "quiet", "q", false, "quiet mode")
	rootCmd.AddCommand(serviceCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(docsCmd)
}
