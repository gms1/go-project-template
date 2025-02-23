package cmd

import (
	"context"
	"log/slog"
	"time"

	"github.com/gms1/go-project-template/pkg/common/core"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/trace"
)

var (
	sighupFunc      *core.SigHupFunc
	serviceMainFunc core.ServiceMainFunc = serviceMain
	serviceSpanName                      = "main"
	tick                                 = 10 * time.Second
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Run a service",
	Long:  `Run as a service`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return core.RunService(serviceMainFunc, sighupFunc, serviceSpanName)
	},
}

//nolint:unparam
func serviceMain(ctx context.Context, cancel context.CancelFunc, span trace.Span) error {
	ticked := false
	for {
		select {
		case <-ctx.Done():
			slog.DebugContext(ctx, "Done in main")
			return nil
		case <-time.Tick(tick):
			if ticked {
				cancel()
				slog.DebugContext(ctx, "Cancelled in main")
				return nil
			} else {
				slog.InfoContext(ctx, "tick")
				ticked = true
			}
		}
	}
}
