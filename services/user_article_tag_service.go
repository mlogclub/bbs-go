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

type UserArticleTagService struct {
	UserArticleTagRepository *repositories.UserArticleTagRepository
	TagRepository            *repositories.TagRepository
}

func NewUserArticleTagService() *UserArticleTagService {
	return &UserArticleTagService{
		UserArticleTagRepository: repositories.NewUserArticleTagRepository(),
	}
}

func (this *UserArticleTagService) Get(id int64) *model.UserArticleTag {
	return this.UserArticleTagRepository.Get(simple.GetDB(), id)
}

func (this *UserArticleTagService) Take(where ...interface{}) *model.UserArticleTag {
	return this.UserArticleTagRepository.Take(simple.GetDB(), where...)
}

func (this *UserArticleTagService) QueryCnd(cnd *simple.QueryCnd) (list []model.UserArticleTag, err error) {
	return this.UserArticleTagRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *UserArticleTagService) Query(queries *simple.ParamQueries) (list []model.UserArticleTag, paging *simple.Paging) {
	return this.UserArticleTagRepository.Query(simple.GetDB(), queries)
}

func (this *UserArticleTagService) Create(t *model.UserArticleTag) error {
	return this.UserArticleTagRepository.Create(simple.GetDB(), t)
}

func (this *UserArticleTagService) Update(t *model.UserArticleTag) error {
	return this.UserArticleTagRepository.Update(simple.GetDB(), t)
}

func (this *UserArticleTagService) Updates(id int64, columns map[string]interface{}) error {
	return this.UserArticleTagRepository.Updates(simple.GetDB(), id, columns)
}

func (this *UserArticleTagService) UpdateColumn(id int64, name string, value interface{}) error {
	return this.UserArticleTagRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *UserArticleTagService) Delete(id int64) {
	this.UserArticleTagRepository.Delete(simple.GetDB(), id)
}

func (this *UserArticleTagService) GetBy(userId, tagId int64) *model.UserArticleTag {
	return this.UserArticleTagRepository.Take(simple.GetDB(), "user_id = ? and tag_id = ?", userId, tagId)
}

func (this *UserArticleTagService) GetUserTags(userId int64) (tags []model.Tag) {
	list, err := this.UserArticleTagRepository.QueryCnd(simple.GetDB(), simple.NewQueryCnd("user_id = ?", userId).Order("id desc"))
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

func (this *UserArticleTagService) AddUserTag(userId int64, name string) error {
	return simple.Tx(simple.GetDB(), func(tx *gorm.DB) error {
		tag, err := this.TagRepository.GetOrCreate(tx, name)
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
		return this.UserArticleTagRepository.Create(tx, userArticleTag)
	})
}
