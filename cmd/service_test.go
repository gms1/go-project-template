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
	assert.NoError(t, rootCmd.Execute())
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
	givenError := errors.New("test init failed")

	assert.Equal(t, givenError, runService(
		func(ctx context.Context, cancel context.CancelFunc) error {
			return givenError
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
	givenError := errors.New("test main failed")

	assert.Equal(t, givenError, runService(
		func(ctx context.Context, cancel context.CancelFunc) error {
			return nil
		},
		func(ctx context.Context, cancel context.CancelFunc) error {
			return givenError
		},
	))
}

func TestServiceInit(t *testing.T) {
	SpanName = t.Name()
	ServiceInstanceId = "Test"
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		common.StopSignalHandling(ctx)
	}()

	assert.NoError(t, serviceInit(ctx, cancel))
	assert.Error(t, serviceInit(ctx, cancel))
}

func TestServiceMainTick(t *testing.T) {
	SpanName = t.Name()
	ServiceInstanceId = "Test"
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		common.StopSignalHandling(ctx)
	}()

	assert.NoError(t, serviceInit(ctx, cancel))

	timoutTimer := time.AfterFunc(time.Millisecond*250, func() {
		cancel()
		assert.Fail(t, "timeout waiting for done")
	})
	defer timoutTimer.Stop()

	Tick = time.Millisecond * 50
	assert.NoError(t, serviceMain(ctx, cancel))
}

func TestServiceMainCancel(t *testing.T) {
	SpanName = t.Name()
	ServiceInstanceId = "Test"
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		common.StopSignalHandling(ctx)
	}()

	assert.NoError(t, serviceInit(ctx, cancel))

	sigintTimer := time.AfterFunc(time.Millisecond*50, func() {
		cancel()
	})
	defer sigintTimer.Stop()

	Tick = time.Millisecond * 250
	assert.NoError(t, serviceMain(ctx, cancel))
}
