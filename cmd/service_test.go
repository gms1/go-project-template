package cmd

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/gms1/go-project-template/pkg/common"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/assert"
)

func TestServiceCmd(t *testing.T) {
	SpanName = t.Name()
	ServiceInstanceId = "Test"
	stubs := gostub.New()
	defer stubs.Reset()
	stubs.StubFunc(&ServiceInitFunc, nil)
	stubs.Stub(&ServiceMainFunc, func(ctx context.Context, cancel context.CancelFunc) error {
		slog.InfoContext(ctx, "ok")
		return nil
	})

	rootCmd.SetArgs([]string{"service"})
	err := rootCmd.Execute()
	assert.NoError(t, err)
}

func TestRunServiceOk(t *testing.T) {
	SpanName = t.Name()
	ServiceInstanceId = "Test"

	assert.NoError(t, runService(
		func(ctx context.Context, cancel context.CancelFunc) error {
			return nil
		},
		func(ctx context.Context, cancel context.CancelFunc) error {
			// NOTE: testing span withoug any log message
			return nil
		},
	))
}

func TestRunServiceFailingInit(t *testing.T) {
	SpanName = t.Name()
	ServiceInstanceId = "Test"

	assert.Error(t, runService(
		func(ctx context.Context, cancel context.CancelFunc) error {
			return errors.New("test init failed")
		},
		func(ctx context.Context, cancel context.CancelFunc) error {
			slog.InfoContext(ctx, "ok")
			return nil
		},
	))
}

func TestRunServiceFailingMain(t *testing.T) {
	SpanName = t.Name()
	ServiceInstanceId = "Test"

	assert.Error(t, runService(
		func(ctx context.Context, cancel context.CancelFunc) error {
			return nil
		},
		func(ctx context.Context, cancel context.CancelFunc) error {
			return errors.New("test main failed")
		},
	))
}

func TestServiceTick(t *testing.T) {
	SpanName = t.Name()
	ServiceInstanceId = "Test"

	ctx, cancel := context.WithCancel(context.Background())
	err := serviceInit(ctx, cancel)
	assert.NoError(t, err)
	err = serviceInit(ctx, cancel)
	assert.Error(t, err)

	timoutTimer := time.AfterFunc(time.Millisecond*250, func() {
		cancel()
		assert.Fail(t, "timeout waiting for done")
	})
	defer timoutTimer.Stop()

	Tick = time.Millisecond * 50
	err = serviceMain(ctx, cancel)
	common.StopSignalHandling(ctx)
	assert.NoError(t, err)
}

func TestServiceCanceled(t *testing.T) {
	SpanName = t.Name()
	ctx, cancel := context.WithCancel(context.Background())
	err := serviceInit(ctx, cancel)
	assert.NoError(t, err)

	sigintTimer := time.AfterFunc(time.Millisecond*50, func() {
		cancel()
	})
	defer sigintTimer.Stop()

	Tick = time.Millisecond * 250
	err = serviceMain(ctx, cancel)
	common.StopSignalHandling(ctx)
	assert.NoError(t, err)
}
