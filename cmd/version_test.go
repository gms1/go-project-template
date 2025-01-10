package cmd

import (
	"log/slog"
	"testing"

	"github.com/gms1/go-project-template/pkg/common"
	"github.com/gms1/go-project-template/test"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/assert"
)

func TestVersionCmd(t *testing.T) {
	testCases := []struct {
		name             string
		verbose          bool
		quiet            bool
		expectedLogLevel slog.Level
	}{
		{"no option", false, false, slog.LevelInfo},
		{"verbose", true, false, slog.LevelDebug},
		{"quiet", false, true, slog.LevelWarn},
	}
	loglevelOri := common.LogLevelVar.Level()
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			common.LogLevelVar.Set(slog.LevelInfo)
			defer func() {
				common.LogLevelVar.Set(loglevelOri)
			}()
			stubs := gostub.New()
			defer stubs.Reset()
			stubs.Stub(&Verbose, false)
			stubs.Stub(&Quiet, false)

			args := []string{"version"}
			if testCase.verbose {
				args = append(args, "-v")
			}
			if testCase.quiet {
				args = append(args, "-q")
			}
			rootCmd.SetArgs(args)
			stdout, _, err := test.CaptureOutput(func() error { return rootCmd.Execute() })
			assert.NoError(t, err)
			assert.Equal(t, common.Version+"\n", stdout)
			assert.Equal(t, testCase.expectedLogLevel, common.LogLevelVar.Level())
		})
	}
}
