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
	assert.NoError(t, rootCmd.Execute())
}

func TestDocsCmdFailing(t *testing.T) {
	dir, _ := os.MkdirTemp("", "gotest")
	defer os.RemoveAll(dir)
	stubs := gostub.New()
	defer stubs.Reset()

	givenError := errors.New("test generating docs failed")
	stubs.StubFunc(&generateDocsFunc, givenError)

	rootCmd.SetArgs([]string{"docs", dir})
	assert.Equal(t, givenError, rootCmd.Execute())
}
