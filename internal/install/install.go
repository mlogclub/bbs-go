package install

import (
	"bbs-go/internal/handlers/render"
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	modelreq "bbs-go/internal/models/req"
	"bbs-go/internal/pkg/config"
	"bbs-go/internal/pkg/idcodec"
	"bbs-go/internal/pkg/iplocator"
	"bbs-go/internal/pkg/locales"
	"bbs-go/internal/pkg/search"
	"bbs-go/internal/scheduler"
	"bbs-go/internal/services"
	"bbs-go/migrations"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/passwd"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"

	// "gorm.io/driver/sqlite" // Sqlite driver based on CGO
	"github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	installMx sync.Mutex
)

const (
	// MySQL Docker 内置环境变量
	DockerBuiltinMySQLEnv         = "BBSGO_INSTALL_DOCKER_BUILTIN_MYSQL"
	DockerBuiltinMySQLHostEnv     = "BBSGO_DOCKER_BUILTIN_MYSQL_HOST"
	DockerBuiltinMySQLPortEnv     = "BBSGO_DOCKER_BUILTIN_MYSQL_PORT"
	DockerBuiltinMySQLDatabaseEnv = "BBSGO_DOCKER_BUILTIN_MYSQL_DATABASE"
	DockerBuiltinMySQLUsernameEnv = "BBSGO_DOCKER_BUILTIN_MYSQL_USERNAME"
	DockerBuiltinMySQLPasswordEnv = "BBSGO_DOCKER_BUILTIN_MYSQL_PASSWORD"

	// PostgreSQL Docker 内置环境变量
	DockerBuiltinPostgreSQLEnv         = "BBSGO_INSTALL_DOCKER_BUILTIN_POSTGRESQL"
	DockerBuiltinPostgreSQLHostEnv     = "BBSGO_DOCKER_BUILTIN_POSTGRESQL_HOST"
	DockerBuiltinPostgreSQLPortEnv     = "BBSGO_DOCKER_BUILTIN_POSTGRESQL_PORT"
	DockerBuiltinPostgreSQLDatabaseEnv = "BBSGO_DOCKER_BUILTIN_POSTGRESQL_DATABASE"
	DockerBuiltinPostgreSQLUsernameEnv = "BBSGO_DOCKER_BUILTIN_POSTGRESQL_USERNAME"
	DockerBuiltinPostgreSQLPasswordEnv = "BBSGO_DOCKER_BUILTIN_POSTGRESQL_PASSWORD"
)

// 测试数据库连接
type DbConfigReq struct {
	Type     string `json:"type" form:"type"`                   // mysql | postgresql | sqlite
	Host     string `json:"host,omitempty" form:"host"`         // mysql | postgresql
	Port     string `json:"port,omitempty" form:"port"`         // mysql | postgresql
	Database string `json:"database,omitempty" form:"database"` // mysql | postgresql
	Username string `json:"username,omitempty" form:"username"` // mysql | postgresql
	Password string `json:"password,omitempty" form:"password"` // mysql | postgresql
}

// 执行安装
type InstallReq struct {
	SiteTitle       string          `json:"siteTitle" form:"siteTitle"`
	SiteDescription string          `json:"siteDescription" form:"siteDescription"`
	BaseURL         string          `json:"baseURL" form:"baseURL"`
	DbConfig        DbConfigReq     `json:"dbConfig" form:"dbConfig"`
	Username        string          `json:"username" form:"username"`
	Password        string          `json:"password" form:"password"`
	Avatar          string          `json:"avatar" form:"avatar"`
	Language        config.Language `json:"language" form:"language"`
}

func (r DbConfigReq) GetConnStr() string {
	if r.Type == "" {
		r.Type = config.DbTypeMySQL
	}
	switch r.Type {
	case config.DbTypeMySQL:
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&multiStatements=true&loc=Local", r.Username, r.Password, r.Host, r.Port, r.Database)
	case config.DbTypePostgreSQL:
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", r.Host, r.Port, r.Username, r.Password, r.Database)
	case config.DbTypeSQLite:
		return buildSqliteDSN()
	default:
		return ""
	}
}

