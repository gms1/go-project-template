package cmd

import (
	"log/slog"
	"testing"

	"github.com/gms1/go-project-template/pkg/common/core"
	"github.com/gms1/go-project-template/test"
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
	loglevelOri := core.LogLevelVar.Level()
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			core.LogLevelVar.Set(slog.LevelInfo)
			defer core.LogLevelVar.Set(loglevelOri)

			args := []string{"version"}
			if testCase.verbose {
				args = append(args, "-v")
			}
			if testCase.quiet {
				args = append(args, "-q")
			}
			rootCmd.SetArgs(args)
			stdout, _, err := test.CaptureOutput(Execute)
			assert.NoError(t, err)
			assert.Equal(t, core.Version+"\n", stdout)
			assert.Equal(t, testCase.expectedLogLevel, core.LogLevelVar.Level())

			v, err := rootCmd.Flags().GetBool(FLAG_VERBOSE_NAME)
			assert.Nil(t, err)

			q, err := rootCmd.Flags().GetBool(FLAG_QUIET_NAME)
			assert.Nil(t, err)

			assert.Equal(t, testCase.verbose, v)
			assert.Equal(t, testCase.quiet, q)
		})
	}
}
