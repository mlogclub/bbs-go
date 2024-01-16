package repositories

import (
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"

	"bbs-go/internal/models"
)

var UserLikeRepository = newUserLikeRepository()

func newUserLikeRepository() *userLikeRepository {
	return &userLikeRepository{}
}

type userLikeRepository struct {
}

func (r *userLikeRepository) Get(db *gorm.DB, id int64) *models.UserLike {
	ret := &models.UserLike{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userLikeRepository) Take(db *gorm.DB, where ...interface{}) *models.UserLike {
	ret := &models.UserLike{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userLikeRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.UserLike) {
	cnd.Find(db, &list)
	return
}

func (r *userLikeRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.UserLike {
	ret := &models.UserLike{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *userLikeRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.UserLike, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *userLikeRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.UserLike, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.UserLike{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *userLikeRepository) Create(db *gorm.DB, t *models.UserLike) (err error) {
	err = db.Create(t).Error
	return
}

func (r *userLikeRepository) Update(db *gorm.DB, t *models.UserLike) (err error) {
	err = db.Save(t).Error
	return
}

func (r *userLikeRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.UserLike{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *userLikeRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.UserLike{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *userLikeRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.UserLike{}, "id = ?", id)
}
