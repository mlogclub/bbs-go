package main

import (
	"flag"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"github.com/mlogclub/bbs-go/app"
	"github.com/mlogclub/bbs-go/common/config"
	"github.com/mlogclub/bbs-go/model"
)

var configFile = flag.String("config", "./bbs-go.yaml", "配置文件路径")

func init() {
	flag.Parse()

	config.InitConfig(*configFile)                                                          // 初始化配置
	initLogrus()                                                                            // 初始化日志
	err := simple.OpenDB(config.Conf.MySqlUrl, 5, 20, config.Conf.ShowSql, model.Models...) // 连接数据库
	if err != nil {
		logrus.Error(err)
	}
}

func initLogrus() {
	// output, err := simple.NewLogWriter(config.Conf.LogFile)
	// if err == nil {
	// 	logrus.SetLevel(logrus.InfoLevel)
	// 	logrus.SetOutput(output)
	// } else {
	// 	logrus.Error(err)
	// }
}

func main() {
	app.StartOn()
	app.InitIris()
}
