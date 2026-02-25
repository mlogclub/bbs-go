package migrations

import (
	"bbs-go/internal/models"
	"bbs-go/internal/services"
	"errors"
	"log/slog"
	"sync"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"github.com/spf13/cast"
)

var migrationFuncs = make(map[int64]MigrationFunc)
var versions = make([]int64, 0)
var migrations = make(map[int64]models.Migration, 0)
var mu sync.Mutex

type MigrationFunc struct {
	Version int64
	Remark  string
	Fn      func() error
}

func Migrate() error {
	mu.Lock()
	defer mu.Unlock()

	if list := services.MigrationService.Find(sqls.NewCnd().Asc("version")); len(list) > 0 {
		for _, element := range list {
			migrations[element.Version] = element
		}
	}

	for _, version := range versions {
		if err := runMigration(version); err != nil {
			slog.Error("migrate failed", "version", version, "error", err)
			return err
		}
	}
	return nil
}

func register(version int64, remark string, fn func() error) {
	if len(versions) == 0 || version > versions[len(versions)-1] {
		versions = append(versions, version)
		migrationFuncs[version] = MigrationFunc{
			Version: version,
			Remark:  remark,
			Fn:      fn,
		}
	} else {
		slog.Error("register migration failed, version is less than latest version", slog.Any("version", version))
		panic(errors.New("register migration failed, version is less than latest version. version: " + cast.ToString(version)))
	}
}

func runMigration(version int64) error {
	migration, found := migrations[version]
	if found && migration.Success {
		return nil
	}

	f, ok := migrationFuncs[version]
	if !ok {
		return errors.New("migration function not found")
	}

	err := f.Fn()

	if !found {
		migration = models.Migration{
			Version:    f.Version,
			Remark:     f.Remark,
			Success:    false,
			RetryCount: 0,
			CreateTime: dates.NowTimestamp(),
			UpdateTime: dates.NowTimestamp(),
		}
	}
	if err == nil {
		migration.Success = true
	} else {
		migration.Success = false
		migration.ErrorInfo = err.Error()
	}
	migration.RetryCount++
	migration.UpdateTime = dates.NowTimestamp()
	if found {
		if e := services.MigrationService.Update(&migration); e != nil {
			slog.Error("update migration failed", "version", version, "error", err)
		}
	} else {
		if e := services.MigrationService.Create(&migration); e != nil {
			slog.Error("create migration failed", "version", version, "error", err)
		}
	}

	return err
}

func init() {
	register(1, "init migration", migrate_init)
	register(2, "init task data", migrate_init_task_data)
	register(3, "add email_code biz_type", migrate_add_email_code_biz_type)
	register(4, "add email log menu", migrate_add_email_log_menu)
	register(5, "migrate smtp config to sys config", migrate_smtp_config_to_sys_config)
}
