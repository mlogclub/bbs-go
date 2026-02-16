package install

import (
	"bbs-go/internal/controllers/render"
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
	"errors"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/passwd"
	"github.com/mlogclub/simple/common/strs"
	"github.com/mlogclub/simple/sqls"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	installMx sync.Mutex
)

// 测试数据库连接
type DbConfigReq struct {
	Type     string `json:"type"`               // mysql | sqlite
	Host     string `json:"host,omitempty"`     // mysql
	Port     string `json:"port,omitempty"`     // mysql
	Database string `json:"database,omitempty"` // mysql
	Username string `json:"username,omitempty"` // mysql
	Password string `json:"password,omitempty"` // mysql
}

// 执行安装
type InstallReq struct {
	SiteTitle       string          `json:"siteTitle"`
	SiteDescription string          `json:"siteDescription"`
	DbConfig        DbConfigReq     `json:"dbConfig"`
	Username        string          `json:"username"`
	Password        string          `json:"password"`
	Avatar          string          `json:"avatar"`
	Language        config.Language `json:"language"`
}

func (r DbConfigReq) GetConnStr() string {
	if r.Type == "" {
		r.Type = config.DbTypeMySQL
	}
	switch r.Type {
	case config.DbTypeSQLite:
		return buildSqliteDSN()
	default:
		return r.Username + ":" + r.Password + "@tcp(" + r.Host + ":" + r.Port + ")/" + r.Database + "?charset=utf8mb4&parseTime=True&multiStatements=true&loc=Local"
	}
}

func TestDbConnection(req DbConfigReq) error {
	if req.Type == "" {
		req.Type = config.DbTypeMySQL
	}
	dsn := req.GetConnStr()

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
		if err = sqlDB.Ping(); err != nil {
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
	default:
		db, err := gorm.Open(mysql.Open(dsn))
		if err != nil {
			return err
		}
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		if err = sqlDB.Ping(); err != nil {
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

	if err := TestDbConnection(req.DbConfig); err != nil {
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
	if err := InitData(req); err != nil {
		return err
	}
	if err := InitOthers(); err != nil {
		return err
	}
	if err := InitLocales(); err != nil {
		return err
	}
	return WriteInstallSuccess()
}

func WriteConfig(req InstallReq) error {
	cfg := config.Instance
	cfg.Language = req.Language
	cfg.IDCodec.Key = idcodec.GenerateRandomKey()
	if req.DbConfig.Type == "" {
		req.DbConfig.Type = config.DbTypeMySQL
	}
	cfg.DB.Type = req.DbConfig.Type
	cfg.DB.Url = req.DbConfig.GetConnStr()
	if strs.IsBlank(cfg.Search.IndexPath) {
		cfg.Search.IndexPath = filepath.Join(config.GetConfigDir(), "data", "topic_index")
	}
	return config.WriteConfig(cfg)
}

func WriteInstallSuccess() error {
	cfg := config.Instance
	cfg.Installed = true
	return config.WriteConfig(cfg)
}

func InitConfig() {
	cfg, _, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}
	config.Instance = cfg
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
			LogLevel:                  logger.Info,
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
	return nil
}

func InitData(req InstallReq) error {
	avatar := req.Avatar
	if strs.IsBlank(avatar) {
		avatar = render.RandomAvatar(time.Now().Unix())
	}
	// 初始化用户
	user := &models.User{
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

	// 初始化用户角色
	role := services.RoleService.GetByCode(constants.RoleOwner)
	if role == nil {
		return errors.New("please configure the super administrator first")
	}
	services.UserRoleService.UpdateUserRoles(user.Id, []int64{role.Id})

	// 初始化系统配置
	services.SysConfigService.Set(constants.SysConfigSiteTitle, req.SiteTitle)
	services.SysConfigService.Set(constants.SysConfigSiteDescription, req.SiteDescription)

	// 初始化默认欢迎帖子（按安装语言）
	if err := initWelcomeTopic(req.Language, user.Id); err != nil {
		return err
	}

	return nil
}

func initWelcomeTopic(language config.Language, userId int64) error {
	title, content := getWelcomeTopicContent(language)
	if services.TopicService.Take("user_id = ? and title = ?", userId, title) != nil {
		return nil
	}

	nodeId := services.SysConfigService.GetDefaultNodeId()
	if nodeId <= 0 {
		nodes := services.TopicNodeService.GetNodes()
		if len(nodes) == 0 {
			return errors.New("no topic node found for welcome topic initialization")
		}
		nodeId = nodes[0].Id
	}

	_, err := services.TopicPublishService.Publish(userId, modelreq.CreateTopicForm{
		Type:        constants.TopicTypeTopic,
		NodeId:      nodeId,
		Title:       title,
		Content:     content,
		ContentType: constants.ContentTypeMarkdown,
	})
	return err
}

func getWelcomeTopicContent(language config.Language) (string, string) {
	if language == config.LanguageZhCN {
		return "欢迎来到 BBS-GO 社区", `欢迎使用 **BBS-GO**！

这是一个轻量、高性能、易扩展的社区系统。

你可以发帖、评论、点赞，并通过任务系统获得积分与成长奖励。

现在就开始发布你的第一篇帖子吧。`
	}

	return "Welcome to the BBS-GO Community", `Welcome to **BBS-GO**!

A lightweight, high-performance, and extensible community platform.

You can create topics, comment, like, and earn points through the task system.

Start by publishing your first post.`
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
