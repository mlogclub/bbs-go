
package services

import (
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
	"github.com/mlogclub/simple"
)

var CollectRuleService = newCollectRuleService()

func newCollectRuleService() *collectRuleService {
	return &collectRuleService {}
}

type collectRuleService struct {
}

func (this *collectRuleService) Get(id int64) *model.CollectRule {
	return repositories.CollectRuleRepository.Get(simple.GetDB(), id)
}

func (this *collectRuleService) Take(where ...interface{}) *model.CollectRule {
	return repositories.CollectRuleRepository.Take(simple.GetDB(), where...)
}

func (this *collectRuleService) QueryCnd(cnd *simple.SqlCnd) (list []model.CollectRule, err error) {
	return repositories.CollectRuleRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *collectRuleService) Query(params *simple.QueryParams) (list []model.CollectRule, paging *simple.Paging) {
	return repositories.CollectRuleRepository.Query(simple.GetDB(), queries)
}

func (this *collectRuleService) Create(t *model.CollectRule) error {
	return repositories.CollectRuleRepository.Create(simple.GetDB(), t)
}

func (this *collectRuleService) Update(t *model.CollectRule) error {
	return repositories.CollectRuleRepository.Update(simple.GetDB(), t)
}

func (this *collectRuleService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.CollectRuleRepository.Updates(simple.GetDB(), id, columns)
}

func (this *collectRuleService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.CollectRuleRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *collectRuleService) Delete(id int64) {
	repositories.CollectRuleRepository.Delete(simple.GetDB(), id)
}

