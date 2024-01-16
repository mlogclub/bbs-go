package services

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/event"
	"errors"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web/params"
	"gorm.io/gorm"

	"bbs-go/internal/models"
	"bbs-go/internal/repositories"
)

var UserLikeService = newUserLikeService()

func newUserLikeService() *userLikeService {
	return &userLikeService{}
}

type userLikeService struct {
}

func (s *userLikeService) Get(id int64) *models.UserLike {
	return repositories.UserLikeRepository.Get(sqls.DB(), id)
}

func (s *userLikeService) Take(where ...interface{}) *models.UserLike {
	return repositories.UserLikeRepository.Take(sqls.DB(), where...)
}

func (s *userLikeService) Find(cnd *sqls.Cnd) []models.UserLike {
	return repositories.UserLikeRepository.Find(sqls.DB(), cnd)
}

func (s *userLikeService) FindOne(cnd *sqls.Cnd) *models.UserLike {
	return repositories.UserLikeRepository.FindOne(sqls.DB(), cnd)
}

func (s *userLikeService) FindPageByParams(params *params.QueryParams) (list []models.UserLike, paging *sqls.Paging) {
	return repositories.UserLikeRepository.FindPageByParams(sqls.DB(), params)
}

func (s *userLikeService) FindPageByCnd(cnd *sqls.Cnd) (list []models.UserLike, paging *sqls.Paging) {
	return repositories.UserLikeRepository.FindPageByCnd(sqls.DB(), cnd)
}

func (s *userLikeService) Create(t *models.UserLike) error {
	return repositories.UserLikeRepository.Create(sqls.DB(), t)
}

func (s *userLikeService) Update(t *models.UserLike) error {
	return repositories.UserLikeRepository.Update(sqls.DB(), t)
}

func (s *userLikeService) Updates(id int64, columns map[string]interface{}) error {
	return repositories.UserLikeRepository.Updates(sqls.DB(), id, columns)
}

func (s *userLikeService) UpdateColumn(id int64, name string, value interface{}) error {
	return repositories.UserLikeRepository.UpdateColumn(sqls.DB(), id, name, value)
}

func (s *userLikeService) Delete(id int64) {
	repositories.UserLikeRepository.Delete(sqls.DB(), id)
}

// 统计数量
func (s *userLikeService) Count(entityType string, entityId int64) int64 {
	var count int64 = 0
	sqls.DB().Model(&models.UserLike{}).Where("entity_id = ?", entityId).Where("entity_type = ?", entityType).Count(&count)
	return count
}

// 最近点赞
func (s *userLikeService) Recent(entityType string, entityId int64, count int) []models.UserLike {
	return s.Find(sqls.NewCnd().Eq("entity_id", entityId).Eq("entity_type", entityType).Desc("id").Limit(count))
}

// Exists 是否点赞
func (s *userLikeService) Exists(userId int64, entityType string, entityId int64) bool {
	return repositories.UserLikeRepository.FindOne(sqls.DB(), sqls.NewCnd().Eq("user_id", userId).
		Eq("entity_id", entityId).Eq("entity_type", entityType)) != nil
}

// 是否点赞，返回已点赞实体编号
func (s *userLikeService) IsLiked(userId int64, entityType string, entityIds []int64) (likedEntityIds []int64) {
	list := repositories.UserLikeRepository.Find(sqls.DB(), sqls.NewCnd().Eq("user_id", userId).
		In("entity_id", entityIds).Eq("entity_type", entityType))
	for _, like := range list {
		likedEntityIds = append(likedEntityIds, like.EntityId)
	}
	return
}

// TopicLike 话题点赞
func (s *userLikeService) TopicLike(userId int64, topicId int64) error {
	topic := repositories.TopicRepository.Get(sqls.DB(), topicId)
	if topic == nil || topic.Status != constants.StatusOk {
		return errors.New("话题不存在")
	}

	if err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := s.like(tx, userId, constants.EntityTopic, topicId); err != nil {
			return err
		}
		// 更新点赞数
		return tx.Exec("update t_topic set like_count = like_count + 1 where id = ?", topicId).Error
	}); err != nil {
		return err
	}

	// 发送事件
	event.Send(event.UserLikeEvent{
		UserId:     userId,
		EntityId:   topicId,
		EntityType: constants.EntityTopic,
	})

	return nil
}

