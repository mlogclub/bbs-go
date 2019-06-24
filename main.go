package main

import (
	"flag"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/mlogclub/mlog/app"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/utils"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
)

var config = flag.String("config", "./mlog.json", "配置文件路径")

func init() {
	flag.Parse()

	initLogrus()

	utils.InitConfig(*config)  // 初始化配置
	utils.InitSessionManager() // 初始化sessionManager
	initDB()                   // 初始化数据库

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
		Url:            utils.Conf.MySqlUrl,
		MaxIdle:        5,
		MaxActive:      20,
		EnableLogModel: utils.Conf.ShowSql,
		Models:         model.Models,
	})
}

func main() {
	app.StartOn()
	app.InitIris()
}
