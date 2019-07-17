package services

import (
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
	"github.com/mlogclub/simple"
)

var SysConfigService = newSysConfigService()

func newSysConfigService() *sysConfigService {
	return &sysConfigService{}
}

type sysConfigService struct {
}

func (this *sysConfigService) Get(id int64) *model.SysConfig {
	return repositories.SysConfigRepository.Get(simple.GetDB(), id)
}

func (this *sysConfigService) Take(where ...interface{}) *model.SysConfig {
	return repositories.SysConfigRepository.Take(simple.GetDB(), where...)
}

func (this *sysConfigService) QueryCnd(cnd *simple.QueryCnd) (list []model.SysConfig, err error) {
	return repositories.SysConfigRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *sysConfigService) Query(queries *simple.ParamQueries) (list []model.SysConfig, paging *simple.Paging) {
	return repositories.SysConfigRepository.Query(simple.GetDB(), queries)
}

func (this *sysConfigService) Create(t *model.SysConfig) error {
	return repositories.SysConfigRepository.Create(simple.GetDB(), t)
}

func (this *sysConfigService) Update(t *model.SysConfig) error {
	return repositories.SysConfigRepository.Update(simple.GetDB(), t)
}

func (this *sysConfigService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.SysConfigRepository.Updates(simple.GetDB(), id, columns)
}

func (this *sysConfigService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.SysConfigRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *sysConfigService) Delete(id int64) {
	repositories.SysConfigRepository.Delete(simple.GetDB(), id)
}
