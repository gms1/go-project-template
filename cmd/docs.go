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
		if err := generateDocsFunc(rootCmd, args[0]); err != nil {
			return err
		}
		slog.InfoContext(cmd.Context(), "Generated docs", slog.String("directory", args[0]))
		return nil
	},
}
