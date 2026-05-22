package install

import (
	"testing"

	"gorm.io/gorm/logger"
)

func TestResolveGormLogLevel(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  logger.LogLevel
	}{
		{name: "empty defaults to info", input: "", want: logger.Info},
		{name: "silent", input: "silent", want: logger.Silent},
		{name: "error", input: "error", want: logger.Error},
		{name: "warn", input: "warn", want: logger.Warn},
		{name: "info", input: "info", want: logger.Info},
		{name: "case insensitive", input: "WARN", want: logger.Warn},
		{name: "invalid defaults to info", input: "unknown", want: logger.Info},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveGormLogLevel(tt.input); got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
