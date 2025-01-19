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

func TestServiceCmd(t *testing.T) {
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

	tick = time.Millisecond * 50
	assert.NoError(t, serviceMain(ctx, cancel, nil))
}

func TestServiceMainCancel(t *testing.T) {
	common.ServiceInstanceId = "Test"
	common.SpanName = t.Name()
	ctx, cancel := context.WithCancel(context.Background())

	assert.NoError(t, common.InitSignalHandler(ctx, cancel, nil))
	defer common.StopSignalHandling(ctx)

	sigintTimer := time.AfterFunc(time.Millisecond*50, func() {
		cancel()
	})
	defer sigintTimer.Stop()

	tick = time.Millisecond * 250
	assert.NoError(t, serviceMain(ctx, cancel, nil))
}
