package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/model"
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

func (this *categoryRepository) QueryCnd(db *gorm.DB, cnd *simple.QueryCnd) (list []model.Category, err error) {
	err = cnd.DoQuery(db).Find(&list).Error
	return
}

func (this *categoryRepository) Query(db *gorm.DB, queries *simple.ParamQueries) (list []model.Category, paging *simple.Paging) {
	queries.StartQuery(db).Find(&list)
	queries.StartCount(db).Model(&model.Category{}).Count(&queries.Paging.Total)
	paging = queries.Paging
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
	return this.QueryCnd(simple.GetDB(), simple.NewQueryCnd("status = ?", model.CategoryStatusOk))
}
