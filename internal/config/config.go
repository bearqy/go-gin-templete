package config

import (
	"fmt"
	"go-gin-templete/docs"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"go-gin-templete/internal/cli"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Log LogConfig            `yaml:"log"`
	Web WebConfig            `yaml:"web"`
	DB  DBConfig             `yaml:"db"`
	Job map[string]JobConfig `yaml:"job"`
}

type LogConfig struct {
	Level          string `yaml:"level"`
	AccessLogfile  string `yaml:"accessLogfile"`
	RuntimeLogfile string `yaml:"runtimeLogfile"`
}

type WebConfig struct {
	Address            string `yaml:"address"`
	ReadTimeoutSecond  int    `yaml:"read_timeout_second"`
	WriteTimeoutSecond int    `yaml:"write_timeout_second"`
	IdleTimeoutSecond  int    `yaml:"idle_timeout_second"`
}

type DBConfig struct {
	ConnStr               string `yaml:"conn_str"`
	MaxOpenConns          int    `yaml:"max_open_conns"`
	MaxIdleConns          int    `yaml:"max_idle_conns"`
	ConnMaxLifetimeSecond int    `yaml:"conn_max_lifetime_second"`
}

type JobConfig struct {
	Cron string `yaml:"cron"`
}

var Main = &Config{}

func Init() error {
	cfg, err := Load(cli.ConfigFilePath)
	if err != nil {
		return err
	}
	Main = cfg
	return nil
}

func Load(configFilePath string) (*Config, error) {
	cfg := Default()

	yamlData, err := os.ReadFile(configFilePath)
	if err != nil {
		slog.Error("read config file error", slog.Any("error", err))
		return nil, err
	}

	if err := yaml.Unmarshal(yamlData, cfg); err != nil {
		slog.Error("yaml unmarshal error", slog.Any("error", err))
		return nil, err
	}

	applyEnv(cfg)
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	slog.Info("load config successfully", "path", configFilePath, "address", cfg.Web.Address, "log_level", cfg.Log.Level, "db_enabled", cfg.DB.ConnStr != "")
	return cfg, nil
}

func Default() *Config {
	return &Config{
		Log: LogConfig{
			Level:          "info",
			AccessLogfile:  "log/access.log",
			RuntimeLogfile: "log/runtime.log",
		},
		Web: WebConfig{
			Address:            ":8080",
			ReadTimeoutSecond:  10,
			WriteTimeoutSecond: 30,
			IdleTimeoutSecond:  60,
		},
		DB: DBConfig{
			MaxOpenConns:          25,
			MaxIdleConns:          5,
			ConnMaxLifetimeSecond: 300,
		},
		Job: map[string]JobConfig{},
	}
}

func (c *Config) Validate() error {
	if strings.TrimSpace(c.Web.Address) == "" {
		return fmt.Errorf("web.address is required")
	}
	if _, ok := ParseLogLevel(c.Log.Level); !ok {
		return fmt.Errorf("unknown log.level %q", c.Log.Level)
	}
	if c.Web.ReadTimeoutSecond <= 0 || c.Web.WriteTimeoutSecond <= 0 || c.Web.IdleTimeoutSecond <= 0 {
		return fmt.Errorf("web timeouts must be greater than zero")
	}
	return nil
}

func ParseLogLevel(level string) (slog.Level, bool) {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return slog.LevelDebug, true
	case "info", "":
		return slog.LevelInfo, true
	case "warn":
		return slog.LevelWarn, true
	case "error":
		return slog.LevelError, true
	default:
		return slog.LevelInfo, false
	}
}

func applyEnv(cfg *Config) {
	overrideString("GIN_TEMPLATE_WEB_ADDRESS", &cfg.Web.Address)
	overrideString("GIN_TEMPLATE_LOG_LEVEL", &cfg.Log.Level)
	overrideString("GIN_TEMPLATE_ACCESS_LOG_FILE", &cfg.Log.AccessLogfile)
	overrideString("GIN_TEMPLATE_RUNTIME_LOG_FILE", &cfg.Log.RuntimeLogfile)
	overrideString("GIN_TEMPLATE_DB_CONN_STR", &cfg.DB.ConnStr)
	overrideInt("GIN_TEMPLATE_DB_MAX_OPEN_CONNS", &cfg.DB.MaxOpenConns)
	overrideInt("GIN_TEMPLATE_DB_MAX_IDLE_CONNS", &cfg.DB.MaxIdleConns)
	overrideInt("GIN_TEMPLATE_DB_CONN_MAX_LIFETIME_SECOND", &cfg.DB.ConnMaxLifetimeSecond)
}

func overrideString(key string, target *string) {
	if value, ok := os.LookupEnv(key); ok {
		*target = value
	}
}

func overrideInt(key string, target *int) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		slog.Warn("ignore invalid integer environment variable", "key", key, "value", value)
		return
	}
	*target = parsed
}

func InitSwagger() {
	docs.SwaggerInfo.Title = "go-gin-templete"
	docs.SwaggerInfo.Version = "v0.1.0"
	docs.SwaggerInfo.Description = "Gin web service scaffold"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.BasePath = ""
	slog.Info("swagger config successfully")
}
