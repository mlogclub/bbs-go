package config

import (
	"strings"

	"github.com/mlogclub/simple/sqls"
)

var Instance *Config

type Config struct {
	Env        string // 环境
	BaseUrl    string // base url
	Port       string // 端口
	IpDataPath string // IP数据文件

	// 日志配置
	Logger struct {
		Filename   string // 日志文件的位置
		MaxSize    int    // 文件最大尺寸（以MB为单位）
		MaxAge     int    // 保留旧文件的最大天数
		MaxBackups int    // 保留的最大旧文件数量
	}

	// 数据库配置
	DB sqls.DbConfig

	// 阿里云oss配置
	Uploader struct {
		Enable    string
		AliyunOss struct {
			Host          string
			Bucket        string
			Endpoint      string
			AccessId      string
			AccessSecret  string
			StyleSplitter string
			StyleAvatar   string
			StylePreview  string
			StyleSmall    string
			StyleDetail   string
		}
		Local struct {
			Host string
			Path string
		}
	}

	// 百度SEO相关配置
	// 文档：https://ziyuan.baidu.com/college/courseinfo?id=267&page=2#h2_article_title14
	BaiduSEO struct {
		Site  string
		Token string
	}

	// 神马搜索SEO相关
	// 文档：https://zhanzhang.sm.cn/open/mip
	SmSEO struct {
		Site     string
		UserName string
		Token    string
	}

	// smtp
	Smtp struct {
		Host     string
		Port     string
		Username string
		Password string
		SSL      bool
	}

	Search struct {
		IndexPath string
	}
}

func (c *Config) IsProd() bool {
	e := strings.ToLower(c.Env)
	return e == "prod" || e == "production"
}
