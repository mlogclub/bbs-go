package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
)

var ThirdAccountRepository = newThirdAccountRepository()

func newThirdAccountRepository() *thirdAccountRepository {
	return &thirdAccountRepository{}
}

type thirdAccountRepository struct {
}

func (this *thirdAccountRepository) Get(db *gorm.DB, id int64) *model.ThirdAccount {
	ret := &model.ThirdAccount{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *thirdAccountRepository) Take(db *gorm.DB, where ...interface{}) *model.ThirdAccount {
	ret := &model.ThirdAccount{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *thirdAccountRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.ThirdAccount) {
	cnd.Find(db, &list)
	return
}

func (this *thirdAccountRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.ThirdAccount) {
	cnd.FindOne(db, &ret)
	return
}

func (this *thirdAccountRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.ThirdAccount, paging *simple.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *thirdAccountRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.ThirdAccount, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.ThirdAccount{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *thirdAccountRepository) Create(db *gorm.DB, t *model.ThirdAccount) (err error) {
	err = db.Create(t).Error
	return
}

func (this *thirdAccountRepository) Update(db *gorm.DB, t *model.ThirdAccount) (err error) {
	err = db.Save(t).Error
	return
}

func (this *thirdAccountRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.ThirdAccount{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *thirdAccountRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.ThirdAccount{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *thirdAccountRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.ThirdAccount{}, "id = ?", id)
}
