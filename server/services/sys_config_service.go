package services

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
	"github.com/mlogclub/bbs-go/services/cache"
)

var SysConfigService = newSysConfigService()

func newSysConfigService() *sysConfigService {
	return &sysConfigService{}
}

type sysConfigService struct {
}

func (this *sysConfigService) Get(id int64) *model.SysConfig {
	return repositories.SysConfigRepository.Get(simple.DB(), id)
}

func (this *sysConfigService) Take(where ...interface{}) *model.SysConfig {
	return repositories.SysConfigRepository.Take(simple.DB(), where...)
}

func (this *sysConfigService) Find(cnd *simple.SqlCnd) []model.SysConfig {
	return repositories.SysConfigRepository.Find(simple.DB(), cnd)
}

func (this *sysConfigService) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.SysConfig) {
	cnd.FindOne(db, &ret)
	return
}

func (this *sysConfigService) FindPageByParams(params *simple.QueryParams) (list []model.SysConfig, paging *simple.Paging) {
	return repositories.SysConfigRepository.FindPageByParams(simple.DB(), params)
}

func (this *sysConfigService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.SysConfig, paging *simple.Paging) {
	return repositories.SysConfigRepository.FindPageByCnd(simple.DB(), cnd)
}

func (this *sysConfigService) GetAll() []model.SysConfig {
	return repositories.SysConfigRepository.Find(simple.DB(), simple.NewSqlCnd().Asc("id"))
}

func (this *sysConfigService) SetAll(configs map[string]string) error {
	if len(configs) == 0 {
		return nil
	}
	return simple.Tx(simple.DB(), func(tx *gorm.DB) error {
		for k, v := range configs {
			if err := this.setSingle(tx, k, v, "", ""); err != nil {
				return err
			}
		}
		return nil
	})
}

// 设置配置，如果配置不存在，那么创建
func (this *sysConfigService) Set(key, value, name, description string) error {
	return simple.Tx(simple.DB(), func(tx *gorm.DB) error {
		if err := this.setSingle(tx, key, value, name, description); err != nil {
			return err
		}
		return nil
	})
}

func (this *sysConfigService) setSingle(db *gorm.DB, key, value, name, description string) error {
	if len(key) == 0 {
		return errors.New("sys config key is null")
	}
	sysConfig := repositories.SysConfigRepository.GetByKey(simple.DB(), key)
	if sysConfig == nil {
		sysConfig = &model.SysConfig{
			CreateTime: simple.NowTimestamp(),
		}
	}
	sysConfig.Key = key
	sysConfig.Value = value
	sysConfig.UpdateTime = simple.NowTimestamp()

	if len(name) > 0 {
		sysConfig.Name = name
	}
	if len(description) > 0 {
		sysConfig.Description = description
	}

	var err error
	if sysConfig.Id > 0 {
		err = repositories.SysConfigRepository.Update(simple.DB(), sysConfig)
	} else {
		err = repositories.SysConfigRepository.Create(simple.DB(), sysConfig)
	}
	if err != nil {
		return err
	} else {
		cache.SysConfigCache.Invalidate(key)
		return nil
	}
}
