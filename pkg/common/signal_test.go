package common

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitSignalHandler(t *testing.T) {
	testCases := []struct {
		name       string
		withSighup bool
	}{
		{"with sighup", true},
		{"without sighup", false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			sighupCounter := 0

			if testCase.withSighup {
				sighupFunc := func() {
					sighupCounter++
				}
				assert.NoError(t, InitSignalHandler(ctx, cancel, &sighupFunc), "init with sighup")
			} else {
				assert.NoError(t, InitSignalHandler(ctx, cancel, nil), "init without sighup")
			}

			err := InitSignalHandler(ctx, cancel, nil)
			assert.Error(t, err, "init twice")

			var sighupTimer, sigintTimer, timoutTimer *time.Timer

			duration := time.Millisecond * 50

			if testCase.withSighup {
				sighupTimer = time.AfterFunc(duration, func() {
					t.Log("SENDING SIGHUP")
					assert.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGHUP), "send SIGHUP signal")
				})
				duration += time.Millisecond * 200
			}

			sigintTimer = time.AfterFunc(duration, func() {
				t.Log("SENDING SIGINT")
				assert.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGINT), "send SIGINT signal")
			})
			duration += time.Millisecond * 200

			timoutTimer = time.AfterFunc(duration, func() {
				StopSignalHandling(ctx)
				cancel()
				assert.Fail(t, "timeout waiting for done")
			})

			defer func() {
				if testCase.withSighup {
					sighupTimer.Stop()
				}
				sigintTimer.Stop()
				timoutTimer.Stop()
			}()

		out:
			//nolint:gosimple
			for {
				select {
				case <-ctx.Done():
					break out
				}
			}

			if testCase.withSighup {
				assert.Equal(t, 1, sighupCounter)
			} else {
				assert.Equal(t, 0, sighupCounter)
			}
		})
	}
}
