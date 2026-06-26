package job

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"go-gin-templete/internal/config"

	"github.com/robfig/cron/v3"
)

var jobMap = make(map[string]func())

type Manager struct {
	cron *cron.Cron
}

func Init() error {
	manager, err := NewManager(config.Main.Job)
	if err != nil {
		return err
	}
	manager.Start()
	return nil
}

func NewManager(jobs map[string]config.JobConfig) (*Manager, error) {
	c := cron.New(cron.WithSeconds(), cron.WithChain(cron.Recover(cron.DefaultLogger)))

	for name, jobConfig := range jobs {
		spec := strings.TrimSpace(jobConfig.Cron)
		if spec == "" {
			return nil, fmt.Errorf("job %q cron expression is empty", name)
		}
		fn, ok := jobMap[name]
		if !ok || fn == nil {
			return nil, fmt.Errorf("job %q is configured but not registered", name)
		}
		_, err := c.AddFunc(spec, namedJob(name, fn))
		if err != nil {
			slog.Error("AddFunc error", slog.Any("error", err))
			return nil, err
		}
		slog.Info("job registered", "name", name, "cron", spec)
	}
	return &Manager{cron: c}, nil
}

func (m *Manager) Start() {
	if m == nil || m.cron == nil {
		return
	}
	m.cron.Start()
}

func (m *Manager) Stop(ctx context.Context) error {
	if m == nil || m.cron == nil {
		return nil
	}
	stopCtx := m.cron.Stop()
	select {
	case <-stopCtx.Done():
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func Register(name string, fn func()) {
	jobMap[name] = fn
}

func namedJob(name string, fn func()) func() {
	return func() {
		slog.Info("job started", "name", name)
		fn()
		slog.Info("job finished", "name", name)
	}
}
