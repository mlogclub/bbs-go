package repositories

import (
	"bbs-go/internal/models"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"
)

var ForbiddenWordRepository = newForbiddenWordRepository()

func newForbiddenWordRepository() *forbiddenWordRepository {
	return &forbiddenWordRepository{}
}

type forbiddenWordRepository struct {
}

func (r *forbiddenWordRepository) Get(db *gorm.DB, id int64) *models.ForbiddenWord {
	ret := &models.ForbiddenWord{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *forbiddenWordRepository) Take(db *gorm.DB, where ...interface{}) *models.ForbiddenWord {
	ret := &models.ForbiddenWord{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *forbiddenWordRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.ForbiddenWord) {
	cnd.Find(db, &list)
	return
}

func (r *forbiddenWordRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.ForbiddenWord {
	ret := &models.ForbiddenWord{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *forbiddenWordRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.ForbiddenWord, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *forbiddenWordRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.ForbiddenWord, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.ForbiddenWord{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *forbiddenWordRepository) FindBySql(db *gorm.DB, sqlStr string, paramArr ...interface{}) (list []models.ForbiddenWord) {
	db.Raw(sqlStr, paramArr...).Scan(&list)
	return
}

func (r *forbiddenWordRepository) CountBySql(db *gorm.DB, sqlStr string, paramArr ...interface{}) (count int64) {
	db.Raw(sqlStr, paramArr...).Count(&count)
	return
}

func (r *forbiddenWordRepository) Count(db *gorm.DB, cnd *sqls.Cnd) int64 {
	return cnd.Count(db, &models.ForbiddenWord{})
}

func (r *forbiddenWordRepository) Create(db *gorm.DB, t *models.ForbiddenWord) (err error) {
	err = db.Create(t).Error
	return
}

func (r *forbiddenWordRepository) Update(db *gorm.DB, t *models.ForbiddenWord) (err error) {
	err = db.Save(t).Error
	return
}

func (r *forbiddenWordRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.ForbiddenWord{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *forbiddenWordRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.ForbiddenWord{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *forbiddenWordRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.ForbiddenWord{}, "id = ?", id)
}
