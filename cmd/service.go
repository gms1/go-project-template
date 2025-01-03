package cmd

import (
	"context"
	"log/slog"
	"time"

	"github.com/gms1/go-project-template/pkg/common"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
)

var (
	ServiceInitFunc   func(ctx context.Context, cancel context.CancelFunc) error = serviceInit
	ServiceMainFunc   func(ctx context.Context, cancel context.CancelFunc) error = serviceMain
	SpanName                                                                     = "RunService"
	ServiceInstanceId                                                            = "Default"
	Tick                                                                         = 10 * time.Second
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Run a service",
	Long:  `Run as a service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runService(ServiceInitFunc, ServiceMainFunc)
	},
}

func runService(init func(ctx context.Context, cancel context.CancelFunc) error, main func(ctx context.Context, cancel context.CancelFunc) error) error {
	common.InitServiceLogging()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	traceProvider := common.NewOtelDefaultTraceProvider(ctx, ServiceInstanceId)
	defer func() { _ = traceProvider.Shutdown(ctx) }()

	otel.SetTracerProvider(traceProvider)

	tracer := traceProvider.Tracer(common.Package)

	err := init(ctx, cancel)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to initialize", slog.Any("error", err))
		return err
	}

	{
		ctx, span := tracer.Start(ctx, SpanName)
		defer span.End()

		err = main(ctx, cancel)
		if err != nil {
			slog.ErrorContext(ctx, "Running service failed", slog.Any("error", err))
			return err
		}
	}
	slog.DebugContext(ctx, "end")
	return nil
}

func serviceInit(ctx context.Context, cancel context.CancelFunc) error {
	if err := common.InitSignalHandler(ctx, cancel, nil); err != nil {
		return err
	}
	return nil
}

//nolint:unparam
func serviceMain(ctx context.Context, cancel context.CancelFunc) error {
	ticked := false
	for {
		select {
		case <-ctx.Done():
			slog.DebugContext(ctx, "Done in main")
			return nil
		case <-time.Tick(Tick):
			if ticked {
				cancel()
			} else {
				slog.InfoContext(ctx, "tick")
			}
			return nil
		}
	}
}
