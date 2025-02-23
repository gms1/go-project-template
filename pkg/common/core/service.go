package core

import (
	"context"
	"log"
	"log/slog"
	"runtime"

	"github.com/KimMachineGun/automemlimit/memlimit"
	"github.com/google/uuid"
	gomaxecs "github.com/rdforte/gomaxecs/maxprocs"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/automaxprocs/maxprocs"
)

type (
	ServiceMainFunc func(ctx context.Context, cancel context.CancelFunc, span trace.Span) error
)

type ServiceContext struct {
	ctx           context.Context
	cancel        context.CancelFunc
	traceProvider *sdktrace.TracerProvider
	tracer        trace.Tracer
}

func NewServiceContext(sigHupFunc *SigHupFunc) ServiceContext {
	service := ServiceContext{}
	service.ctx, service.cancel = context.WithCancel(context.Background())
	InitSignalHandler(service.ctx, service.cancel, sigHupFunc)
	service.traceProvider = NewOtelDefaultTraceProvider(service.ctx, uuid.NewString())
	otel.SetTracerProvider(service.traceProvider)
	service.tracer = service.traceProvider.Tracer(Package)
	return service
}

func (service *ServiceContext) Run(main ServiceMainFunc, serviceSpanName string) error {
	ctx, span := service.tracer.Start(service.ctx, serviceSpanName)
	defer span.End()

	err := main(ctx, service.cancel, span)
	if err != nil {
		slog.ErrorContext(ctx, "Running service failed", slog.Any("error", err))
		return err
	}
	slog.DebugContext(service.ctx, "end")
	return nil
}

func (service *ServiceContext) Shutdown() {
	StopSignalHandling(service.ctx)
	service.cancel()
	_ = service.traceProvider.Shutdown(service.ctx)
}

func RunService(main ServiceMainFunc, sigHupFunc *SigHupFunc, serviceSpanName string) error {
	InitServiceLogging()
	InitServiceRuntime()

	service := NewServiceContext(sigHupFunc)
	defer service.Shutdown()

	return service.Run(main, serviceSpanName)
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
