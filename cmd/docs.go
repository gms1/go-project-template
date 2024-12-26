package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var generateDocsFunc = doc.GenMarkdownTree

var docsCmd = &cobra.Command{
	Use:   "docs <path-to-docs-folder>",
	Short: "Generate docs",
	Long:  `Generate docs`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := generateDocsFunc(rootCmd, args[0])
		if err != nil {
			slog.Error("Failed to generate docs", "error", err)
			return err
		}
		for _, c := range rootCmd.Commands() {
			slog.Info("command", "command", c.Name())
		}
		slog.Info("Generated docs", slog.String("directory", args[0]))
		return nil
	},
}
