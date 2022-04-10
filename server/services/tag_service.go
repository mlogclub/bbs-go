package services

import (
	"bbs-go/model/constants"
	"strings"

	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"

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
	return repositories.TagRepository.Get(sqls.DB(), id)
}

func (s *tagService) Take(where ...interface{}) *model.Tag {
	return repositories.TagRepository.Take(sqls.DB(), where...)
}

func (s *tagService) Find(cnd *sqls.Cnd) []model.Tag {
	return repositories.TagRepository.Find(sqls.DB(), cnd)
}

func (s *tagService) FindOne(cnd *sqls.Cnd) *model.Tag {
	return repositories.TagRepository.FindOne(sqls.DB(), cnd)
}

func (s *tagService) FindPageByParams(params *params.QueryParams) (list []model.Tag, paging *sqls.Paging) {
	return repositories.TagRepository.FindPageByParams(sqls.DB(), params)
}

func (s *tagService) FindPageByCnd(cnd *sqls.Cnd) (list []model.Tag, paging *sqls.Paging) {
	return repositories.TagRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *tagService) Create(t *model.Tag) error {
	return repositories.TagRepository.Create(sqls.DB(), t)
}

func (s *tagService) Update(t *model.Tag) error {
	if err := repositories.TagRepository.Update(sqls.DB(), t); err != nil {
		return err
	}
	cache.TagCache.Invalidate(t.Id)
	return nil
}

// func (s *tagService) Updates(id int64, columns map[string]interface{}) error {
// 	return repositories.TagRepository.Updates(sqls.DB(), id, columns)
// }
//
// func (s *tagService) UpdateColumn(id int64, name string, value interface{}) error {
// 	return repositories.TagRepository.UpdateColumn(sqls.DB(), id, name, value)
// }
//
// func (s *tagService) Delete(id int64) {
// 	repositories.TagRepository.Delete(sqls.DB(), id)
// }

// 自动完成
func (s *tagService) Autocomplete(input string) []model.Tag {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil
	}
	return repositories.TagRepository.Find(sqls.DB(), sqls.NewCnd().Where("status = ? and name like ?",
		constants.StatusOk, "%"+input+"%").Limit(6))
}

func (s *tagService) GetOrCreate(name string) (*model.Tag, error) {
	return repositories.TagRepository.GetOrCreate(sqls.DB(), name)
}

func (s *tagService) GetByName(name string) *model.Tag {
	return repositories.TagRepository.GetByName(name)
}

func (s *tagService) GetTags() []model.TagResponse {
	list := repositories.TagRepository.Find(sqls.DB(), sqls.NewCnd().Where("status = ?", constants.StatusOk))

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
		list := repositories.TagRepository.Find(sqls.DB(), sqls.NewCnd().Where("id > ?", cursor).Asc("id").Limit(100))
		if len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		callback(list)
	}
}
