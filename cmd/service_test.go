package cmd

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/gms1/go-project-template/pkg/common"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func AssertNoSignalHandler(t *testing.T) {
	t.Helper()
	assert.False(t, common.HasSignalHandler(), "a signal handler is registered")
}

func TestServiceCmd(t *testing.T) {
	defer AssertNoSignalHandler(t)
	common.ServiceInstanceId = "Test"
	common.SpanName = t.Name()
	stubs := gostub.New()
	defer stubs.Reset()
	stubs.Stub(&serviceMainFunc, func(ctx context.Context, cancel context.CancelFunc, span trace.Span) error {
		slog.InfoContext(ctx, "ok")
		return nil
	})

	rootCmd.SetArgs([]string{"service"})
	assert.NoError(t, Execute())
}

func TestServiceMainTick(t *testing.T) {
	defer AssertNoSignalHandler(t)
	common.ServiceInstanceId = "Test"
	common.SpanName = t.Name()
	ctx, cancel := context.WithCancel(context.Background())

	assert.NoError(t, common.InitSignalHandler(ctx, cancel, nil))
	defer common.StopSignalHandling(ctx)

	timoutTimer := time.AfterFunc(time.Millisecond*250, func() {
		cancel()
		assert.Fail(t, "timeout waiting for done")
	})
	defer timoutTimer.Stop()

	stubs := gostub.New()
	stubs.Stub(&tick, time.Millisecond*50)
	defer stubs.Reset()

	assert.NoError(t, serviceMain(ctx, cancel, nil))
}

func TestServiceMainCancel(t *testing.T) {
	defer AssertNoSignalHandler(t)
	common.ServiceInstanceId = "Test"
	common.SpanName = t.Name()
	ctx, cancel := context.WithCancel(context.Background())

	assert.NoError(t, common.InitSignalHandler(ctx, cancel, nil))
	defer common.StopSignalHandling(ctx)

	sigintTimer := time.AfterFunc(time.Millisecond*50, func() {
		cancel()
	})
	defer sigintTimer.Stop()

	stubs := gostub.New()
	stubs.Stub(&tick, time.Millisecond*250)
	defer stubs.Reset()

	assert.NoError(t, serviceMain(ctx, cancel, nil))
}
