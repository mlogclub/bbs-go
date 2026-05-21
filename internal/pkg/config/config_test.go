package config

import "testing"

func TestSetDbDefaultsSetsDefaultLogLevel(t *testing.T) {
	var cfg DBConfig

	SetDbDefaults(&cfg)

	if cfg.LogLevel != DefaultDBLogLevel {
		t.Fatalf("expected default db log level %q, got %q", DefaultDBLogLevel, cfg.LogLevel)
	}
}
