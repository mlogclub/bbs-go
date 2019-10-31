package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
)

var CategoryRepository = newCategoryRepository()

func newCategoryRepository() *categoryRepository {
	return &categoryRepository{}
}

type categoryRepository struct {
}

func (this *categoryRepository) Get(db *gorm.DB, id int64) *model.Category {
	ret := &model.Category{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *categoryRepository) Take(db *gorm.DB, where ...interface{}) *model.Category {
	ret := &model.Category{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *categoryRepository) Find(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Category, err error) {
	err = cnd.Find(db, &list)
	return
}

func (this *categoryRepository) FindPageByParams(db *gorm.DB, params *simple.QueryParams) (list []model.Category, paging *simple.Paging) {
	return this.FindPageByCnd(db, &params.SqlCnd)
}

func (this *categoryRepository) FindPageByCnd(db *gorm.DB, cnd *simple.SqlCnd) (list []model.Category, paging *simple.Paging) {
	err := cnd.Find(db, &list)
	if err != nil {
		return
	}

	count, err := cnd.Count(db, &model.Category{})
	if err != nil {
		return
	}
	paging = &simple.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}

func (this *categoryRepository) Create(db *gorm.DB, t *model.Category) (err error) {
	err = db.Create(t).Error
	return
}

func (this *categoryRepository) Update(db *gorm.DB, t *model.Category) (err error) {
	err = db.Save(t).Error
	return
}

func (this *categoryRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Category{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *categoryRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Category{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *categoryRepository) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.Category{}, "id = ?", id)
}

func (this *categoryRepository) GetCategories() ([]model.Category, error) {
	return this.Find(simple.DB(), simple.NewSqlCnd().Where("status = ?", model.CategoryStatusOk))
}
