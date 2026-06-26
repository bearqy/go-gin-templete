package job

import (
	"strings"
	"testing"

	"go-gin-templete/internal/config"
)

func TestNewManagerRejectsUnregisteredJob(t *testing.T) {
	_, err := NewManager(map[string]config.JobConfig{
		"missing": {Cron: "* * * * * *"},
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "not registered") {
		t.Fatalf("expected not registered error, got %v", err)
	}
}

func TestNewManagerAcceptsRegisteredSecondsCron(t *testing.T) {
	Register("test_seconds_cron", func() {})

	manager, err := NewManager(map[string]config.JobConfig{
		"test_seconds_cron": {Cron: "*/5 * * * * *"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if manager == nil {
		t.Fatal("expected manager")
	}
}
