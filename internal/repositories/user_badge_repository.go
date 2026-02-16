package repositories

import (
	"bbs-go/internal/models"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var UserBadgeRepository = newUserBadgeRepository()

func newUserBadgeRepository() *userBadgeRepository {
	return &userBadgeRepository{}
}

type userBadgeRepository struct {
}

func (r *userBadgeRepository) Get(db *gorm.DB, id int64) *models.UserBadge {
	ret := &models.UserBadge{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userBadgeRepository) Take(db *gorm.DB, where ...interface{}) *models.UserBadge {
	ret := &models.UserBadge{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *userBadgeRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.UserBadge) {
	cnd.Find(db, &list)
	return
}

func (r *userBadgeRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.UserBadge {
	ret := &models.UserBadge{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *userBadgeRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.UserBadge, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *userBadgeRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.UserBadge, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.UserBadge{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *userBadgeRepository) FindBySql(db *gorm.DB, sqlStr string, paramArr ...interface{}) (list []models.UserBadge) {
	db.Raw(sqlStr, paramArr...).Scan(&list)
	return
}

func (r *userBadgeRepository) CountBySql(db *gorm.DB, sqlStr string, paramArr ...interface{}) (count int64) {
	db.Raw(sqlStr, paramArr...).Count(&count)
	return
}

func (r *userBadgeRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.UserBadge{})
}

func (r *userBadgeRepository) Create(db *gorm.DB, t *models.UserBadge) (err error) {
	err = db.Create(t).Error
	return
}

func (r *userBadgeRepository) Update(db *gorm.DB, t *models.UserBadge) (err error) {
	err = db.Save(t).Error
	return
}

func (r *userBadgeRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.UserBadge{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *userBadgeRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.UserBadge{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *userBadgeRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.UserBadge{}, "id = ?", id)
}

func (r *userBadgeRepository) GetBy(db *gorm.DB, userId int64, badgeId int64) *models.UserBadge {
	return r.FindOne(db, sqls.NewCnd().Where("user_id = ? AND badge_id = ?", userId, badgeId))
}
