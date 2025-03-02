package cmd

import (
	"context"
	"log/slog"
	"runtime/debug"

	"github.com/gms1/go-project-template/pkg/common/core"
	"github.com/spf13/cobra"
)

var (
	verbose        bool
	quiet          bool
	silentEnd      bool // Commands that handle logging themselves should set this to true
	commandStarted bool // A command has started after the argument/option parsing has been succeeded
)

func Execute() error {
	verbose = false
	quiet = false
	silentEnd = false
	commandStarted = false
	err := executeCommand()
	if !silentEnd && commandStarted {
		if err != nil {
			core.LogErrorAndStackTrace(rootCmd.Context(), "Failed", err)
			return err
		}
		slog.InfoContext(rootCmd.Context(), "Succeeded")
	}
	return err
}

func executeCommand() (err error) {
	// NOTE: this method should not panic, so we are able to log all errors returned
	defer func() {
		if r := recover(); r != nil {
			err = core.NewStackTraceError(debug.Stack(), r)
		}
	}()
	err = rootCmd.ExecuteContext(context.Background())
	return err
}

func init() { //nolint:gochecknoinits
	cobra.EnableCommandSorting = false
	rootCmd.PersistentFlags().BoolVarP(&verbose, FLAG_VERBOSE_NAME, FLAG_VERBOSE_SHORTHAND, false, "verbose mode")
	rootCmd.PersistentFlags().BoolVarP(&quiet, FLAG_QUIET_NAME, FLAG_QUIET_SHORTHAND, false, "quiet mode")
	rootCmd.AddCommand(serviceCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(docsCmd)
}
