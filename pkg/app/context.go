package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// NewAppContext creates a cancellable context that listens to OS signals (SIGINT, SIGTERM).
// Useful for graceful shutdown across the app.
func NewAppContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		fmt.Println("\nReceived shutdown signal...")
		cancel()
	}()

	return ctx, cancel
}
