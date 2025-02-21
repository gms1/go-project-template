package cmd

import "github.com/spf13/cobra"

var (
	verbose bool
	quiet   bool
)

func Execute() error {
	verbose = false
	quiet = false
	return rootCmd.Execute()
}

func init() { //nolint:gochecknoinits
	cobra.EnableCommandSorting = false
	rootCmd.PersistentFlags().BoolVarP(&verbose, FLAG_VERBOSE_NAME, FLAG_VERBOSE_SHORTHAND, false, "verbose mode")
	rootCmd.PersistentFlags().BoolVarP(&quiet, FLAG_QUIET_NAME, FLAG_QUIET_SHORTHAND, false, "quiet mode")
	rootCmd.AddCommand(serviceCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(docsCmd)
}
