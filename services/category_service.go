package services

import (
	"github.com/mlogclub/simple"

	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
)

type CategoryService struct {
	CategoryRepository *repositories.CategoryRepository
}

func NewCategoryService() *CategoryService {
	return &CategoryService{
		CategoryRepository: repositories.NewCategoryRepository(),
	}
}

func (this *CategoryService) Get(id int64) *model.Category {
	return this.CategoryRepository.Get(simple.GetDB(), id)
}

func (this *CategoryService) Take(where ...interface{}) *model.Category {
	return this.CategoryRepository.Take(simple.GetDB(), where...)
}

func (this *CategoryService) QueryCnd(cnd *simple.QueryCnd) (list []model.Category, err error) {
	return this.CategoryRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *CategoryService) Query(queries *simple.ParamQueries) (list []model.Category, paging *simple.Paging) {
	return this.CategoryRepository.Query(simple.GetDB(), queries)
}

func (this *CategoryService) Create(t *model.Category) error {
	return this.CategoryRepository.Create(simple.GetDB(), t)
}

func (this *CategoryService) Update(t *model.Category) error {
	return this.CategoryRepository.Update(simple.GetDB(), t)
}

func (this *CategoryService) Updates(id int64, columns map[string]interface{}) error {
	return this.CategoryRepository.Updates(simple.GetDB(), id, columns)
}

func (this *CategoryService) UpdateColumn(id int64, name string, value interface{}) error {
	return this.CategoryRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *CategoryService) Delete(id int64) {
	this.CategoryRepository.Delete(simple.GetDB(), id)
}

func (this *CategoryService) GetOrCreate(name string) *model.Category {
	category := this.FindByName(name)
	if category != nil {
		return category
	} else {
		category = &model.Category{
			Name:       name,
			Status:     model.CategoryStatusOk,
			CreateTime: simple.NowTimestamp(),
			UpdateTime: simple.NowTimestamp(),
		}
		this.Create(category)
		return category
	}
}

func (this *CategoryService) FindByName(name string) *model.Category {
	if len(name) == 0 {
		return nil
	}
	return this.Take("name = ?", name)
}

func (this *CategoryService) GetCategories() ([]model.Category, error) {
	return this.CategoryRepository.GetCategories()
}
