package cmd

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gms1/go-project-template/pkg/common"
	"github.com/stretchr/testify/assert"
)

func TestServiceCmd(t *testing.T) {
	serviceInitFuncOri := ServiceInitFunc
	serviceMainFuncOri := ServiceMainFunc

	ServiceInitFunc = func() (context.Context, context.CancelFunc, error) {
		ctx, cancel := context.WithCancel(context.Background())
		return ctx, cancel, nil
	}

	ServiceMainFunc = func(ctx context.Context, cancel context.CancelFunc) error {
		return nil
	}

	rootCmd.SetArgs([]string{"service"})

	err := rootCmd.Execute()

	ServiceInitFunc = serviceInitFuncOri
	ServiceMainFunc = serviceMainFuncOri

	assert.NoError(t, err)
}

func TestRunServiceOk(t *testing.T) {
	assert.NoError(t, runService(
		func() (context.Context, context.CancelFunc, error) {
			ctx, cancel := context.WithCancel(context.Background())
			return ctx, cancel, nil
		},
		func(ctx context.Context, cancel context.CancelFunc) error {
			return nil
		},
	))
}

func TestRunServiceFailingInit(t *testing.T) {
	assert.Error(t, runService(
		func() (context.Context, context.CancelFunc, error) {
			return nil, nil, errors.New("test init failed")
		},
		func(ctx context.Context, cancel context.CancelFunc) error {
			return nil
		},
	))
}

func TestRunServiceFailingMain(t *testing.T) {
	assert.Error(t, runService(
		func() (context.Context, context.CancelFunc, error) {
			ctx, cancel := context.WithCancel(context.Background())
			return ctx, cancel, nil
		},
		func(ctx context.Context, cancel context.CancelFunc) error {
			return errors.New("test main failed")
		},
	))
}

func TestServiceTick(t *testing.T) {
	ctx, cancel, err := serviceInit()
	assert.NoError(t, err)
	_, _, err = serviceInit()
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
	ctx, cancel, err := serviceInit()
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
