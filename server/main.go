package main

import (
	"flag"
	"gorm.io/gorm"
	"io"
	"log"
	"os"
	"server/controllers"
	"server/model"
	"server/pkg/common"
	"server/pkg/config"
	"server/scheduler"
	_ "server/services/eventhandler"
	"time"

	"github.com/mlogclub/simple/sqls"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

var configFile = flag.String("config", "./server.yaml", "配置文件路径")

func init() {
	flag.Parse()

	// 初始化配置
	conf := config.Init(*configFile)
	//alipay.Init()

	// 初始化日志
	if file, err := os.OpenFile(conf.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
		logrus.SetOutput(io.MultiWriter(os.Stdout, file))
	} else {
		logrus.SetOutput(os.Stdout)
		logrus.Error(err)
	}

	// 连接数据库
	gormConf := &gorm.Config{
		Logger: logger.New(logrus.StandardLogger(), logger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  true,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
		}),
	}
	if err := sqls.Open(conf.DB.Url, gormConf, conf.DB.MaxIdleConns, conf.DB.MaxOpenConns, model.Models...); err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	if common.IsProd() {
		// 开启定时任务
		scheduler.Start()
	}
	controllers.Router()
}
