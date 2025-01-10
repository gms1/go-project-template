package cmd

import (
	"errors"
	"os"
	"testing"

	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/assert"
)

func TestDocsCmdOk(t *testing.T) {
	dir, _ := os.MkdirTemp("", "gotest")
	defer os.RemoveAll(dir)
	stubs := gostub.New()
	defer stubs.Reset()
	stubs.StubFunc(&generateDocsFunc, nil)

	rootCmd.SetArgs([]string{"docs", dir})
	err := rootCmd.Execute()
	assert.NoError(t, err)
}

func TestDocsCmdFailing(t *testing.T) {
	dir, _ := os.MkdirTemp("", "gotest")
	defer os.RemoveAll(dir)
	stubs := gostub.New()
	defer stubs.Reset()
	stubs.StubFunc(&generateDocsFunc, errors.New("Generating docs failed"))

	rootCmd.SetArgs([]string{"docs", dir})
	err := rootCmd.Execute()
	assert.Error(t, err)
}
