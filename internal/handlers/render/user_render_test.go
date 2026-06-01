package render

import (
	"fmt"
	"os"
	"testing"
	"time"

	"bbs-go/internal/models"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/pkg/idcodec"

	"github.com/glebarez/sqlite"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func TestMain(m *testing.M) {
	config.Instance = &config.Config{Language: config.LanguageEnUS}
	idcodec.Init(1)

	dsn := fmt.Sprintf("file:render_user_test_%d?mode=memory&cache=shared&_fk=1", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqls.SetDB(db)
	if err := db.AutoMigrate(&models.LevelConfig{}); err != nil {
		panic(err)
	}

	code := m.Run()
	_ = sqlDB.Close()
	os.Exit(code)
}

func TestBuildUserProfileKeepsEmptyDescriptionRaw(t *testing.T) {
	user := &models.User{}

	profile := BuildUserProfile(user)

	if profile == nil {
		t.Fatal("expected profile")
	}
	if profile.Description != "" {
		t.Fatalf("expected raw empty description, got %q", profile.Description)
	}
}

func TestBuildUserDetailUsesDefaultDescriptionForDisplay(t *testing.T) {
	user := &models.User{}

	detail := BuildUserDetail(user)

	if detail == nil {
		t.Fatal("expected detail")
	}
	if detail.Description == "" {
		t.Fatal("expected display description fallback")
	}
}
