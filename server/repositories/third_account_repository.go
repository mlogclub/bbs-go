
package repositories

import (
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/simple"
	"github.com/jinzhu/gorm"
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

func (this *thirdAccountRepository) QueryCnd(db *gorm.DB, cnd *simple.QueryCnd) (list []model.ThirdAccount, err error) {
	err = cnd.DoQuery(db).Find(&list).Error
	return
}

func (this *thirdAccountRepository) Query(db *gorm.DB, params *simple.ParamQueries) (list []model.ThirdAccount, paging *simple.Paging) {
	params.StartQuery(db).Find(&list)
    params.StartCount(db).Model(&model.ThirdAccount{}).Count(&params.Paging.Total)
	paging = params.Paging
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

