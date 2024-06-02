package server

import (
	"bbs-go/internal/models"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/pkg/gormlogs"
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
	"gorm.io/gorm"
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
	env := os.Getenv("ENV")
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
	gormConf := &gorm.Config{
		Logger: gormlogs.New(
			gormlogs.SetLogLevel(gormlogs.SlowQueryLogType, slog.LevelWarn),
			gormlogs.WithSlowThreshold(time.Second),
		),
	}
	if err := sqls.Open(config.Instance.DB, gormConf, models.Models...); err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
	}
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
