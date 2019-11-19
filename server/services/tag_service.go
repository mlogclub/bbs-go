package services

import (
	"math"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
	"github.com/mlogclub/bbs-go/services/cache"
)

type ScanTagCallback func(tags []model.Tag) bool

var TagService = newTagService()

func newTagService() *tagService {
	return &tagService{}
}

type tagService struct {
}

func (this *tagService) Get(id int64) *model.Tag {
	return repositories.TagRepository.Get(simple.DB(), id)
}

func (this *tagService) Take(where ...interface{}) *model.Tag {
	return repositories.TagRepository.Take(simple.DB(), where...)
}

func (this *tagService) Find(cnd *simple.SqlCnd) []model.Tag {
	return repositories.TagRepository.Find(simple.DB(), cnd)
}

func (this *tagService) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.Tag) {
	cnd.FindOne(db, &ret)
	return
}

func (this *tagService) FindPageByParams(params *simple.QueryParams) (list []model.Tag, paging *simple.Paging) {
	return repositories.TagRepository.FindPageByParams(simple.DB(), params)
}

func (this *tagService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.Tag, paging *simple.Paging) {
	return repositories.TagRepository.FindPageByCnd(simple.DB(), cnd)
}

func (this *tagService) Create(t *model.Tag) error {
	return repositories.TagRepository.Create(simple.DB(), t)
}

func (this *tagService) Update(t *model.Tag) error {
	if err := repositories.TagRepository.Update(simple.DB(), t); err != nil {
		return err
	}
	cache.TagCache.Invalidate(t.Id)
	return nil
}

// func (this *tagService) Updates(id int64, columns map[string]interface{}) error {
// 	return repositories.TagRepository.Updates(simple.DB(), id, columns)
// }
//
// func (this *tagService) UpdateColumn(id int64, name string, value interface{}) error {
// 	return repositories.TagRepository.UpdateColumn(simple.DB(), id, name, value)
// }
//
// func (this *tagService) Delete(id int64) {
// 	repositories.TagRepository.Delete(simple.DB(), id)
// }

// 自动完成
func (this *tagService) Autocomplete(input string) []model.Tag {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return nil
	}
	return repositories.TagRepository.Find(simple.DB(), simple.NewSqlCnd().Where("status = ? and name like ?",
		model.TagStatusOk, "%"+input+"%").Limit(6))
}

func (this *tagService) GetOrCreate(name string) (*model.Tag, error) {
	return repositories.TagRepository.GetOrCreate(simple.DB(), name)
}

func (this *tagService) GetByName(name string) *model.Tag {
	return repositories.TagRepository.GetByName(name)
}

func (this *tagService) GetTags() []model.TagResponse {
	list := repositories.TagRepository.Find(simple.DB(), simple.NewSqlCnd().Where("status = ?", model.TagStatusOk))

	var tags []model.TagResponse
	for _, tag := range list {
		tags = append(tags, model.TagResponse{TagId: tag.Id, TagName: tag.Name})
	}
	return tags
}

func (this *tagService) GetTagInIds(tagIds []int64) []model.Tag {
	return repositories.TagRepository.GetTagInIds(tagIds)
}

// 倒序扫描
func (this *tagService) ScanDesc(cb ScanTagCallback) {
	var cursor int64 = math.MaxInt64
	for {
		list := repositories.TagRepository.Find(simple.DB(), simple.NewSqlCnd().Where("id < ?", cursor).Desc("id").Limit(100))
		if list == nil || len(list) == 0 {
			break
		}
		cursor = list[len(list)-1].Id
		if !cb(list) {
			break
		}
	}
}
