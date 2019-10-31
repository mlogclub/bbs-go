package services

import (
	"github.com/mlogclub/simple"
	"strings"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
)

var TagService = newTagService()

func newTagService() *tagService {
	return &tagService{}
}

type tagService struct {
}

func (this *tagService) Get(id int64) *model.Tag {
	return repositories.TagRepository.Get(simple.GetDB(), id)
}

func (this *tagService) Take(where ...interface{}) *model.Tag {
	return repositories.TagRepository.Take(simple.GetDB(), where...)
}

func (this *tagService) QueryCnd(cnd *simple.SqlCnd) (list []model.Tag, err error) {
	return repositories.TagRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *tagService) Query(params *simple.QueryParams) (list []model.Tag, paging *simple.Paging) {
	return repositories.TagRepository.Query(simple.GetDB(), queries)
}

func (this *tagService) Create(t *model.Tag) error {
	return repositories.TagRepository.Create(simple.GetDB(), t)
}

func (this *tagService) Update(t *model.Tag) error {
	return repositories.TagRepository.Update(simple.GetDB(), t)
}

func (this *tagService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.TagRepository.Updates(simple.GetDB(), id, columns)
}

func (this *tagService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.TagRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *tagService) Delete(id int64) {
	repositories.TagRepository.Delete(simple.GetDB(), id)
}

// 自动完成
func (this *tagService) Autocomplete(input string) []model.Tag {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil
	}
	list, _ := repositories.TagRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("status = ? and name like ?",
		model.TagStatusOk, "%"+input+"%").Size(6))
	return list
}

func (this *tagService) GetOrCreate(name string) (*model.Tag, error) {
	return repositories.TagRepository.GetOrCreate(simple.GetDB(), name)
}

func (this *tagService) GetByName(name string) *model.Tag {
	return repositories.TagRepository.FindByName(name)
}

func (this *tagService) ListAll(categoryId int64) ([]model.Tag, error) {
	return repositories.TagRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("category_id = ? and status = ?", categoryId, model.TagStatusOk))
}

func (this *tagService) GetTags() []model.TagResponse {
	list, err := repositories.TagRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("status = ?", model.TagStatusOk))
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
	return repositories.TagRepository.GetTagInIds(tagIds)
}
