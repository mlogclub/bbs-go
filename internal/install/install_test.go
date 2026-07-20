package install

import (
	"bbs-go/internal/pkg/config"
	"net/url"
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

func TestDockerBuiltinDBType(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		t.Setenv(DockerBuiltinMySQLEnv, "")
		t.Setenv(DockerBuiltinPostgreSQLEnv, "")

		if got := DockerBuiltinDBType(); got != "" {
			t.Fatalf("expected empty builtin db type, got %q", got)
		}
		if IsDockerBuiltinDBInstall() {
			t.Fatal("expected builtin db install to be disabled")
		}
	})

	t.Run("mysql", func(t *testing.T) {
		t.Setenv(DockerBuiltinMySQLEnv, "true")
		t.Setenv(DockerBuiltinPostgreSQLEnv, "")

		if got := DockerBuiltinDBType(); got != config.DbTypeMySQL {
			t.Fatalf("expected %q, got %q", config.DbTypeMySQL, got)
		}
	})

	t.Run("postgresql takes precedence", func(t *testing.T) {
		t.Setenv(DockerBuiltinMySQLEnv, "true")
		t.Setenv(DockerBuiltinPostgreSQLEnv, "true")

		if got := DockerBuiltinDBType(); got != config.DbTypePostgreSQL {
			t.Fatalf("expected %q, got %q", config.DbTypePostgreSQL, got)
		}
	})
}

func TestDbConfigReqGetConnStrPostgreSQLEscapesValues(t *testing.T) {
	req := DbConfigReq{
		Type:     config.DbTypePostgreSQL,
		Host:     "localhost",
		Port:     "5432",
		Database: "bbs go",
		Username: "bbs user",
		Password: `pa ss'word\`,
	}

	dsn := req.GetConnStr()
	parsed, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("expected valid postgresql dsn, got %q: %v", dsn, err)
	}
	if parsed.Scheme != "postgres" {
		t.Fatalf("expected postgres scheme, got %q", parsed.Scheme)
	}
	if parsed.Hostname() != req.Host {
		t.Fatalf("expected host %q, got %q", req.Host, parsed.Hostname())
	}
	if parsed.Port() != req.Port {
		t.Fatalf("expected port %q, got %q", req.Port, parsed.Port())
	}
	if parsed.User.Username() != req.Username {
		t.Fatalf("expected username %q, got %q", req.Username, parsed.User.Username())
	}
	password, ok := parsed.User.Password()
	if !ok {
		t.Fatal("expected password in postgresql dsn")
	}
	if password != req.Password {
		t.Fatalf("expected password %q, got %q", req.Password, password)
	}
	if parsed.Path != "/"+req.Database {
		t.Fatalf("expected database path %q, got %q", "/"+req.Database, parsed.Path)
	}
	if parsed.Query().Get("sslmode") != "disable" {
		t.Fatalf("expected sslmode=disable, got query %q", parsed.RawQuery)
	}
}

func TestDbConfigReqGetConnStrPostgreSQLIPv6Host(t *testing.T) {
	req := DbConfigReq{
		Type:     config.DbTypePostgreSQL,
		Host:     "::1",
		Port:     "5432",
		Database: "bbsgo",
		Username: "bbsgo",
		Password: "secret",
	}

	parsed, err := url.Parse(req.GetConnStr())
	if err != nil {
		t.Fatalf("expected valid postgresql dsn: %v", err)
	}
	if parsed.Hostname() != req.Host {
		t.Fatalf("expected host %q, got %q", req.Host, parsed.Hostname())
	}
	if parsed.Port() != req.Port {
		t.Fatalf("expected port %q, got %q", req.Port, parsed.Port())
	}
}

func TestApplyDockerBuiltinPostgreSQLConfig(t *testing.T) {
	previousConfig := config.Instance
	t.Cleanup(func() {
		config.Instance = previousConfig
	})

	t.Setenv(DockerBuiltinPostgreSQLEnv, "true")
	t.Setenv(DockerBuiltinPostgreSQLHostEnv, "pg")
	t.Setenv(DockerBuiltinPostgreSQLPortEnv, "15432")
	t.Setenv(DockerBuiltinPostgreSQLDatabaseEnv, "bbsgo_test")
	t.Setenv(DockerBuiltinPostgreSQLUsernameEnv, "bbsgo_user")
	t.Setenv(DockerBuiltinPostgreSQLPasswordEnv, "bbsgo_secret")

	config.Instance = &config.Config{}
	ApplyDockerBuiltinPostgreSQLConfig()

	if got := config.Instance.DB.Type; got != config.DbTypePostgreSQL {
		t.Fatalf("expected db type %q, got %q", config.DbTypePostgreSQL, got)
	}
	wantDSN := "postgres://bbsgo_user:bbsgo_secret@pg:15432/bbsgo_test?sslmode=disable"
	if got := config.Instance.DB.Url; got != wantDSN {
		t.Fatalf("expected dsn %q, got %q", wantDSN, got)
	}
}

func TestRuntimeConfigForWriteClearsBuiltinDBUrl(t *testing.T) {
	cfg := &config.Config{
		DB: config.DBConfig{
			Type: config.DbTypePostgreSQL,
			Url:  "host=pg port=5432 user=bbsgo password=secret dbname=bbsgo sslmode=disable",
		},
	}

	t.Run("postgresql", func(t *testing.T) {
		t.Setenv(DockerBuiltinMySQLEnv, "")
		t.Setenv(DockerBuiltinPostgreSQLEnv, "true")

		got := runtimeConfigForWrite(cfg)
		if got == cfg {
			t.Fatal("expected builtin runtime config to be copied before write")
		}
		if got.DB.Type != config.DbTypePostgreSQL {
			t.Fatalf("expected db type %q, got %q", config.DbTypePostgreSQL, got.DB.Type)
		}
		if got.DB.Url != "" {
			t.Fatalf("expected persisted db url to be empty, got %q", got.DB.Url)
		}
		if cfg.DB.Url == "" {
			t.Fatal("expected runtime config db url to stay available")
		}
	})

	t.Run("mysql", func(t *testing.T) {
		t.Setenv(DockerBuiltinMySQLEnv, "true")
		t.Setenv(DockerBuiltinPostgreSQLEnv, "")

		got := runtimeConfigForWrite(cfg)
		if got.DB.Type != config.DbTypeMySQL {
			t.Fatalf("expected db type %q, got %q", config.DbTypeMySQL, got.DB.Type)
		}
		if got.DB.Url != "" {
			t.Fatalf("expected persisted db url to be empty, got %q", got.DB.Url)
		}
	})

	t.Run("none", func(t *testing.T) {
		t.Setenv(DockerBuiltinMySQLEnv, "")
		t.Setenv(DockerBuiltinPostgreSQLEnv, "")

		got := runtimeConfigForWrite(cfg)
		if got != cfg {
			t.Fatal("expected non-builtin runtime config to be written directly")
		}
	})
}
