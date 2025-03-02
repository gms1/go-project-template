package cmd

import (
	"fmt"

	"github.com/gms1/go-project-template/pkg/common/core"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Long:  fmt.Sprintf("Print the version of %s", core.Program),
	Run: func(cmd *cobra.Command, args []string) {
		silentEnd = true
		fmt.Println(core.Version)
	},
}
