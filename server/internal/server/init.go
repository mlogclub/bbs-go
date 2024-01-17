package server

import (
	"bbs-go/internal/models"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/pkg/iplocator"
	"bbs-go/internal/scheduler"
	"fmt"
	"log/slog"
	"time"

	"github.com/mlogclub/simple/sqls"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init() {
	initConfig()
	initLogger()
	initDB()
	initCron()
	initIpLocator()
}

func initConfig() {
	viper.SetConfigName("bbs-go")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.bbs-go")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../../")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if err := viper.Unmarshal(&config.Instance); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func initDB() {
	// 连接数据库
	gormConf := &gorm.Config{
		Logger: logger.New(logrus.StandardLogger(), logger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  true,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
		}),
	}
	if err := sqls.Open(config.Instance.DB, gormConf, models.Models...); err != nil {
		slog.Error(err.Error())
	}
}

func initCron() {
	if common.IsProd() {
		// 开启定时任务
		scheduler.Start()
	}
}

func initIpLocator() {
	iplocator.InitIpLocator(config.Instance.IpDataPath)
}
