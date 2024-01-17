package server

import (
	"bbs-go/internal/models"
	"bbs-go/internal/pkg/common"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/pkg/iplocator"
	"bbs-go/internal/scheduler"
	"fmt"
	"io"
	"log/slog"
	"os"
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

func initLogger() {
	// conf := config.Instance
	// // 初始化日志
	// if file, err := os.OpenFile(conf.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
	// 	logrus.SetOutput(io.MultiWriter(os.Stdout, file))
	// } else {
	// 	logrus.SetOutput(os.Stdout)
	// 	logrus.Error(err)
	// }

	logger := slog.New(slog.NewTextHandler(io.MultiWriter(os.Stdout)))
	slog.SetDefault(logger)
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
		logrus.Error(err)
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
