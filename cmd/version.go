package cmd

import (
	"fmt"

	"github.com/gms1/go-project-template/pkg/common"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Long:  fmt.Sprintf("Print the version of %s", common.Program),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(common.Version)
	},
}
