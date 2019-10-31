
package repositories

import (
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/simple"
	"github.com/jinzhu/gorm"
)

var CollectRuleRepository = newCollectRuleRepository()

func newCollectRuleRepository() *collectRuleRepository {
	return &collectRuleRepository{}
}

type collectRuleRepository struct {
}

func (this *collectRuleRepository) Get(db *gorm.DB, id int64) *model.CollectRule {
	ret := &model.CollectRule{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *collectRuleRepository) Take(db *gorm.DB, where ...interface{}) *model.CollectRule {
	ret := &model.CollectRule{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *collectRuleRepository) QueryCnd(db *gorm.DB, cnd *simple.QueryCnd) (list []model.CollectRule, err error) {
	err = cnd.DoQuery(db).Find(&list).Error
	return
}

func (this *collectRuleRepository) Query(db *gorm.DB, params *simple.ParamQueries) (list []model.CollectRule, paging *simple.Paging) {
	params.StartQuery(db).Find(&list)
    params.StartCount(db).Model(&model.CollectRule{}).Count(&params.Paging.Total)
	paging = params.Paging
	return
}

func (this *collectRuleRepository) Create(db *gorm.DB, t *model.CollectRule) (err error) {
	err = db.Create(t).Error
	return
}

func (this *collectRuleRepository) Update(db *gorm.DB, t *model.CollectRule) (err error) {
	err = db.Save(t).Error
	return
}

func (this *collectRuleRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.CollectRule{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *collectRuleRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.CollectRule{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *collectRuleRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.CollectRule{}, "id = ?", id)
}

