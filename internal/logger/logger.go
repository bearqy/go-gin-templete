package logger

import (
	"go-gin-templete/internal/config"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

var AccessFile *os.File
var RuntimeFile *os.File
var initLogLevel = slog.LevelDebug

type Handles struct {
	AccessFile  *os.File
	RuntimeFile *os.File
}

func InitConsole() {
	l := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource:   true,
		Level:       initLogLevel,
		ReplaceAttr: nil,
	}))

	slog.SetDefault(l)
}

func Init() error {
	InitConsole()
	return nil
}

func Init2(logLevel string) error {
	handles, err := Setup(config.Main.Log)
	if err != nil {
		return err
	}
	AccessFile = handles.AccessFile
	RuntimeFile = handles.RuntimeFile
	return nil
}

func Setup(logConfig config.LogConfig) (*Handles, error) {
	var err error
	handles := &Handles{}
	handles.AccessFile, err = createLogFile(logConfig.AccessLogfile)
	if err != nil {
		return nil, err
	}
	handles.RuntimeFile, err = createLogFile(logConfig.RuntimeLogfile)
	if err != nil {
		handles.Close()
		return nil, err
	}

	l := slog.New(slog.NewTextHandler(handles.RuntimeWriter(), &slog.HandlerOptions{
		AddSource:   true,
		Level:       Level2Level(logConfig.Level),
		ReplaceAttr: nil,
	}))

	slog.SetDefault(l)

	return handles, nil
}

func (h *Handles) AccessWriter() io.Writer {
	if h == nil || h.AccessFile == nil {
		return os.Stdout
	}
	return h.AccessFile
}

func (h *Handles) RuntimeWriter() io.Writer {
	if h == nil || h.RuntimeFile == nil {
		return os.Stderr
	}
	return h.RuntimeFile
}

func (h *Handles) Close() error {
	if h == nil {
		return nil
	}
	var closeErr error
	if h.AccessFile != nil {
		closeErr = h.AccessFile.Close()
	}
	if h.RuntimeFile != nil && h.RuntimeFile != h.AccessFile {
		if err := h.RuntimeFile.Close(); err != nil && closeErr == nil {
			closeErr = err
		}
	}
	return closeErr
}

func createLogFile(file string) (*os.File, error) {
	if strings.TrimSpace(file) == "" {
		return nil, nil
	}

	//确保存放日志文件的目录始终存在
	logPathDir := filepath.Dir(file)
	if logPathDir != "." {
		if err := os.MkdirAll(logPathDir, 0755); err != nil {
			slog.Error("创建日志目录失败", slog.Any("error", err))
			return nil, err
		}
	}
	if logPathDir == "." {
		slog.Debug("日志路径为当前目录")
	} else {
		slog.Debug("日志路径为", slog.Any("logPathDir", logPathDir))
	}

	file1, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		slog.Error("创建日志目录失败", slog.Any("error", err))
		return nil, err
	}
	return file1, err
}

func Level2Level(level string) slog.Level {
	slog.Debug("日志等级", slog.Any("level", level))
	if parsed, ok := config.ParseLogLevel(level); ok {
		return parsed
	}
	slog.Error("unknown log level: " + strings.ToUpper(level) + " 现在用的是info类型")
	return slog.LevelInfo
}

func Close() error {
	return (&Handles{AccessFile: AccessFile, RuntimeFile: RuntimeFile}).Close()
}
