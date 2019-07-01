package services

import (
	"errors"
	"github.com/mlogclub/mlog/model"
	"github.com/mlogclub/mlog/repositories"
	"github.com/mlogclub/simple"
)

type FavoriteService struct {
	FavoriteRepository *repositories.FavoriteRepository
	ArticleRepository  *repositories.ArticleRepository
	TopicRepository    *repositories.TopicRepository
}

func NewFavoriteService() *FavoriteService {
	return &FavoriteService{
		FavoriteRepository: repositories.NewFavoriteRepository(),
	}
}

func (this *FavoriteService) Get(id int64) *model.Favorite {
	return this.FavoriteRepository.Get(simple.GetDB(), id)
}

func (this *FavoriteService) Take(where ...interface{}) *model.Favorite {
	return this.FavoriteRepository.Take(simple.GetDB(), where...)
}

func (this *FavoriteService) QueryCnd(cnd *simple.QueryCnd) (list []model.Favorite, err error) {
	return this.FavoriteRepository.QueryCnd(simple.GetDB(), cnd)
}

func (this *FavoriteService) Query(queries *simple.ParamQueries) (list []model.Favorite, paging *simple.Paging) {
	return this.FavoriteRepository.Query(simple.GetDB(), queries)
}

func (this *FavoriteService) Create(t *model.Favorite) error {
	return this.FavoriteRepository.Create(simple.GetDB(), t)
}

func (this *FavoriteService) Update(t *model.Favorite) error {
	return this.FavoriteRepository.Update(simple.GetDB(), t)
}

func (this *FavoriteService) Updates(id int64, columns map[string]interface{}) error {
	return this.FavoriteRepository.Updates(simple.GetDB(), id, columns)
}

func (this *FavoriteService) UpdateColumn(id int64, name string, value interface{}) error {
	return this.FavoriteRepository.UpdateColumn(simple.GetDB(), id, name, value)
}

func (this *FavoriteService) Delete(id int64) {
	this.FavoriteRepository.Delete(simple.GetDB(), id)
}

func (this *FavoriteService) GetBy(entityType string, entityId int64) *model.Favorite {
	return this.FavoriteRepository.Take(simple.GetDB(), "entity_type = ? and entity_id = ?", entityType, entityId)
}

// 收藏文章
func (this *FavoriteService) AddArticleFavorite(userId, articleId int64) error {
	article := this.ArticleRepository.Get(simple.GetDB(), articleId)
	if article == nil || article.Status != model.ArticleStatusPublished {
		return errors.New("收藏的文章不存在")
	}
	temp := this.GetBy(model.EntityTypeArticle, articleId)
	if temp != nil { // 已经收藏
		return nil
	}
	return this.FavoriteRepository.Create(simple.GetDB(), &model.Favorite{
		UserId:     userId,
		EntityType: model.EntityTypeArticle,
		EntityId:   articleId,
		CreateTime: simple.NowTimestamp(),
	})
}

// 收藏主题
func (this *FavoriteService) AddTopicFavorite(userId, topicId int64) error {
	topic := this.TopicRepository.Get(simple.GetDB(), topicId)
	if topic == nil || topic.Status != model.TopicStatusOk {
		return errors.New("收藏的话题不存在")
	}
	temp := this.GetBy(model.EntityTypeTopic, topicId)
	if temp != nil { // 已经收藏
		return nil
	}
	return this.FavoriteRepository.Create(simple.GetDB(), &model.Favorite{
		UserId:     userId,
		EntityType: model.EntityTypeTopic,
		EntityId:   topicId,
		CreateTime: simple.NowTimestamp(),
	})
}
