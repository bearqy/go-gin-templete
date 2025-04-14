package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go-gin-templete/internal/api"
	"go-gin-templete/internal/cli"
	"go-gin-templete/internal/config"
	"go-gin-templete/internal/db"
	"go-gin-templete/internal/job"
	"go-gin-templete/internal/logger"
	"go-gin-templete/internal/webserver"
)

func main() {

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	err := logger.Init()
	if err != nil {
		os.Exit(5)
	}

	cli.Init()

	err = config.Init()
	if err != nil {
		os.Exit(5)
	}

	err = logger.Init2(config.Main.Log.Level)
	if err != nil {
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

	fmt.Println("应用启动成功，请直接访问swagger界面。http://127.0.0.1" + config.Main.Web.Address + "/wiki")
	// Ignore errors; 出错自动os.Exit(5)
	webserver.Run(ctx, api.Router)
}
