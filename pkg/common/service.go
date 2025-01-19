package common

import (
	"context"
	"log"
	"log/slog"
	"runtime"

	"github.com/KimMachineGun/automemlimit/memlimit"
	gomaxecs "github.com/rdforte/gomaxecs/maxprocs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/automaxprocs/maxprocs"
)

type (
	ServiceInitFunc func(ctx context.Context, cancel context.CancelFunc) error
	ServiceMainFunc func(ctx context.Context, cancel context.CancelFunc, span trace.Span) error
)

func RunService(main ServiceMainFunc, sigHupFunc *SigHupFunc) error {
	InitServiceLogging()
	InitServiceRuntime()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	traceProvider := NewOtelDefaultTraceProvider(ctx, ServiceInstanceId)
	defer func() { _ = traceProvider.Shutdown(ctx) }()

	otel.SetTracerProvider(traceProvider)

	tracer := traceProvider.Tracer(Package)

	err := InitSignalHandler(ctx, cancel, sigHupFunc)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to initialize", slog.Any("error", err))
		return err
	}
	defer StopSignalHandling(ctx)

	{
		ctx, span := tracer.Start(ctx, SpanName)
		defer span.End()

		err = main(ctx, cancel, span)
		if err != nil {
			slog.ErrorContext(ctx, "Running service failed", slog.Any("error", err))
			return err
		}
	}
	slog.DebugContext(ctx, "end")
	return nil
}

// InitServiceRuntime set defaults for GOMAXPROCS and GOMEMLIMIT if running in cgroup
// since currently the go runtime is not container/cgroup-aware (please see e.g https://github.com/golang/go/issues/33803)
func InitServiceRuntime() {
	// NOTE: maxprocs.Set honors GOMAXPROCS environment variable if present
	if gomaxecs.IsECS() {
		//nolint:errcheck
		gomaxecs.Set(gomaxecs.WithLogger(log.Printf))
	} else {
		//nolint:errcheck
		maxprocs.Set(maxprocs.Logger(log.Printf))
	}
	slog.Info("CPU:", slog.Int("GOMAXPROCS", runtime.GOMAXPROCS(0)), slog.Int("NumCPU", runtime.NumCPU()))

	// NOTE: memlimit.SetGoMemLimitWithOpts honors GOMEMLIMIT environment variable if present
	//nolint:errcheck
	memlimit.SetGoMemLimitWithOpts(
		memlimit.WithLogger(slog.Default()),
	)
}
