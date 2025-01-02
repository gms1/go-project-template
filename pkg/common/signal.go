package common

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	signalMutex sync.Mutex
	signalChan  chan os.Signal
)

func InitSignalHandler(ctx context.Context, cancel func(), sighupFunc *func()) error {
	signalMutex.Lock()
	defer signalMutex.Unlock()
	if signalChan != nil {
		return errors.New("signal handler is already initialized")
	}
	signalChan = make(chan os.Signal, 1)
	if sighupFunc != nil {
		signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP)
		slog.DebugContext(ctx, "Initialize signal handler for SIGINT, SIGKILL and SIGHUP")
	} else {
		signal.Notify(signalChan, os.Interrupt)
		slog.DebugContext(ctx, "Initialize signal handler for SIGINT and SIGKILL")
	}

	go signalHandler(signalChan, ctx, cancel, sighupFunc)
	return nil
}

func StopSignalHandling(ctx context.Context) {
	signalMutex.Lock()
	defer signalMutex.Unlock()
	if signalChan == nil {
		return
	}
	signal.Stop(signalChan)
	signalChan = nil
	slog.DebugContext(ctx, "Stopped signal handling")
}

func signalHandler(channel chan os.Signal, ctx context.Context, cancel func(), sighupFunc *func()) {
	for {
		select {
		case s := <-channel:
			switch s {
			case syscall.SIGHUP:
				slog.DebugContext(ctx, "Got SIGHUP")
				if sighupFunc != nil {
					(*sighupFunc)()
				}
			case os.Interrupt:
				slog.DebugContext(ctx, "Got SIGINT or SIGKILL")
				StopSignalHandling(ctx)
				cancel()
			}
		case <-ctx.Done():
			StopSignalHandling(ctx)
			slog.DebugContext(ctx, "Done in signal handler")
			return
		}
	}
}
