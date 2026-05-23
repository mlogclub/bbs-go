package repositories

import (
	"bbs-go/internal/pkg/params"

	"github.com/mlogclub/simple/sqls"
	"gorm.io/gorm"

	"bbs-go/internal/models"
)

var CategoryRepository = newCategoryRepository()

func newCategoryRepository() *categoryRepository {
	return &categoryRepository{}
}

type categoryRepository struct {
}

func (r *categoryRepository) Get(db *gorm.DB, id int64) *models.Category {
	ret := &models.Category{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (r *categoryRepository) Take(db *gorm.DB, where ...interface{}) *models.Category {
	ret := &models.Category{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (r *categoryRepository) Find(db *gorm.DB, cnd *sqls.Cnd) (list []models.Category) {
	cnd.Find(db, &list)
	return
}

func (r *categoryRepository) FindOne(db *gorm.DB, cnd *sqls.Cnd) *models.Category {
	ret := &models.Category{}
	if err := cnd.FindOne(db, &ret); err != nil {
		return nil
	}
	return ret
}

func (r *categoryRepository) FindPageByParams(db *gorm.DB, params *params.QueryParams) (list []models.Category, paging *sqls.Paging) {
	return r.FindPageByCnd(db, &params.Cnd)
}

func (r *categoryRepository) FindPageByCnd(db *gorm.DB, cnd *sqls.Cnd) (list []models.Category, paging *sqls.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &models.Category{})

	paging = &sqls.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (r *categoryRepository) Create(db *gorm.DB, t *models.Category) (err error) {
	err = db.Create(t).Error
	return
}

func (r *categoryRepository) Update(db *gorm.DB, t *models.Category) (err error) {
	err = db.Save(t).Error
	return
}

func (r *categoryRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&models.Category{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (r *categoryRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&models.Category{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (r *categoryRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&models.Category{}, "id = ?", id)
}
