package services

import (
	"bbs-go/model/constants"
	"strings"

	"github.com/mlogclub/simple"

	"bbs-go/cache"
	"bbs-go/model"
	"bbs-go/repositories"
)

var TagService = newTagService()

func newTagService() *tagService {
	return &tagService{}
}

type tagService struct {
}

func (s *tagService) Get(id int64) *model.Tag {
	return repositories.TagRepository.Get(simple.DB(), id)
}

func (s *tagService) Take(where ...interface{}) *model.Tag {
	return repositories.TagRepository.Take(simple.DB(), where...)
}

func (s *tagService) Find(cnd *simple.SqlCnd) []model.Tag {
	return repositories.TagRepository.Find(simple.DB(), cnd)
}

func (s *tagService) FindOne(cnd *simple.SqlCnd) *model.Tag {
	return repositories.TagRepository.FindOne(simple.DB(), cnd)
}

func (s *tagService) FindPageByParams(params *simple.QueryParams) (list []model.Tag, paging *simple.Paging) {
	return repositories.TagRepository.FindPageByParams(simple.DB(), params)
}

func (s *tagService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.Tag, paging *simple.Paging) {
	return repositories.TagRepository.FindPageByCnd(simple.DB(), cnd)
}

func (s *tagService) Create(t *model.Tag) error {
	return repositories.TagRepository.Create(simple.DB(), t)
}

func (s *tagService) Update(t *model.Tag) error {
	if err := repositories.TagRepository.Update(simple.DB(), t); err != nil {
		return err
	}
	cache.TagCache.Invalidate(t.Id)
	return nil
}

// func (s *tagService) Updates(id int64, columns map[string]interface{}) error {
// 	return repositories.TagRepository.Updates(simple.DB(), id, columns)
// }
//
// func (s *tagService) UpdateColumn(id int64, name string, value interface{}) error {
// 	return repositories.TagRepository.UpdateColumn(simple.DB(), id, name, value)
// }
//
// func (s *tagService) Delete(id int64) {
// 	repositories.TagRepository.Delete(simple.DB(), id)
// }

// 自动完成
func (s *tagService) Autocomplete(input string) []model.Tag {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil
	}
	return repositories.TagRepository.Find(simple.DB(), simple.NewSqlCnd().Where("status = ? and name like ?",
		constants.StatusOk, "%"+input+"%").Limit(6))
}

func (s *tagService) GetOrCreate(name string) (*model.Tag, error) {
	return repositories.TagRepository.GetOrCreate(simple.DB(), name)
}

func (s *tagService) GetByName(name string) *model.Tag {
	return repositories.TagRepository.GetByName(name)
}

func (s *tagService) GetTags() []model.TagResponse {
	list := repositories.TagRepository.Find(simple.DB(), simple.NewSqlCnd().Where("status = ?", constants.StatusOk))

	var tags []model.TagResponse
	for _, tag := range list {
		tags = append(tags, model.TagResponse{TagId: tag.Id, TagName: tag.Name})
	}
	return tags
}

func (s *tagService) GetTagInIds(tagIds []int64) []model.Tag {
	return repositories.TagRepository.GetTagInIds(tagIds)
}

// 扫描
func (s *tagService) Scan(callback func(tags []model.Tag)) {
	var cursor int64
	for {
		list := repositories.TagRepository.Find(simple.DB(), simple.NewSqlCnd().Where("id > ?", cursor).Asc("id").Limit(100))
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}
