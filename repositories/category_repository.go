package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/model"
)

type CategoryRepository struct {
}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

func (this *CategoryRepository) Get(db *gorm.DB, id int64) *model.Category {
	ret := &model.Category{}
	if err := db.First(ret, "id = ?", id).Error; err != nil {
		return nil
	}
	return ret
}

func (this *CategoryRepository) Take(db *gorm.DB, where ...interface{}) *model.Category {
	ret := &model.Category{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (this *CategoryRepository) QueryCnd(db *gorm.DB, cnd *simple.QueryCnd) (list []model.Category, err error) {
	err = cnd.DoQuery(db).Find(&list).Error
	return
}

func (this *CategoryRepository) Query(db *gorm.DB, queries *simple.ParamQueries) (list []model.Category, paging *simple.Paging) {
	queries.StartQuery(db).Find(&list)
	queries.StartCount(db).Model(&model.Category{}).Count(&queries.Paging.Total)
	paging = queries.Paging
	return
}

func (this *CategoryRepository) Create(db *gorm.DB, t *model.Category) (err error) {
	err = db.Create(t).Error
	return
}

func (this *CategoryRepository) Update(db *gorm.DB, t *model.Category) (err error) {
	err = db.Save(t).Error
	return
}

func (this *CategoryRepository) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Category{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (this *CategoryRepository) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Category{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (this *CategoryRepository) Delete(db *gorm.DB, id int64) {
	db.Model(&model.Category{}).Delete("id", id)
}

func (this *CategoryRepository) GetCategories() ([]model.Category, error) {
	return this.QueryCnd(simple.GetDB(), simple.NewQueryCnd("status = ?", model.CategoryStatusOk))
}
