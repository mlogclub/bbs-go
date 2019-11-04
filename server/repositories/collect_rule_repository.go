package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
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

func (this *collectRuleRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.CollectRule) {
	cnd.Find(db, &list)
	return
}

func (this *collectRuleRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.CollectRule) {
	cnd.FindOne(db, &ret)
	return
}

func (this *collectRuleRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.CollectRule, paging *simple.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *collectRuleRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.CollectRule, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.CollectRule{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
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
