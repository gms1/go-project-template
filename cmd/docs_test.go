package cmd

import (
	"errors"
	"os"
	"testing"

	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/assert"
)

func TestDocsCmd(t *testing.T) {
	testCases := []struct {
		name string
		err  error
	}{
		{"without error", nil},
		{"with error", errors.New("test generating docs failed")},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			dir, _ := os.MkdirTemp("", "gotest")
			defer os.RemoveAll(dir)
			stubs := gostub.New()
			defer stubs.Reset()
			stubs.StubFunc(&generateDocsFunc, testCase.err)

			rootCmd.SetArgs([]string{"docs", dir})
			assert.Equal(t, testCase.err, Execute())
		})
	}
}
