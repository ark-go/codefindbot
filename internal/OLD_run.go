package internal

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

// Run runs f callback with context and logger, panics on error.
func Run(f func(ctx context.Context, log *zap.Logger) error) {
	//	log, err := zap.NewDevelopment()
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() { _ = log.Sync() }()
	// Никакого изящного выключения..
	//ctx := context.Background()
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
		//  os.Exit(1)
	}()
	//	go func() {
	if err := f(ctx, log); err != nil {
		log.Fatal("Run failed", zap.Error(err))
	}
	log.Info("Выход")
	//	}()

	// Done.
}
