package core

import (
	"context"
	"log"
	"log/slog"
	"runtime"

	"github.com/KimMachineGun/automemlimit/memlimit"
	gomaxecs "github.com/rdforte/gomaxecs/maxprocs"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/automaxprocs/maxprocs"
)

type (
	ServiceMainFunc func(ctx context.Context, cancel context.CancelFunc, span trace.Span) error
)

type Service struct {
	main          ServiceMainFunc
	sigHupFunc    *SigHupFunc
	ctx           context.Context
	cancel        context.CancelFunc
	traceProvider *sdktrace.TracerProvider
	tracer        trace.Tracer
}

func NewService(main ServiceMainFunc, sigHupFunc *SigHupFunc) Service {
	service := Service{main: main, sigHupFunc: sigHupFunc}
	service.ctx, service.cancel = context.WithCancel(context.Background())
	return service
}

func (service *Service) Init() error {
	if err := InitSignalHandler(service.ctx, service.cancel, service.sigHupFunc); err != nil {
		slog.ErrorContext(service.ctx, "Failed to initialize", slog.Any("error", err))
		return err
	}
	service.traceProvider = NewOtelDefaultTraceProvider(service.ctx, ServiceInstanceId)
	otel.SetTracerProvider(service.traceProvider)
	service.tracer = service.traceProvider.Tracer(Package)
	return nil
}

func (service *Service) Run() error {
	ctx, span := service.tracer.Start(service.ctx, SpanName)
	defer span.End()

	err := service.main(ctx, service.cancel, span)
	if err != nil {
		slog.ErrorContext(ctx, "Running service failed", slog.Any("error", err))
		return err
	}
	return nil
}

func RunService(main ServiceMainFunc, sigHupFunc *SigHupFunc) error {
	InitServiceLogging()
	InitServiceRuntime()

	service := NewService(main, sigHupFunc)
	defer service.cancel()

	err := service.Init()
	if err != nil {
		return err
	}
	defer func() {
		StopSignalHandling(service.ctx)
		_ = service.traceProvider.Shutdown(service.ctx)
	}()

	err = service.Run()
	if err != nil {
		return err
	}
	slog.DebugContext(service.ctx, "end")
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
