package services

import (
	"github.com/mlogclub/bbs-go/services/cache"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
)

var CategoryService = newCategoryService()

func newCategoryService() *categoryService {
	return &categoryService{}
}

type categoryService struct {
}

func (this *categoryService) Get(id int64) *model.Category {
	return repositories.CategoryRepository.Get(simple.GetDB(), id)
}

func (this *categoryService) Take(where ...interface{}) *model.Category {
	return repositories.CategoryRepository.Take(simple.GetDB(), where...)
}

func (this *categoryService) QueryCnd(cnd *simple.QueryCnd) (list []model.Category, err error) {
	return repositories.CategoryRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *categoryService) Query(queries *simple.ParamQueries) (list []model.Category, paging *simple.Paging) {
	return repositories.CategoryRepository.Query(simple.GetDB(), queries)
}

func (this *categoryService) Create(t *model.Category) error {
	err := repositories.CategoryRepository.Create(simple.GetDB(), t)
	if err == nil {
		cache.CategoryCache.Invalidate(t.Id)
		cache.CategoryCache.InvalidateAll()
	}
	return err
}

func (this *categoryService) Update(t *model.Category) error {
	err := repositories.CategoryRepository.Update(simple.GetDB(), t)
	if err == nil {
		cache.CategoryCache.Invalidate(t.Id)
		cache.CategoryCache.InvalidateAll()
	}
	return err
}

func (this *categoryService) Updates(id int64, columns map[string]interface{}) error {
	err := repositories.CategoryRepository.Updates(simple.GetDB(), id, columns)
	if err == nil {
		cache.CategoryCache.Invalidate(id)
		cache.CategoryCache.InvalidateAll()
	}
	return err
}

func (this *categoryService) UpdateColumn(id int64, name string, value interface{}) error {
	err := repositories.CategoryRepository.UpdateColumn(simple.GetDB(), id, name, value)
	if err == nil {
		cache.CategoryCache.Invalidate(id)
		cache.CategoryCache.InvalidateAll()
	}
	return err
}

func (this *categoryService) Delete(id int64) {
	repositories.CategoryRepository.Delete(simple.GetDB(), id)
	cache.CategoryCache.Invalidate(id)
	cache.CategoryCache.InvalidateAll()
}

func (this *categoryService) GetOrCreate(name string) (*model.Category, error) {
	category := this.FindByName(name)
	if category != nil {
		return category, nil
	} else {
		category = &model.Category{
			Name:       name,
			Status:     model.CategoryStatusOk,
			CreateTime: simple.NowTimestamp(),
			UpdateTime: simple.NowTimestamp(),
		}
		err := this.Create(category)
		if err != nil {
			return nil, err
		}
		return category, nil
	}
}

func (this *categoryService) FindByName(name string) *model.Category {
	if len(name) == 0 {
		return nil
	}
	return this.Take("name = ?", name)
}

func (this *categoryService) GetCategories() ([]model.Category, error) {
	return repositories.CategoryRepository.GetCategories()
}
