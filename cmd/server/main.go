package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"go-gin-templete/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx, os.Args[1:]); err != nil {
		slog.Error("application stopped", slog.Any("error", err))
		os.Exit(1)
	}
}
