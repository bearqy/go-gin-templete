package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAppliesDefaultsAndEnvOverrides(t *testing.T) {
	t.Setenv("GIN_TEMPLATE_WEB_ADDRESS", ":9090")
	t.Setenv("GIN_TEMPLATE_DB_MAX_OPEN_CONNS", "12")

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(`log:
  level: debug
`), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Web.Address != ":9090" {
		t.Fatalf("expected env address override, got %q", cfg.Web.Address)
	}
	if cfg.Log.Level != "debug" {
		t.Fatalf("expected yaml log level, got %q", cfg.Log.Level)
	}
	if cfg.DB.MaxOpenConns != 12 {
		t.Fatalf("expected env max open conns override, got %d", cfg.DB.MaxOpenConns)
	}
	if cfg.Web.ReadTimeoutSecond == 0 || cfg.Web.WriteTimeoutSecond == 0 {
		t.Fatal("expected default web timeouts")
	}
}

func TestValidateRejectsUnknownLogLevel(t *testing.T) {
	cfg := Default()
	cfg.Log.Level = "trace"

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error")
	}
}
