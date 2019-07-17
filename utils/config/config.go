package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var Conf *Config

type Config struct {
	Env            string `yaml:"Env"`            // 环境：prod、dev
	BaseUrl        string `yaml:"BaseUrl"`        // base url
	Port           string `yaml:"Port"`           // 端口
	ShowSql        bool   `yaml:"ShowSql"`        // 是否显示日志
	ViewsPath      string `yaml:"ViewsPath"`      // 模板路径
	RootStaticPath string `yaml:"RootStaticPath"` // 根路径下的静态文件目录
	StaticPath     string `yaml:"StaticPath"`     // 静态文件路径

	MySqlUrl  string `yaml:"MySqlUrl"`  // 数据库连接地址
	RedisAddr string `yaml:"RedisAddr"` // redis

	// redis
	Redis struct {
		Addr     string `yaml:"Addr"`     // redis链接
		Password string `yaml:"Password"` // redis 密码
	} `yaml:"Redis"`

	// oauth server
	OauthServer struct {
		AuthUrl  string `yaml:"AuthUrl"`
		TokenUrl string `yaml:"TokenUrl"`
	} `yaml:"OauthServer"`

	// oauth client
	OauthClient struct {
		ClientId          string `yaml:"ClientId"`          // oauth2客户端编号
		ClientSecret      string `yaml:"ClientSecret"`      // oauth2客户端秘钥
		ClientRedirectUrl string `yaml:"ClientRedirectUrl"` // oauth2客户端回调地址
		ClientSuccessUrl  string `yaml:"ClientSuccessUrl"`  // oauth2客户端登录成功之后跳转到的页面地址
	} `yaml:"OauthClient"`

	// Github
	Github struct {
		ClientID     string `yaml:"ClientID"`
		ClientSecret string `yaml:"ClientSecret"`
	} `yaml:"Github"`

	// 阿里云oss配置
	AliyunOss struct {
		Host         string `yaml:"Host"`
		Bucket       string `yaml:"Bucket"`
		Endpoint     string `yaml:"Endpoint"`
		AccessId     string `yaml:"AccessId"`
		AccessSecret string `yaml:"AccessSecret"`
	} `yaml:"AliyunOss"`

	// smtp
	Smtp struct {
		Addr     string `yaml:"Addr"`
		Port     string `yaml:"Port"`
		Username string `yaml:"Username"`
		Password string `yaml:"Password"`
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
