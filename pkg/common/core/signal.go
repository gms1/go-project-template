package core

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type SigHupFunc func()

var (
	signalMutex                        sync.Mutex     //nolint:gochecknoglobals
	signalChan                         chan os.Signal //nolint:gochecknoglobals
	ErrSignalHandlerAlreadyInitialized = errors.New("signal handler is already initialized")
)

func InitSignalHandler(ctx context.Context, cancel func(), sighupFunc *SigHupFunc) {
	signalMutex.Lock()
	defer signalMutex.Unlock()
	if signalChan != nil {
		panic(ErrSignalHandlerAlreadyInitialized)
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

func HasSignalHandler() bool {
	signalMutex.Lock()
	defer signalMutex.Unlock()
	return signalChan != nil
}

func signalHandler(channel chan os.Signal, ctx context.Context, cancel func(), sighupFunc *SigHupFunc) {
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
				slog.DebugContext(ctx, "Cancelled in signal handler")
			}
		case <-ctx.Done():
			slog.DebugContext(ctx, "Done in signal handler")
			return
		}
	}
}
