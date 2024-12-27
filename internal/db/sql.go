package db

import (
	"log/slog"

	"github.com/bearqy/go-gin-templete/internal/config"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var Engine *xorm.Engine

func Init() error {
	var err error
	Engine, err = xorm.NewEngine("mysql", config.Main.DB.ConnStr)
	if err != nil {
		slog.Error("init xorm error", slog.Any("error", err))
		return err
	}
	return nil
}
