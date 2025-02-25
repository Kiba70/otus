package process

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var (
	ctx       context.Context
	ctxCancel context.CancelFunc
)

func Start() (context.Context, context.CancelFunc) {
	ctx, ctxCancel = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)
	return ctx, ctxCancel
}

func Stop() {
	ctxCancel()
}