func TestDbConnection(ctx context.Context, req DbConfigReq) error {
	var dsn string
	if IsDockerBuiltinMySQLInstall() {
		req = DbConfigReq{
			Type: config.DbTypeMySQL,
		}
		dsn = config.Instance.DB.Url
	}
	if req.Type == "" {
		req.Type = config.DbTypeMySQL
	}
	if dsn == "" {
		dsn = req.GetConnStr()
	}

	switch req.Type {
	case config.DbTypeSQLite:
		db, err := gorm.Open(sqlite.Open(dsn))
		if err != nil {
			return err
		}
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		if err = sqlDB.PingContext(ctx); err != nil {
			return err
		}
		var name string
		err = db.Raw("SELECT name FROM sqlite_master WHERE type='table' LIMIT 1").Scan(&name).Error
		if err != nil {
			return err
		}
		if name != "" {
			return errors.New("please use an empty database for installation")
		}

	case config.DbTypePostgreSQL:
		db, err := gorm.Open(postgres.Open(dsn))
		if err != nil {
			return err
		}
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		if err = sqlDB.PingContext(ctx); err != nil {
			return err
		}
		var tables []string
		if err := db.Raw("SELECT tablename FROM pg_tables WHERE schemaname='public'").Scan(&tables).Error; err != nil {
			return err
		}
		if len(tables) > 0 {
			return errors.New("please use an empty database for installation")
		}

	default: // MySQL
		db, err := gorm.Open(mysql.Open(dsn))
		if err != nil {
			return err
		}
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		if err = sqlDB.PingContext(ctx); err != nil {
			return err
		}
		var tables []string
		if err := db.Raw("SHOW TABLES").Scan(&tables).Error; err != nil {
			return err
		}
		if len(tables) > 0 {
			return errors.New("please use an empty database for installation")
		}
	}
	return nil
}

func Install(req InstallReq) error {
	installMx.Lock()
	defer installMx.Unlock()

	if IsDockerBuiltinMySQLInstall() {
		req.DbConfig = DbConfigReq{
			Type: config.DbTypeMySQL,
		}
	}
	// 传入 context
	if err := TestDbConnection(context.Background(), req.DbConfig); err != nil {
		return err
	}
	if err := WriteConfig(req); err != nil {
		return err
	}
	if err := InitDB(); err != nil {
		return err
	}
	if err := InitMigrations(); err != nil {
		return err
	}
	if err := InitOthers(); err != nil {
		return err
	}
	if err := InitLocales(); err != nil {
		return err
	}
	if err := InitData(req); err != nil {
		return err
	}
	return WriteInstallSuccess()
}

func WriteConfig(req InstallReq) error {
	cfg := config.Instance
	cfg.Language = req.Language
	cfg.IDCodec.Key = idcodec.GenerateRandomKey()
	if !IsDockerBuiltinMySQLInstall() {
		if req.DbConfig.Type == "" {
			req.DbConfig.Type = config.DbTypeMySQL
		}
		cfg.DB.Type = req.DbConfig.Type
		cfg.DB.Url = req.DbConfig.GetConnStr()
	}
	if strs.IsBlank(cfg.Search.IndexPath) {
		cfg.Search.IndexPath = filepath.Join(config.GetConfigDir(), "data", "topic_index")
	}
	return WriteRuntimeConfig(cfg)
}

func WriteInstallSuccess() error {
	cfg := config.Instance
	cfg.Installed = true
	return WriteRuntimeConfig(cfg)
}

func InitConfig() {
	cfg, _, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}
	config.Instance = cfg
	ApplyDockerBuiltinMySQLConfig()
	ApplyDockerBuiltinPostgreSQLConfig()
}

// ResolveSqlitePath 将配置的 sqlite 相对路径转换为绝对路径（相对 bbs-go.yaml 所在目录）
func buildSqliteDSN() string {
	filepath := filepath.Join(config.GetConfigDir(), "bbs-go.db")
	return "file:" + filepath + "?cache=shared&_fk=0&mode=rwc&_journal_mode=wal"
}

