package services

import (
	"github.com/mlogclub/simple"
	"strings"

	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
)

var TagService = newTagService()

func newTagService() *tagService {
	return &tagService{
		TagRepository:      repositories.NewTagRepository(),
		CategoryRepository: repositories.NewCategoryRepository(),
	}
}

type tagService struct {
	TagRepository      *repositories.TagRepository
	CategoryRepository *repositories.CategoryRepository
}

func (this *tagService) Get(id int64) *model.Tag {
	return this.TagRepository.Get(simple.GetDB(), id)
}

func (this *tagService) Take(where ...interface{}) *model.Tag {
	return this.TagRepository.Take(simple.GetDB(), where...)
}

func (this *tagService) QueryCnd(cnd *simple.QueryCnd) (list []model.Tag, err error) {
	return this.TagRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *tagService) Query(queries *simple.ParamQueries) (list []model.Tag, paging *simple.Paging) {
	return this.TagRepository.Query(simple.GetDB(), queries)
}

func (this *tagService) Create(t *model.Tag) error {
	return this.TagRepository.Create(simple.GetDB(), t)
}

func (this *tagService) Update(t *model.Tag) error {
	return this.TagRepository.Update(simple.GetDB(), t)
}

func (this *tagService) Updates(id int64, columns map[string]interface{}) error {
	return this.TagRepository.Updates(simple.GetDB(), id, columns)
}

func (this *tagService) UpdateColumn(id int64, name string, value interface{}) error {
	return this.TagRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *tagService) Delete(id int64) {
	this.TagRepository.Delete(simple.GetDB(), id)
}

// 自动完成
func (this *tagService) Autocomplete(input string) []model.Tag {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil
	}
	list, _ := this.TagRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("status = ? and name like ?",
		model.TagStatusOk, "%"+input+"%").Size(6))
	return list
}

func (this *tagService) GetOrCreate(name string) (*model.Tag, error) {
	return this.TagRepository.GetOrCreate(simple.GetDB(), name)
}

func (this *tagService) GetByName(name string) *model.Tag {
	return this.TagRepository.FindByName(name)
}

func (this *tagService) ListAll(categoryId int64) ([]model.Tag, error) {
	return this.TagRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("category_id = ? and status = ?", categoryId, model.TagStatusOk))
}

func (this *tagService) GetTags() []model.TagResponse {
	list, err := this.TagRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("status = ?", model.TagStatusOk))
	if err != nil {
		return nil
	}

	var tags []model.TagResponse
	for _, tag := range list {
		tags = append(tags, model.TagResponse{TagId: tag.Id, TagName: tag.Name})
	}
	return tags
}

func (this *tagService) GetTagInIds(tagIds []int64) []model.Tag {
	return this.TagRepository.GetTagInIds(tagIds)
}
