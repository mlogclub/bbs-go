package repositories

import (
	"github.com/mlogclub/simple"
	"gorm.io/gorm"

	"bbs-go/model"
)

var ThirdAccountRepository = newThirdAccountRepository()

func newThirdAccountRepository() *thirdAccountRepository {
	return &thirdAccountRepository{}
}

type thirdAccountRepository struct {
}

func (r *thirdAccountRepository) Get(db *gorm.DB, id int64) *model.ThirdAccount {
	ret := &model.ThirdAccount{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *thirdAccountRepository) Take(db *gorm.DB, where ...interface{}) *model.ThirdAccount {
	ret := &model.ThirdAccount{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *thirdAccountRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.ThirdAccount) {
	cnd.Find(db, &list)
	return
}

func (r *thirdAccountRepository) FindOne(db *gorm.DB, cnd *simple.SqlCnd) *model.ThirdAccount {
	ret := &model.ThirdAccount{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *thirdAccountRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.ThirdAccount, paging *simple.Paging) {
	return r.FindPageByCnd(db, &params.SqlCnd)
}

func (r *thirdAccountRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.ThirdAccount, paging *simple.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.ThirdAccount{})

	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *thirdAccountRepository) Create(db *gorm.DB, t *model.ThirdAccount) (err error) {
	err = db.Create(t).Error
	return
}

func (r *thirdAccountRepository) Update(db *gorm.DB, t *model.ThirdAccount) (err error) {
	err = db.Save(t).Error
	return
}

func (r *thirdAccountRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.ThirdAccount{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *thirdAccountRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.ThirdAccount{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *thirdAccountRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.ThirdAccount{}, "id = ?", id)
}
