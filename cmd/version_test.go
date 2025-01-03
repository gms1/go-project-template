package cmd

import (
	"log/slog"
	"os"
	"testing"

	"github.com/gms1/go-project-template/pkg/common"
	"github.com/gms1/go-project-template/test"
	"github.com/stretchr/testify/assert"
)

func TestVersionCmd(t *testing.T) {
	testCases := []struct {
		verbose          bool
		quiet            bool
		expectedLogLevel slog.Level
	}{
		{false, false, slog.LevelInfo},
		{true, false, slog.LevelDebug},
		{false, true, slog.LevelWarn},
	}
	loglevelOri := common.LogLevelVar.Level()
	for _, testCase := range testCases {
		Verbose = false
		Quiet = false
		args := []string{"version"}
		if os.Getenv(common.LOG_LEVEL_NAME) == "" {
			common.LogLevelVar.Set(slog.LevelInfo)
			if testCase.verbose {
				args = append(args, "-v")
			}
			if testCase.quiet {
				args = append(args, "-q")
			}
		}
		rootCmd.SetArgs(args)
		stdout, _, err := test.CaptureOutput(func() error { return rootCmd.Execute() })
		assert.NoError(t, err)
		assert.Equal(t, common.Version+"\n", stdout)
		if os.Getenv(common.LOG_LEVEL_NAME) == "" {
			assert.Equal(t, testCase.expectedLogLevel, common.LogLevelVar.Level())
			common.LogLevelVar.Set(loglevelOri)
		}
	}
}
