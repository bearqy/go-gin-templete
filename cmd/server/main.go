package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/bearqy/go-gin-templete/internal/api"
	"github.com/bearqy/go-gin-templete/internal/cli"
	"github.com/bearqy/go-gin-templete/internal/config"
	"github.com/bearqy/go-gin-templete/internal/db"
	"github.com/bearqy/go-gin-templete/internal/job"
	"github.com/bearqy/go-gin-templete/internal/logger"
	"github.com/bearqy/go-gin-templete/internal/web"
)

func main() {

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	logger.Init()

	cli.Init()

	err := config.Init()
	if err != nil {
		os.Exit(5)
	}

	oldLevel := logger.SetLevel(config.Main.Log.Level)
	if oldLevel == "" {
		os.Exit(5)
	}

	err = db.Init()
	if err != nil {
		os.Exit(5)
	}

	err = job.Init()
	if err != nil {
		os.Exit(5)
	}

	// Ignore errors; 出错自动os.Exit(5)
	web.Run(ctx, api.Router)
}
