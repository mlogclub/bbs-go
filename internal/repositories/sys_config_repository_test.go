package repositories

import (
	"strings"
	"testing"

	"bbs-go/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func TestSysConfigRepositoryGetByKey(t *testing.T) {
	db := openSysConfigSQLiteDB(t)

	cfg := &models.SysConfig{
		Key:         "siteTitle",
		Value:       "BBS-GO",
		Name:        "Site Title",
		Description: "Site Title",
		CreateTime:  1,
		UpdateTime:  1,
	}
	if err := db.Create(cfg).Error; err != nil {
		t.Fatalf("create sys config: %v", err)
	}

	got := SysConfigRepository.GetByKey(db, "siteTitle")
	if got == nil {
		t.Fatal("expected sys config, got nil")
	}
	if got.Value != cfg.Value {
		t.Fatalf("expected value %q, got %q", cfg.Value, got.Value)
	}
	if got := SysConfigRepository.GetByKey(db, ""); got != nil {
		t.Fatalf("expected nil for empty key, got %#v", got)
	}
}

func TestSysConfigRepositoryGetByKeySQLUsesDialectQuoting(t *testing.T) {
	tests := []struct {
		name       string
		db         *gorm.DB
		wantClause string
		badClause  string
	}{
		{
			name:       "mysql",
			db:         openSysConfigDryRunDB(t, mysql.New(mysql.Config{DSN: "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local", SkipInitializeWithVersion: true})),
			wantClause: "`key` = 'siteTitle'",
			badClause:  `"key" = 'siteTitle'`,
		},
		{
			name:       "postgres",
			db:         openSysConfigDryRunDB(t, postgres.Open("host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable")),
			wantClause: `"key" = 'siteTitle'`,
			badClause:  "`key` = 'siteTitle'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := tt.db.ToSQL(func(tx *gorm.DB) *gorm.DB {
				return tx.Take(&models.SysConfig{}, &models.SysConfig{Key: "siteTitle"})
			})
			if !strings.Contains(sql, tt.wantClause) {
				t.Fatalf("expected SQL %q to contain %q", sql, tt.wantClause)
			}
			if strings.Contains(sql, tt.badClause) {
				t.Fatalf("expected SQL %q not to contain %q", sql, tt.badClause)
			}
		})
	}
}

func openSysConfigSQLiteDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), sysConfigGormConfig())
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}
	if err := db.AutoMigrate(&models.SysConfig{}); err != nil {
		t.Fatalf("auto migrate sys config: %v", err)
	}
	return db
}

func openSysConfigDryRunDB(t *testing.T, dialector gorm.Dialector) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(dialector, sysConfigGormConfig())
	if err != nil {
		t.Fatalf("open dry run db: %v", err)
	}
	return db.Session(&gorm.Session{DryRun: true})
}

func sysConfigGormConfig() *gorm.Config {
	return &gorm.Config{
		DisableAutomaticPing: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		},
	}
}
