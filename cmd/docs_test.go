package cmd

import (
	"errors"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestDocsCmdOk(t *testing.T) {
	dir, _ := os.MkdirTemp("", "gotest")
	defer os.RemoveAll(dir)
	rootCmd.SetArgs([]string{"docs", dir})
	generateDocsFuncOrig := generateDocsFunc
	generateDocsFunc = func(cmd *cobra.Command, dir string) error {
		return nil
	}
	err := rootCmd.Execute()
	generateDocsFunc = generateDocsFuncOrig

	assert.NoError(t, err)
}

func TestDocsCmdFailing(t *testing.T) {
	dir, _ := os.MkdirTemp("", "gotest")
	defer os.RemoveAll(dir)
	rootCmd.SetArgs([]string{"docs", dir})
	generateDocsFuncOrig := generateDocsFunc
	generateDocsFunc = func(cmd *cobra.Command, dir string) error {
		return errors.New("Generating docs failed")
	}
	err := rootCmd.Execute()
	generateDocsFunc = generateDocsFuncOrig

	assert.Error(t, err)
}
