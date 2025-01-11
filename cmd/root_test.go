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
	givenError := errors.New("test test-failing-command failed")
	failingCmd := &cobra.Command{
		Use:   "test-failing-command",
		Short: "failing test command",
		Long:  `test command that always fails`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return givenError
		},
	}
	rootCmd.AddCommand(failingCmd)
	defer func() {
		rootCmd.RemoveCommand(failingCmd)
	}()

	rootCmd.SetArgs([]string{"test-failing-command"})
	assert.Equal(t, givenError, Execute())
}
