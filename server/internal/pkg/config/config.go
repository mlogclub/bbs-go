package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"bbs-go/internal/pkg/simple/common/strs"
	"bbs-go/internal/pkg/simple/sqls"
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
	Language       Language          `yaml:"language"`       // 语言
	BaseURL        string            `yaml:"baseURL"`        // baseURL
	Port           int               `yaml:"port"`           // 端口
	IpDataPath     string            `yaml:"ipDataPath"`     // IP数据文件
	AllowedOrigins []string          `yaml:"allowedOrigins"` // 跨域白名单
	Installed      bool              `yaml:"installed"`      // 是否已安装
	Logger         LoggerConfig      `yaml:"logger"`         // 日志配置
	DB             sqls.DbConfig     `yaml:"db"`             // 数据库配置
	Smtp           SmtpConfig        `yaml:"smtp"`           // smtp
	Search         SearchConfig      `yaml:"search"`         // 搜索配置
	MeiliSearch    MeiliSearchConfig `yaml:"meilisearch"`    // MeiliSearch配置
	BaiduSEO       BaiduSEOConfig    `yaml:"baiduSEO"`       // 百度SEO配置
	SmSEO          SmSEOConfig       `yaml:"smSEO"`          // 神马搜索SEO配置
	Redis          RedisConfig       `yaml:"redis"`          // Redis配置
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

type MeiliSearchConfig struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	APIKey  string `yaml:"apiKey"`
	Index   string `yaml:"index"`
	Enabled bool   `yaml:"enabled"`
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

// Redis配置
type RedisConfig struct {
	Host         string `yaml:"host" default:"localhost"`          // 服务器IP地址
	Port         int    `yaml:"port" default:"6379"`               // 服务器端口号
	Password     string `yaml:"password"`                          // 密码
	DB           int    `yaml:"db" default:"0"`                    // 数据库
	PoolSize     int    `yaml:"pool_size" default:"100"`           // 连接池大小
	MinIdleConns int    `yaml:"min_idle_conns" default:"10"`       // 最小空闲连接
	MaxRetries   int    `yaml:"max_retries" default:"3"`           // 最大重试次数
	DialTimeout  int    `yaml:"dial_timeout" default:"5"`          // 连接超时时间(秒)
	ReadTimeout  int    `yaml:"read_timeout" default:"3"`          // 读取超时时间(秒)
	WriteTimeout int    `yaml:"write_timeout" default:"3"`         // 写入超时时间(秒)
	Enabled      bool   `yaml:"enabled" default:"true"`            // 是否启用Redis
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
			Redis: RedisConfig{
				Host:         "localhost",
				Port:         6379,
				DB:           0,
				PoolSize:     100,
				MinIdleConns: 10,
				MaxRetries:   3,
				DialTimeout:  5,
				ReadTimeout:  3,
				WriteTimeout: 3,
				Enabled:      true,
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
