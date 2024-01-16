package repositories

import (
	"bbs-go/internal/models"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var UserFollowRepository = newUserFollowRepository()

func newUserFollowRepository() *userFollowRepository {
	return &userFollowRepository{}
}

type userFollowRepository struct {
}

func (r *userFollowRepository) Get(db *gorm.DB, id int64) *models.UserFollow {
	ret := &models.UserFollow{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userFollowRepository) Take(db *gorm.DB, where ...interface{}) *models.UserFollow {
	ret := &models.UserFollow{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userFollowRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.UserFollow) {
	cnd.Find(db, &list)
	return
}

func (r *userFollowRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.UserFollow {
	ret := &models.UserFollow{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *userFollowRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.UserFollow, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *userFollowRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.UserFollow, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.UserFollow{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *userFollowRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.UserFollow{})
}

func (r *userFollowRepository) Create(db *gorm.DB, t *models.UserFollow) (err error) {
	err = db.Create(t).Error
	return
}

func (r *userFollowRepository) Update(db *gorm.DB, t *models.UserFollow) (err error) {
	err = db.Save(t).Error
	return
}

func (r *userFollowRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.UserFollow{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *userFollowRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.UserFollow{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *userFollowRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.UserFollow{}, "id = ?", id)
}
