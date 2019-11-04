package services

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
)

var CollectRuleService = newCollectRuleService()

func newCollectRuleService() *collectRuleService {
	return &collectRuleService{}
}

type collectRuleService struct {
}

func (this *collectRuleService) Get(id int64) *model.CollectRule {
	return repositories.CollectRuleRepository.Get(simple.DB(), id)
}

func (this *collectRuleService) Take(where ...interface{}) *model.CollectRule {
	return repositories.CollectRuleRepository.Take(simple.DB(), where...)
}

func (this *collectRuleService) Find(cnd *simple.SqlCnd) []model.CollectRule {
	return repositories.CollectRuleRepository.Find(simple.DB(), cnd)
}

func (this *collectRuleService) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.CollectRule) {
	cnd.FindOne(db, &ret)
	return
}

func (this *collectRuleService) FindPageByParams(params *simple.QueryParams) (list []model.CollectRule, paging *simple.Paging) {
	return repositories.CollectRuleRepository.FindPageByParams(simple.DB(), params)
}

func (this *collectRuleService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.CollectRule, paging *simple.Paging) {
	return repositories.CollectRuleRepository.FindPageByCnd(simple.DB(), cnd)
}
func (this *collectRuleService) Create(t *model.CollectRule) error {
	return repositories.CollectRuleRepository.Create(simple.DB(), t)
}

func (this *collectRuleService) Update(t *model.CollectRule) error {
	return repositories.CollectRuleRepository.Update(simple.DB(), t)
}

func (this *collectRuleService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.CollectRuleRepository.Updates(simple.DB(), id, columns)
}

func (this *collectRuleService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.CollectRuleRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (this *collectRuleService) Delete(id int64) {
	repositories.CollectRuleRepository.Delete(simple.DB(), id)
}