func (s *userLikeService) TopicUnLike(userId int64, topicId int64) error {
	topic := repositories.TopicRepository.Get(sqls.DB(), topicId)
	if topic == nil || topic.Status != constants.StatusOk {
		return errors.New("话题不存在")
	}

	if err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := s.unlike(tx, userId, constants.EntityTopic, topicId); err != nil {
			return err
		}
		// 更新点赞数
		return repositories.TopicRepository.UpdateColumn(tx, topicId, "like_count", gorm.Expr("like_count - 1"))
	}); err != nil {
		return err
	}

	// 发送事件
	event.Send(event.UserUnLikeEvent{
		UserId:     userId,
		EntityId:   topicId,
		EntityType: constants.EntityTopic,
	})

	return nil
}

func (s *userLikeService) ArticleLike(userId int64, articleId int64) error {
	article := repositories.ArticleRepository.Get(sqls.DB(), articleId)
	if article == nil || article.Status != constants.StatusOk {
		return errors.New("文章不存在")
	}
	if err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := s.like(tx, userId, constants.EntityArticle, articleId); err != nil {
			return err
		}
		// 更新点赞数
		return repositories.ArticleRepository.UpdateColumn(tx, articleId, "like_count", gorm.Expr("like_count - 1"))
	}); err != nil {
		return err
	}

	// 发送事件
	event.Send(event.UserLikeEvent{
		UserId:     userId,
		EntityId:   articleId,
		EntityType: constants.EntityArticle,
	})
	return nil
}

func (s *userLikeService) ArticleUnLike(userId int64, articleId int64) error {
	article := repositories.ArticleRepository.Get(sqls.DB(), articleId)
	if article == nil || article.Status != constants.StatusOk {
		return errors.New("文章不存在")
	}

	if err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := s.unlike(tx, userId, constants.EntityArticle, articleId); err != nil {
			return err
		}
		// 更新点赞数
		return repositories.ArticleRepository.UpdateColumn(tx, articleId, "like_count", gorm.Expr("like_count - 1"))
	}); err != nil {
		return err
	}

	// 发送事件
	event.Send(event.UserUnLikeEvent{
		UserId:     userId,
		EntityId:   articleId,
		EntityType: constants.EntityArticle,
	})

	return nil
}

// CommentLike 话题点赞
func (s *userLikeService) CommentLike(userId int64, commentId int64) error {
	comment := repositories.CommentRepository.Get(sqls.DB(), commentId)
	if comment == nil || comment.Status != constants.StatusOk {
		return errors.New("评论不存在")
	}

	if err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := s.like(tx, userId, constants.EntityComment, commentId); err != nil {
			return err
		}
		// 更新点赞数
		return repositories.CommentRepository.UpdateColumn(tx, commentId, "like_count", gorm.Expr("like_count + 1"))
	}); err != nil {
		return err
	}

	// 发送事件
	event.Send(event.UserLikeEvent{
		UserId:     userId,
		EntityId:   commentId,
		EntityType: constants.EntityComment,
	})

	return nil
}

// CommentLike 话题点赞
func (s *userLikeService) CommentUnLike(userId int64, commentId int64) error {
	comment := repositories.CommentRepository.Get(sqls.DB(), commentId)
	if comment == nil || comment.Status != constants.StatusOk {
		return errors.New("评论不存在")
	}

	if err := sqls.DB().Transaction(func(tx *gorm.DB) error {
		if err := s.unlike(tx, userId, constants.EntityComment, commentId); err != nil {
			return err
		}
		// 更新点赞数
		return repositories.CommentRepository.UpdateColumn(tx, commentId, "like_count", gorm.Expr("like_count - 1"))
	}); err != nil {
		return err
	}

	// 发送事件
	event.Send(event.UserUnLikeEvent{
		UserId:     userId,
		EntityId:   commentId,
		EntityType: constants.EntityComment,
	})

	return nil
}

func (s *userLikeService) like(tx *gorm.DB, userId int64, entityType string, entityId int64) error {
	// 判断是否已经点赞了
	if s.Exists(userId, entityType, entityId) {
		return errors.New("已点赞")
	}
	// 点赞
	return repositories.UserLikeRepository.Create(tx, &models.UserLike{
		UserId:     userId,
		EntityType: entityType,
		EntityId:   entityId,
		CreateTime: dates.NowTimestamp(),
	})
}

func (s userLikeService) unlike(tx *gorm.DB, userId int64, entityType string, entityId int64) error {
	return tx.Delete(&models.UserLike{}, "user_id = ? and entity_id = ? and entity_type = ?", userId, entityId, entityType).Error
}
