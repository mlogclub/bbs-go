package repositories

import (
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"

	"bbs-go/internal/models"
)

var UserTokenRepository = newUserTokenRepository()

func newUserTokenRepository() *userTokenRepository {
	return &userTokenRepository{}
}

type userTokenRepository struct {
}

func (r *userTokenRepository) GetByToken(db *gorm.DB, token string) *models.UserToken {
	if len(token) == 0 {
		return nil
	}
	return r.Take(db, "token = ?", token)
}

func (r *userTokenRepository) Get(db *gorm.DB, id int64) *models.UserToken {
	ret := &models.UserToken{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userTokenRepository) Take(db *gorm.DB, where ...interface{}) *models.UserToken {
	ret := &models.UserToken{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userTokenRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.UserToken) {
	cnd.Find(db, &list)
	return
}

func (r *userTokenRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.UserToken {
	ret := &models.UserToken{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *userTokenRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.UserToken, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *userTokenRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.UserToken, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.UserToken{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *userTokenRepository) Create(db *gorm.DB, t *models.UserToken) (err error) {
	err = db.Create(t).Error
	return
}

func (r *userTokenRepository) Update(db *gorm.DB, t *models.UserToken) (err error) {
	err = db.Save(t).Error
	return
}

func (r *userTokenRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.UserToken{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *userTokenRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.UserToken{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *userTokenRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.UserToken{}, "id = ?", id)
}
