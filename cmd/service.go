package cmd

import (
	"context"
	"log/slog"
	"time"

	"github.com/gms1/go-project-template/pkg/common"
	"github.com/spf13/cobra"
)

var (
	ServiceInitFunc func() (context.Context, context.CancelFunc, error)        = serviceInit
	ServiceMainFunc func(ctx context.Context, cancel context.CancelFunc) error = serviceMain
	Tick                                                                       = 10 * time.Second
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Run a service",
	Long:  `Run as a service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runService(ServiceInitFunc, ServiceMainFunc)
	},
}

func runService(init func() (context.Context, context.CancelFunc, error), main func(ctx context.Context, cancel context.CancelFunc) error) error {
	ctx, cancel, err := init()
	if err != nil {
		slog.Error("Failed to initialize", slog.Any("error", err))
		return err
	}
	defer cancel()
	err = main(ctx, cancel)
	if err != nil {
		slog.Error("Running service failed", slog.Any("error", err))
		return err
	}
	return nil
}

func serviceInit() (context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())
	if err := common.InitSignalHandler(ctx, cancel, nil); err != nil {
		return nil, nil, err
	}
	return ctx, cancel, nil
}

//nolint:unparam
func serviceMain(ctx context.Context, cancel context.CancelFunc) error {
	for {
		select {
		case <-ctx.Done():
			slog.Debug("Done in main")
			return nil
		case <-time.Tick(Tick):
			slog.Info("tick")
			cancel()
			return nil
		}
	}
}
