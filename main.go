package main

import (
	"flag"

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

	initLogrus()

	config.InitConfig(*configFile) // 初始化配置
	utils.InitEmail()              // 初始化邮件
	initDB()                       // 初始化数据库
}

func initLogrus() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
	logrus.SetLevel(logrus.InfoLevel)
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
