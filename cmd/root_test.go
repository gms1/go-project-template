package cmd

import (
	"errors"
	"testing"

	"github.com/gms1/go-project-template/test"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestRootCmdOk(t *testing.T) {
	stdout, _, err := test.CaptureOutput(func() error {
		rootCmd.SetArgs([]string{"help"})
		return Execute()
	})
	assert.NoError(t, err)
	assert.Contains(t, stdout, "Usage:")
}

func TestRootCmdFailing(t *testing.T) {
	failingError := errors.New("test-failing-command failed")
	failingCmd := &cobra.Command{
		Use:   "test-failing-command",
		Short: "failing test command",
		Long:  `test command that always fails`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return failingError
		},
	}
	rootCmd.AddCommand(failingCmd)

	rootCmd.SetArgs([]string{"test-failing-command"})
	err := Execute()
	rootCmd.RemoveCommand(failingCmd)
	t.Logf("got error '%v'", err)
	assert.Error(t, err)
	assert.Equal(t, failingError.Error(), err.Error())
}
