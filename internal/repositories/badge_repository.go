package repositories

import (
	"bbs-go/internal/models"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var BadgeRepository = newBadgeRepository()

func newBadgeRepository() *badgeRepository {
	return &badgeRepository{}
}

type badgeRepository struct {
}

func (r *badgeRepository) Get(db *gorm.DB, id int64) *models.Badge {
	ret := &models.Badge{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *badgeRepository) Take(db *gorm.DB, where ...interface{}) *models.Badge {
	ret := &models.Badge{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *badgeRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.Badge) {
	cnd.Find(db, &list)
	return
}

func (r *badgeRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.Badge {
	ret := &models.Badge{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *badgeRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.Badge, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *badgeRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.Badge, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.Badge{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *badgeRepository) FindBySql(db *gorm.DB, sqlStr string, paramArr... interface{}) (list []models.Badge) {
	db.Raw(sqlStr, paramArr...).Scan(&list)
	return
}

func (r *badgeRepository) CountBySql(db *gorm.DB, sqlStr string, paramArr... interface{}) (count int64) {
	db.Raw(sqlStr, paramArr...).Count(&count)
	return
}

func (r *badgeRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.Badge{})
}

func (r *badgeRepository) Create(db *gorm.DB, t *models.Badge) (err error) {
	err = db.Create(t).Error
	return
}

func (r *badgeRepository) Update(db *gorm.DB, t *models.Badge) (err error) {
	err = db.Save(t).Error
	return
}

func (r *badgeRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.Badge{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *badgeRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.Badge{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *badgeRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.Badge{}, "id = ?", id)
}

