package cmd

import (
	"errors"
	"os"
	"testing"

	"github.com/gms1/go-project-template/pkg/common/core"
	"github.com/gms1/go-project-template/test"
	"github.com/prashantv/gostub"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestExecuteWithFailingCommand(t *testing.T) {
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
	defer rootCmd.RemoveCommand(failingCmd)

	rootCmd.SetArgs([]string{"test-failing-command"})
	assert.Equal(t, givenError, core.ErrorRootCause(Execute()))
}

func TestExecuteHelpCommand(t *testing.T) {
	stdout, _, err := test.CaptureOutput(func() error {
		rootCmd.SetArgs([]string{"help"})
		return Execute()
	})
	assert.NoError(t, err)
	assert.Contains(t, stdout, "Usage:")
}

func TestExecuteDocsCommand(t *testing.T) {
	dir, _ := os.MkdirTemp("", "gotest")
	defer os.RemoveAll(dir)

	testCases := []struct {
		name string
		args []string
		err  error
	}{
		{"docs command without error", []string{"docs", dir}, nil},
		{"docs command with error", []string{"docs", dir}, errors.New("test generating docs failed")},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			stubs := gostub.New()
			defer stubs.Reset()
			stubs.StubFunc(&generateDocsFunc, testCase.err)

			rootCmd.SetArgs(testCase.args)
			assert.Equal(t, testCase.err, core.ErrorRootCause(Execute()))
		})
	}
}

func TestExecuteDocsCommandPanics(t *testing.T) {
	givenErrorText := "i am panicing"

	dir, _ := os.MkdirTemp("", "gotest")
	defer os.RemoveAll(dir)
	stubs := gostub.New()
	defer stubs.Reset()
	stubs.Stub(&generateDocsFunc, func(*cobra.Command, string) error {
		panic(givenErrorText)
	})

	rootCmd.SetArgs([]string{"docs", dir})
	err := Execute()
	assert.Contains(t, givenErrorText, err.Error())
}
