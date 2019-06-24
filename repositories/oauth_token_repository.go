package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/simple"
)

type OauthTokenRepository struct {
}

func NewOauthTokenRepository() *OauthTokenRepository {
	return &OauthTokenRepository{}
}

func (this *OauthTokenRepository) Get(db *gorm.DB, id int64) *model.OauthToken {
	ret := &model.OauthToken{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *OauthTokenRepository) Take(db *gorm.DB, where ...interface{}) *model.OauthToken {
	ret := &model.OauthToken{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *OauthTokenRepository) QueryCnd(db *gorm.DB, cnd *simple.QueryCnd) (list []model.OauthToken, err error) {
	err = cnd.DoQuery(db).Find(&list).Error
	return
}

func (this *OauthTokenRepository) Query(db *gorm.DB, queries *simple.ParamQueries) (list []model.OauthToken, paging *simple.Paging) {
	queries.StartQuery(db).Find(&list)
	queries.StartCount(db).Model(&model.OauthToken{}).Count(&queries.Paging.Total)
	paging = queries.Paging
	return
}

func (this *OauthTokenRepository) Create(db *gorm.DB, t *model.OauthToken) (err error) {
	err = db.Create(t).Error
	return
}

func (this *OauthTokenRepository) Update(db *gorm.DB, t *model.OauthToken) (err error) {
	err = db.Save(t).Error
	return
}

func (this *OauthTokenRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.OauthToken{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *OauthTokenRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.OauthToken{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *OauthTokenRepository) Delete(db *gorm.DB, id int64) {
	db.Model(&model.OauthToken{}).Delete("id", id)
}

func (OauthTokenRepository) RemoveByCode(db *gorm.DB, code string) {
	db.Where("code = ?", code).Delete(&model.OauthToken{})
}

func (OauthTokenRepository) RemoveByAccessToken(db *gorm.DB, accessToken string) {
	db.Where("access_token = ?", accessToken).Delete(&model.OauthToken{})
}

func (OauthTokenRepository) RemoveByRefreshToken(db *gorm.DB, refreshToken string) {
	db.Where("refresh_token = ?", refreshToken).Delete(&model.OauthToken{})
}

func (OauthTokenRepository) GetByCode(db *gorm.DB, code string) *model.OauthToken {
	var oauthToken model.OauthToken
	if err := db.First(&oauthToken, "code = ?", code).Error; err != nil {
		return nil
	}
	return &oauthToken
}

func (OauthTokenRepository) GetByAccessToken(db *gorm.DB, accessToken string) *model.OauthToken {
	var oauthToken model.OauthToken
	if err := db.First(&oauthToken, "access_token = ?", accessToken).Error; err != nil {
		return nil
	}
	return &oauthToken
}

func (OauthTokenRepository) GetByRefreshToken(db *gorm.DB, refreshToken string) *model.OauthToken {
	var oauthToken model.OauthToken
	if err := db.First(&oauthToken, "refresh_token = ?", refreshToken).Error; err != nil {
		return nil
	}
	return &oauthToken
}
