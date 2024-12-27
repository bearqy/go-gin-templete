package config

import (
	"log/slog"
	"os"

	"github.com/bearqy/go-gin-templete/internal/cli"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`
	Web struct {
		Address string `yaml:"address"`
	} `yaml:"web"`
	DB struct {
		ConnStr string `yaml:"conn_str"`
	} `yaml:"db"`
	Job map[string]struct {
		Cron string `yaml:"cron"`
	} `yaml:"job"`
}

var Main = &Config{}

func Init() error {
	yamlData, err := os.ReadFile(cli.ConfigFilePath)
	if err != nil {
		slog.Error("read config file error", slog.Any("error", err))
		return err
	}

	err = yaml.Unmarshal(yamlData, Main)
	if err != nil {
		slog.Error("yaml unmarshal error", slog.Any("error", err))
		return err
	}
	slog.Info("load config successfully", slog.Any("config", Main))
	return nil
}
