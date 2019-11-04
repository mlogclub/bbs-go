package services

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/mlogclub/simple"

	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/repositories"
)

var FavoriteService = newFavoriteService()

func newFavoriteService() *favoriteService {
	return &favoriteService{
	}
}

type favoriteService struct {
}

func (this *favoriteService) Get(id int64) *model.Favorite {
	return repositories.FavoriteRepository.Get(simple.DB(), id)
}

func (this *favoriteService) Take(where ...interface{}) *model.Favorite {
	return repositories.FavoriteRepository.Take(simple.DB(), where...)
}

func (this *favoriteService) Find(cnd *simple.SqlCnd) []model.Favorite {
	return repositories.FavoriteRepository.Find(simple.DB(), cnd)
}

func (this *favoriteService) FindOne(db *gorm.DB, cnd *simple.SqlCnd) (ret *model.Favorite) {
	cnd.FindOne(db, &ret)
	return
}

func (this *favoriteService) FindPageByParams(params *simple.QueryParams) (list []model.Favorite, paging *simple.Paging) {
	return repositories.FavoriteRepository.FindPageByParams(simple.DB(), params)
}

func (this *favoriteService) FindPageByCnd(cnd *simple.SqlCnd) (list []model.Favorite, paging *simple.Paging) {
	return repositories.FavoriteRepository.FindPageByCnd(simple.DB(), cnd)
}

func (this *favoriteService) Create(t *model.Favorite) error {
	return repositories.FavoriteRepository.Create(simple.DB(), t)
}

func (this *favoriteService) Update(t *model.Favorite) error {
	return repositories.FavoriteRepository.Update(simple.DB(), t)
}

func (this *favoriteService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.FavoriteRepository.Updates(simple.DB(), id, columns)
}

func (this *favoriteService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.FavoriteRepository.UpdateColumn(simple.DB(), id, name, value)
}

func (this *favoriteService) Delete(id int64) {
	repositories.FavoriteRepository.Delete(simple.DB(), id)
}

func (this *favoriteService) GetBy(userId int64, entityType string, entityId int64) *model.Favorite {
	return repositories.FavoriteRepository.Take(simple.DB(), "user_id = ? and entity_type = ? and entity_id = ?",
		userId, entityType, entityId)
}

// 收藏文章
func (this *favoriteService) AddArticleFavorite(userId, articleId int64) error {
	article := repositories.ArticleRepository.Get(simple.DB(), articleId)
	if article == nil || article.Status != model.ArticleStatusPublished {
		return errors.New("收藏的文章不存在")
	}
	return this.addFavorite(userId, model.EntityTypeArticle, articleId)
}

// 收藏主题
func (this *favoriteService) AddTopicFavorite(userId, topicId int64) error {
	topic := repositories.TopicRepository.Get(simple.DB(), topicId)
	if topic == nil || topic.Status != model.TopicStatusOk {
		return errors.New("收藏的话题不存在")
	}
	return this.addFavorite(userId, model.EntityTypeTopic, topicId)
}

func (this *favoriteService) addFavorite(userId int64, entityType string, entityId int64) error {
	temp := this.GetBy(userId, entityType, entityId)
	if temp != nil { // 已经收藏
		return nil
	}
	return repositories.FavoriteRepository.Create(simple.DB(), &model.Favorite{
		UserId:     userId,
		EntityType: entityType,
		EntityId:   entityId,
		CreateTime: simple.NowTimestamp(),
	})
}
