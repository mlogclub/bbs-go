package config

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var Conf *Config

type Config struct {
	Env        string `yaml:"Env"`        // 环境：prod、dev
	BaseUrl    string `yaml:"BaseUrl"`    // base url
	Port       string `yaml:"Port"`       // 端口
	LogFile    string `yaml:"LogFile"`    // 日志文件
	ShowSql    bool   `yaml:"ShowSql"`    // 是否显示日志
	StaticPath string `yaml:"StaticPath"` // 静态文件目录

	MySqlUrl string `yaml:"MySqlUrl"` // 数据库连接地址

	// Github
	Github struct {
		ClientID     string `yaml:"ClientID"`
		ClientSecret string `yaml:"ClientSecret"`
	} `yaml:"Github"`

	// QQ登录
	QQConnect struct {
		AppId  string `yaml:"AppId"`
		AppKey string `yaml:"AppKey"`
	} `yaml:"QQConnect"`

	// 阿里云oss配置
	AliyunOss struct {
		Host         string `yaml:"Host"`
		Bucket       string `yaml:"Bucket"`
		Endpoint     string `yaml:"Endpoint"`
		AccessId     string `yaml:"AccessId"`
		AccessSecret string `yaml:"AccessSecret"`
	} `yaml:"AliyunOss"`

	// 百度ai
	BaiduAi struct {
		ApiKey    string `yaml:"ApiKey"`
		SecretKey string `yaml:"SecretKey"`
	} `yaml:"BaiduAi"`

	// 百度SEO相关配置
	// 文档：https://ziyuan.baidu.com/college/courseinfo?id=267&page=2#h2_article_title14
	BaiduSEO struct {
		Site  string `yaml:"Site"`
		Token string `yaml:"Token"`
	} `yaml:"BaiduSEO"`

	// smtp
	Smtp struct {
		Host     string `yaml:"Host"`
		Port     string `yaml:"Port"`
		Username string `yaml:"Username"`
		Password string `yaml:"Password"`
		SSL      bool   `yaml:"SSL"`
	} `yaml:"Smtp"`
}

func InitConfig(filename string) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		logrus.Error(err)
		return
	}

	Conf = &Config{}
	err = yaml.Unmarshal(yamlFile, Conf)
	if err != nil {
		logrus.Error(err)
	}
}
