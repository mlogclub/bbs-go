package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
)

var UserTokenRepository = newUserTokenRepository()

func newUserTokenRepository() *userTokenRepository {
	return &userTokenRepository{}
}

type userTokenRepository struct {
}

func (this *userTokenRepository) GetByToken(db *gorm.DB, token string) *model.UserToken {
	if len(token) == 0 {
		return nil
	}
	return this.Take(db, "token = ?", token)
}

func (this *userTokenRepository) Get(db *gorm.DB, id int64) *model.UserToken {
	ret := &model.UserToken{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *userTokenRepository) Take(db *gorm.DB, where ...interface{}) *model.UserToken {
	ret := &model.UserToken{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *userTokenRepository) QueryCnd(db *gorm.DB, cnd *simple.QueryCnd) (list []model.UserToken, err error) {
	err = cnd.DoQuery(db).Find(&list).Error
	return
}

func (this *userTokenRepository) Query(db *gorm.DB, queries *simple.ParamQueries) (list []model.UserToken, paging *simple.Paging) {
	queries.StartQuery(db).Find(&list)
	queries.StartCount(db).Model(&model.UserToken{}).Count(&queries.Paging.Total)
	paging = queries.Paging
	return
}

func (this *userTokenRepository) Create(db *gorm.DB, t *model.UserToken) (err error) {
	err = db.Create(t).Error
	return
}

func (this *userTokenRepository) Update(db *gorm.DB, t *model.UserToken) (err error) {
	err = db.Save(t).Error
	return
}

func (this *userTokenRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.UserToken{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *userTokenRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.UserToken{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *userTokenRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.UserToken{}, "id = ?", id)
}
