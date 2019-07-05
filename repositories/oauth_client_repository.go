package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/simple"
)

var OauthClientRepository = newOauthClientRepository()

func newOauthClientRepository() *oauthClientRepository {
	return &oauthClientRepository{}
}

type oauthClientRepository struct {
}

func (this *oauthClientRepository) GetByClientId(db *gorm.DB, clientId string) *model.OauthClient {
	var oauthClient model.OauthClient
	if err := db.First(&oauthClient, "client_id = ?", clientId).Error; err != nil {
		return nil
	}
	return &oauthClient
}

func (this *oauthClientRepository) Get(db *gorm.DB, id int64) *model.OauthClient {
	ret := &model.OauthClient{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *oauthClientRepository) Take(db *gorm.DB, where ...interface{}) *model.OauthClient {
	ret := &model.OauthClient{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *oauthClientRepository) QueryCnd(db *gorm.DB, cnd *simple.QueryCnd) (list []model.OauthClient, err error) {
	err = cnd.DoQuery(db).Find(&list).Error
	return
}

func (this *oauthClientRepository) Query(db *gorm.DB, queries *simple.ParamQueries) (list []model.OauthClient, paging *simple.Paging) {
	queries.StartQuery(db).Find(&list)
	queries.StartCount(db).Model(&model.OauthClient{}).Count(&queries.Paging.Total)
	paging = queries.Paging
	return
}

func (this *oauthClientRepository) Create(db *gorm.DB, t *model.OauthClient) (err error) {
	err = db.Create(t).Error
	return
}

func (this *oauthClientRepository) Update(db *gorm.DB, t *model.OauthClient) (err error) {
	err = db.Save(t).Error
	return
}

func (this *oauthClientRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.OauthClient{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *oauthClientRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.OauthClient{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *oauthClientRepository) Delete(db *gorm.DB, id int64) {
	db.Model(&model.OauthClient{}).Delete("id", id)
}
