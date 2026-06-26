package app

import (
	"context"
	"log/slog"
	"time"

	"go-gin-templete/internal/api"
	"go-gin-templete/internal/cli"
	"go-gin-templete/internal/config"
	"go-gin-templete/internal/db"
	"go-gin-templete/internal/job"
	"go-gin-templete/internal/logger"
	"go-gin-templete/internal/webserver"
)

func Run(ctx context.Context, args []string) error {
	logger.InitConsole()

	configPath, err := cli.Parse(args)
	if err != nil {
		return err
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}
	config.Main = cfg
	config.InitSwagger()

	logHandles, err := logger.Setup(cfg.Log)
	if err != nil {
		return err
	}
	defer closeWithLog("logger", logHandles.Close)

	store, err := db.InitWithConfig(cfg.DB)
	if err != nil {
		return err
	}
	defer closeWithLog("database", store.Close)

	job.RegisterDefaults()
	jobManager, err := job.NewManager(cfg.Job)
	if err != nil {
		return err
	}
	jobManager.Start()
	defer stopJobs(jobManager)

	slog.Info("application ready", "addr", cfg.Web.Address, "swagger", "http://127.0.0.1"+cfg.Web.Address+"/wiki")
	return webserver.Run(ctx, cfg.Web, logHandles.AccessWriter(), api.Router)
}

func stopJobs(manager *job.Manager) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := manager.Stop(ctx); err != nil {
		slog.Error("stop jobs error", slog.Any("error", err))
	}
}

func closeWithLog(name string, closeFn func() error) {
	if closeFn == nil {
		return
	}
	if err := closeFn(); err != nil {
		slog.Error("close resource error", "name", name, slog.Any("error", err))
	}
}