func InitDB() error {
	conf := config.Instance.DB
	config.SetDbDefaults(&conf)
	config.Instance.DB = conf

	var dialector gorm.Dialector
	switch conf.Type {
	case config.DbTypePostgreSQL:
		dialector = postgres.Open(conf.Url)
	case config.DbTypeSQLite:
		dsn := conf.Url
		dsn = buildSqliteDSN()
		dialector = sqlite.Open(dsn)
	default:
		dialector = mysql.Open(conf.Url)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		},
		Logger: logger.New(log.New(os.Stdout, "", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  resolveGormLogLevel(conf.LogLevel),
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		}),
	})

	if err != nil {
		slog.Error(err.Error())
		return err
	}

	if sqlDB, err := db.DB(); err == nil {
		if conf.Type == config.DbTypeSQLite {
			sqlDB.SetMaxOpenConns(1)
			sqlDB.SetMaxIdleConns(1)
		} else {
			sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
			sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
			sqlDB.SetConnMaxIdleTime(time.Duration(conf.ConnMaxIdleTimeSeconds) * time.Second)
			sqlDB.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetimeSeconds) * time.Second)
		}
	}

	sqls.SetDB(db)
	return nil
}

func InitMigrations() error {
	// migrate schema
	if err := sqls.DB().AutoMigrate(models.Models...); err != nil {
		slog.Error("auto migrate tables failed", slog.Any("error", err))
		return err
	}

	// migrate data
	if err := migrations.Migrate(); err != nil {
		slog.Error("migrate failed", slog.Any("error", err))
		return err
	}

	// sysc permissions
	if err := services.PermissionService.SyncDefinitions(); err != nil {
		slog.Error("sync permissions failed", slog.Any("error", err))
		return err
	}
	return nil
}

func InitData(req InstallReq) error {
	avatar := req.Avatar
	if strs.IsBlank(avatar) {
		avatar = render.RandomAvatar(time.Now().Unix())
	}
	// 初始化用户
	user := services.UserService.GetByUsername(req.Username)
	if user == nil {
		user = &models.User{
			Nickname:   req.Username,
			Username:   sqls.SqlNullString(req.Username),
			Password:   passwd.EncodePassword(req.Password),
			Avatar:     avatar,
			Status:     constants.StatusOk,
			CreateTime: dates.NowTimestamp(),
			UpdateTime: dates.NowTimestamp(),
		}
		if err := services.UserService.Create(user); err != nil {
			return err
		}
	}

	// 初始化用户角色
	role := services.RoleService.GetByCode(constants.RoleOwner)
	if role == nil {
		return errors.New("please configure the super administrator first")
	}
	services.UserRoleService.UpdateUserRoles(user.Id, []int64{role.Id})

	// 初始化系统配置
	baseURL := req.BaseURL
	if strs.IsBlank(baseURL) {
		baseURL = "/"
	}
	services.SysConfigService.Set(constants.SysConfigSiteTitle, req.SiteTitle)
	services.SysConfigService.Set(constants.SysConfigSiteDescription, req.SiteDescription)
	services.SysConfigService.Set(constants.SysConfigBaseURL, baseURL)

	// 初始化默认欢迎帖子（按安装语言）
	if err := initWelcomeTopic(req.Language, user.Id); err != nil {
		return err
	}

	return nil
}

