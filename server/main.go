package main

import (
	"flag"
	"os"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"

	"bbs-go/app"
	"bbs-go/common/config"
	"bbs-go/model"
	"reflect"
)

var configFile = flag.String("config", "./bbs-go.yaml", "配置文件路径")

func init() {
	flag.Parse()

	config.InitConfig(*configFile)
	v := reflect.ValueOf(config.Conf).Elem()
	t := reflect.TypeOf(*config.Conf)
	for i := 0; i < t.NumField(); i++ {
		fieldName := t.Field(i).Name
		if os.Getenv(fieldName) != "" {
			env := os.Getenv("Env")
			v.Field(i).SetString(env)
		}

	} // 初始化配置
	// 初始化日志
	err := simple.OpenMySql(config.Conf.MySqlUrl, 10, 20, config.Conf.ShowSql, model.Models...) // 连接数据库
	if err != nil {
		logrus.Error(err)
	}
}

func initLogrus() {
	file, err := os.OpenFile(config.Conf.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logrus.SetOutput(file)
	} else {
		logrus.Error(err)
	}
}

func main() {
	app.StartOn()
	app.InitIris()
}
