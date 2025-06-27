package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	BBSGO_ENV  = "BBSGO_ENV"
	ENV_PREFIX = "BBSGO"

	EnvDev  = "dev"
	EnvTest = "test"
	EnvProd = "prod"
)

type Language string

const (
	LanguageZhCN Language = "zh-CN"
	LanguageEnUS Language = "en-US"

	DefaultLanguage = LanguageEnUS
)

var (
	Instance   *Config
	v          *viper.Viper
	configFile string
	writeMx    sync.Mutex
)

func init() {
	var (
		configFileName = "bbs-go.yaml"
	)
	v = viper.New()
	v.SetConfigFile(configFileName)
	v.AddConfigPath(".")
	if workDir, err := os.Executable(); err == nil {
		v.AddConfigPath(filepath.Dir(workDir))
	}
	v.AutomaticEnv()
	v.SetEnvPrefix(ENV_PREFIX)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	configFile = getConfigFilePath(configFileName)
}

type Config struct {
	Language       Language       `yaml:"language"`       // 语言
	BaseURL        string         `yaml:"baseURL"`        // baseURL
	Port           int            `yaml:"port"`           // 端口
	IpDataPath     string         `yaml:"ipDataPath"`     // IP数据文件
	AllowedOrigins []string       `yaml:"allowedOrigins"` // 跨域白名单
	Installed      bool           `yaml:"installed"`      // 是否已安装
	Logger         LoggerConfig   `yaml:"logger"`         // 日志配置
	DB             sqls.DbConfig  `yaml:"db"`             // 数据库配置
	Smtp           SmtpConfig     `yaml:"smtp"`           // smtp
	Search         SearchConfig   `yaml:"search"`         // 搜索配置
	BaiduSEO       BaiduSEOConfig `yaml:"baiduSEO"`       // 百度SEO配置
	SmSEO          SmSEOConfig    `yaml:"smSEO"`          // 神马搜索SEO配置
}

type LoggerConfig struct {
	Filename   string `yaml:"filename"`   // 日志文件的位置
	MaxSize    int    `yaml:"maxSize"`    // 文件最大尺寸（以MB为单位）
	MaxAge     int    `yaml:"maxAge"`     // 保留旧文件的最大天数
	MaxBackups int    `yaml:"maxBackups"` // 保留的最大旧文件数量
}

type SmtpConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	SSL      bool   `yaml:"ssl"`
}

type SearchConfig struct {
	IndexPath string `yaml:"indexPath"`
}

// 百度SEO配置
// 文档：https://ziyuan.baidu.com/college/courseinfo?id=267&page=2#h2_article_title14
type BaiduSEOConfig struct {
	Site  string `yaml:"site"`
	Token string `yaml:"token"`
}

// 神马搜索SEO配置
// 文档：https://zhanzhang.sm.cn/open/mip
type SmSEOConfig struct {
	Site     string `yaml:"site"`
	UserName string `yaml:"userName"`
	Token    string `yaml:"token"`
}

func ReadConfig() (cfg *Config, exists bool, err error) {
	exists = true
	if e := v.ReadInConfig(); e != nil {
		exists = false
		slog.Warn("Config file not found, use default", slog.Any("error", e))
	}

	if exists {
		if e := v.Unmarshal(&cfg); e != nil {
			err = fmt.Errorf("fatal error unmarshal config: %w", err)
			return
		}
		// 如果配置文件存在但没有语言设置，使用默认语言
		if strs.IsBlank(string(cfg.Language)) {
			cfg.Language = DefaultLanguage
		}
	} else {
		// default config
		cfg = &Config{
			Language:  DefaultLanguage,
			Port:      8082,
			Installed: false,
			Logger: LoggerConfig{
				Filename:   getLogFilename(),
				MaxSize:    10,
				MaxAge:     10,
				MaxBackups: 10,
			},
			DB: sqls.DbConfig{
				MaxIdleConns:           50,
				MaxOpenConns:           200,
				ConnMaxIdleTimeSeconds: 300,
				ConnMaxLifetimeSeconds: 3600,
			},
		}
	}

	slog.Info("Load config", slog.String("ENV", GetEnv()))
	return cfg, exists, nil
}

func WriteConfig(cfg *Config) error {
	if !writeMx.TryLock() {
		return errors.New("config is being written, please try again later")
	}
	defer writeMx.Unlock()

	yamlData, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFile, yamlData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func IsProd() bool {
	e := strings.ToLower(GetEnv())
	return e == "prod" || e == "production"
}

func GetEnv() string {
	env := os.Getenv("BBSGO_ENV")
	if strs.IsBlank(env) {
		env = EnvDev
	}
	return env
}

func getConfigFilePath(configName string) string {
	workDir, err := os.Executable()
	if err != nil {
		slog.Error("Failed to get working directory", slog.Any("error", err))
		return ""
	}
	return filepath.Join(filepath.Dir(workDir), configName)
}

func getLogFilename() string {
	// workDir, err := os.Getwd()
	// if err != nil {
	// 	slog.Error("Failed to get working directory", slog.Any("error", err))
	// 	return ""
	// }
	return filepath.Join("./", "logs", "bbs-go.log")
}
