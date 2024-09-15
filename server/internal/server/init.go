package server

import (
	"bbs-go/internal/models"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/pkg/iplocator"
	"bbs-go/internal/pkg/search"
	"bbs-go/internal/scheduler"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func Init() {
	initConfig()
	initLogger()
	initDB()
	initCron()
	initIpLocator()
	initSearch()
}

func initConfig() {
	env := os.Getenv("BBSGO_ENV")
	if strs.IsBlank(env) {
		env = "dev"
	}

	slog.Info("Load config", slog.String("ENV", env))

	viper.SetConfigName("bbs-go." + env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.bbs-go")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../../")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("BBSGO")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := viper.Unmarshal(&config.Instance); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	config.Instance.Env = env

	slog.Info("Load config", slog.String("ENV", env), slog.String("config", jsons.ToJsonStr(config.Instance)))
}

func initDB() {
	conf := config.Instance.DB
	db, err := gorm.Open(mysql.Open(conf.Url), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		},
	})

	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
		sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
		sqlDB.SetConnMaxIdleTime(time.Duration(conf.ConnMaxIdleTimeSeconds) * time.Second)
		sqlDB.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetimeSeconds) * time.Second)
	}

	if err := db.AutoMigrate(models.Models...); nil != err {
		slog.Error(err.Error(), slog.Any("error", err))
	}

	sqls.SetDB(db)
}

func initCron() {
	if config.Instance.IsProd() {
		scheduler.Start()
	}
}

func initIpLocator() {
	iplocator.InitIpLocator(config.Instance.IpDataPath)
}

func initSearch() {
	search.Init(config.Instance.Search.IndexPath)
}
