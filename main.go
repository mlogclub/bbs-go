package main

import (
	"flag"
	"os"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/mlog/app"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/utils"
	"github.com/mlogclub/mlog/utils/config"
)

var configFile = flag.String("config", "./mlog.yaml", "配置文件路径")

func init() {
	flag.Parse()

	config.InitConfig(*configFile) // 初始化配置
	initLogrus()                   // 初始化日志
	utils.InitEmail()              // 初始化邮件
	initDB()                       // 初始化数据库
}

func initLogrus() {
	file, err := os.OpenFile(config.Conf.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		logrus.SetLevel(logrus.InfoLevel)
		logrus.SetOutput(file)
	} else {
		logrus.Error(err)
	}
}

func initDB() {
	// 连接数据库
	simple.OpenDB(&simple.DBConfiguration{
		Dialect:        "mysql",
		Url:            config.Conf.MySqlUrl,
		MaxIdle:        5,
		MaxActive:      20,
		EnableLogModel: config.Conf.ShowSql,
		Models:         model.Models,
	})
}

func main() {
	app.StartOn()
	app.InitIris()
}
