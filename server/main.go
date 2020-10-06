package main

import (
	"bbs-go/app"
	"bbs-go/config"
	"bbs-go/model"
	"flag"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var configFile = flag.String("config", "./bbs-go.yaml", "配置文件路径")

func init() {
	flag.Parse()

	// 初始化配置
	config.Init(*configFile)

	// 初始化日志
	if file, err := os.OpenFile(config.Instance.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err == nil {
		logrus.SetOutput(file)
	} else {
		logrus.Error(err)
	}

	// 连接数据库
	gormConf := &gorm.Config{Logger: logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			Colorful:      true,
			LogLevel:      logger.Info,
		}),
	}
	if err := simple.OpenDB(config.Instance.MySqlUrl, gormConf, 10, 20, model.Models...); err != nil {
		logrus.Error(err)
	}
}

func main() {
	app.StartOn()
	app.InitIris()
}
