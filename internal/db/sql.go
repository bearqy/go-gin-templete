package db

import (
	"context"
	"go-gin-templete/internal/config"
	"go-gin-templete/internal/entity"
	"log/slog"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

var Engine *xorm.Engine

type Store struct {
	Engine *xorm.Engine
}

func Init() error {
	store, err := InitWithConfig(config.Main.DB)
	if err != nil {
		return err
	}
	if store != nil {
		Engine = store.Engine
	}
	return nil
}

func InitWithConfig(dbConfig config.DBConfig) (*Store, error) {
	if dbConfig.ConnStr == "" {
		slog.Info("database connection string is empty, skip database initialization")
		return nil, nil
	}
	engine, err := xorm.NewEngine("mysql", dbConfig.ConnStr)
	if err != nil {
		slog.Error("init xorm error", slog.Any("error", err))
		return nil, err
	}
	engine.SetMaxOpenConns(dbConfig.MaxOpenConns)
	engine.SetMaxIdleConns(dbConfig.MaxIdleConns)
	engine.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetimeSecond) * time.Second)

	// 使用 context 设置 ping 超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// 验证数据库连接是否真的建立成功
	if err := engine.PingContext(ctx); err != nil {
		slog.Error("database connection test failed", slog.Any("error", err))
		engine.Close()
		return nil, err
	}

	// 在控制台打印生成的SQL语句
	// engine.ShowSQL(true)

	// 初始化数据库并自动创建表
	if err := InitializeDB(engine); err != nil {
		engine.Close()
		return nil, err
	}

	Engine = engine
	slog.Info("init xorm success")
	return &Store{Engine: engine}, nil
}

// InitializeDB 初始化数据库并自动创建表
func InitializeDB(engine *xorm.Engine) error {
	if engine == nil {
		return nil
	}
	// 检查并自动创建表
	err := engine.Sync2(
		new(entity.User),
	)
	if err != nil {
		slog.Error("failed to sync database tables", slog.Any("error", err))
		return err
	}
	slog.Info("database tables synced successfully")
	return nil
}

func (s *Store) Close() error {
	if s == nil || s.Engine == nil {
		return nil
	}
	return s.Engine.Close()
}
