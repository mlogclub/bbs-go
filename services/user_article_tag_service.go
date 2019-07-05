package services

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
	"github.com/mlogclub/mlog/services/cache"
	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
)

var UserArticleTagService = newUserArticleTagService()

func newUserArticleTagService() *userArticleTagService {
	return &userArticleTagService{}
}

type userArticleTagService struct {
}

func (this *userArticleTagService) Get(id int64) *model.UserArticleTag {
	return repositories.UserArticleTagRepository.Get(simple.GetDB(), id)
}

func (this *userArticleTagService) Take(where ...interface{}) *model.UserArticleTag {
	return repositories.UserArticleTagRepository.Take(simple.GetDB(), where...)
}

func (this *userArticleTagService) QueryCnd(cnd *simple.QueryCnd) (list []model.UserArticleTag, err error) {
	return repositories.UserArticleTagRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *userArticleTagService) Query(queries *simple.ParamQueries) (list []model.UserArticleTag, paging *simple.Paging) {
	return repositories.UserArticleTagRepository.Query(simple.GetDB(), queries)
}

func (this *userArticleTagService) Create(t *model.UserArticleTag) error {
	return repositories.UserArticleTagRepository.Create(simple.GetDB(), t)
}

func (this *userArticleTagService) Update(t *model.UserArticleTag) error {
	return repositories.UserArticleTagRepository.Update(simple.GetDB(), t)
}

func (this *userArticleTagService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.UserArticleTagRepository.Updates(simple.GetDB(), id, columns)
}

func (this *userArticleTagService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.UserArticleTagRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *userArticleTagService) Delete(id int64) {
	repositories.UserArticleTagRepository.Delete(simple.GetDB(), id)
}

func (this *userArticleTagService) GetBy(userId, tagId int64) *model.UserArticleTag {
	return repositories.UserArticleTagRepository.Take(simple.GetDB(), "user_id = ? and tag_id = ?", userId, tagId)
}

func (this *userArticleTagService) GetUserTags(userId int64) (tags []model.Tag) {
	list, err := repositories.UserArticleTagRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("user_id = ?", userId).Order("id desc"))
	if err != nil {
		logrus.Error(err)
		return
	}
	for _, userArticleTag := range list {
		tag := cache.TagCache.Get(userArticleTag.TagId)
		if tag != nil {
			tags = append(tags, *tag)
		}
	}
	return
}

func (this *userArticleTagService) AddUserTag(userId int64, name string) error {
	return simple.Tx(simple.GetDB(), func(tx *gorm.DB) error {
		tag, err := repositories.TagRepository.GetOrCreate(tx, name)
		if err != nil {
			return err
		}

		userArticleTag := this.GetBy(userId, tag.Id)
		if userArticleTag != nil {
			return errors.New("标签已存在")
		}
		userArticleTag = &model.UserArticleTag{
			UserId: userId,
			TagId:  tag.Id,
		}
		return repositories.UserArticleTagRepository.Create(tx, userArticleTag)
	})
}
