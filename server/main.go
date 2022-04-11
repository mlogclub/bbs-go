package main

import (
	"bbs-go/controllers"
	"bbs-go/model"
	"bbs-go/pkg/common"
	"bbs-go/pkg/config"
	"bbs-go/scheduler"
	_ "bbs-go/services/eventhandler"
	"flag"
	"io"
	"log"
	"os"
	"time"

	"github.com/mlogclub/simple/sqls"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var configFile = flag.String("config", "./bbs-go.yaml", "配置文件路径")

func init() {
	flag.Parse()

	// 初始化配置
	conf := config.Init(*configFile)

	// gorm配置
	gormConf := &gorm.Config{}

	// 初始化日志
	if file, err := os.OpenFile(conf.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
		logrus.SetOutput(io.MultiWriter(os.Stdout, file))
		if conf.ShowSql {
			gormConf.Logger = logger.New(log.New(file, "\r\n", log.LstdFlags), logger.Config{
				SlowThreshold: time.Second,
				Colorful:      true,
				LogLevel:      logger.Info,
			})
		}
	} else {
		logrus.SetOutput(os.Stdout)
		logrus.Error(err)
	}

	// 连接数据库
	if err := sqls.Open(conf.DB.Url, gormConf, conf.DB.MaxIdleConns, conf.DB.MaxOpenConns, model.Models...); err != nil {
		logrus.Error(err)
	}
}

func main() {
	if common.IsProd() {
		// 开启定时任务
		scheduler.Start()
	}
	controllers.Router()
}
