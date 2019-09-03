package main

import (
	"flag"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/mlog/app"
	"github.com/mlogclub/mlog/common/config"
	"github.com/mlogclub/mlog/model"
)

var configFile = flag.String("config", "./mlog.yaml", "配置文件路径")

func init() {
	flag.Parse()

	config.InitConfig(*configFile) // 初始化配置
	initLogrus()                   // 初始化日志
	initDB()                       // 初始化数据库
}

func initLogrus() {
	output, err := simple.NewLogWriter(config.Conf.LogFile)
	if err == nil {
		logrus.SetLevel(logrus.InfoLevel)
		logrus.SetOutput(output)
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
