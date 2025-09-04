package sqls

import (
	"gorm.io/gorm"
)

type GormModel struct {
	Id int64 `gorm:"primaryKey;autoIncrement" json:"id" form:"id"`
}

type DbConfig struct {
	Url                    string `yaml:"url"`
	MaxIdleConns           int    `yaml:"maxIdleConns"`
	MaxOpenConns           int    `yaml:"maxOpenConns"`
	ConnMaxIdleTimeSeconds int    `yaml:"connMaxIdleTimeSeconds"`
	ConnMaxLifetimeSeconds int    `yaml:"connMaxLifetimeSeconds"`
}

var (
	_db *gorm.DB
)

// func Open(dbConfig DbConfig, config *gorm.Config, models ...interface{}) (err error) {
// 	if config == nil {
// 		config = &gorm.Config{}
// 	}

// 	if config.NamingStrategy == nil {
// 		config.NamingStrategy = schema.NamingStrategy{
// 			TablePrefix:   "t_",
// 			SingularTable: true,
// 		}
// 	}

// 	if _db, err = gorm.Open(mysql.Open(dbConfig.Url), config); err != nil {
// 		slog.Error("opens database failed", slog.Any("error", err))
// 		return
// 	}

// 	if sqlDB, err = _db.DB(); err == nil {
// 		sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
// 		sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
// 		sqlDB.SetConnMaxIdleTime(time.Duration(dbConfig.ConnMaxIdleTimeSeconds) * time.Second)
// 		sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetimeSeconds) * time.Second)
// 	} else {
// 		slog.Error(err.Error(), slog.Any("error", err))
// 	}

// 	if err = _db.AutoMigrate(models...); nil != err {
// 		slog.Error("auto migrate tables failed", slog.Any("error", err))
// 	}
// 	return
// }

func DB() *gorm.DB {
	return _db
}

func SetDB(gormDB *gorm.DB) {
	_db = gormDB
}
