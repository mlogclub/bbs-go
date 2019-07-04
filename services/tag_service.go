package services

import (
	"github.com/mlogclub/simple"
	"strings"

	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
)

type TagService struct {
	TagRepository      *repositories.TagRepository
	CategoryRepository *repositories.CategoryRepository
}

func NewTagService() *TagService {
	return &TagService{
		TagRepository:      repositories.NewTagRepository(),
		CategoryRepository: repositories.NewCategoryRepository(),
	}
}

func (this *TagService) Get(id int64) *model.Tag {
	return this.TagRepository.Get(simple.GetDB(), id)
}

func (this *TagService) Take(where ...interface{}) *model.Tag {
	return this.TagRepository.Take(simple.GetDB(), where...)
}

func (this *TagService) QueryCnd(cnd *simple.QueryCnd) (list []model.Tag, err error) {
	return this.TagRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *TagService) Query(queries *simple.ParamQueries) (list []model.Tag, paging *simple.Paging) {
	return this.TagRepository.Query(simple.GetDB(), queries)
}

func (this *TagService) Create(t *model.Tag) error {
	return this.TagRepository.Create(simple.GetDB(), t)
}

func (this *TagService) Update(t *model.Tag) error {
	return this.TagRepository.Update(simple.GetDB(), t)
}

func (this *TagService) Updates(id int64, columns map[string]interface{}) error {
	return this.TagRepository.Updates(simple.GetDB(), id, columns)
}

func (this *TagService) UpdateColumn(id int64, name string, value interface{}) error {
	return this.TagRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *TagService) Delete(id int64) {
	this.TagRepository.Delete(simple.GetDB(), id)
}

// 自动完成
func (this *TagService) Autocomplete(input string) []model.Tag {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil
	}
	list, _ := this.TagRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("status = ? and name like ?",
		model.TagStatusOk, "%"+input+"%").Size(6))
	return list
}

func (this *TagService) GetOrCreate(name string) (*model.Tag, error) {
	return this.TagRepository.GetOrCreate(simple.GetDB(), name)
}

func (this *TagService) GetByName(name string) *model.Tag {
	return this.TagRepository.FindByName(name)
}

func (this *TagService) ListAll(categoryId int64) ([]model.Tag, error) {
	return this.TagRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("category_id = ? and status = ?", categoryId, model.TagStatusOk))
}

func (this *TagService) GetTags() []model.TagResponse {
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

func (this *TagService) GetTagInIds(tagIds []int64) []model.Tag {
	return this.TagRepository.GetTagInIds(tagIds)
}
