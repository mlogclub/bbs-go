package utils

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

var Conf *Config

type Config struct {
	Env        string // 环境：prod、dev
	BaseUrl    string // base url
	SiteTitle  string // 网站标题
	Port       string // 端口
	MySqlUrl   string // 数据库连接地址
	ShowSql    bool   // 是否显示日志
	ViewsPath  string // 模板路径
	StaticPath string // 静态文件路径

	// oauth
	OauthAuthUrl  string
	OauthTokenUrl string

	// oauth client
	OauthClientId          string // oauth2客户端编号
	OauthClientSecret      string // oauth2客户端秘钥
	OauthClientRedirectUrl string // oauth2客户端回调地址
	OauthClientAuthUrl     string // oauth2客户端的授权地址
	OauthClientSuccessUrl  string // oauth2客户端登录成功之后跳转到的页面地址

	// Redis
	RedisAddr string

	// Github
	GithubClientID     string
	GithubClientSecret string

	// 阿里云oss配置
	AliyunOssHost         string
	AliyunOssBucket       string
	AliyunOssEndpoint     string
	AliyunOssAccessId     string
	AliyunOssAccessSecret string

	// smtp
	SmtpAddr     string
	SmtpPort     string
	SmtpUsername string
	SmtpPassword string
}

func InitConfig(config string) {
	initConfig(config)
	initGithubConfig()
	initEmail()
}

func initConfig(config string) {
	bytes, err := ioutil.ReadFile(config)
	if err != nil {
		logrus.Error("ReadFile: ", err.Error())
		os.Exit(-1)
	}

	Conf = &Config{}
	if err := json.Unmarshal(bytes, &Conf); err != nil {
		logrus.Error("invalid config: ", err.Error())
		os.Exit(-1)
	}
}