func initWelcomeTopic(language config.Language, userId int64) error {
	var (
		title   string
		content string
	)
	if language == config.LanguageZhCN {
		title = "欢迎来到 BBS-GO 社区"
		content = `欢迎使用 **BBS-GO**！

这是一个轻量、高性能、易扩展的社区系统。

你可以发帖、评论、点赞，并通过任务系统获得积分与成长奖励。

现在就开始发布你的第一篇帖子吧。`
	} else {
		title = "Welcome to the BBS-GO Community"
		content = `Welcome to **BBS-GO**!

A lightweight, high-performance, and extensible community platform.

You can create topics, comment, like, and earn points through the task system.

Start by publishing your first post.`
	}

	if services.TopicService.Take("user_id = ? and title = ?", userId, title) != nil {
		return nil
	}

	categoryId := services.SysConfigService.GetDefaultCategoryId()
	if categoryId <= 0 {
		categories := services.CategoryService.GetCategories()
		if len(categories) == 0 {
			return errors.New("no category found for welcome topic initialization")
		}
		categoryId = categories[0].Id
	}

	_, err := services.TopicPublishService.Publish(userId, modelreq.CreateTopicReq{
		Type:        constants.TopicTypeTopic,
		CategoryId:  categoryId,
		Title:       title,
		Content:     content,
		ContentType: constants.ContentTypeMarkdown,
	})
	return err
}

func InitLocales() error {
	return locales.Init()
}

func InitOthers() error {
	idcodec.Init(config.Instance.IDCodec.Key)
	if config.IsProd() {
		scheduler.Start()
	}
	iplocator.InitIpLocator()
	search.Init()
	return nil
}

func resolveGormLogLevel(level string) logger.LogLevel {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info", "":
		return logger.Info
	default:
		return logger.Info
	}
}

func IsDockerBuiltinMySQLInstall() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(DockerBuiltinMySQLEnv))) {
	case "1", "true", "yes", "y", "on":
		return true
	default:
		return false
	}
}

func IsDockerBuiltinPostgreSQLInstall() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(DockerBuiltinPostgreSQLEnv))) {
	case "1", "true", "yes", "y", "on":
		return true
	default:
		return false
	}
}

func ApplyDockerBuiltinMySQLConfig() {
	if !IsDockerBuiltinMySQLInstall() {
		return
	}
	config.Instance.DB.Type = config.DbTypeMySQL
	config.Instance.DB.Url = DbConfigReq{
		Type:     config.DbTypeMySQL,
		Host:     dockerBuiltinMySQLEnv(DockerBuiltinMySQLHostEnv, "mysql"),
		Port:     dockerBuiltinMySQLEnv(DockerBuiltinMySQLPortEnv, "3306"),
		Database: dockerBuiltinMySQLEnv(DockerBuiltinMySQLDatabaseEnv, "bbsgo"),
		Username: dockerBuiltinMySQLEnv(DockerBuiltinMySQLUsernameEnv, "bbsgo"),
		Password: dockerBuiltinMySQLEnv(DockerBuiltinMySQLPasswordEnv, "bbsgo_password"),
	}.GetConnStr()
}

func ApplyDockerBuiltinPostgreSQLConfig() {
	if !IsDockerBuiltinPostgreSQLInstall() {
		return
	}

	config.Instance.DB.Type = config.DbTypePostgreSQL
	config.Instance.DB.Url = DbConfigReq{
		Type:     config.DbTypePostgreSQL,
		Host:     dockerBuiltinPostgreSQLEnv(DockerBuiltinPostgreSQLHostEnv, "postgresql"),
		Port:     dockerBuiltinPostgreSQLEnv(DockerBuiltinPostgreSQLPortEnv, "5432"),
		Database: dockerBuiltinPostgreSQLEnv(DockerBuiltinPostgreSQLDatabaseEnv, "bbsgo"),
		Username: dockerBuiltinPostgreSQLEnv(DockerBuiltinPostgreSQLUsernameEnv, "bbsgo"),
		Password: dockerBuiltinPostgreSQLEnv(DockerBuiltinPostgreSQLPasswordEnv, "bbsgo_password"),
	}.GetConnStr()
}

func WriteRuntimeConfig(cfg *config.Config) error {
	if !IsDockerBuiltinMySQLInstall() {
		return config.WriteConfig(cfg)
	}
	persisted := *cfg
	persisted.DB.Type = config.DbTypeMySQL
	persisted.DB.Url = ""
	return config.WriteConfig(&persisted)
}

func dockerBuiltinMySQLEnv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func dockerBuiltinPostgreSQLEnv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
